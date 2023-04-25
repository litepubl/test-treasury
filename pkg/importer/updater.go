package importer

import (
	"context"
	"errors"
	"sync"

	"github.com/rs/zerolog/log"
)

type Updater struct {
	importer DataImporter
	state    State
	mutex    sync.RWMutex
}

var (
	ErrStateAlreadyUpdating = errors.New("data is already updating")
	ErrNotUpdated           = errors.New("data has been not updated")
)

func NewUpdater(importer DataImporter) *Updater {
	return &Updater{
		importer: importer,
		state:    Empty,
	}
}

func (updater *Updater) Update(ctx context.Context) error {
	if updater.GetState() == Updating {
		return ErrStateAlreadyUpdating
	}

	updater.setState(Updating)

	signal := make(chan any, 1)
	go func() {
		err := updater.importer.Import(ctx)
		if err == nil {
			updater.setState(Ok)
		} else {
			updater.setState(Empty)
			log.Info().Err(err).Msg("Updater: importer return error")
		}

		signal <- nil
	}()

	select {
	case <-ctx.Done():
		updater.setState(Empty)
		log.Info().Msg("Updater: context done")

	case <-signal:
		log.Info().Msg("Updater: importer finished")
	}

	if updater.GetState() == Empty {
		return ErrNotUpdated
	}

	return nil
}

func (updater *Updater) GetState() State {
	updater.mutex.RLock()
	defer updater.mutex.RUnlock()

	return updater.state
}

func (updater *Updater) setState(state State) {
	log.Info().Msg("Updater: state changed to " + state.String())
	updater.mutex.Lock()
	defer updater.mutex.Unlock()

	updater.state = state
}
