package storage

import (
	"github.com/google/uuid"
	"time"
)

// Slot - место на сайте, на котором мы показываем баннер.
type Slot struct {
	Id          string `json:"id"`          // Id - уникальный идентификатор слота (UUID)
	Description string `json:"description"` // Описание слота
}

// Banner - рекламный/информационный элемент, который показывается в слоте.
type Banner struct {
	Id          string `json:"id"`          // Id - уникальный идентификатор баннера (UUID)
	Description string `json:"description"` // Описание баннера
}

// Segment - группа пользователей сайта со схожими интересами, например "девушки 20-25" или "дедушки 80+".
type Segment struct {
	Id          string `json:"id"`          // Id - уникальный идентификатор сегмента (UUID)
	Description string `json:"description"` // Описание сегмента
}

// Rotation - баннер в ротации в данном слоте.
type Rotation struct {
	IdSlot   string `json:"idSlot"`   // ID слота
	IdBanner string `json:"idBanner"` // ID баннера
}

// Event - событие по переходу или показу баннера.
type Event struct {
	Action    ActionType `json:"action"`    // Действие: клик или показ
	IdSlot    string     `json:"idSlot"`    // ID слота
	IdBanner  string     `json:"idBanner"`  // ID баннера
	IdSegment string     `json:"idSegment"` // ID сегмента
	Date      time.Time  `json:"date"`      // Дата и время события
}

type ActionType int

const (
	Click ActionType = 1
	Show  ActionType = 2
)

func NewUuidId() string {
	return uuid.New().String()
}
