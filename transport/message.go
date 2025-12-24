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
	State     *int    `xml:"VACATION_STATE,omitempty"`
	StartDate *string `xml:"START_DATE,omitempty"`
	StartTime *string `xml:"START_TIME,omitempty"`
	EndDate   *string `xml:"END_DATE,omitempty"`
	EndTime   *string `xml:"END_TIME,omitempty"`
}

type Network struct {
	MAC           *string `xml:"MAC,omitempty"`
	DHCP          *int    `xml:"DHCP,omitempty"`
	IPv6Active    *int    `xml:"IPV6ACTIVE,omitempty"`
	IPv4Actual    *string `xml:"IPV4ACTUAL,omitempty"`
	IPv4Set       *string `xml:"IPV4SET,omitempty"`
	IPv6Actual    *string `xml:"IPV6ACTUAL,omitempty"`
	IPv6Set       *string `xml:"IPV6SET,omitempty"`
	NetmaskActual *string `xml:"NETMASKACTUAL,omitempty"`
	NetmaskSet    *string `xml:"NETMASKSET,omitempty"`
	DNS           *string `xml:"DNS,omitempty"`
	Gateway       *string `xml:"GATEWAY,omitempty"`
}

type Cloud struct {
	UserID           *string `xml:"USERID,omitempty"`
	Password         *string `xml:"PASSWORD,omitempty"`
	M2MServerPort    *int    `xml:"M2MSERVERPORT,omitempty"`
	M2MLocalPort     *int    `xml:"M2MLOCALPORT,omitempty"`
	M2MHTTPPort      *int    `xml:"M2MHTTPPORT,omitempty"`
	M2MHTTPSPort     *int    `xml:"M2MHTTPSPORT,omitempty"`
	M2MServerAddress *string `xml:"M2MSERVERADDRESS,omitempty"`
	M2MActive        *int    `xml:"M2MACTIVE,omitempty"`
	M2MState         *string `xml:"M2MSTATE,omitempty"`
}

type KWLCtrl struct {
	Visible    *int    `xml:"KWL_CONTROL_VISIBLE,omitempty"`
	Present    *int    `xml:"KWL_PRESENT,omitempty"`
	Connection *int    `xml:"KWL_CONNECTION,omitempty"`
	URL        *string `xml:"KWL_URL,omitempty"`
	Port       *int    `xml:"KWL_PORT,omitempty"`
	Status     *int    `xml:"KWL_STATUS,omitempty"`
	FlowCtrl   *int    `xml:"KWL_FLOWCTRL,omitempty"`
}

type Code struct {
	Expert *string `xml:"EXPERT,omitempty"`
}

type Program struct {
	ShiftPrograms *[]ShiftProgram `xml:"SHIFT_PROGRAM,omitempty"`
}

type ShiftProgram struct {
	Nr           *int    `xml:"nr,attr,omitempty"`
	ShiftingTime *int    `xml:"shiftingtime,attr,omitempty"`
	Start        *string `xml:"START,omitempty"`
	End          *string `xml:"END,omitempty"`
}

type PumpOutput struct {
	LocalGlobal   *int `xml:"LOCALGLOBAL,omitempty"`
	Type          *int `xml:"PUMP_OUTPUT_TYPE,omitempty"`
	LeadTime      *int `xml:"PUMP_LEADTIME,omitempty"`
	StoppingTime  *int `xml:"PUMP_STOPPINGTIME,omitempty"`
	OperationMode *int `xml:"PUMP_OPERATIONMODE,omitempty"`
	MinRuntime    *int `xml:"MINRUNTIME,omitempty"`
	MinStandstill *int `xml:"MINSTANDSTILL,omitempty"`
}

type Relais struct {
	Function      *int `xml:"FUNCTION,omitempty"`
	LeadTime      *int `xml:"RELAIS_LEADTIME,omitempty"`
	StoppingTime  *int `xml:"RELAIS_STOPPINGTIME,omitempty"`
	OperationMode *int `xml:"RELAIS_OPERATIONMODE,omitempty"`
}

