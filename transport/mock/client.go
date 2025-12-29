package mock

import (
	"github.com/chrishrb/ezr2mqtt/transport"
)

type MockClient struct {
	currentMessage *transport.Message
}

func NewMockClient() *MockClient {
	return &MockClient{
		currentMessage: createMockMessage(),
	}
}

func (c *MockClient) Connect() (*transport.Message, error) {
	return c.currentMessage, nil
}

func (c *MockClient) Send(msg *transport.Message) error {
	if c.currentMessage == nil {
		c.currentMessage = &transport.Message{}
	}

	// Merge the sent message with the current message
	// Only update fields that are present in the sent message
	if msg.Device.ID != nil {
		c.currentMessage.Device.ID = msg.Device.ID
	}
	if msg.Device.Type != nil {
		c.currentMessage.Device.Type = msg.Device.Type
	}
	if msg.Device.Name != nil {
		c.currentMessage.Device.Name = msg.Device.Name
	}
	if msg.Device.Origin != nil {
		c.currentMessage.Device.Origin = msg.Device.Origin
	}
	if msg.Device.ErrorCount != nil {
		c.currentMessage.Device.ErrorCount = msg.Device.ErrorCount
	}
	if msg.Device.DateTime != nil {
		c.currentMessage.Device.DateTime = msg.Device.DateTime
	}
	if msg.Device.DayOfWeek != nil {
		c.currentMessage.Device.DayOfWeek = msg.Device.DayOfWeek
	}
	if msg.Device.TimeZone != nil {
		c.currentMessage.Device.TimeZone = msg.Device.TimeZone
	}
	if msg.Device.NTPSync != nil {
		c.currentMessage.Device.NTPSync = msg.Device.NTPSync
	}
	if msg.Device.VersSWSTM != nil {
		c.currentMessage.Device.VersSWSTM = msg.Device.VersSWSTM
	}
	if msg.Device.VersSWETH != nil {
		c.currentMessage.Device.VersSWETH = msg.Device.VersSWETH
	}
	if msg.Device.VersHW != nil {
		c.currentMessage.Device.VersHW = msg.Device.VersHW
	}
	if msg.Device.TemperatureUnit != nil {
		c.currentMessage.Device.TemperatureUnit = msg.Device.TemperatureUnit
	}
	if msg.Device.SummerWinter != nil {
		c.currentMessage.Device.SummerWinter = msg.Device.SummerWinter
	}
	if msg.Device.TPS != nil {
		c.currentMessage.Device.TPS = msg.Device.TPS
	}
	if msg.Device.Limiter != nil {
		c.currentMessage.Device.Limiter = msg.Device.Limiter
	}
	if msg.Device.MasterID != nil {
		c.currentMessage.Device.MasterID = msg.Device.MasterID
	}
	if msg.Device.Changeover != nil {
		c.currentMessage.Device.Changeover = msg.Device.Changeover
	}
	if msg.Device.Cooling != nil {
		c.currentMessage.Device.Cooling = msg.Device.Cooling
	}
	if msg.Device.Mode != nil {
		c.currentMessage.Device.Mode = msg.Device.Mode
	}
	if msg.Device.OperationModeActor != nil {
		c.currentMessage.Device.OperationModeActor = msg.Device.OperationModeActor
	}
	if msg.Device.Antifreeze != nil {
		c.currentMessage.Device.Antifreeze = msg.Device.Antifreeze
	}
	if msg.Device.AntifreezeTemp != nil {
		c.currentMessage.Device.AntifreezeTemp = msg.Device.AntifreezeTemp
	}
	if msg.Device.FirstOpenTime != nil {
		c.currentMessage.Device.FirstOpenTime = msg.Device.FirstOpenTime
	}
	if msg.Device.SmartStart != nil {
		c.currentMessage.Device.SmartStart = msg.Device.SmartStart
	}
	if msg.Device.EcoDiff != nil {
		c.currentMessage.Device.EcoDiff = msg.Device.EcoDiff
	}
	if msg.Device.EcoInputMode != nil {
		c.currentMessage.Device.EcoInputMode = msg.Device.EcoInputMode
	}
	if msg.Device.EcoInputState != nil {
		c.currentMessage.Device.EcoInputState = msg.Device.EcoInputState
	}
	if msg.Device.THeatVacation != nil {
		c.currentMessage.Device.THeatVacation = msg.Device.THeatVacation
	}

	// Update nested structures - merge fields instead of replacing entire structures
	if msg.Device.Vacation != nil {
		if c.currentMessage.Device.Vacation == nil {
			c.currentMessage.Device.Vacation = &transport.Vacation{}
		}
		if msg.Device.Vacation.State != nil {
			c.currentMessage.Device.Vacation.State = msg.Device.Vacation.State
		}
		if msg.Device.Vacation.StartDate != nil {
			c.currentMessage.Device.Vacation.StartDate = msg.Device.Vacation.StartDate
		}
		if msg.Device.Vacation.StartTime != nil {
			c.currentMessage.Device.Vacation.StartTime = msg.Device.Vacation.StartTime
		}
		if msg.Device.Vacation.EndDate != nil {
			c.currentMessage.Device.Vacation.EndDate = msg.Device.Vacation.EndDate
		}
		if msg.Device.Vacation.EndTime != nil {
			c.currentMessage.Device.Vacation.EndTime = msg.Device.Vacation.EndTime
		}
	}

	if msg.Device.Network != nil {
		if c.currentMessage.Device.Network == nil {
			c.currentMessage.Device.Network = &transport.Network{}
		}
		if msg.Device.Network.MAC != nil {
			c.currentMessage.Device.Network.MAC = msg.Device.Network.MAC
		}
		if msg.Device.Network.DHCP != nil {
			c.currentMessage.Device.Network.DHCP = msg.Device.Network.DHCP
		}
		if msg.Device.Network.IPv6Active != nil {
			c.currentMessage.Device.Network.IPv6Active = msg.Device.Network.IPv6Active
		}
		if msg.Device.Network.IPv4Actual != nil {
			c.currentMessage.Device.Network.IPv4Actual = msg.Device.Network.IPv4Actual
		}
		if msg.Device.Network.IPv4Set != nil {
			c.currentMessage.Device.Network.IPv4Set = msg.Device.Network.IPv4Set
		}
		if msg.Device.Network.IPv6Actual != nil {
			c.currentMessage.Device.Network.IPv6Actual = msg.Device.Network.IPv6Actual
		}
		if msg.Device.Network.IPv6Set != nil {
			c.currentMessage.Device.Network.IPv6Set = msg.Device.Network.IPv6Set
		}
		if msg.Device.Network.NetmaskActual != nil {
			c.currentMessage.Device.Network.NetmaskActual = msg.Device.Network.NetmaskActual
		}
		if msg.Device.Network.NetmaskSet != nil {
			c.currentMessage.Device.Network.NetmaskSet = msg.Device.Network.NetmaskSet
		}
		if msg.Device.Network.DNS != nil {
			c.currentMessage.Device.Network.DNS = msg.Device.Network.DNS
		}
		if msg.Device.Network.Gateway != nil {
			c.currentMessage.Device.Network.Gateway = msg.Device.Network.Gateway
		}
	}

	if msg.Device.Cloud != nil {
		if c.currentMessage.Device.Cloud == nil {
			c.currentMessage.Device.Cloud = &transport.Cloud{}
		}
		if msg.Device.Cloud.UserID != nil {
			c.currentMessage.Device.Cloud.UserID = msg.Device.Cloud.UserID
		}
		if msg.Device.Cloud.Password != nil {
			c.currentMessage.Device.Cloud.Password = msg.Device.Cloud.Password
		}
		if msg.Device.Cloud.M2MServerPort != nil {
			c.currentMessage.Device.Cloud.M2MServerPort = msg.Device.Cloud.M2MServerPort
		}
		if msg.Device.Cloud.M2MLocalPort != nil {
			c.currentMessage.Device.Cloud.M2MLocalPort = msg.Device.Cloud.M2MLocalPort
		}
		if msg.Device.Cloud.M2MHTTPPort != nil {
			c.currentMessage.Device.Cloud.M2MHTTPPort = msg.Device.Cloud.M2MHTTPPort
		}
		if msg.Device.Cloud.M2MHTTPSPort != nil {
			c.currentMessage.Device.Cloud.M2MHTTPSPort = msg.Device.Cloud.M2MHTTPSPort
		}
		if msg.Device.Cloud.M2MServerAddress != nil {
			c.currentMessage.Device.Cloud.M2MServerAddress = msg.Device.Cloud.M2MServerAddress
		}
		if msg.Device.Cloud.M2MActive != nil {
			c.currentMessage.Device.Cloud.M2MActive = msg.Device.Cloud.M2MActive
		}
		if msg.Device.Cloud.M2MState != nil {
			c.currentMessage.Device.Cloud.M2MState = msg.Device.Cloud.M2MState
		}
	}

	// For complex nested structures, merge them field-by-field
	if msg.Device.KWLCtrl != nil {
		if c.currentMessage.Device.KWLCtrl == nil {
			c.currentMessage.Device.KWLCtrl = &transport.KWLCtrl{}
		}
		if msg.Device.KWLCtrl.Visible != nil {
			c.currentMessage.Device.KWLCtrl.Visible = msg.Device.KWLCtrl.Visible
		}
		if msg.Device.KWLCtrl.Present != nil {
			c.currentMessage.Device.KWLCtrl.Present = msg.Device.KWLCtrl.Present
		}
		if msg.Device.KWLCtrl.Connection != nil {
			c.currentMessage.Device.KWLCtrl.Connection = msg.Device.KWLCtrl.Connection
		}
		if msg.Device.KWLCtrl.URL != nil {
			c.currentMessage.Device.KWLCtrl.URL = msg.Device.KWLCtrl.URL
		}
		if msg.Device.KWLCtrl.Port != nil {
			c.currentMessage.Device.KWLCtrl.Port = msg.Device.KWLCtrl.Port
		}
		if msg.Device.KWLCtrl.Status != nil {
			c.currentMessage.Device.KWLCtrl.Status = msg.Device.KWLCtrl.Status
		}
		if msg.Device.KWLCtrl.FlowCtrl != nil {
			c.currentMessage.Device.KWLCtrl.FlowCtrl = msg.Device.KWLCtrl.FlowCtrl
		}
	}

	if msg.Device.Code != nil {
		if c.currentMessage.Device.Code == nil {
			c.currentMessage.Device.Code = &transport.Code{}
		}
		if msg.Device.Code.Expert != nil {
			c.currentMessage.Device.Code.Expert = msg.Device.Code.Expert
		}
	}

	if msg.Device.Program != nil {
		if c.currentMessage.Device.Program == nil {
			c.currentMessage.Device.Program = &transport.Program{}
		}
		if msg.Device.Program.ShiftPrograms != nil {
			c.currentMessage.Device.Program.ShiftPrograms = msg.Device.Program.ShiftPrograms
		}
	}

	if msg.Device.PumpOutput != nil {
		if c.currentMessage.Device.PumpOutput == nil {
			c.currentMessage.Device.PumpOutput = &transport.PumpOutput{}
		}
		if msg.Device.PumpOutput.LocalGlobal != nil {
			c.currentMessage.Device.PumpOutput.LocalGlobal = msg.Device.PumpOutput.LocalGlobal
		}
		if msg.Device.PumpOutput.Type != nil {
			c.currentMessage.Device.PumpOutput.Type = msg.Device.PumpOutput.Type
		}
		if msg.Device.PumpOutput.LeadTime != nil {
			c.currentMessage.Device.PumpOutput.LeadTime = msg.Device.PumpOutput.LeadTime
		}
		if msg.Device.PumpOutput.StoppingTime != nil {
			c.currentMessage.Device.PumpOutput.StoppingTime = msg.Device.PumpOutput.StoppingTime
		}
		if msg.Device.PumpOutput.OperationMode != nil {
			c.currentMessage.Device.PumpOutput.OperationMode = msg.Device.PumpOutput.OperationMode
		}
		if msg.Device.PumpOutput.MinRuntime != nil {
			c.currentMessage.Device.PumpOutput.MinRuntime = msg.Device.PumpOutput.MinRuntime
		}
		if msg.Device.PumpOutput.MinStandstill != nil {
			c.currentMessage.Device.PumpOutput.MinStandstill = msg.Device.PumpOutput.MinStandstill
		}
	}

	if msg.Device.Relais != nil {
		if c.currentMessage.Device.Relais == nil {
			c.currentMessage.Device.Relais = &transport.Relais{}
		}
		if msg.Device.Relais.Function != nil {
			c.currentMessage.Device.Relais.Function = msg.Device.Relais.Function
		}
		if msg.Device.Relais.LeadTime != nil {
			c.currentMessage.Device.Relais.LeadTime = msg.Device.Relais.LeadTime
		}
		if msg.Device.Relais.StoppingTime != nil {
			c.currentMessage.Device.Relais.StoppingTime = msg.Device.Relais.StoppingTime
		}
		if msg.Device.Relais.OperationMode != nil {
			c.currentMessage.Device.Relais.OperationMode = msg.Device.Relais.OperationMode
		}
	}

	if msg.Device.ChangeoverFunc != nil {
		if c.currentMessage.Device.ChangeoverFunc == nil {
			c.currentMessage.Device.ChangeoverFunc = &transport.ChangeoverFunc{}
		}
		if msg.Device.ChangeoverFunc.Mode != nil {
			c.currentMessage.Device.ChangeoverFunc.Mode = msg.Device.ChangeoverFunc.Mode
		}
	}

	if msg.Device.EmergencyMode != nil {
		if c.currentMessage.Device.EmergencyMode == nil {
			c.currentMessage.Device.EmergencyMode = &transport.EmergencyMode{}
		}
		if msg.Device.EmergencyMode.Time != nil {
			c.currentMessage.Device.EmergencyMode.Time = msg.Device.EmergencyMode.Time
		}
		if msg.Device.EmergencyMode.PWMCycle != nil {
			c.currentMessage.Device.EmergencyMode.PWMCycle = msg.Device.EmergencyMode.PWMCycle
		}
		if msg.Device.EmergencyMode.PWMHeat != nil {
			c.currentMessage.Device.EmergencyMode.PWMHeat = msg.Device.EmergencyMode.PWMHeat
		}
		if msg.Device.EmergencyMode.PWMCool != nil {
			c.currentMessage.Device.EmergencyMode.PWMCool = msg.Device.EmergencyMode.PWMCool
		}
	}

	if msg.Device.ValveProtect != nil {
		if c.currentMessage.Device.ValveProtect == nil {
			c.currentMessage.Device.ValveProtect = &transport.ValveProtect{}
		}
		if msg.Device.ValveProtect.Time != nil {
			c.currentMessage.Device.ValveProtect.Time = msg.Device.ValveProtect.Time
		}
		if msg.Device.ValveProtect.Duration != nil {
			c.currentMessage.Device.ValveProtect.Duration = msg.Device.ValveProtect.Duration
		}
	}

	if msg.Device.PumpProtect != nil {
		if c.currentMessage.Device.PumpProtect == nil {
			c.currentMessage.Device.PumpProtect = &transport.PumpProtect{}
		}
		if msg.Device.PumpProtect.Time != nil {
			c.currentMessage.Device.PumpProtect.Time = msg.Device.PumpProtect.Time
		}
		if msg.Device.PumpProtect.Duration != nil {
			c.currentMessage.Device.PumpProtect.Duration = msg.Device.PumpProtect.Duration
		}
	}

	// For arrays, merge by matching the 'Nr' attribute
	if msg.Device.HeatAreas != nil {
		if c.currentMessage.Device.HeatAreas == nil {
			c.currentMessage.Device.HeatAreas = &[]transport.HeatArea{}
		}
		for _, msgArea := range *msg.Device.HeatAreas {
			mergeHeatArea(c.currentMessage.Device.HeatAreas, &msgArea)
		}
	}

	if msg.Device.HeatCtrls != nil {
		if c.currentMessage.Device.HeatCtrls == nil {
			c.currentMessage.Device.HeatCtrls = &[]transport.HeatCtrl{}
		}
		for _, msgCtrl := range *msg.Device.HeatCtrls {
			mergeHeatCtrl(c.currentMessage.Device.HeatCtrls, &msgCtrl)
		}
	}

	if msg.Device.IODevices != nil {
		if c.currentMessage.Device.IODevices == nil {
			c.currentMessage.Device.IODevices = &[]transport.IODevice{}
		}
		for _, msgDevice := range *msg.Device.IODevices {
			mergeIODevice(c.currentMessage.Device.IODevices, &msgDevice)
		}
	}

	return nil
}

