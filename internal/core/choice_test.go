package core

import (
	"testing"

	"github.com/astrviktor/banner-rotation/internal/storage"
	memorystorage "github.com/astrviktor/banner-rotation/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestGetBanner(t *testing.T) {
	t.Run("test 2 banners and no clicks", func(t *testing.T) {
		s := memorystorage.New()
		err := s.Connect()
		require.NoError(t, err)

		segment, err := s.CreateSegment("men")
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

	t.Run("test 3 banners and no clicks", func(t *testing.T) {
		s := memorystorage.New()

		segment, err := s.CreateSegment("men")
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
	})

	t.Run("test 2 banners and clicks on bannerA", func(t *testing.T) {
		s := memorystorage.New()

		segment, err := s.CreateSegment("men")
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
	})
}
