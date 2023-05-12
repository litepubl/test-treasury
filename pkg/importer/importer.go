package importer

import (
	"context"
	"fmt"
	"strconv"

	"github.com/litepubl/test-treasury/pkg/entity"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const individual = "Individual"

// Importer непосредственно скачивает данные xml и обновляет БД
type Importer struct {
	repo       PersonRepo
	downloader Downloader
}

// NewImporter  - конструктор.
func NewImporter(r PersonRepo, d Downloader) *Importer {
	return &Importer{
		repo:       r,
		downloader: d,
	}
}

// Import возвращает ошибку, если импорт не удался
func (im *Importer) Import(ctx context.Context) error {
	persons, err := im.repo.All(context.Background())
	if err != nil {
		return fmt.Errorf("Importer - Person - s.repo.All: %w", err)
	}

	sdnList, err := im.downloader.Download(ctx)
	if err != nil {
		return fmt.Errorf("Importer - Update - downloader.Download: %w", err)
	}

	for _, sdnEntry := range sdnList.SdnEntry {
		if sdnEntry.SdnType == individual {
			uid, err := strconv.Atoi(sdnEntry.Uid)
			if err != nil {
				continue
			}

			person := &entity.Person{
				Uid:       uid,
				FirstName: cases.Title(language.Und).String(sdnEntry.FirstName),
				LastName:  cases.Title(language.Und).String(sdnEntry.LastName),
			}

			err = im.updateEntity(ctx, persons, person)
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
		err := im.repo.DeleteByIDList(context.Background(), idList)
		if err != nil {
			return fmt.Errorf("Importer - Update - s.repo.Store: %w", err)
		}
	}

	return nil
}

func (im *Importer) updateEntity(ctx context.Context, persons map[int]entity.Person, person *entity.Person) error {
	oldPerson, ok := persons[person.Uid]
	if ok {
		delete(persons, person.Uid)
		if person.FirstName != oldPerson.FirstName || person.LastName != oldPerson.LastName {
			return im.repo.Update(context.Background(), person)
		}
	} else {
		return im.repo.Store(context.Background(), person)
	}

	return nil
}
