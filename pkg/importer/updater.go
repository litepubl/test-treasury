package importer

import (
	"context"
	"errors"
	"sync"

	"github.com/rs/zerolog/log"
)

// Updater  отвечает за обновление и состояние данных
type Updater struct {
	importer   DataImporter
	state      State
	stateMutex sync.RWMutex
}

var (
	// ErrStateAlreadyUpdating ошибка возникает, когда во время обновления пытаются  еще раз вызвать метод обновления
	ErrStateAlreadyUpdating = errors.New("data is already updating")
	ErrNotUpdated           = errors.New("data has been not updated")
)

// NewUpdater конструктор
func NewUpdater(importer DataImporter) *Updater {
	return &Updater{
		importer: importer,
		state:    Empty,
	}
}

// Update обновляет данные в БД из xml
func (u *Updater) Update(ctx context.Context) error {
	if u.State() == Updating {
		return ErrStateAlreadyUpdating
	}

	u.setState(Updating)

	importerFinished := func() <-chan State {
		out := make(chan State, 1)
		go func() {
			defer close(out)
			err := u.importer.Import(ctx)
			if err != nil {
				out <- Empty
				log.Info().Err(err).Msg("Updater: importer return error")
				return
			}

			out <- Ok
		}()

		return out
	}()

	select {
	case <-ctx.Done():
		u.setState(Empty)
		log.Info().Msg("Updater: context done")

	case state := <-importerFinished:
		u.setState(state)
		log.Info().Msg("Updater: importer finished with result " + state.String())
	}

	if u.State() == Empty {
		return ErrNotUpdated
	}

	return nil
}

// State возвращает текущее состояние, потокобезопасен
func (u *Updater) State() State {
	u.stateMutex.RLock()
	defer u.stateMutex.RUnlock()

	return u.state
}

func (u *Updater) setState(state State) {
	u.stateMutex.Lock()
	defer u.stateMutex.Unlock()

	u.state = state
	log.Info().Msg("Updater: state changed to " + state.String())
}
