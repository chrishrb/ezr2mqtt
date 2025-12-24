package api

type HAComponent string

var (
	HAComponentSensor HAComponent = "sensor"
	HAComponentNumber HAComponent = "number"
	HAComponentSelect HAComponent = "select"
)

type HASensorDiscovery struct {
	Name                string    `json:"name,omitempty"`
	UniqueID            string    `json:"unique_id,omitempty"`
	StateTopic          string    `json:"state_topic"`
	UnitOfMeasurement   string    `json:"unit_of_measurement,omitempty"`
	DeviceClass         string    `json:"device_class,omitempty"`
	StateClass          string    `json:"state_class,omitempty"`
	ValueTemplate       string    `json:"value_template,omitempty"`
	AvailabilityTopic   string    `json:"availability_topic,omitempty"`
	PayloadAvailable    string    `json:"payload_available,omitempty"`
	PayloadNotAvailable string    `json:"payload_not_available,omitempty"`
	Device              *HADevice `json:"device,omitempty"`
	ExpireAfter         int       `json:"expire_after,omitempty"`
	Icon                string    `json:"icon,omitempty"`
	JSONAttributesTopic string    `json:"json_attributes_topic,omitempty"`
	CommandTopic        string    `json:"command_topic,omitempty"`
	Minimum             float64   `json:"min,omitempty"`
	Maximum             float64   `json:"max,omitempty"`
	Step                float64   `json:"step,omitempty"`
	Mode                string    `json:"mode,omitempty"`
	Options             []string  `json:"options,omitempty"`
}

type HADevice struct {
	Identifiers  []string `json:"identifiers,omitempty"`
	Name         string   `json:"name,omitempty"`
	Manufacturer string   `json:"manufacturer,omitempty"`
	Model        string   `json:"model,omitempty"`
	SwVersion    string   `json:"sw_version,omitempty"`
}