func createMockMessage() *transport.Message {
	// Helper function to create pointers
	strPtr := func(s string) *string { return &s }
	intPtr := func(i int) *int { return &i }
	floatPtr := func(f float64) *float64 { return &f }

	return &transport.Message{
		Device: transport.Device{
			ID:     strPtr("MOCK-12345"),
			Type:   strPtr("EZR"),
			Name:   strPtr("Mock Device"),
			Origin: strPtr("mock"),

			ErrorCount: intPtr(0),
			DateTime:   strPtr("2025-12-29T12:00:00"),
			DayOfWeek:  intPtr(7),
			TimeZone:   intPtr(1),
			NTPSync:    intPtr(1),

			VersSWSTM: strPtr("1.0.0"),
			VersSWETH: strPtr("1.0.0"),
			VersHW:    strPtr("1.0"),

			TemperatureUnit: intPtr(0),
			SummerWinter:    intPtr(0),
			TPS:             intPtr(0),
			Limiter:         intPtr(0),

			MasterID:   strPtr("MASTER-001"),
			Changeover: intPtr(0),
			Cooling:    intPtr(0),
			Mode:       intPtr(1),

			OperationModeActor: intPtr(0),
			Antifreeze:         intPtr(1),
			AntifreezeTemp:     floatPtr(5.0),

			FirstOpenTime: intPtr(60),
			SmartStart:    intPtr(1),

			EcoDiff:       floatPtr(2.0),
			EcoInputMode:  intPtr(0),
			EcoInputState: intPtr(0),

			THeatVacation: floatPtr(15.0),

			Vacation: &transport.Vacation{
				State:     intPtr(0),
				StartDate: strPtr("2025-01-01"),
				StartTime: strPtr("00:00"),
				EndDate:   strPtr("2025-01-07"),
				EndTime:   strPtr("23:59"),
			},

			Network: &transport.Network{
				MAC:           strPtr("00:11:22:33:44:55"),
				DHCP:          intPtr(1),
				IPv6Active:    intPtr(0),
				IPv4Actual:    strPtr("192.168.1.100"),
				IPv4Set:       strPtr("192.168.1.100"),
				NetmaskActual: strPtr("255.255.255.0"),
				NetmaskSet:    strPtr("255.255.255.0"),
				DNS:           strPtr("192.168.1.1"),
				Gateway:       strPtr("192.168.1.1"),
			},

			HeatAreas: &[]transport.HeatArea{
				{
					Nr:          intPtr(1),
					Name:        strPtr("Living Room"),
					Mode:        intPtr(0),
					TActual:     floatPtr(22.5),
					TTarget:     floatPtr(22.0),
					TTargetBase: floatPtr(20.0),
					TTargetMin:  floatPtr(5.0),
					TTargetMax:  floatPtr(30.0),
					State:       intPtr(1),
					THeatDay:    floatPtr(22.0),
					THeatNight:  floatPtr(18.0),
				},
				{
					Nr:          intPtr(2),
					Name:        strPtr("Bedroom"),
					Mode:        intPtr(0),
					TActual:     floatPtr(19.5),
					TTarget:     floatPtr(20.0),
					TTargetBase: floatPtr(18.0),
					TTargetMin:  floatPtr(5.0),
					TTargetMax:  floatPtr(30.0),
					State:       intPtr(1),
					THeatDay:    floatPtr(20.0),
					THeatNight:  floatPtr(16.0),
				},
			},
		},
	}
}

