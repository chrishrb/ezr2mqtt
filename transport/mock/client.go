package mock

import (
	"encoding/xml"
	"reflect"

	"github.com/chrishrb/ezr2mqtt/transport"
)

type MockClient struct {
	// currentMessage stores the current state of the device
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
	// Mutate the current message by updating only the fields that are present in msg
	mergeMessages(c.currentMessage, msg)
	return nil
}

// createMockMessage creates a transport.Message with mock data
func createMockMessage() *transport.Message {
	return &transport.Message{
		XMLName: xml.Name{Local: "Devices"},
		Device: transport.Device{
			// Identification
			ID:     "MOCK-12345",
			Type:   "EZR",
			Name:   "Mock Device",
			Origin: "Mock",

			// System
			ErrorCount: 0,
			DateTime:   "2025-12-23 10:00:00",
			DayOfWeek:  1,
			TimeZone:   1,
			NTPSync:    1,

			VersSWSTM: "1.0.0",
			VersSWETH: "1.0.0",
			VersHW:    "1.0",

			TemperatureUnit: 0,
			SummerWinter:    0,
			TPS:             0,
			Limiter:         0,

			MasterID:   "MOCK-MASTER",
			Changeover: 0,
			Cooling:    0,
			Mode:       1,

			OperationModeActor: 0,

			Antifreeze:     0,
			AntifreezeTemp: 5.0,

			FirstOpenTime: 30,
			SmartStart:    0,

			EcoDiff:       2.0,
			EcoInputMode:  0,
			EcoInputState: 0,

			THeatVacation: 15.0,

			Vacation: transport.Vacation{
				State:     0,
				StartDate: "01.01.2025",
				StartTime: "00:00",
				EndDate:   "01.01.2025",
				EndTime:   "00:00",
			},

			Network: transport.Network{
				MAC:           "00:11:22:33:44:55",
				DHCP:          1,
				IPv6Active:    0,
				IPv4Actual:    "192.168.1.100",
				IPv4Set:       "192.168.1.100",
				IPv6Actual:    "",
				IPv6Set:       "",
				NetmaskActual: "255.255.255.0",
				NetmaskSet:    "255.255.255.0",
				DNS:           "192.168.1.1",
				Gateway:       "192.168.1.1",
			},

			Cloud: transport.Cloud{
				UserID:           "mock-user",
				Password:         "mock-password",
				M2MServerPort:    8080,
				M2MLocalPort:     8081,
				M2MHTTPPort:      80,
				M2MHTTPSPort:     443,
				M2MServerAddress: "mock.server.com",
				M2MActive:        0,
				M2MState:         "disconnected",
			},

			KWLCtrl: transport.KWLCtrl{
				Visible:    0,
				Present:    0,
				Connection: 0,
				URL:        "",
				Port:       0,
				Status:     0,
				FlowCtrl:   0,
			},

			Code: transport.Code{
				Expert: "0000",
			},

			Program: transport.Program{
				ShiftPrograms: []transport.ShiftProgram{
					{Nr: 1, ShiftingTime: 0, Start: "06:00", End: "22:00"},
					{Nr: 2, ShiftingTime: 0, Start: "08:00", End: "18:00"},
				},
			},

			PumpOutput: transport.PumpOutput{
				LocalGlobal:   0,
				Type:          0,
				LeadTime:      0,
				StoppingTime:  0,
				OperationMode: 0,
				MinRuntime:    0,
				MinStandstill: 0,
			},

			Relais: transport.Relais{
				Function:      0,
				LeadTime:      0,
				StoppingTime:  0,
				OperationMode: 0,
			},

			ChangeoverFunc: transport.ChangeoverFunc{
				Mode: 0,
			},

			EmergencyMode: transport.EmergencyMode{
				Time:     0,
				PWMCycle: 0,
				PWMHeat:  0,
				PWMCool:  0,
			},

			ValveProtect: transport.ValveProtect{
				Time:     0,
				Duration: 0,
			},

			PumpProtect: transport.PumpProtect{
				Time:     0,
				Duration: 0,
			},

			HeatAreas: []transport.HeatArea{
				{
					Nr:   1,
					Name: "Living Room",
					Mode: 1,

					TActual:     22.5,
					TActualExt:  20.0,
					TTarget:     22.0,
					TTargetBase: 20.0,

					State: 1,

					ProgramSource:  1,
					ProgramWeek:    1,
					ProgramWeekend: 2,

					Party:              0,
					PartyRemainingTime: 0,
					Presence:           1,

					TTargetMin: 15.0,
					TTargetMax: 30.0,

					RPMMotor: 0,
					Offset:   0.0,

					THeatDay:   22.0,
					THeatNight: 18.0,
					TCoolDay:   24.0,
					TCoolNight: 26.0,
					TFloorDay:  25.0,

					HeatingSystem: 0,

					BlockHC: 0,

					IsLocked:      0,
					LockCode:      "",
					LockAvailable: 0,

					Light:      0,
					SensorExt:  0,
					Adjustable: 1,
				},
				{
					Nr:   2,
					Name: "Bedroom",
					Mode: 1,

					TActual:     19.5,
					TActualExt:  18.0,
					TTarget:     20.0,
					TTargetBase: 18.0,

					State: 1,

					ProgramSource:  1,
					ProgramWeek:    1,
					ProgramWeekend: 2,

					Party:              0,
					PartyRemainingTime: 0,
					Presence:           1,

					TTargetMin: 15.0,
					TTargetMax: 25.0,

					RPMMotor: 0,
					Offset:   0.0,

					THeatDay:   20.0,
					THeatNight: 17.0,
					TCoolDay:   23.0,
					TCoolNight: 25.0,
					TFloorDay:  24.0,

					HeatingSystem: 0,

					BlockHC: 0,

					IsLocked:      0,
					LockCode:      "",
					LockAvailable: 0,

					Light:      0,
					SensorExt:  0,
					Adjustable: 1,
				},
			},

			HeatCtrls: []transport.HeatCtrl{
				{Nr: 1, InUse: 1, HeatAreaNr: 1, Actor: 50, ActorPercent: 50, State: 1},
				{Nr: 2, InUse: 1, HeatAreaNr: 2, Actor: 30, ActorPercent: 30, State: 1},
			},

			IODevices: []transport.IODevice{
				{
					Nr:             1,
					Type:           1,
					ID:             101,
					VersHW:         "1.0",
					VersSW:         "1.0.0",
					HeatAreaNr:     1,
					SignalStrength: 85,
					Battery:        100,
					State:          1,
					ComError:       0,
					IsOn:           1,
				},
			},
		},
	}
}

