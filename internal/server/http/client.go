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
		return errors.New("service status is not OK")
	}
	return nil
}

func (c *Client) CreateBanner(description string) (string, error) {
	return c.CreateItem(Banner, description)
}

func (c *Client) CreateSlot(description string) (string, error) {
	return c.CreateItem(Slot, description)
}

func (c *Client) CreateSegment(description string) (string, error) {
	return c.CreateItem(Segment, description)
}

func (c *Client) CreateItem(item ItemType, description string) (string, error) {
	desc := Description{Description: description}
	b, err := json.Marshal(desc)
	if err != nil {
		return "", err
	}

	body := bytes.NewReader(b)

	var url string
	switch item {
	case Banner:
		url = "http://" + c.addr + "/banner"
	case Slot:
		url = "http://" + c.addr + "/slot"
	case Segment:
		url = "http://" + c.addr + "/segment"
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", url, body)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: c.timeout}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	responseID := ResponseID{}
	err = json.Unmarshal(responseBody, &responseID)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("error while creating")
	}
	return responseID.ID, nil
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
		return errors.New("error when adding rotation")
	}
	return nil
}

func (c *Client) Click(slotID, bannerID, segmentID string) error {
	url := "http://" + c.addr + "/click/" + slotID + "/" + bannerID + "/" + segmentID

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
		return errors.New("error when counting click for the banner")
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

	responseID := ResponseID{}
	err = json.Unmarshal(body, &responseID)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("error getting banner to show")
	}
	return responseID.ID, nil
}

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
		return ResponseStat{}, errors.New("error while getting stat")
	}
	return responseStat, nil
}
