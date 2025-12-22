package http

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/chrishrb/ezr2mqtt/transport"
)

type HTTPClient struct {
	Hostname string
	Client   *http.Client
}

func NewHTTPClient(hostname string) *HTTPClient {
	return &HTTPClient{
		Hostname: hostname,
		Client:   &http.Client{},
	}
}

func (c *HTTPClient) Connect() (*transport.Message, error) {
	url := fmt.Sprintf("http://%s/data/static.xml", c.Hostname)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Connection", "close")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var msg transport.Message
	err = xml.NewDecoder(resp.Body).Decode(&msg)
	if err != nil {
		return nil, fmt.Errorf("failed to decode XML: %w", err)
	}
	return &msg, nil
}

func (c *HTTPClient) Send(msg *transport.Message) error {
	out, err := xml.MarshalIndent(msg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal XML: %w", err)
	}
	xmlData := []byte(xml.Header + string(out))

	url := fmt.Sprintf("http://%s/data/changes.xml", c.Hostname)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(xmlData))
	if err != nil {
		return err
	}
	req.Header.Set("Connection", "close")
	req.Header.Set("Content-Type", "application/xml")

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	return nil
}
