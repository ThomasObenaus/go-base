package main

import (
	"fmt"
	"reflect"

	"github.com/ThomasObenaus/go-base/config"
	"github.com/pkg/errors"
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

		// apply the value
		val := provider.Get(cfgTag.Name)
		newValue := reflect.ValueOf(val)
		typeOfValue := reflect.TypeOf(val)
		debug("%s applied value '%v' (type=%v) to '%s' based on config '%s'\n", logPrefix, newValue, typeOfValue, fieldName, cfgTag.Name)
		fieldValue.Set(newValue)
		debug("%s applied value '%v' (type=%v) to '%s' based on config '%s'\n", logPrefix, newValue, typeOfValue, fieldName, cfgTag.Name)
		return nil
	})

	if err != nil {
		return errors.Wrapf(err, "Applying config to %v", targetType)
	}
	return nil
}
