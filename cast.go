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
	castedToTargetType, err := createAndFillStruct(targetType, parsedDefinitionCastedToMap)
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

	typeOfRawValue := reflect.TypeOf(rawValue)
	if typeOfRawValue.Kind() != reflect.Slice {
		return nil, fmt.Errorf("Can't cast to slice since the given raw value is no slice. Instead it is %v", typeOfRawValue)
	}

	sliceValue := reflect.ValueOf(rawValue)

	// obtain the type of the slices elements
	elementType := targetType.Elem()
	sliceInTargetType := reflect.MakeSlice(targetType, 0, sliceValue.Len())

	for i := 0; i < sliceValue.Len(); i++ {
		rawDefaultValueElement := sliceValue.Index(i).Interface()
		switch castedRawElement := rawDefaultValueElement.(type) {
		case map[string]interface{}:
			// handles structs
			castedToTargetType, err := createAndFillStruct(elementType, castedRawElement)
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

// castToTargetType casts the given raw value to the given target type.
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

func getConfigTagDefinition(fieldDeclaration reflect.StructField) (string, bool) {
	return fieldDeclaration.Tag.Lookup("cfg")
}

// createAndFillStruct creates a struct based on the given type and fills its fields based on the given data.
// For being able to fill the struct the given datas keys have to match the config tags that are defined on the target type.
//
// e.g. for type
//
//	type my struct {
//		Field1 string `cfg:"{'name':'field_1'}"`
//	}
//
// the data map should contain an entry with name 'field_1'
// 	data := map[string]interface{}{"field_1":"a value"}
func createAndFillStruct(targetTypeOfStruct reflect.Type, data map[string]interface{}) (reflect.Value, error) {
	if targetTypeOfStruct.Kind() != reflect.Struct {
		return reflect.Zero(targetTypeOfStruct), fmt.Errorf("The target type must be a struct")
	}

	newStruct := reflect.New(targetTypeOfStruct)
	newStructValue := newStruct.Elem()

	for i := 0; i < targetTypeOfStruct.NumField(); i++ {
		fieldDeclaration := targetTypeOfStruct.Field(i)
		fieldValue := newStructValue.FieldByName(fieldDeclaration.Name)
		fieldType := fieldDeclaration.Type
		configTag, hasConfig := getConfigTagDefinition(fieldDeclaration)
		if !hasConfig {
			continue
		}

		entry, err := parseConfigTagDefinition(configTag, fieldType, "")
		if err != nil {
			return reflect.Zero(targetTypeOfStruct), errors.Wrapf(err, "Parsing configTag '%s'", configTag)
		}
		val, ok := data[entry.Name]
		if !ok {
			if entry.IsRequired() {
				return reflect.Zero(targetTypeOfStruct), fmt.Errorf("Missing value for required field (struct-field='%s',expected-key='%s')", fieldDeclaration.Name, entry.Name)
			}

			// take the default value
			val = entry.Def
		}

		// cast the parsed default value to the target type
		castedToTargetTypeIf, err := castToTargetType(val, fieldType)
		if err != nil {
			return reflect.Zero(targetTypeOfStruct), errors.Wrapf(err, "Casting to target type")
		}
		castedToTargetType := reflect.ValueOf(castedToTargetTypeIf)

		// ensure that the casted value can be set
		if !isFieldExported(fieldDeclaration) {
			return reflect.Zero(targetTypeOfStruct), fmt.Errorf("Can't set value for unexported field (struct-field='%s',key='%s').", fieldDeclaration.Name, entry.Name)
		}
		if !fieldValue.CanSet() {
			return reflect.Zero(targetTypeOfStruct), fmt.Errorf("Can't set value for field (struct-field='%s',key='%s').", fieldDeclaration.Name, entry.Name)
		}
		fieldValue.Set(castedToTargetType)
	}

	return newStructValue, nil
}

func fullFieldName(nameOfParent string, fieldName string) string {
	if len(nameOfParent) == 0 {
		return fieldName
	}
	return fmt.Sprintf("%s.%s", nameOfParent, fieldName)
}
