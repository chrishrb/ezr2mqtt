package http

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chrishrb/ezr2mqtt/transport"
	"github.com/stretchr/testify/assert"
)

func TestNewHTTPClient(t *testing.T) {
	hostname := "192.168.1.100"
	client := NewHTTPClient(hostname)

	assert.NotNil(t, client)
	assert.Equal(t, hostname, client.Hostname)
	assert.NotNil(t, client.Client)
}

func TestHTTPClient_Connect_Success(t *testing.T) {
	// Create mock response
	mockMessage := &transport.Message{
		XMLName: xml.Name{Local: "Devices"},
		Device: transport.Device{
			ID:   "TEST-123",
			Type: "EZR",
			Name: "Test Device",
			HeatAreas: []transport.HeatArea{
				{Nr: 1, Name: "Room 1", TTarget: 22.0, TActual: 21.5},
				{Nr: 2, Name: "Room 2", TTarget: 20.0, TActual: 19.5},
			},
		},
	}

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/data/static.xml", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "close", r.Header.Get("Connection"))

		w.Header().Set("Content-Type", "application/xml")
		xmlData, _ := xml.Marshal(mockMessage)
		_, _ = w.Write([]byte(xml.Header))
		_, _ = w.Write(xmlData)
	}))
	defer server.Close()

	// Extract hostname from server URL (remove http://)
	hostname := server.URL[7:] // Remove "http://"
	client := NewHTTPClient(hostname)

	result, err := client.Connect()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "TEST-123", result.Device.ID)
	assert.Equal(t, "EZR", result.Device.Type)
	assert.Equal(t, "Test Device", result.Device.Name)
	assert.Len(t, result.Device.HeatAreas, 2)
}

func TestHTTPClient_Connect_InvalidXML(t *testing.T) {
	// Create test server that returns invalid XML
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		_, _ = w.Write([]byte("invalid xml content"))
	}))
	defer server.Close()

	hostname := server.URL[7:]
	client := NewHTTPClient(hostname)

	result, err := client.Connect()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to decode XML")
}

func TestHTTPClient_Connect_ServerError(t *testing.T) {
	// Create test server that returns 500 error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	hostname := server.URL[7:]
	client := NewHTTPClient(hostname)

	// Even with 500 error, Connect tries to decode the body
	// This will fail because there's no valid XML
	result, err := client.Connect()
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestHTTPClient_Send_Success(t *testing.T) {
	var receivedBody []byte
	var receivedHeaders http.Header

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/data/changes.xml", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "close", r.Header.Get("Connection"))
		assert.Equal(t, "application/xml", r.Header.Get("Content-Type"))

		receivedHeaders = r.Header

		// Read the body
		var buf [1024]byte
		n, _ := r.Body.Read(buf[:])
		receivedBody = buf[:n]

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	hostname := server.URL[7:]
	client := NewHTTPClient(hostname)

	msg := &transport.Message{
		Device: transport.Device{
			ID: "TEST-123",
			HeatAreas: []transport.HeatArea{
				{Nr: 1, TTarget: 23.0},
			},
		},
	}

	err := client.Send(msg)

	assert.NoError(t, err)
	assert.NotNil(t, receivedBody)
	assert.Contains(t, string(receivedBody), "<?xml version")
	assert.Contains(t, string(receivedBody), "TEST-123")
	assert.Equal(t, "application/xml", receivedHeaders.Get("Content-Type"))
}

func TestHTTPClient_Send_VerifyXMLStructure(t *testing.T) {
	var receivedMessage *transport.Message

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Decode the received XML
		receivedMessage = &transport.Message{}
		err := xml.NewDecoder(r.Body).Decode(receivedMessage)
		assert.NoError(t, err)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	hostname := server.URL[7:]
	client := NewHTTPClient(hostname)

	sentMsg := &transport.Message{
		Device: transport.Device{
			ID: "DEVICE-456",
			HeatAreas: []transport.HeatArea{
				{Nr: 2, TTarget: 24.5},
				{Nr: 3, TTarget: 19.0},
			},
		},
	}

	err := client.Send(sentMsg)

	assert.NoError(t, err)
	assert.NotNil(t, receivedMessage)
	assert.Equal(t, "DEVICE-456", receivedMessage.Device.ID)
	assert.Len(t, receivedMessage.Device.HeatAreas, 2)
	assert.Equal(t, 2, receivedMessage.Device.HeatAreas[0].Nr)
	assert.Equal(t, 24.5, receivedMessage.Device.HeatAreas[0].TTarget)
	assert.Equal(t, 3, receivedMessage.Device.HeatAreas[1].Nr)
	assert.Equal(t, 19.0, receivedMessage.Device.HeatAreas[1].TTarget)
}

func TestHTTPClient_Connect_InvalidHostname(t *testing.T) {
	client := NewHTTPClient("invalid-hostname-that-does-not-exist:9999")

	result, err := client.Connect()

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestHTTPClient_Send_InvalidHostname(t *testing.T) {
	client := NewHTTPClient("invalid-hostname-that-does-not-exist:9999")

	msg := &transport.Message{
		Device: transport.Device{
			ID: "TEST",
		},
	}

	err := client.Send(msg)

	assert.Error(t, err)
}

func TestHTTPClient_Connect_EmptyResponse(t *testing.T) {
	// Create test server that returns empty response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		// Return nothing
	}))
	defer server.Close()

	hostname := server.URL[7:]
	client := NewHTTPClient(hostname)

	result, err := client.Connect()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to decode XML")
}

func TestHTTPClient_Send_ServerUnavailable(t *testing.T) {
	// Create and immediately close server to simulate unavailability
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	hostname := server.URL[7:]
	server.Close() // Close immediately

	client := NewHTTPClient(hostname)

	msg := &transport.Message{
		Device: transport.Device{ID: "TEST"},
	}

	err := client.Send(msg)

	assert.Error(t, err)
}

func TestHTTPClient_Connect_ComplexDevice(t *testing.T) {
	// Create a more complex mock message
	mockMessage := &transport.Message{
		XMLName: xml.Name{Local: "Devices"},
		Device: transport.Device{
			ID:       "COMPLEX-123",
			Type:     "EZR",
			Name:     "Complex Device",
			DateTime: "2025-12-23 10:00:00",
			HeatAreas: []transport.HeatArea{
				{
					Nr:         1,
					Name:       "Living Room",
					TTarget:    22.0,
					TActual:    21.5,
					TTargetMin: 15.0,
					TTargetMax: 30.0,
					Mode:       1,
					State:      1,
				},
			},
			HeatCtrls: []transport.HeatCtrl{
				{Nr: 1, InUse: 1, HeatAreaNr: 1, Actor: 50},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		xmlData, _ := xml.Marshal(mockMessage)
		_, _ = w.Write([]byte(xml.Header))
		_, _ = w.Write(xmlData)
	}))
	defer server.Close()

	hostname := server.URL[7:]
	client := NewHTTPClient(hostname)

	result, err := client.Connect()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "COMPLEX-123", result.Device.ID)
	assert.Len(t, result.Device.HeatAreas, 1)
	assert.Equal(t, "Living Room", result.Device.HeatAreas[0].Name)
	assert.Len(t, result.Device.HeatCtrls, 1)
}
