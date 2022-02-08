//go:build integration

package integration_test

import (
	"os"
	"testing"
	"time"

	internalhttp "github.com/astrviktor/banner-rotation/internal/server/http"
	"github.com/stretchr/testify/suite"
)

type BannerRotationSuite struct {
	suite.Suite
	client *internalhttp.Client
}

func (s *BannerRotationSuite) SetupSuite() {
	// wait project up
	time.Sleep(15 * time.Second)

	host := os.Getenv("SERVICE_HOST")
	port := os.Getenv("SERVICE_PORT")

	if port == "" {
		port = "8888"
	}
	if host == "" {
		host = "127.0.0.1"
	}

	s.client = internalhttp.NewClient(host, port, time.Second)
}

func (s *BannerRotationSuite) SetupTest() {
}

func (s *BannerRotationSuite) TestTwoBannersAndNoClicks() {
	err := s.client.GetStatus()
	s.Require().NoError(err)

	segment, err := s.client.CreateSegment("segment")
	s.Require().NoError(err)

	slot, err := s.client.CreateSlot("slot")
	s.Require().NoError(err)

	bannerA, err := s.client.CreateBanner("bannerA")
	s.Require().NoError(err)
	bannerB, err := s.client.CreateBanner("bannerB")
	s.Require().NoError(err)

	err = s.client.CreateRotation(slot, bannerA)
	s.Require().NoError(err)
	err = s.client.CreateRotation(slot, bannerB)
	s.Require().NoError(err)

	for i := 0; i < 1000; i++ {
		_, err := s.client.Choice(slot, segment)
		s.Require().NoError(err)
	}

	statBannerA, err := s.client.GetStat(bannerA, segment)
	s.Require().NoError(err)
	statBannerB, err := s.client.GetStat(bannerB, segment)
	s.Require().NoError(err)

	s.Require().Equal(500, statBannerA.ShowCount)
	s.Require().Equal(500, statBannerB.ShowCount)

	s.Require().Equal(0, statBannerA.ClickCount)
	s.Require().Equal(0, statBannerB.ClickCount)
}

func (s *BannerRotationSuite) TestTwoBannersAndAllClicks() {
	err := s.client.GetStatus()
	s.Require().NoError(err)

	segment, err := s.client.CreateSegment("segment")
	s.Require().NoError(err)

	slot, err := s.client.CreateSlot("slot")
	s.Require().NoError(err)

	bannerA, err := s.client.CreateBanner("bannerA")
	s.Require().NoError(err)
	bannerB, err := s.client.CreateBanner("bannerB")
	s.Require().NoError(err)

	err = s.client.CreateRotation(slot, bannerA)
	s.Require().NoError(err)
	err = s.client.CreateRotation(slot, bannerB)
	s.Require().NoError(err)

	for i := 0; i < 1000; i++ {
		bannerID, err := s.client.Choice(slot, segment)
		s.Require().NoError(err)

		err = s.client.Click(slot, bannerID, segment)
		s.Require().NoError(err)
	}

	statBannerA, err := s.client.GetStat(bannerA, segment)
	s.Require().NoError(err)
	statBannerB, err := s.client.GetStat(bannerB, segment)
	s.Require().NoError(err)

	s.Require().Equal(500, statBannerA.ShowCount)
	s.Require().Equal(500, statBannerB.ShowCount)

	s.Require().Equal(500, statBannerA.ClickCount)
	s.Require().Equal(500, statBannerB.ClickCount)
}

