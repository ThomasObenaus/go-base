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
	isNotSupportedField := targetValue.Kind() != reflect.Ptr
	if isNotSupportedField || isNilPtr {
		return nil, reflect.Zero(targetType), fmt.Errorf("Can't handle %v (kind=%s,value=%v) (probably the type is not supported)", targetType, targetType.Kind(), targetValue)
	}

	// use the element type since we have a pointer
	targetType = targetType.Elem()
	targetValue = targetValue.Elem()

	return targetType, targetValue, nil
}

// applyConfig applies the config that is stored in the given provider. The config will be used to fill the given target type.
func applyConfig(provider config.Provider, target interface{}, nameOfParentType string, parent configTag) error {

	targetType, targetValue, err := getTargetTypeAndValue(target)
	if err != nil {
		return errors.Wrapf(err, "Applying config target=%v,nameOfParentType=%s,parent=%s,", target, nameOfParentType, parent)
	}

	debug("[Apply-(%s)] structure-type=%v state of structure-type=%v\n", nameOfParentType, targetType, targetValue)

	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		fType := field.Type

		fieldName := fullFieldName(nameOfParentType, field.Name)
		logPrefix := fmt.Sprintf("[Apply-(%s)]", fieldName)
		debug("%s field-type=%s\n", logPrefix, fType)

		// find out if we already have a primitive type
		isPrimitive, err := isOfPrimitiveType(fType)
		if err != nil {
			return errors.Wrapf(err, "Checking for primitive type failed for field '%s'", fieldName)
		}
		debug("%s is primitive=%t\n", logPrefix, isPrimitive)

		cfgSetting, hasCfgTag := getConfigTagDeclaration(field)
		// ignore fields without a config tag
		if !hasCfgTag {
			debug("%s no tag found entry will be skipped\n", logPrefix)
			continue
		}
		debug("%s tag found cfgSetting=%v\n", logPrefix, cfgSetting)

		eDef, err := parseConfigTag(cfgSetting, fType, parent.Name)
		if err != nil {
			return errors.Wrapf(err, "Parsing the config definition failed for field '%s'", fieldName)
		}
		debug("%s parsed config entry=%v\n", logPrefix, eDef)

		v := targetValue.Field(i)
		fieldValue := v.Addr().Interface()
		debug("%s field-type=%s field-value=%v\n", logPrefix, fType, v)

		// handling of non primitives (stucts)
		if !isPrimitive {
			if err := applyConfig(provider, fieldValue, nameOfParentType, eDef); err != nil {
				return errors.Wrap(err, "Applying non primitive")
			}
			debug("%s applied non primitive %v\n", logPrefix, fieldValue)
			continue
		}

		if !provider.IsSet(eDef.Name) {
			debug("%s parameter not provided, nothing will be applied\n", logPrefix)
			continue
		}

		// apply the value
		val := provider.Get(eDef.Name)
		newValue := reflect.ValueOf(val)
		typeOfValue := reflect.TypeOf(val)
		debug("%s applied value '%v' (type=%v) to '%s' based on config '%s'\n", logPrefix, newValue, typeOfValue, fieldName, eDef.Name)
		v.Set(newValue)
		debug("%s applied value '%v' (type=%v) to '%s' based on config '%s'\n", logPrefix, newValue, typeOfValue, fieldName, eDef.Name)

	}
	return nil
}
