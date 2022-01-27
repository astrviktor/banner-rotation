package core

import "errors"

var ErrBannerClicksMoreThenShows = errors.New("для баннера количество кликов больше чем количество показов")
var ErrSlotTooFewBanners = errors.New("для слота недостаточное количество баннеров в ротации")