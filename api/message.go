package api

type Message struct {
	Room int
	Type string
	Data any
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
