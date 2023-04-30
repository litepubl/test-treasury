package importer

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/litepubl/test-treasury/pkg/entity"
)

const _individual = "Individual"

type Importer struct {
	repo       PersonRepo
	downloader Downloader
}

func NewImporter(r PersonRepo, d Downloader) *Importer {
	return &Importer{
		repo:       r,
		downloader: d,
	}
}

func (importer *Importer) Import(ctx context.Context) error {
	persons, err := importer.repo.GetAll(context.Background())
	if err != nil {
		return fmt.Errorf("Importer - Person - s.repo.GetAll: %w", err)
	}

	sdnList, err := importer.downloader.Download(ctx)
	if err != nil {
		return fmt.Errorf("Importer - Update - downloader.Download: %w", err)
	}

	for _, sdnEntry := range sdnList.SdnEntry {
		if sdnEntry.SdnType == _individual {
			uid, err := strconv.Atoi(sdnEntry.Uid)
			if err != nil {
				continue
			}

			person := &entity.Person{
				Uid:       uid,
				FirstName: strings.Title(sdnEntry.FirstName),
				LastName:  strings.Title(sdnEntry.LastName),
			}

			err = importer.updateEntity(ctx, persons, person)
			if err != nil {
				return fmt.Errorf("Importer - Update - s.repo.Store: %w", err)
			}
		}
	}

	// оставшиеся элементы в persons не были найдены в xml и следовательно должны быть удалены из БД
	if len(persons) > 0 {
		idList := make([]int, 0, len(persons))
		for id := range persons {
			idList = append(idList, id)
		}
		err := importer.repo.DeleteById(context.Background(), idList)
		if err != nil {
			return fmt.Errorf("Importer - Update - s.repo.Store: %w", err)
		}
	}

	return nil
}

func (importer *Importer) updateEntity(ctx context.Context, persons map[int]entity.Person, person *entity.Person) error {
	oldPerson, ok := persons[person.Uid]
	if ok {
		delete(persons, person.Uid)
		if person.FirstName != oldPerson.FirstName || person.LastName != oldPerson.LastName {
			return importer.repo.Update(context.Background(), person)
		}
	} else {
		return importer.repo.Store(context.Background(), person)
	}

	return nil
}
