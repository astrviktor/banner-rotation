package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/astrviktor/banner-rotation/internal/storage"
)

type Client struct {
	addr    string
	timeout time.Duration
}

func NewClient(host string, port string, timeout time.Duration) *Client {
	return &Client{net.JoinHostPort(host, port), timeout}
}

func (c *Client) GetStatus() error {
	req, err := http.NewRequestWithContext(context.Background(), "GET", "http://"+c.addr+"/status", nil)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: c.timeout}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK || string(body) != "OK" {
		return errors.New("status сервиса не OK")
	}
	return nil
}

func (c *Client) CreateBanner(id, description string) error {
	banner := storage.Banner{ID: id, Description: description}
	b, err := json.Marshal(banner)
	if err != nil {
		return err
	}

	body := bytes.NewReader(b)

	req, err := http.NewRequestWithContext(context.Background(), "POST", "http://"+c.addr+"/banner", body)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: c.timeout}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("ошибка при добавлении баннера")
	}
	return nil
}

func (c *Client) CreateSlot(id, description string) error {
	slot := storage.Slot{ID: id, Description: description}
	b, err := json.Marshal(slot)
	if err != nil {
		return err
	}

	body := bytes.NewReader(b)

	req, err := http.NewRequestWithContext(context.Background(), "POST", "http://"+c.addr+"/slot", body)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: c.timeout}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("ошибка при добавлении слота")
	}
	return nil
}

func (c *Client) CreateSegment(id, description string) error {
	segment := storage.Segment{ID: id, Description: description}
	b, err := json.Marshal(segment)
	if err != nil {
		return err
	}

	body := bytes.NewReader(b)

	req, err := http.NewRequestWithContext(context.Background(), "POST", "http://"+c.addr+"/segment", body)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: c.timeout}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("ошибка при добавлении сегмента")
	}
	return nil
}

func (c *Client) CreateRotation(slotID, bannerID string) error {
	url := "http://" + c.addr + "/rotation/" + slotID + "/" + bannerID

	req, err := http.NewRequestWithContext(context.Background(), "POST", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: c.timeout}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("ошибка при добавлении ротации")
	}
	return nil
}

func (c *Client) Choice(slotID, segmentID string) (string, error) {
	url := "http://" + c.addr + "/choice/" + slotID + "/" + segmentID

	req, err := http.NewRequestWithContext(context.Background(), "POST", url, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: c.timeout}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	responseBannerID := ResponseBannerID{}
	err = json.Unmarshal(body, &responseBannerID)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("ошибка при получении баннера для показа")
	}
	return responseBannerID.BannerID, nil
}

// mux.HandleFunc("/click/", Logging(s.handleClick))

func (c *Client) GetStat(bannerID, segmentID string) (ResponseStat, error) {
	url := "http://" + c.addr + "/stat/" + bannerID + "/" + segmentID

	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return ResponseStat{}, err
	}

	client := &http.Client{Timeout: c.timeout}

	resp, err := client.Do(req)
	if err != nil {
		return ResponseStat{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ResponseStat{}, err
	}

	responseStat := ResponseStat{}
	err = json.Unmarshal(body, &responseStat)
	if err != nil {
		return ResponseStat{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return ResponseStat{}, errors.New("ошибка при получении статистики")
	}
	return responseStat, nil
}