func (s *BannerRotationSuite) TestTwoBannersAndOneClicks() {
	err := s.client.GetStatus()
	s.Require().NoError(err)

	segment, err := s.client.CreateSegment("segment")
	s.Require().NoError(err)

	slot, err := s.client.CreateSlot("slot")
	s.Require().NoError(err)

	bannerA, err := s.client.CreateBanner("bannerA")
	s.Require().NoError(err)
	bannerB, err := s.client.CreateBanner("bannerB")
	s.Require().NoError(err)

	err = s.client.CreateRotation(slot, bannerA)
	s.Require().NoError(err)
	err = s.client.CreateRotation(slot, bannerB)
	s.Require().NoError(err)

	for i := 0; i < 1000; i++ {
		bannerID, err := s.client.Choice(slot, segment)
		s.Require().NoError(err)

		if bannerID == bannerA {
			err = s.client.Click(slot, bannerID, segment)
			s.Require().NoError(err)
		}
	}

	statBannerA, err := s.client.GetStat(bannerA, segment)
	s.Require().NoError(err)
	statBannerB, err := s.client.GetStat(bannerB, segment)
	s.Require().NoError(err)

	s.Require().Equal(988, statBannerA.ShowCount)
	s.Require().Equal(12, statBannerB.ShowCount)

	s.Require().Equal(988, statBannerA.ClickCount)
	s.Require().Equal(0, statBannerB.ClickCount)
}

func (s *BannerRotationSuite) TestThreeBannersAndNoClicks() {
	err := s.client.GetStatus()
	s.Require().NoError(err)

	segment, err := s.client.CreateSegment("segment")
	s.Require().NoError(err)

	slot, err := s.client.CreateSlot("slot")
	s.Require().NoError(err)

	bannerA, err := s.client.CreateBanner("bannerA")
	s.Require().NoError(err)
	bannerB, err := s.client.CreateBanner("bannerB")
	s.Require().NoError(err)
	bannerC, err := s.client.CreateBanner("bannerC")
	s.Require().NoError(err)

	err = s.client.CreateRotation(slot, bannerA)
	s.Require().NoError(err)
	err = s.client.CreateRotation(slot, bannerB)
	s.Require().NoError(err)
	err = s.client.CreateRotation(slot, bannerC)
	s.Require().NoError(err)

	for i := 0; i < 900; i++ {
		_, err := s.client.Choice(slot, segment)
		s.Require().NoError(err)
	}

	statBannerA, err := s.client.GetStat(bannerA, segment)
	s.Require().NoError(err)
	statBannerB, err := s.client.GetStat(bannerB, segment)
	s.Require().NoError(err)
	statBannerC, err := s.client.GetStat(bannerC, segment)
	s.Require().NoError(err)

	s.Require().Equal(300, statBannerA.ShowCount)
	s.Require().Equal(300, statBannerB.ShowCount)
	s.Require().Equal(300, statBannerC.ShowCount)

	s.Require().Equal(0, statBannerA.ClickCount)
	s.Require().Equal(0, statBannerB.ClickCount)
	s.Require().Equal(0, statBannerC.ClickCount)
}

func (s *BannerRotationSuite) TestThreeBannersAndAllClicks() {
	err := s.client.GetStatus()
	s.Require().NoError(err)

	segment, err := s.client.CreateSegment("segment")
	s.Require().NoError(err)

	slot, err := s.client.CreateSlot("slot")
	s.Require().NoError(err)

	bannerA, err := s.client.CreateBanner("bannerA")
	s.Require().NoError(err)
	bannerB, err := s.client.CreateBanner("bannerB")
	s.Require().NoError(err)
	bannerC, err := s.client.CreateBanner("bannerC")
	s.Require().NoError(err)

	err = s.client.CreateRotation(slot, bannerA)
	s.Require().NoError(err)
	err = s.client.CreateRotation(slot, bannerB)
	s.Require().NoError(err)
	err = s.client.CreateRotation(slot, bannerC)
	s.Require().NoError(err)

	for i := 0; i < 900; i++ {
		bannerID, err := s.client.Choice(slot, segment)
		s.Require().NoError(err)

		err = s.client.Click(slot, bannerID, segment)
		s.Require().NoError(err)
	}

	statBannerA, err := s.client.GetStat(bannerA, segment)
	s.Require().NoError(err)
	statBannerB, err := s.client.GetStat(bannerB, segment)
	s.Require().NoError(err)
	statBannerC, err := s.client.GetStat(bannerC, segment)
	s.Require().NoError(err)

	s.Require().Equal(300, statBannerA.ShowCount)
	s.Require().Equal(300, statBannerB.ShowCount)
	s.Require().Equal(300, statBannerC.ShowCount)

	s.Require().Equal(300, statBannerA.ClickCount)
	s.Require().Equal(300, statBannerB.ClickCount)
	s.Require().Equal(300, statBannerC.ClickCount)
}

