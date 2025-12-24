package mock

import (
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
	return &transport.Message{}
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