// mergeHeatArea merges a HeatArea into the existing array by Nr
func mergeHeatArea(current *[]transport.HeatArea, incoming *transport.HeatArea) {
	if incoming.Nr == nil {
		// If no Nr is specified, append as new
		*current = append(*current, *incoming)
		return
	}

	// Find existing area with same Nr
	for i := range *current {
		if (*current)[i].Nr != nil && *(*current)[i].Nr == *incoming.Nr {
			// Merge fields
			if incoming.Name != nil {
				(*current)[i].Name = incoming.Name
			}
			if incoming.Mode != nil {
				(*current)[i].Mode = incoming.Mode
			}
			if incoming.TActual != nil {
				(*current)[i].TActual = incoming.TActual
			}
			if incoming.TActualExt != nil {
				(*current)[i].TActualExt = incoming.TActualExt
			}
			if incoming.TTarget != nil {
				(*current)[i].TTarget = incoming.TTarget
			}
			if incoming.TTargetBase != nil {
				(*current)[i].TTargetBase = incoming.TTargetBase
			}
			if incoming.State != nil {
				(*current)[i].State = incoming.State
			}
			if incoming.ProgramSource != nil {
				(*current)[i].ProgramSource = incoming.ProgramSource
			}
			if incoming.ProgramWeek != nil {
				(*current)[i].ProgramWeek = incoming.ProgramWeek
			}
			if incoming.ProgramWeekend != nil {
				(*current)[i].ProgramWeekend = incoming.ProgramWeekend
			}
			if incoming.Party != nil {
				(*current)[i].Party = incoming.Party
			}
			if incoming.PartyRemainingTime != nil {
				(*current)[i].PartyRemainingTime = incoming.PartyRemainingTime
			}
			if incoming.Presence != nil {
				(*current)[i].Presence = incoming.Presence
			}
			if incoming.TTargetMin != nil {
				(*current)[i].TTargetMin = incoming.TTargetMin
			}
			if incoming.TTargetMax != nil {
				(*current)[i].TTargetMax = incoming.TTargetMax
			}
			if incoming.RPMMotor != nil {
				(*current)[i].RPMMotor = incoming.RPMMotor
			}
			if incoming.Offset != nil {
				(*current)[i].Offset = incoming.Offset
			}
			if incoming.THeatDay != nil {
				(*current)[i].THeatDay = incoming.THeatDay
			}
			if incoming.THeatNight != nil {
				(*current)[i].THeatNight = incoming.THeatNight
			}
			if incoming.TCoolDay != nil {
				(*current)[i].TCoolDay = incoming.TCoolDay
			}
			if incoming.TCoolNight != nil {
				(*current)[i].TCoolNight = incoming.TCoolNight
			}
			if incoming.TFloorDay != nil {
				(*current)[i].TFloorDay = incoming.TFloorDay
			}
			if incoming.HeatingSystem != nil {
				(*current)[i].HeatingSystem = incoming.HeatingSystem
			}
			if incoming.BlockHC != nil {
				(*current)[i].BlockHC = incoming.BlockHC
			}
			if incoming.IsLocked != nil {
				(*current)[i].IsLocked = incoming.IsLocked
			}
			if incoming.LockCode != nil {
				(*current)[i].LockCode = incoming.LockCode
			}
			if incoming.LockAvailable != nil {
				(*current)[i].LockAvailable = incoming.LockAvailable
			}
			if incoming.Light != nil {
				(*current)[i].Light = incoming.Light
			}
			if incoming.SensorExt != nil {
				(*current)[i].SensorExt = incoming.SensorExt
			}
			if incoming.Adjustable != nil {
				(*current)[i].Adjustable = incoming.Adjustable
			}
			return
		}
	}

	// Not found, append as new
	*current = append(*current, *incoming)
}

