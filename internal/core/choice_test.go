package core

import (
	"testing"

	"github.com/astrviktor/banner-rotation/internal/storage"
	memorystorage "github.com/astrviktor/banner-rotation/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestGetBanner(t *testing.T) {
	t.Run("simple test 2 banners and no clicks", func(t *testing.T) {
		s := memorystorage.New()

		segment := storage.Segment{ID: "men", Description: "men"}
		err := s.CreateSegment(segment)
		require.NoError(t, err)

		slot := storage.Slot{ID: "slot", Description: "slot"}
		err = s.CreateSlot(slot)
		require.NoError(t, err)

		bannerA := storage.Banner{ID: "bannerA", Description: "bannerA"}
		err = s.CreateBanner(bannerA)
		require.NoError(t, err)

		bannerB := storage.Banner{ID: "bannerB", Description: "bannerB"}
		err = s.CreateBanner(bannerB)
		require.NoError(t, err)

		err = s.CreateRotation(storage.Rotation{IDSlot: slot.ID, IDBanner: bannerA.ID})
		require.NoError(t, err)
		err = s.CreateRotation(storage.Rotation{IDSlot: slot.ID, IDBanner: bannerB.ID})
		require.NoError(t, err)

		for i := 0; i < 1000; i++ {
			IDBanner, err := GetBanner(s, slot.ID, segment.ID)
			require.NoError(t, err)

			err = s.AddEvent(slot.ID, IDBanner, segment.ID, storage.Show)
			require.NoError(t, err)
		}

		showCountBannerA := s.GetCountActionsForBannerAndSegment(bannerA.ID, segment.ID, storage.Show)
		showCountBannerB := s.GetCountActionsForBannerAndSegment(bannerB.ID, segment.ID, storage.Show)

		require.Equal(t, 500, showCountBannerA)
		require.Equal(t, 500, showCountBannerB)
	})

	t.Run("simple test 3 banners and no clicks", func(t *testing.T) {
		s := memorystorage.New()

		segment := storage.Segment{ID: "men", Description: "men"}
		err := s.CreateSegment(segment)
		require.NoError(t, err)

		slot := storage.Slot{ID: "slot", Description: "slot"}
		err = s.CreateSlot(slot)
		require.NoError(t, err)

		bannerA := storage.Banner{ID: "bannerA", Description: "bannerA"}
		err = s.CreateBanner(bannerA)
		require.NoError(t, err)

		bannerB := storage.Banner{ID: "bannerB", Description: "bannerB"}
		err = s.CreateBanner(bannerB)
		require.NoError(t, err)

		bannerC := storage.Banner{ID: "bannerC", Description: "bannerC"}
		err = s.CreateBanner(bannerB)
		require.NoError(t, err)

		err = s.CreateRotation(storage.Rotation{IDSlot: slot.ID, IDBanner: bannerA.ID})
		require.NoError(t, err)
		err = s.CreateRotation(storage.Rotation{IDSlot: slot.ID, IDBanner: bannerB.ID})
		require.NoError(t, err)
		err = s.CreateRotation(storage.Rotation{IDSlot: slot.ID, IDBanner: bannerC.ID})
		require.NoError(t, err)

		for i := 0; i < 900; i++ {
			IDBanner, err := GetBanner(s, slot.ID, segment.ID)
			require.NoError(t, err)

			err = s.AddEvent(slot.ID, IDBanner, segment.ID, storage.Show)
			require.NoError(t, err)
		}

		showCountBannerA := s.GetCountActionsForBannerAndSegment(bannerA.ID, segment.ID, storage.Show)
		showCountBannerB := s.GetCountActionsForBannerAndSegment(bannerB.ID, segment.ID, storage.Show)
		showCountBannerC := s.GetCountActionsForBannerAndSegment(bannerC.ID, segment.ID, storage.Show)

		require.Equal(t, 300, showCountBannerA)
		require.Equal(t, 300, showCountBannerB)
		require.Equal(t, 300, showCountBannerC)
	})

	t.Run("simple test 2 banners and clicks on bannerA", func(t *testing.T) {
		s := memorystorage.New()

		segment := storage.Segment{ID: "men", Description: "men"}
		err := s.CreateSegment(segment)
		require.NoError(t, err)

		slot := storage.Slot{ID: "slot", Description: "slot"}
		err = s.CreateSlot(slot)
		require.NoError(t, err)

		bannerA := storage.Banner{ID: "bannerA", Description: "bannerA"}
		err = s.CreateBanner(bannerA)
		require.NoError(t, err)

		bannerB := storage.Banner{ID: "bannerB", Description: "bannerB"}
		err = s.CreateBanner(bannerB)
		require.NoError(t, err)

		err = s.CreateRotation(storage.Rotation{IDSlot: slot.ID, IDBanner: bannerA.ID})
		require.NoError(t, err)
		err = s.CreateRotation(storage.Rotation{IDSlot: slot.ID, IDBanner: bannerB.ID})
		require.NoError(t, err)

		for i := 0; i < 1000; i++ {
			IDBanner, err := GetBanner(s, slot.ID, segment.ID)
			require.NoError(t, err)

			err = s.AddEvent(slot.ID, IDBanner, segment.ID, storage.Show)
			require.NoError(t, err)

			if IDBanner == bannerA.ID {
				err = s.AddEvent(slot.ID, IDBanner, segment.ID, storage.Click)
				require.NoError(t, err)
			}
		}

		showCountBannerA := s.GetCountActionsForBannerAndSegment(bannerA.ID, segment.ID, storage.Show)
		showCountBannerB := s.GetCountActionsForBannerAndSegment(bannerB.ID, segment.ID, storage.Show)

		require.Equal(t, 988, showCountBannerA)
		require.Equal(t, 12, showCountBannerB)
	})
}
