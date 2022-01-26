package core

import (
	"github.com/astrviktor/banner-rotation/internal/storage"
	memorystorage "github.com/astrviktor/banner-rotation/internal/storage/memory"
	"math"
)

func GetBanner(s *memorystorage.Storage, idSlot, idSegment string) (string, error) {
	// 1. получить список баннеров в ротации с IdSlot
	idBanners, err := s.GetBannersForRotations(idSlot)
	if err != nil {
		return "", err
	}

	// 2. для каждого баннера посчитать количество показов для сегмента (независимо от слотов)
	// если для баннера 0 показов, можно сразу его вернуть для показа

	// 3. для каждого баннера посчитать количество переходов для сегмента (независимо от слотов)

	// 4. посчитать количество показов всех баннеров в слоте для сегмента (сумма из п.2)

	shows := make(map[string]int)
	clicks := make(map[string]int)
	showsSum := 0
	for _, idBanner := range idBanners {
		showCount := s.GetCountActionsForBannerAndSegment(idBanner, idSegment, storage.Show)
		clickCount := s.GetCountActionsForBannerAndSegment(idBanner, idSegment, storage.Click)

		if showCount == 0 {
			return idBanner, nil
		}

		if clickCount > showCount {
			return "", ErrClicksMoreThenShows
		}

		shows[idBanner] = showCount
		clicks[idBanner] = clickCount
		showsSum += showCount
	}

	// 5. weight = xi + sqrt(2 * Ln(n) / ni)
	// нужно взять баннер с максимальным весом

	ln := math.Log(float64(showsSum)) / math.Log(math.E)

	var maxWeight float64
	resultIdBanner := ""

	for _, idBanner := range idBanners {
		showCount := shows[idBanner]
		clickCount := clicks[idBanner]

		weight := float64(clickCount)/float64(showCount) + math.Sqrt(2*ln/float64(showCount))

		if weight > maxWeight {
			maxWeight = weight
			resultIdBanner = idBanner
		}
	}

	return resultIdBanner, nil
}
