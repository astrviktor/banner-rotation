package core

import "errors"

var (
	ErrTooFewBannersForSlot      = errors.New("для слота недостаточное количество баннеров в ротации")
	ErrBannerClicksMoreThenShows = errors.New("для баннера количество кликов больше чем количество показов")
)
