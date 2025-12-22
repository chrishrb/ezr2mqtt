package transport

import (
	"encoding/xml"
)

type Message struct {
	XMLName xml.Name `xml:"Devices"`
	Device  Device   `xml:"Device"`
}

func NewMessage(device Device) *Message {
	return &Message{
		Device: device,
	}
}

type Vacation struct {
	State     int    `xml:"VACATION_STATE"`
	StartDate string `xml:"START_DATE"`
	StartTime string `xml:"START_TIME"`
	EndDate   string `xml:"END_DATE"`
	EndTime   string `xml:"END_TIME"`
}

type Network struct {
	MAC           string `xml:"MAC"`
	DHCP          int    `xml:"DHCP"`
	IPv6Active    int    `xml:"IPV6ACTIVE"`
	IPv4Actual    string `xml:"IPV4ACTUAL"`
	IPv4Set       string `xml:"IPV4SET"`
	IPv6Actual    string `xml:"IPV6ACTUAL"`
	IPv6Set       string `xml:"IPV6SET"`
	NetmaskActual string `xml:"NETMASKACTUAL"`
	NetmaskSet    string `xml:"NETMASKSET"`
	DNS           string `xml:"DNS"`
	Gateway       string `xml:"GATEWAY"`
}

type Cloud struct {
	UserID           string `xml:"USERID"`
	Password         string `xml:"PASSWORD"`
	M2MServerPort    int    `xml:"M2MSERVERPORT"`
	M2MLocalPort     int    `xml:"M2MLOCALPORT"`
	M2MHTTPPort      int    `xml:"M2MHTTPPORT"`
	M2MHTTPSPort     int    `xml:"M2MHTTPSPORT"`
	M2MServerAddress string `xml:"M2MSERVERADDRESS"`
	M2MActive        int    `xml:"M2MACTIVE"`
	M2MState         string `xml:"M2MSTATE"`
}

type KWLCtrl struct {
	Visible    int    `xml:"KWL_CONTROL_VISIBLE"`
	Present    int    `xml:"KWL_PRESENT"`
	Connection int    `xml:"KWL_CONNECTION"`
	URL        string `xml:"KWL_URL"`
	Port       int    `xml:"KWL_PORT"`
	Status     int    `xml:"KWL_STATUS"`
	FlowCtrl   int    `xml:"KWL_FLOWCTRL"`
}

type Code struct {
	Expert string `xml:"EXPERT"`
}

type Program struct {
	ShiftPrograms []ShiftProgram `xml:"SHIFT_PROGRAM"`
}

type ShiftProgram struct {
	Nr           int    `xml:"nr,attr"`
	ShiftingTime int    `xml:"shiftingtime,attr"`
	Start        string `xml:"START"`
	End          string `xml:"END"`
}

type PumpOutput struct {
	LocalGlobal   int `xml:"LOCALGLOBAL"`
	Type          int `xml:"PUMP_OUTPUT_TYPE"`
	LeadTime      int `xml:"PUMP_LEADTIME"`
	StoppingTime  int `xml:"PUMP_STOPPINGTIME"`
	OperationMode int `xml:"PUMP_OPERATIONMODE"`
	MinRuntime    int `xml:"MINRUNTIME"`
	MinStandstill int `xml:"MINSTANDSTILL"`
}

type Relais struct {
	Function      int `xml:"FUNCTION"`
	LeadTime      int `xml:"RELAIS_LEADTIME"`
	StoppingTime  int `xml:"RELAIS_STOPPINGTIME"`
	OperationMode int `xml:"RELAIS_OPERATIONMODE"`
}

type ChangeoverFunc struct {
	Mode int `xml:"CHANGEOVER_FUNC_MODE"`
}

type EmergencyMode struct {
	Time     int `xml:"EMERGENCYMODE_TIME"`
	PWMCycle int `xml:"PWMCYCLE"`
	PWMHeat  int `xml:"PWMHEAT"`
	PWMCool  int `xml:"PWMCOOL"`
}

type ValveProtect struct {
	Time     int `xml:"VALVEPROTECT_TIME"`
	Duration int `xml:"DURATION"`
}

type PumpProtect struct {
	Time     int `xml:"PUMPPROTECT_TIME"`
	Duration int `xml:"DURATION"`
}

