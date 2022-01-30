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

	err = s.client.CreateSegment("men", "men")
	s.Require().NoError(err)

	err = s.client.CreateSlot("slot", "slot")
	s.Require().NoError(err)

	err = s.client.CreateBanner("bannerA", "bannerA")
	s.Require().NoError(err)

	err = s.client.CreateBanner("bannerB", "bannerB")
	s.Require().NoError(err)

	err = s.client.CreateRotation("slot", "bannerA")
	s.Require().NoError(err)

	err = s.client.CreateRotation("slot", "bannerB")
	s.Require().NoError(err)

	for i := 0; i < 1000; i++ {
		_, err := s.client.Choice("slot", "men")
		s.Require().NoError(err)
	}

	statA, err := s.client.GetStat("bannerA", "men")
	s.Require().NoError(err)
	statB, err := s.client.GetStat("bannerB", "men")
	s.Require().NoError(err)

	s.Require().Equal(500, statA.Show)
	s.Require().Equal(500, statB.Show)
}

func (s *BannerRotationSuite) TearDownTest() {
}

func (s *BannerRotationSuite) TearDownSuite() {
}

func TestBannerRotationSuite(t *testing.T) {
	suite.Run(t, new(BannerRotationSuite))
}