type ChangeoverFunc struct {
	Mode *int `xml:"CHANGEOVER_FUNC_MODE,omitempty"`
}

type EmergencyMode struct {
	Time     *int `xml:"EMERGENCYMODE_TIME,omitempty"`
	PWMCycle *int `xml:"PWMCYCLE,omitempty"`
	PWMHeat  *int `xml:"PWMHEAT,omitempty"`
	PWMCool  *int `xml:"PWMCOOL,omitempty"`
}

type ValveProtect struct {
	Time     *int `xml:"VALVEPROTECT_TIME,omitempty"`
	Duration *int `xml:"DURATION,omitempty"`
}

type PumpProtect struct {
	Time     *int `xml:"PUMPPROTECT_TIME,omitempty"`
	Duration *int `xml:"DURATION,omitempty"`
}

type HeatArea struct {
	Nr   *int    `xml:"nr,attr,omitempty"`
	Name *string `xml:"HEATAREA_NAME,omitempty"`

	Mode *int `xml:"HEATAREA_MODE,omitempty"`

	TActual     *float64 `xml:"T_ACTUAL,omitempty"`
	TActualExt  *float64 `xml:"T_ACTUAL_EXT,omitempty"`
	TTarget     *float64 `xml:"T_TARGET,omitempty"`
	TTargetBase *float64 `xml:"T_TARGET_BASE,omitempty"`

	State *int `xml:"HEATAREA_STATE,omitempty"`

	ProgramSource  *int `xml:"PROGRAM_SOURCE,omitempty"`
	ProgramWeek    *int `xml:"PROGRAM_WEEK,omitempty"`
	ProgramWeekend *int `xml:"PROGRAM_WEEKEND,omitempty"`

	Party              *int `xml:"PARTY,omitempty"`
	PartyRemainingTime *int `xml:"PARTY_REMAININGTIME,omitempty"`
	Presence           *int `xml:"PRESENCE,omitempty"`

	TTargetMin *float64 `xml:"T_TARGET_MIN,omitempty"`
	TTargetMax *float64 `xml:"T_TARGET_MAX,omitempty"`

	RPMMotor *int     `xml:"RPM_MOTOR,omitempty"`
	Offset   *float64 `xml:"OFFSET,omitempty"`

	THeatDay   *float64 `xml:"T_HEAT_DAY,omitempty"`
	THeatNight *float64 `xml:"T_HEAT_NIGHT,omitempty"`
	TCoolDay   *float64 `xml:"T_COOL_DAY,omitempty"`
	TCoolNight *float64 `xml:"T_COOL_NIGHT,omitempty"`
	TFloorDay  *float64 `xml:"T_FLOOR_DAY,omitempty"`

	HeatingSystem *int `xml:"HEATINGSYSTEM,omitempty"`

	BlockHC *int `xml:"BLOCK_HC,omitempty"`

	IsLocked      *int    `xml:"ISLOCKED,omitempty"`
	LockCode      *string `xml:"LOCK_CODE,omitempty"`
	LockAvailable *int    `xml:"LOCK_AVAILABLE,omitempty"`

	Light      *int `xml:"LIGHT,omitempty"`
	SensorExt  *int `xml:"SENSOR_EXT,omitempty"`
	Adjustable *int `xml:"T_TARGET_ADJUSTABLE,omitempty"`
}

type HeatCtrl struct {
	Nr           *int `xml:"nr,attr,omitempty"`
	InUse        *int `xml:"INUSE,omitempty"`
	HeatAreaNr   *int `xml:"HEATAREA_NR,omitempty"`
	Actor        *int `xml:"ACTOR,omitempty"`
	ActorPercent *int `xml:"ACTOR_PERCENT,omitempty"`
	State        *int `xml:"HEATCTRL_STATE,omitempty"`
}