type HeatArea struct {
	Nr   int    `xml:"nr,attr"`
	Name string `xml:"HEATAREA_NAME"`

	Mode int `xml:"HEATAREA_MODE"`

	TActual     float64 `xml:"T_ACTUAL"`
	TActualExt  float64 `xml:"T_ACTUAL_EXT"`
	TTarget     float64 `xml:"T_TARGET"`
	TTargetBase float64 `xml:"T_TARGET_BASE"`

	State int `xml:"HEATAREA_STATE"`

	ProgramSource  int `xml:"PROGRAM_SOURCE"`
	ProgramWeek    int `xml:"PROGRAM_WEEK"`
	ProgramWeekend int `xml:"PROGRAM_WEEKEND"`

	Party              int `xml:"PARTY"`
	PartyRemainingTime int `xml:"PARTY_REMAININGTIME"`
	Presence           int `xml:"PRESENCE"`

	TTargetMin float64 `xml:"T_TARGET_MIN"`
	TTargetMax float64 `xml:"T_TARGET_MAX"`

	RPMMotor int     `xml:"RPM_MOTOR"`
	Offset   float64 `xml:"OFFSET"`

	THeatDay   float64 `xml:"T_HEAT_DAY"`
	THeatNight float64 `xml:"T_HEAT_NIGHT"`
	TCoolDay   float64 `xml:"T_COOL_DAY"`
	TCoolNight float64 `xml:"T_COOL_NIGHT"`
	TFloorDay  float64 `xml:"T_FLOOR_DAY"`

	HeatingSystem int `xml:"HEATINGSYSTEM"`

	BlockHC int `xml:"BLOCK_HC"`

	IsLocked      int    `xml:"ISLOCKED"`
	LockCode      string `xml:"LOCK_CODE"`
	LockAvailable int    `xml:"LOCK_AVAILABLE"`

	Light      int `xml:"LIGHT"`
	SensorExt  int `xml:"SENSOR_EXT"`
	Adjustable int `xml:"T_TARGET_ADJUSTABLE"`
}

type HeatCtrl struct {
	Nr           int `xml:"nr,attr"`
	InUse        int `xml:"INUSE"`
	HeatAreaNr   int `xml:"HEATAREA_NR"`
	Actor        int `xml:"ACTOR"`
	ActorPercent int `xml:"ACTOR_PERCENT"`
	State        int `xml:"HEATCTRL_STATE"`
}

type IODevice struct {
	Nr int `xml:"nr,attr"`

	Type int `xml:"IODEVICE_TYPE"`
	ID   int `xml:"IODEVICE_ID"`

	VersHW string `xml:"IODEVICE_VERS_HW"`
	VersSW string `xml:"IODEVICE_VERS_SW"`

	HeatAreaNr int `xml:"HEATAREA_NR"`

	SignalStrength int `xml:"SIGNALSTRENGTH"`
	Battery        int `xml:"BATTERY"`

	State    int `xml:"IODEVICE_STATE"`
	ComError int `xml:"IODEVICE_COMERROR"`
	IsOn     int `xml:"ISON"`
}

type Device struct {
	// Identification
	ID     string `xml:"ID"`
	Type   string `xml:"TYPE"`
	Name   string `xml:"NAME"`
	Origin string `xml:"ORIGIN"`

	// System
	ErrorCount int    `xml:"ERRORCOUNT"`
	DateTime   string `xml:"DATETIME"`
	DayOfWeek  int    `xml:"DAYOFWEEK"`
	TimeZone   int    `xml:"TIMEZONE"`

	NTPSync int `xml:"NTPTIMESYNC"`

	VersSWSTM string `xml:"VERS_SW_STM"`
	VersSWETH string `xml:"VERS_SW_ETH"`
	VersHW    string `xml:"VERS_HW"`

	TemperatureUnit int `xml:"TEMPERATUREUNIT"`
	SummerWinter    int `xml:"SUMMERWINTER"`
	TPS             int `xml:"TPS"`
	Limiter         int `xml:"LIMITER"`

	MasterID   string `xml:"MASTERID"`
	Changeover int    `xml:"CHANGEOVER"`
	Cooling    int    `xml:"COOLING"`
	Mode       int    `xml:"MODE"`

	OperationModeActor int `xml:"OPERATIONMODE_ACTOR"`

	Antifreeze     int     `xml:"ANTIFREEZE"`
	AntifreezeTemp float64 `xml:"ANTIFREEZE_TEMP"`

	FirstOpenTime int `xml:"FIRSTOPEN_TIME"`
	SmartStart    int `xml:"SMARTSTART"`

	EcoDiff       float64 `xml:"ECO_DIFF"`
	EcoInputMode  int     `xml:"ECO_INPUTMODE"`
	EcoInputState int     `xml:"ECO_INPUT_STATE"`

	THeatVacation float64 `xml:"T_HEAT_VACATION"`

	Vacation Vacation `xml:"VACATION"`

	Network Network `xml:"NETWORK"`
	Cloud   Cloud   `xml:"CLOUD"`
	KWLCtrl KWLCtrl `xml:"KWLCTRL"`

	Code    Code    `xml:"CODE"`
	Program Program `xml:"PROGRAM"`

	PumpOutput     PumpOutput     `xml:"PUMP_OUTPUT"`
	Relais         Relais         `xml:"RELAIS"`
	ChangeoverFunc ChangeoverFunc `xml:"CHANGEOVER_FUNC"`
	EmergencyMode  EmergencyMode  `xml:"EMERGENCYMODE"`
	ValveProtect   ValveProtect   `xml:"VALVEPROTECT"`
	PumpProtect    PumpProtect    `xml:"PUMPPROTECT"`

	HeatAreas []HeatArea `xml:"HEATAREA"`
	HeatCtrls []HeatCtrl `xml:"HEATCTRL"`
	IODevices []IODevice `xml:"IODEVICE"`
}
