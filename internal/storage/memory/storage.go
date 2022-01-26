package memorystorage

import (
	"context"
	"github.com/astrviktor/banner-rotation/internal/storage"
	"sync"
	"time"
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
	s.slots[slot.Id] = slot
	s.mutex.Unlock()
	return nil
}

func (s *Storage) CreateBanner(banner storage.Banner) error {
	s.mutex.Lock()
	s.banners[banner.Id] = banner
	s.mutex.Unlock()
	return nil
}

func (s *Storage) CreateSegment(segment storage.Segment) error {
	s.mutex.Lock()
	s.segments[segment.Id] = segment
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
	for idx, r := range s.rotations {
		if r.IdSlot == rotation.IdSlot && r.IdBanner == rotation.IdBanner {
			s.rotations = append(s.rotations[:idx], s.rotations[idx+1:]...)
			break
		}
	}
	s.mutex.Unlock()
	return nil
}

func (s *Storage) AddEvent(action storage.ActionType, rotation storage.Rotation, IdSegment string) error {
	event := storage.Event{
		Action:    action,
		IdSlot:    rotation.IdSlot,
		IdBanner:  rotation.IdBanner,
		IdSegment: IdSegment,
		Date:      time.Now().UTC(),
	}

	s.mutex.Lock()
	s.events = append(s.events, event)
	s.mutex.Unlock()
	return nil
}
