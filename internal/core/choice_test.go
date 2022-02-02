package core

import (
	"testing"

	"github.com/astrviktor/banner-rotation/internal/storage"
	memorystorage "github.com/astrviktor/banner-rotation/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

// тесты основного алгоритма:
// - 2 баннера, нет переходов - должна быть одинаковая статистика показов
// - 2 баннера, переходы при каждом показе - должна быть одинаковая статистика показов
// - 2 баннера, по 1 переходят, по 2 нет - по первому показов сильно больше, но и по второму есть

// - 3 баннера, нет переходов - должна быть одинаковая статистика показов
// - 3 баннера, переходы при каждом показе - должна быть одинаковая статистика показов
// - 3 баннера, по 1 переходят, по 2 и 3 нет - по первому показов сильно больше, но и по остальным есть

// тесты с ошибками:
// - для слота не создано ротаций
// - количество переходов по какой то причине стало больше чем количество показов

func TestGetBannerTwoOK(t *testing.T) {
	t.Run("test 2 banners and no clicks", func(t *testing.T) {
		s := memorystorage.New()
		err := s.Connect()
		require.NoError(t, err)

		segment, err := s.CreateSegment("segment")
		require.NoError(t, err)

		slot, err := s.CreateSlot("slot")
		require.NoError(t, err)

		bannerA, err := s.CreateBanner("bannerA")
		require.NoError(t, err)
		bannerB, err := s.CreateBanner("bannerB")
		require.NoError(t, err)

		err = s.CreateRotation(storage.Rotation{SlotID: slot, BannerID: bannerA})
		require.NoError(t, err)
		err = s.CreateRotation(storage.Rotation{SlotID: slot, BannerID: bannerB})
		require.NoError(t, err)

		for i := 0; i < 1000; i++ {
			bannerID, err := GetBanner(s, slot, segment)
			require.NoError(t, err)

			err = s.CreateEvent(slot, bannerID, segment, storage.Show)
			require.NoError(t, err)
		}

		statBannerA, err := s.GetStatForBannerAndSegment(bannerA, segment)
		require.NoError(t, err)
		statBannerB, err := s.GetStatForBannerAndSegment(bannerB, segment)
		require.NoError(t, err)

		require.Equal(t, 500, statBannerA.ShowCount)
		require.Equal(t, 500, statBannerB.ShowCount)
		require.Equal(t, 0, statBannerA.ClickCount)
		require.Equal(t, 0, statBannerB.ClickCount)

		s.Close()
	})

	t.Run("test 2 banners and all clicks", func(t *testing.T) {
		s := memorystorage.New()
		err := s.Connect()
		require.NoError(t, err)

		segment, err := s.CreateSegment("segment")
		require.NoError(t, err)

		slot, err := s.CreateSlot("slot")
		require.NoError(t, err)

		bannerA, err := s.CreateBanner("bannerA")
		require.NoError(t, err)
		bannerB, err := s.CreateBanner("bannerB")
		require.NoError(t, err)

		err = s.CreateRotation(storage.Rotation{SlotID: slot, BannerID: bannerA})
		require.NoError(t, err)
		err = s.CreateRotation(storage.Rotation{SlotID: slot, BannerID: bannerB})
		require.NoError(t, err)

		for i := 0; i < 1000; i++ {
			bannerID, err := GetBanner(s, slot, segment)
			require.NoError(t, err)

			err = s.CreateEvent(slot, bannerID, segment, storage.Show)
			require.NoError(t, err)
			err = s.CreateEvent(slot, bannerID, segment, storage.Click)
			require.NoError(t, err)
		}

		statBannerA, err := s.GetStatForBannerAndSegment(bannerA, segment)
		require.NoError(t, err)
		statBannerB, err := s.GetStatForBannerAndSegment(bannerB, segment)
		require.NoError(t, err)

		require.Equal(t, 500, statBannerA.ShowCount)
		require.Equal(t, 500, statBannerB.ShowCount)
		require.Equal(t, 500, statBannerA.ClickCount)
		require.Equal(t, 500, statBannerB.ClickCount)

		s.Close()
	})

	t.Run("test 2 banners and clicks on bannerA", func(t *testing.T) {
		s := memorystorage.New()
		err := s.Connect()
		require.NoError(t, err)

		segment, err := s.CreateSegment("segment")
		require.NoError(t, err)

		slot, err := s.CreateSlot("slot")
		require.NoError(t, err)

		bannerA, err := s.CreateBanner("bannerA")
		require.NoError(t, err)
		bannerB, err := s.CreateBanner("bannerB")
		require.NoError(t, err)

		err = s.CreateRotation(storage.Rotation{SlotID: slot, BannerID: bannerA})
		require.NoError(t, err)
		err = s.CreateRotation(storage.Rotation{SlotID: slot, BannerID: bannerB})
		require.NoError(t, err)

		for i := 0; i < 1000; i++ {
			bannerID, err := GetBanner(s, slot, segment)
			require.NoError(t, err)

			err = s.CreateEvent(slot, bannerID, segment, storage.Show)
			require.NoError(t, err)

			if bannerID == bannerA {
				err = s.CreateEvent(slot, bannerID, segment, storage.Click)
				require.NoError(t, err)
			}
		}

		statBannerA, err := s.GetStatForBannerAndSegment(bannerA, segment)
		require.NoError(t, err)
		statBannerB, err := s.GetStatForBannerAndSegment(bannerB, segment)
		require.NoError(t, err)

		require.Equal(t, 988, statBannerA.ShowCount)
		require.Equal(t, 12, statBannerB.ShowCount)

		s.Close()
	})
}

func TestGetBannerThreeOK(t *testing.T) {
	t.Run("test 3 banners and no clicks", func(t *testing.T) {
		s := memorystorage.New()

		segment, err := s.CreateSegment("segment")
		require.NoError(t, err)

		slot, err := s.CreateSlot("slot")
		require.NoError(t, err)

		bannerA, err := s.CreateBanner("bannerA")
		require.NoError(t, err)
		bannerB, err := s.CreateBanner("bannerB")
		require.NoError(t, err)
		bannerC, err := s.CreateBanner("bannerC")
		require.NoError(t, err)

		err = s.CreateRotation(storage.Rotation{SlotID: slot, BannerID: bannerA})
		require.NoError(t, err)
		err = s.CreateRotation(storage.Rotation{SlotID: slot, BannerID: bannerB})
		require.NoError(t, err)
		err = s.CreateRotation(storage.Rotation{SlotID: slot, BannerID: bannerC})
		require.NoError(t, err)

		for i := 0; i < 900; i++ {
			bannerID, err := GetBanner(s, slot, segment)
			require.NoError(t, err)

			err = s.CreateEvent(slot, bannerID, segment, storage.Show)
			require.NoError(t, err)
		}

		statBannerA, err := s.GetStatForBannerAndSegment(bannerA, segment)
		require.NoError(t, err)
		statBannerB, err := s.GetStatForBannerAndSegment(bannerB, segment)
		require.NoError(t, err)
		statBannerC, err := s.GetStatForBannerAndSegment(bannerC, segment)
		require.NoError(t, err)

		require.Equal(t, 300, statBannerA.ShowCount)
		require.Equal(t, 300, statBannerB.ShowCount)
		require.Equal(t, 300, statBannerC.ShowCount)
		require.Equal(t, 0, statBannerA.ClickCount)
		require.Equal(t, 0, statBannerB.ClickCount)
		require.Equal(t, 0, statBannerC.ClickCount)

		s.Close()
	})

	t.Run("test 3 banners and all clicks", func(t *testing.T) {
		s := memorystorage.New()

		segment, err := s.CreateSegment("segment")
		require.NoError(t, err)

		slot, err := s.CreateSlot("slot")
		require.NoError(t, err)

		bannerA, err := s.CreateBanner("bannerA")
		require.NoError(t, err)
		bannerB, err := s.CreateBanner("bannerB")
		require.NoError(t, err)
		bannerC, err := s.CreateBanner("bannerC")
		require.NoError(t, err)

		err = s.CreateRotation(storage.Rotation{SlotID: slot, BannerID: bannerA})
		require.NoError(t, err)
		err = s.CreateRotation(storage.Rotation{SlotID: slot, BannerID: bannerB})
		require.NoError(t, err)
		err = s.CreateRotation(storage.Rotation{SlotID: slot, BannerID: bannerC})
		require.NoError(t, err)

		for i := 0; i < 900; i++ {
			bannerID, err := GetBanner(s, slot, segment)
			require.NoError(t, err)

			err = s.CreateEvent(slot, bannerID, segment, storage.Show)
			require.NoError(t, err)
			err = s.CreateEvent(slot, bannerID, segment, storage.Click)
			require.NoError(t, err)
		}

		statBannerA, err := s.GetStatForBannerAndSegment(bannerA, segment)
		require.NoError(t, err)
		statBannerB, err := s.GetStatForBannerAndSegment(bannerB, segment)
		require.NoError(t, err)
		statBannerC, err := s.GetStatForBannerAndSegment(bannerC, segment)
		require.NoError(t, err)

		require.Equal(t, 300, statBannerA.ShowCount)
		require.Equal(t, 300, statBannerB.ShowCount)
		require.Equal(t, 300, statBannerC.ShowCount)
		require.Equal(t, 300, statBannerA.ClickCount)
		require.Equal(t, 300, statBannerB.ClickCount)
		require.Equal(t, 300, statBannerC.ClickCount)

		s.Close()
	})

	t.Run("test 3 banners and clicks on bannerA", func(t *testing.T) {
		s := memorystorage.New()

		segment, err := s.CreateSegment("segment")
		require.NoError(t, err)

		slot, err := s.CreateSlot("slot")
		require.NoError(t, err)

		bannerA, err := s.CreateBanner("bannerA")
		require.NoError(t, err)
		bannerB, err := s.CreateBanner("bannerB")
		require.NoError(t, err)
		bannerC, err := s.CreateBanner("bannerC")
		require.NoError(t, err)

		err = s.CreateRotation(storage.Rotation{SlotID: slot, BannerID: bannerA})
		require.NoError(t, err)
		err = s.CreateRotation(storage.Rotation{SlotID: slot, BannerID: bannerB})
		require.NoError(t, err)
		err = s.CreateRotation(storage.Rotation{SlotID: slot, BannerID: bannerC})
		require.NoError(t, err)

		for i := 0; i < 1000; i++ {
			bannerID, err := GetBanner(s, slot, segment)
			require.NoError(t, err)

			err = s.CreateEvent(slot, bannerID, segment, storage.Show)
			require.NoError(t, err)

			if bannerID == bannerA {
				err = s.CreateEvent(slot, bannerID, segment, storage.Click)
				require.NoError(t, err)
			}
		}

		statBannerA, err := s.GetStatForBannerAndSegment(bannerA, segment)
		require.NoError(t, err)
		statBannerB, err := s.GetStatForBannerAndSegment(bannerB, segment)
		require.NoError(t, err)
		statBannerC, err := s.GetStatForBannerAndSegment(bannerC, segment)
		require.NoError(t, err)

		require.Equal(t, 976, statBannerA.ShowCount)
		require.Equal(t, 12, statBannerB.ShowCount)
		require.Equal(t, 12, statBannerC.ShowCount)
		require.Equal(t, 976, statBannerA.ClickCount)
		require.Equal(t, 0, statBannerB.ClickCount)
		require.Equal(t, 0, statBannerC.ClickCount)

		s.Close()
	})
}

func TestGetBannerErrors(t *testing.T) {
	t.Run("no rotations created for slot", func(t *testing.T) {
		s := memorystorage.New()
		err := s.Connect()
		require.NoError(t, err)

		segment, err := s.CreateSegment("segment")
		require.NoError(t, err)

		slot, err := s.CreateSlot("slot")
		require.NoError(t, err)

		_, err = GetBanner(s, slot, segment)
		require.ErrorIs(t, err, ErrTooFewBannersForSlot)

		s.Close()
	})

	t.Run("number of clicks more than number of shows", func(t *testing.T) {
		s := memorystorage.New()
		err := s.Connect()
		require.NoError(t, err)

		segment, err := s.CreateSegment("segment")
		require.NoError(t, err)

		slot, err := s.CreateSlot("slot")
		require.NoError(t, err)

		bannerA, err := s.CreateBanner("bannerA")
		require.NoError(t, err)

		err = s.CreateRotation(storage.Rotation{SlotID: slot, BannerID: bannerA})
		require.NoError(t, err)

		statBannerA, err := s.GetStatForBannerAndSegment(bannerA, segment)
		require.NoError(t, err)

		require.Equal(t, 0, statBannerA.ShowCount)
		require.Equal(t, 0, statBannerA.ClickCount)

		_, err = GetBanner(s, slot, segment)
		require.NoError(t, err)

		err = s.CreateEvent(slot, bannerA, segment, storage.Click)
		require.NoError(t, err)

		statBannerA, err = s.GetStatForBannerAndSegment(bannerA, segment)
		require.NoError(t, err)

		require.Equal(t, 0, statBannerA.ShowCount)
		require.Equal(t, 1, statBannerA.ClickCount)

		_, err = GetBanner(s, slot, segment)
		require.ErrorIs(t, err, ErrBannerClicksMoreThenShows)

		s.Close()
	})
}