// mergeMessages updates target with values from source where source has non-zero values
func mergeMessages(target, source *transport.Message) {
	if source == nil || target == nil {
		return
	}

	// Merge the Device struct
	mergeStructs(reflect.ValueOf(&target.Device).Elem(), reflect.ValueOf(&source.Device).Elem())
}

// mergeStructs recursively merges source struct into target struct
// Only non-zero values from source are copied to target
func mergeStructs(target, source reflect.Value) {
	if !target.IsValid() || !source.IsValid() {
		return
	}

	switch source.Kind() {
	case reflect.Struct:
		for i := 0; i < source.NumField(); i++ {
			sourceField := source.Field(i)
			targetField := target.Field(i)

			if !targetField.CanSet() {
				continue
			}

			mergeStructs(targetField, sourceField)
		}

	case reflect.Slice:
		// For slices, replace the entire slice if source is non-empty
		if source.Len() > 0 {
			// Match slice elements by their 'Nr' field if it exists
			if source.Len() > 0 && source.Index(0).Kind() == reflect.Struct {
				mergeSlicesByNr(target, source)
			} else {
				target.Set(source)
			}
		}

	case reflect.String:
		if source.String() != "" {
			target.SetString(source.String())
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// Only update if source has a non-zero value
		if source.Int() != 0 {
			target.SetInt(source.Int())
		}

	case reflect.Float32, reflect.Float64:
		// Only update if source has a non-zero value
		if source.Float() != 0.0 {
			target.SetFloat(source.Float())
		}
	}
}

// mergeSlicesByNr merges slices by matching elements with the same 'Nr' field
func mergeSlicesByNr(target, source reflect.Value) {
	if source.Len() == 0 {
		return
	}

	// Check if elements have Nr field
	elemType := source.Type().Elem()
	if elemType.Kind() != reflect.Struct {
		return
	}

	nrField, hasNr := elemType.FieldByName("Nr")
	if !hasNr || nrField.Type.Kind() != reflect.Int {
		// If no Nr field, just replace the entire slice
		target.Set(source)
		return
	}

	// Create a map of target elements by Nr
	targetMap := make(map[int]reflect.Value)
	for i := 0; i < target.Len(); i++ {
		elem := target.Index(i)
		nr := int(elem.FieldByName("Nr").Int())
		targetMap[nr] = elem
	}

	// Merge source elements into target
	for i := 0; i < source.Len(); i++ {
		sourceElem := source.Index(i)
		nr := int(sourceElem.FieldByName("Nr").Int())

		if targetElem, exists := targetMap[nr]; exists {
			// Update existing element - use special merge for elements matched by Nr
			mergeStructsByNr(targetElem, sourceElem)
		} else {
			// Append new element
			target.Set(reflect.Append(target, sourceElem))
		}
	}
}

// mergeStructsByNr merges structs that were matched by Nr field
// For these structs, we merge non-zero values, plus specific fields that should allow zero
func mergeStructsByNr(target, source reflect.Value) {
	if !target.IsValid() || !source.IsValid() {
		return
	}

	if source.Kind() != reflect.Struct {
		return
	}

	// Fields that should be merged even when zero
	// These are fields where 0 is a valid meaningful value
	zeroAllowedFields := map[string]bool{
		"Mode":  true, // Heat area mode: 0=auto, 1=day, 2=night
		"State": true, // Various state fields where 0 is valid
		"Actor": true, // Actor percentage can be 0
	}

	for i := 0; i < source.NumField(); i++ {
		sourceField := source.Field(i)
		targetField := target.Field(i)
		fieldName := source.Type().Field(i).Name

		if !targetField.CanSet() {
			continue
		}

		// Skip the Nr field itself to avoid changing the identifier
		if fieldName == "Nr" {
			continue
		}

		// For fields in Nr-matched structs, apply different rules
		switch sourceField.Kind() {
		case reflect.Struct:
			mergeStructs(targetField, sourceField)
		case reflect.Slice:
			if sourceField.Len() > 0 {
				targetField.Set(sourceField)
			}
		case reflect.String:
			if sourceField.String() != "" {
				targetField.SetString(sourceField.String())
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			// Always set if non-zero, or if it's a field that allows zero values
			if sourceField.Int() != 0 || zeroAllowedFields[fieldName] {
				targetField.SetInt(sourceField.Int())
			}
		case reflect.Float32, reflect.Float64:
			// Always set if non-zero, or if it's a field that allows zero values
			if sourceField.Float() != 0.0 || zeroAllowedFields[fieldName] {
				targetField.SetFloat(sourceField.Float())
			}
		}
	}
}
