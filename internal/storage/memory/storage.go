package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/astrviktor/banner-rotation/internal/storage"
)

type Storage struct {
	slots     map[string]storage.Slot
	banners   map[string]storage.Banner
	segments  map[string]storage.Segment
	rotations []storage.Rotation
	events    []storage.Event

	mutex *sync.RWMutex
}

func New() *Storage {
	mutex := sync.RWMutex{}

	return &Storage{
		slots:     make(map[string]storage.Slot),
		banners:   make(map[string]storage.Banner),
		segments:  make(map[string]storage.Segment),
		rotations: nil,
		events:    nil,
		mutex:     &mutex,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return nil
}

func (s *Storage) CreateSlot(slot storage.Slot) error {
	s.mutex.Lock()
	s.slots[slot.ID] = slot
	s.mutex.Unlock()
	return nil
}

func (s *Storage) CreateBanner(banner storage.Banner) error {
	s.mutex.Lock()
	s.banners[banner.ID] = banner
	s.mutex.Unlock()
	return nil
}

func (s *Storage) CreateSegment(segment storage.Segment) error {
	s.mutex.Lock()
	s.segments[segment.ID] = segment
	s.mutex.Unlock()
	return nil
}

func (s *Storage) CreateRotation(rotation storage.Rotation) error {
	s.mutex.Lock()
	s.rotations = append(s.rotations, rotation)
	s.mutex.Unlock()
	return nil
}

func (s *Storage) DeleteRotation(rotation storage.Rotation) error {
	s.mutex.Lock()
	for IDx, r := range s.rotations {
		if r.IDSlot == rotation.IDSlot && r.IDBanner == rotation.IDBanner {
			s.rotations = append(s.rotations[:IDx], s.rotations[IDx+1:]...)
			break
		}
	}
	s.mutex.Unlock()
	return nil
}

func (s *Storage) AddEvent(idSlot, idBanner, idSegment string, action storage.ActionType) error {
	event := storage.Event{
		Action:    action,
		IDSlot:    idSlot,
		IDBanner:  idBanner,
		IDSegment: idSegment,
		Date:      time.Now().UTC(),
	}

	s.mutex.Lock()
	s.events = append(s.events, event)
	s.mutex.Unlock()
	return nil
}

func (s *Storage) GetBannersForRotations(idSlot string) ([]string, error) {
	var res []string

	s.mutex.RLock()
	for _, rotation := range s.rotations {
		if rotation.IDSlot == idSlot {
			res = append(res, rotation.IDBanner)
		}
	}
	s.mutex.RUnlock()

	return res, nil
}

func (s *Storage) GetCountActionsForBannerAndSegment(idBanner, idSegment string, action storage.ActionType) int {
	count := 0

	s.mutex.RLock()
	for _, event := range s.events {
		if event.IDBanner == idBanner && event.IDSegment == idSegment && event.Action == action {
			count++
		}
	}
	s.mutex.RUnlock()

	return count
}
