package main

import (
	"fmt"
	"reflect"

	"github.com/ThomasObenaus/go-base/config"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

func getTargetTypeAndValue(target interface{}) (reflect.Type, reflect.Value, error) {
	targetType := reflect.TypeOf(target)
	targetValue := reflect.ValueOf(target)

	isNilPtr := targetValue.Kind() == reflect.Ptr && targetValue.IsNil()
	isTypeNotSupported := targetValue.Kind() != reflect.Ptr
	if isTypeNotSupported {
		return nil, reflect.Zero(targetType), fmt.Errorf("Can't handle %v (kind=%s,value=%v) (the type has to be a pointer)", targetType, targetType.Kind(), targetValue)
	}
	if isNilPtr {
		return nil, reflect.Zero(targetType), fmt.Errorf("Can't handle %v (kind=%s,value=%v) (probably the type is not supported)", targetType, targetType.Kind(), targetValue)
	}

	// use the element type since we have a pointer
	if targetType.Kind() == reflect.Ptr {
		targetType = targetType.Elem()
		targetValue = targetValue.Elem()
	}

	return targetType, targetValue, nil
}

// applyConfig applies the config that is stored in the given provider. The config will be used to fill the given target type.
func applyConfig(provider config.Provider, target interface{}, nameOfParentType string, parent configTag) error {

	targetType, targetValue, err := getTargetTypeAndValue(target)
	if err != nil {
		return errors.Wrapf(err, "Applying config target=%v,nameOfParentType=%s,parent=%s,", target, nameOfParentType, parent)
	}

	debug("[Apply-(%s)] structure-type=%v state of structure-type=%v\n", nameOfParentType, targetType, targetValue)

	// TODO: move to function factory
	err = processAllConfigTagsOfStruct(target, nameOfParentType, parent, func(fieldName string, isPrimitive bool, fieldType reflect.Type, fieldValue reflect.Value, cfgTag configTag) error {

		logPrefix := fmt.Sprintf("[Apply-(%s)]", fieldName)
		debug("%s field-type=%s field-value=%v\n", logPrefix, fieldType, fieldValue)

		// handling of non primitives (stucts)
		if !isPrimitive {
			fieldValueIf := fieldValue.Addr().Interface()
			if err := applyConfig(provider, fieldValueIf, nameOfParentType, cfgTag); err != nil {
				return errors.Wrap(err, "Applying non primitive")
			}
			debug("%s applied non primitive %v\n", logPrefix, fieldValueIf)
			return nil
		}

		if !provider.IsSet(cfgTag.Name) {
			debug("%s parameter not provided, nothing will be applied\n", logPrefix)
			return nil
		}

		if !fieldValue.CanSet() {
			return fmt.Errorf("Can't set value to field (fieldName=%s,fieldType=%v,fieldValue=%s)", fieldName, fieldType, fieldValue)
		}

		// apply the value
		val := provider.Get(cfgTag.Name)
		typeOfValueFromConfig := reflect.TypeOf(val)

		// Special treatment for slices of structs. This is needed since flag can't handle them instead the value is encodes in a string.
		if typeOfValueFromConfig.Kind() == reflect.String && fieldType.Kind() == reflect.Slice && fieldType.Elem().Kind() == reflect.Struct {

			sliceOfMapsAsString, err := cast.ToStringE(val)
			if err != nil {
				return errors.Wrapf(err, "Casting %v (type=%T) to string", val, val)
			}

			sliceOfMaps, err := parseStringContainingSliceOfMaps(sliceOfMapsAsString)
			if err != nil {
				return errors.Wrapf(err, "Parsing %v (type=%T) to []map[string]interface{}", val, val)
			}
			val = sliceOfMaps
		}

		// cast the parsed default value to the target type
		castedToTargetTypeIf, err := castToTargetType(val, fieldType)
		if err != nil {
			return errors.Wrapf(err, "Casting to target type")
		}
		castedToTargetType := reflect.ValueOf(castedToTargetTypeIf)

		fieldValue.Set(castedToTargetType)
		debug("%s applied value '%v' (type=%v) to '%s' based on config '%s'\n", logPrefix, val, fieldType, fieldName, cfgTag.Name)
		return nil
	})

	if err != nil {
		return errors.Wrapf(err, "Applying config to %v", targetType)
	}
	return nil
}