func (s *BannerRotationSuite) TestThreeBannersAndOneClicks() {
	err := s.client.GetStatus()
	s.Require().NoError(err)

	segment, err := s.client.CreateSegment("segment")
	s.Require().NoError(err)

	slot, err := s.client.CreateSlot("slot")
	s.Require().NoError(err)

	bannerA, err := s.client.CreateBanner("bannerA")
	s.Require().NoError(err)
	bannerB, err := s.client.CreateBanner("bannerB")
	s.Require().NoError(err)
	bannerC, err := s.client.CreateBanner("bannerC")
	s.Require().NoError(err)

	err = s.client.CreateRotation(slot, bannerA)
	s.Require().NoError(err)
	err = s.client.CreateRotation(slot, bannerB)
	s.Require().NoError(err)
	err = s.client.CreateRotation(slot, bannerC)
	s.Require().NoError(err)

	for i := 0; i < 1000; i++ {
		bannerID, err := s.client.Choice(slot, segment)
		s.Require().NoError(err)

		if bannerID == bannerA {
			err = s.client.Click(slot, bannerID, segment)
			s.Require().NoError(err)
		}
	}

	statBannerA, err := s.client.GetStat(bannerA, segment)
	s.Require().NoError(err)
	statBannerB, err := s.client.GetStat(bannerB, segment)
	s.Require().NoError(err)
	statBannerC, err := s.client.GetStat(bannerC, segment)
	s.Require().NoError(err)

	s.Require().Equal(976, statBannerA.ShowCount)
	s.Require().Equal(12, statBannerB.ShowCount)
	s.Require().Equal(12, statBannerC.ShowCount)

	s.Require().Equal(976, statBannerA.ClickCount)
	s.Require().Equal(0, statBannerB.ClickCount)
	s.Require().Equal(0, statBannerC.ClickCount)
}

func (s *BannerRotationSuite) NoRotationsCreatedForSlot() {
	err := s.client.GetStatus()
	s.Require().NoError(err)

	segment, err := s.client.CreateSegment("segment")
	s.Require().NoError(err)

	slot, err := s.client.CreateSlot("slot")
	s.Require().NoError(err)

	bannerA, err := s.client.CreateBanner("bannerA")
	s.Require().NoError(err)

	err = s.client.CreateRotation(slot, bannerA)
	s.Require().NoError(err)

	_, err = s.client.Choice(slot, segment)
	s.Require().NoError(err)

	err = s.client.DeleteRotation(slot, bannerA)
	s.Require().NoError(err)

	_, err = s.client.Choice(slot, segment)
	s.Require().Error(err)
}

func (s *BannerRotationSuite) ClicksMoreThenShows() {
	err := s.client.GetStatus()
	s.Require().NoError(err)

	segment, err := s.client.CreateSegment("segment")
	s.Require().NoError(err)

	slot, err := s.client.CreateSlot("slot")
	s.Require().NoError(err)

	bannerA, err := s.client.CreateBanner("bannerA")
	s.Require().NoError(err)

	err = s.client.CreateRotation(slot, bannerA)
	s.Require().NoError(err)

	statBannerA, err := s.client.GetStat(bannerA, segment)
	s.Require().NoError(err)

	s.Require().Equal(0, statBannerA.ShowCount)
	s.Require().Equal(0, statBannerA.ClickCount)

	err = s.client.Click(slot, bannerA, segment)
	s.Require().NoError(err)

	_, err = s.client.Choice(slot, segment)
	s.Require().Error(err)
}

func (s *BannerRotationSuite) TearDownTest() {
}

func (s *BannerRotationSuite) TearDownSuite() {
}

func TestBannerRotationSuite(t *testing.T) {
	suite.Run(t, new(BannerRotationSuite))
}