type IODevice struct {
	Nr *int `xml:"nr,attr,omitempty"`

	Type *int `xml:"IODEVICE_TYPE,omitempty"`
	ID   *int `xml:"IODEVICE_ID,omitempty"`

	VersHW *string `xml:"IODEVICE_VERS_HW,omitempty"`
	VersSW *string `xml:"IODEVICE_VERS_SW,omitempty"`

	HeatAreaNr *int `xml:"HEATAREA_NR,omitempty"`

	SignalStrength *int `xml:"SIGNALSTRENGTH,omitempty"`
	Battery        *int `xml:"BATTERY,omitempty"`

	State    *int `xml:"IODEVICE_STATE,omitempty"`
	ComError *int `xml:"IODEVICE_COMERROR,omitempty"`
	IsOn     *int `xml:"ISON,omitempty"`
}

type Device struct {
	// Identification
	ID     *string `xml:"ID,omitempty"`
	Type   *string `xml:"TYPE,omitempty"`
	Name   *string `xml:"NAME,omitempty"`
	Origin *string `xml:"ORIGIN,omitempty"`

	// System
	ErrorCount *int    `xml:"ERRORCOUNT,omitempty"`
	DateTime   *string `xml:"DATETIME,omitempty"`
	DayOfWeek  *int    `xml:"DAYOFWEEK,omitempty"`
	TimeZone   *int    `xml:"TIMEZONE,omitempty"`

	NTPSync *int `xml:"NTPTIMESYNC,omitempty"`

	VersSWSTM *string `xml:"VERS_SW_STM,omitempty"`
	VersSWETH *string `xml:"VERS_SW_ETH,omitempty"`
	VersHW    *string `xml:"VERS_HW,omitempty"`

	TemperatureUnit *int `xml:"TEMPERATUREUNIT,omitempty"`
	SummerWinter    *int `xml:"SUMMERWINTER,omitempty"`
	TPS             *int `xml:"TPS,omitempty"`
	Limiter         *int `xml:"LIMITER,omitempty"`

	MasterID   *string `xml:"MASTERID,omitempty"`
	Changeover *int    `xml:"CHANGEOVER,omitempty"`
	Cooling    *int    `xml:"COOLING,omitempty"`
	Mode       *int    `xml:"MODE,omitempty"`

	OperationModeActor *int `xml:"OPERATIONMODE_ACTOR,omitempty"`

	Antifreeze     *int     `xml:"ANTIFREEZE,omitempty"`
	AntifreezeTemp *float64 `xml:"ANTIFREEZE_TEMP,omitempty"`

	FirstOpenTime *int `xml:"FIRSTOPEN_TIME,omitempty"`
	SmartStart    *int `xml:"SMARTSTART,omitempty"`

	EcoDiff       *float64 `xml:"ECO_DIFF,omitempty"`
	EcoInputMode  *int     `xml:"ECO_INPUTMODE,omitempty"`
	EcoInputState *int     `xml:"ECO_INPUT_STATE,omitempty"`

	THeatVacation *float64 `xml:"T_HEAT_VACATION,omitempty"`

	Vacation *Vacation `xml:"VACATION,omitempty"`

	Network *Network `xml:"NETWORK,omitempty"`
	Cloud   *Cloud   `xml:"CLOUD,omitempty"`
	KWLCtrl *KWLCtrl `xml:"KWLCTRL,omitempty"`

	Code    *Code    `xml:"CODE,omitempty"`
	Program *Program `xml:"PROGRAM,omitempty"`

	PumpOutput     *PumpOutput     `xml:"PUMP_OUTPUT,omitempty"`
	Relais         *Relais         `xml:"RELAIS,omitempty"`
	ChangeoverFunc *ChangeoverFunc `xml:"CHANGEOVER_FUNC,omitempty"`
	EmergencyMode  *EmergencyMode  `xml:"EMERGENCYMODE,omitempty"`
	ValveProtect   *ValveProtect   `xml:"VALVEPROTECT,omitempty"`
	PumpProtect    *PumpProtect    `xml:"PUMPPROTECT,omitempty"`

	HeatAreas *[]HeatArea `xml:"HEATAREA,omitempty"`
	HeatCtrls *[]HeatCtrl `xml:"HEATCTRL,omitempty"`
	IODevices *[]IODevice `xml:"IODEVICE,omitempty"`
}