// mergeHeatCtrl merges a HeatCtrl into the existing array by Nr
func mergeHeatCtrl(current *[]transport.HeatCtrl, incoming *transport.HeatCtrl) {
	if incoming.Nr == nil {
		*current = append(*current, *incoming)
		return
	}

	for i := range *current {
		if (*current)[i].Nr != nil && *(*current)[i].Nr == *incoming.Nr {
			if incoming.InUse != nil {
				(*current)[i].InUse = incoming.InUse
			}
			if incoming.HeatAreaNr != nil {
				(*current)[i].HeatAreaNr = incoming.HeatAreaNr
			}
			if incoming.Actor != nil {
				(*current)[i].Actor = incoming.Actor
			}
			if incoming.ActorPercent != nil {
				(*current)[i].ActorPercent = incoming.ActorPercent
			}
			if incoming.State != nil {
				(*current)[i].State = incoming.State
			}
			return
		}
	}

	*current = append(*current, *incoming)
}

// mergeIODevice merges an IODevice into the existing array by Nr
func mergeIODevice(current *[]transport.IODevice, incoming *transport.IODevice) {
	if incoming.Nr == nil {
		*current = append(*current, *incoming)
		return
	}

	for i := range *current {
		if (*current)[i].Nr != nil && *(*current)[i].Nr == *incoming.Nr {
			if incoming.Type != nil {
				(*current)[i].Type = incoming.Type
			}
			if incoming.ID != nil {
				(*current)[i].ID = incoming.ID
			}
			if incoming.VersHW != nil {
				(*current)[i].VersHW = incoming.VersHW
			}
			if incoming.VersSW != nil {
				(*current)[i].VersSW = incoming.VersSW
			}
			if incoming.HeatAreaNr != nil {
				(*current)[i].HeatAreaNr = incoming.HeatAreaNr
			}
			if incoming.SignalStrength != nil {
				(*current)[i].SignalStrength = incoming.SignalStrength
			}
			if incoming.Battery != nil {
				(*current)[i].Battery = incoming.Battery
			}
			if incoming.State != nil {
				(*current)[i].State = incoming.State
			}
			if incoming.ComError != nil {
				(*current)[i].ComError = incoming.ComError
			}
			if incoming.IsOn != nil {
				(*current)[i].IsOn = incoming.IsOn
			}
			return
		}
	}

	*current = append(*current, *incoming)
}
