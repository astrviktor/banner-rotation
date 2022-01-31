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
}

func (s *BannerRotationSuite) TearDownTest() {
}

func (s *BannerRotationSuite) TearDownSuite() {
}

func TestBannerRotationSuite(t *testing.T) {
	suite.Run(t, new(BannerRotationSuite))
}
