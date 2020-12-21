package main

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

func isFieldExported(typeOfField reflect.StructField) bool {
	return typeOfField.PkgPath == ""
}

// castToPrimitive supports casting of primitive types (such as int, string,...) to the given target type.
func castToPrimitive(rawValue interface{}, targetType reflect.Type) (interface{}, error) {
	typeOfValue := reflect.TypeOf(rawValue)
	if !typeOfValue.ConvertibleTo(targetType) {
		return nil, fmt.Errorf("Can't convert %v to %v", typeOfValue, targetType)
	}
	return reflect.ValueOf(rawValue).Convert(targetType).Interface(), nil
}

// castToStruct supports casting of structs (also nested) to the given target type.
func castToStruct(rawValue interface{}, targetType reflect.Type) (interface{}, error) {
	if targetType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("Can't cast to struct since the target type is not a struct. Instead it is %v", targetType)
	}

	parsedDefinitionCastedToMap, ok := rawValue.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Unable to cast %v (type=%T) to %v. Type must be map[string]interface{}", rawValue, rawValue, targetType)
	}
	castedToTargetType, err := createAndMapStruct(targetType, parsedDefinitionCastedToMap)
	if err != nil {
		return nil, errors.Wrap(err, "Handling default value for element in a slice of structs")
	}
	return castedToTargetType.Interface(), nil
}

// castToSlice supports casting of slices (of primitives and structs) to the given target type.
func castToSlice(rawValue interface{}, targetType reflect.Type) (interface{}, error) {
	if targetType.Kind() != reflect.Slice {
		return nil, fmt.Errorf("Can't cast to slice since the target type is not a slice. Instead it is %v", targetType)
	}

	typedDefaultValue, ok := rawValue.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Types does not match. The target type is a slice (type=%v) but the given default value is no slice (type=%T).", targetType, rawValue)
	}

	// obtain the type of the slices elements
	elementType := targetType.Elem()
	sliceInTargetType := reflect.MakeSlice(targetType, 0, len(typedDefaultValue))

	for _, rawDefaultValueElement := range typedDefaultValue {
		switch castedRawElement := rawDefaultValueElement.(type) {
		case map[string]interface{}:
			// handles structs
			castedToTargetType, err := createAndMapStruct(elementType, castedRawElement)
			if err != nil {
				return nil, errors.Wrap(err, "Handling default value for element in a slice of structs")
			}
			sliceInTargetType = reflect.Append(sliceInTargetType, castedToTargetType)
		default:
			// handles primitive elements (int, string, ...)
			castedToTargetType, err := castToPrimitive(rawDefaultValueElement, elementType)
			if err != nil {
				return nil, errors.Wrap(err, "Casting default value to element type")
			}
			sliceInTargetType = reflect.Append(sliceInTargetType, reflect.ValueOf(castedToTargetType))
		}

	}
	return sliceInTargetType.Interface(), nil
}

func castToTargetType(rawUntypedValue interface{}, targetType reflect.Type) (interface{}, error) {
	switch targetType.Kind() {
	case reflect.Struct:
		return castToStruct(rawUntypedValue, targetType)
	case reflect.Slice:
		return castToSlice(rawUntypedValue, targetType)
	default:
		return castToPrimitive(rawUntypedValue, targetType)
	}
}
