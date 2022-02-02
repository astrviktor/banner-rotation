package core

import (
	"math"

	"github.com/astrviktor/banner-rotation/internal/storage"
)

func GetBanner(s storage.Storage, slotID, segmentID string) (string, error) {
	// 1. получить список баннеров в ротации с slotID
	bannersID, err := s.GetBannersForSlot(slotID)
	if err != nil {
		return storage.EmptyID, err
	}

	if len(bannersID) == 0 {
		return storage.EmptyID, ErrTooFewBannersForSlot
	}

	// 2. для каждого баннера посчитать количество показов для сегмента (независимо от слотов)
	// если для баннера 0 показов, можно сразу его вернуть для показа

	// 3. для каждого баннера посчитать количество переходов для сегмента (независимо от слотов)

	// 4. посчитать количество показов всех баннеров в слоте для сегмента (сумма из п.2)

	stats := make(map[string]storage.Stat)
	showsAmount := 0
	for _, bannerID := range bannersID {
		stat, err := s.GetStatForBannerAndSegment(bannerID, segmentID)
		if err != nil {
			return storage.EmptyID, err
		}

		if stat.ClickCount > stat.ShowCount {
			return storage.EmptyID, ErrBannerClicksMoreThenShows
		}

		if stat.ShowCount == 0 {
			return bannerID, nil
		}

		stats[bannerID] = stat
		showsAmount += stat.ShowCount
	}

	// 5. weight = xi + sqrt(2 * Ln(n) / ni)
	// нужно взять баннер с максимальным весом

	var weightMax float64
	resultID := ""

	ln := math.Log(float64(showsAmount)) / math.Log(math.E)
	for _, bannerID := range bannersID {
		stat := stats[bannerID]

		weight := float64(stat.ClickCount)/float64(stat.ShowCount) + math.Sqrt(2*ln/float64(stat.ShowCount))

		if weight > weightMax {
			weightMax = weight
			resultID = bannerID
		}
	}

	return resultID, nil
}
