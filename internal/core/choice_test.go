package core

import (
	"github.com/astrviktor/banner-rotation/internal/storage"
	memorystorage "github.com/astrviktor/banner-rotation/internal/storage/memory"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetBanner(t *testing.T) {
	t.Run("simple test 2 banners and no clicks", func(t *testing.T) {
		s := memorystorage.New()

		segment := storage.Segment{Id: "men", Description: "men"}
		err := s.CreateSegment(segment)
		require.NoError(t, err)

		slot := storage.Slot{Id: "slot", Description: "slot"}
		err = s.CreateSlot(slot)
		require.NoError(t, err)

		bannerA := storage.Banner{Id: "bannerA", Description: "bannerA"}
		err = s.CreateBanner(bannerA)
		require.NoError(t, err)

		bannerB := storage.Banner{Id: "bannerB", Description: "bannerB"}
		err = s.CreateBanner(bannerB)
		require.NoError(t, err)

		err = s.CreateRotation(storage.Rotation{IdSlot: slot.Id, IdBanner: bannerA.Id})
		require.NoError(t, err)
		err = s.CreateRotation(storage.Rotation{IdSlot: slot.Id, IdBanner: bannerB.Id})
		require.NoError(t, err)

		for i := 0; i < 1000; i++ {
			idBanner, err := GetBanner(s, slot.Id, segment.Id)
			require.NoError(t, err)

			err = s.AddEvent(slot.Id, idBanner, segment.Id, storage.Show)
			require.NoError(t, err)
		}

		showCountBannerA := s.GetCountActionsForBannerAndSegment(bannerA.Id, segment.Id, storage.Show)
		showCountBannerB := s.GetCountActionsForBannerAndSegment(bannerB.Id, segment.Id, storage.Show)

		require.Equal(t, 500, showCountBannerA)
		require.Equal(t, 500, showCountBannerB)
	})

	t.Run("simple test 3 banners and no clicks", func(t *testing.T) {
		s := memorystorage.New()

		segment := storage.Segment{Id: "men", Description: "men"}
		err := s.CreateSegment(segment)
		require.NoError(t, err)

		slot := storage.Slot{Id: "slot", Description: "slot"}
		err = s.CreateSlot(slot)
		require.NoError(t, err)

		bannerA := storage.Banner{Id: "bannerA", Description: "bannerA"}
		err = s.CreateBanner(bannerA)
		require.NoError(t, err)

		bannerB := storage.Banner{Id: "bannerB", Description: "bannerB"}
		err = s.CreateBanner(bannerB)
		require.NoError(t, err)

		bannerC := storage.Banner{Id: "bannerC", Description: "bannerC"}
		err = s.CreateBanner(bannerB)
		require.NoError(t, err)

		err = s.CreateRotation(storage.Rotation{IdSlot: slot.Id, IdBanner: bannerA.Id})
		require.NoError(t, err)
		err = s.CreateRotation(storage.Rotation{IdSlot: slot.Id, IdBanner: bannerB.Id})
		require.NoError(t, err)
		err = s.CreateRotation(storage.Rotation{IdSlot: slot.Id, IdBanner: bannerC.Id})
		require.NoError(t, err)

		for i := 0; i < 900; i++ {
			idBanner, err := GetBanner(s, slot.Id, segment.Id)
			require.NoError(t, err)

			err = s.AddEvent(slot.Id, idBanner, segment.Id, storage.Show)
			require.NoError(t, err)
		}

		showCountBannerA := s.GetCountActionsForBannerAndSegment(bannerA.Id, segment.Id, storage.Show)
		showCountBannerB := s.GetCountActionsForBannerAndSegment(bannerB.Id, segment.Id, storage.Show)
		showCountBannerC := s.GetCountActionsForBannerAndSegment(bannerC.Id, segment.Id, storage.Show)

		require.Equal(t, 300, showCountBannerA)
		require.Equal(t, 300, showCountBannerB)
		require.Equal(t, 300, showCountBannerC)
	})

	t.Run("simple test 2 banners and clicks on bannerA", func(t *testing.T) {
		s := memorystorage.New()

		segment := storage.Segment{Id: "men", Description: "men"}
		err := s.CreateSegment(segment)
		require.NoError(t, err)

		slot := storage.Slot{Id: "slot", Description: "slot"}
		err = s.CreateSlot(slot)
		require.NoError(t, err)

		bannerA := storage.Banner{Id: "bannerA", Description: "bannerA"}
		err = s.CreateBanner(bannerA)
		require.NoError(t, err)

		bannerB := storage.Banner{Id: "bannerB", Description: "bannerB"}
		err = s.CreateBanner(bannerB)
		require.NoError(t, err)

		err = s.CreateRotation(storage.Rotation{IdSlot: slot.Id, IdBanner: bannerA.Id})
		require.NoError(t, err)
		err = s.CreateRotation(storage.Rotation{IdSlot: slot.Id, IdBanner: bannerB.Id})
		require.NoError(t, err)

		for i := 0; i < 1000; i++ {
			idBanner, err := GetBanner(s, slot.Id, segment.Id)
			require.NoError(t, err)

			err = s.AddEvent(slot.Id, idBanner, segment.Id, storage.Show)
			require.NoError(t, err)

			if idBanner == bannerA.Id {
				err = s.AddEvent(slot.Id, idBanner, segment.Id, storage.Click)
				require.NoError(t, err)
			}
		}

		showCountBannerA := s.GetCountActionsForBannerAndSegment(bannerA.Id, segment.Id, storage.Show)
		showCountBannerB := s.GetCountActionsForBannerAndSegment(bannerB.Id, segment.Id, storage.Show)

		require.Equal(t, 988, showCountBannerA)
		require.Equal(t, 12, showCountBannerB)
	})

}
