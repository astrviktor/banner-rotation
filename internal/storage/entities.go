package storage

import (
	"time"

	"github.com/google/uuid"
)

type Storage interface {
	Connect() error
	Close() error
	CreateSlot(description string) (string, error)
	CreateBanner(description string) (string, error)
	CreateSegment(description string) (string, error)
	CreateRotation(rotation Rotation) error
	DeleteRotation(rotation Rotation) error
	CreateEvent(slotID, bannerID, segmentID string, action ActionType) error
	GetBannersForSlot(slotID string) ([]string, error)
	GetStatForBannerAndSegment(bannerID, segmentID string) (Stat, error)
}

// Slot - место на сайте, на котором мы показываем баннер.
type Slot struct {
	ID          string `json:"id"`          // ID - уникальный идентификатор слота (UUID)
	Description string `json:"description"` // Описание слота
}

// Banner - рекламный/информационный элемент, который показывается в слоте.
type Banner struct {
	ID          string `json:"id"`          // ID - уникальный идентификатор баннера (UUID)
	Description string `json:"description"` // Описание баннера
}

// Segment - группа пользователей сайта со схожими интересами, например "девушки 20-25" или "дедушки 80+".
type Segment struct {
	ID          string `json:"id"`          // ID - уникальный идентификатор сегмента (UUID)
	Description string `json:"description"` // Описание сегмента
}

// Rotation - баннер в ротации в данном слоте.
type Rotation struct {
	SlotID   string `json:"slotId"`   // ID слота
	BannerID string `json:"bannerId"` // ID баннера
}

// Stat - агрегированная статистика по переходу и показу баннера.
type Stat struct {
	BannerID   string `json:"bannerId"`   // ID баннера
	SegmentID  string `json:"segmentId"`  // ID сегмента
	ShowCount  int    `json:"showCount"`  // количество показов
	ClickCount int    `json:"clickCount"` // количество переходов
}

// Event - событие по переходу или показу баннера.
type Event struct {
	SlotID    string     `json:"slotId"`    // ID слота
	BannerID  string     `json:"bannerId"`  // ID баннера
	SegmentID string     `json:"segmentId"` // ID сегмента
	Action    ActionType `json:"action"`    // Действие: клик или показ
	Date      time.Time  `json:"date"`      // Дата и время события
}

type ActionType int

const (
	Show  ActionType = 1
	Click ActionType = 2
)

func NewID() string {
	return uuid.New().String()
}

const EmptyID string = "00000000-0000-0000-0000-000000000000"
