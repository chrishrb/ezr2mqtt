package api

import "fmt"

type Message struct {
	Room int
	Type string
	Data string
}

type RoomDiscovery struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ClimateDiscovery struct {
	// Identity
	Name string `json:"name"`
	ID   string `json:"id"`
	Type string `json:"type"`

	// Rooms
	Rooms []RoomDiscovery `json:"rooms"`
}

func FormatFloat(f float64) string {
	return fmt.Sprintf("%.2f", f)
}
