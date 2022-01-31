package memorystorage

import (
	"sync"
	"time"

	"github.com/astrviktor/banner-rotation/internal/storage"
)

type Storage struct {
	slots     map[string]storage.Slot
	banners   map[string]storage.Banner
	segments  map[string]storage.Segment
	rotations []storage.Rotation
	stats     []storage.Stat
	events    []storage.Event

	mutex *sync.RWMutex
}

func New() *Storage {
	mutex := sync.RWMutex{}

	return &Storage{
		slots:     make(map[string]storage.Slot),
		banners:   make(map[string]storage.Banner),
		segments:  make(map[string]storage.Segment),
		rotations: make([]storage.Rotation, 0),
		stats:     make([]storage.Stat, 0),
		events:    make([]storage.Event, 0),
		mutex:     &mutex,
	}
}

func (s *Storage) Connect() error {
	return nil
}

func (s *Storage) Close() error {
	return nil
}

func (s *Storage) CreateSlot(description string) (string, error) {
	id := storage.NewID()
	slot := storage.Slot{ID: id, Description: description}
	s.mutex.Lock()
	s.slots[id] = slot
	s.mutex.Unlock()
	return id, nil
}

func (s *Storage) CreateBanner(description string) (string, error) {
	id := storage.NewID()
	banner := storage.Banner{ID: id, Description: description}
	s.mutex.Lock()
	s.banners[id] = banner
	for _, segment := range s.segments {
		stat := storage.Stat{
			BannerID:   id,
			SegmentID:  segment.ID,
			ShowCount:  0,
			ClickCount: 0,
		}
		s.stats = append(s.stats, stat)
	}
	s.mutex.Unlock()
	return id, nil
}

func (s *Storage) CreateSegment(description string) (string, error) {
	id := storage.NewID()
	segment := storage.Segment{ID: id, Description: description}
	s.mutex.Lock()
	s.segments[id] = segment
	for _, banner := range s.banners {
		stat := storage.Stat{
			BannerID:   banner.ID,
			SegmentID:  id,
			ShowCount:  0,
			ClickCount: 0,
		}
		s.stats = append(s.stats, stat)
	}
	s.mutex.Unlock()
	return id, nil
}

func (s *Storage) CreateRotation(rotation storage.Rotation) error {
	s.mutex.Lock()
	s.rotations = append(s.rotations, rotation)
	s.mutex.Unlock()
	return nil
}

func (s *Storage) DeleteRotation(rotation storage.Rotation) error {
	s.mutex.Lock()
	for idx, elem := range s.rotations {
		if elem.SlotID == rotation.SlotID && elem.BannerID == rotation.BannerID {
			s.rotations = append(s.rotations[:idx], s.rotations[idx+1:]...)
			break
		}
	}
	s.mutex.Unlock()
	return nil
}

func (s *Storage) CreateEvent(slotID, bannerID, segmentID string, action storage.ActionType) error {
	event := storage.Event{
		SlotID:    slotID,
		BannerID:  bannerID,
		SegmentID: segmentID,
		Action:    action,
		Date:      time.Now().UTC(),
	}

	s.mutex.Lock()
	s.events = append(s.events, event)

	for idx, stat := range s.stats {
		if stat.BannerID == bannerID && stat.SegmentID == segmentID {
			switch action {
			case storage.Show:
				s.stats[idx].ShowCount++
			case storage.Click:
				s.stats[idx].ClickCount++
			}
			break
		}
	}
	s.mutex.Unlock()
	return nil
}

func (s *Storage) GetBannersForSlot(slotID string) ([]string, error) {
	var bannersID []string

	s.mutex.RLock()
	for _, rotation := range s.rotations {
		if rotation.SlotID == slotID {
			bannersID = append(bannersID, rotation.BannerID)
		}
	}
	s.mutex.RUnlock()

	return bannersID, nil
}

func (s *Storage) GetStatForBannerAndSegment(bannerID, segmentID string) storage.Stat {
	s.mutex.RLock()
	for _, stat := range s.stats {
		if stat.BannerID == bannerID && stat.SegmentID == segmentID {
			s.mutex.RUnlock()
			return stat
		}
	}
	s.mutex.RUnlock()

	return storage.Stat{}
}
