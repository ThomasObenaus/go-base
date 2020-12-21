package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// configTag represents the definition for a config read from the type tag.
// A config tag on a type is expected to be defined as:
//
// 	`cfg:"{'name':'<name of the config>','desc':'<description>','default':<default value>}"`
type configTag struct {
	Name        string      `json:"name,omitempty"`
	Description string      `json:"desc,omitempty"`
	Def         interface{} `json:"default,omitempty"`
	desiredType reflect.Type
}

func (e configTag) String() string {
	return fmt.Sprintf(`n:"%s",d:"%s",df:%v`, e.Name, e.Description, e.Def)
}

func (e configTag) IsRequired() bool {
	return e.Def == nil
}

// parseConfigTag parses a definition like
// 	`cfg:"{'name':'<name of the config>','desc':'<description>','default':<default value>}"`
// to a configTag
func parseConfigTag(configTagStr string, typeOfEntry reflect.Type, nameOfParent string) (configTag, error) {
	configTagStr = strings.TrimSpace(configTagStr)
	// replace all single quotes by double quotes to get a valid json
	configTagStr = strings.ReplaceAll(configTagStr, "'", `"`)

	// parse the config tag
	parsedDefinition := configTag{}
	if err := json.Unmarshal([]byte(configTagStr), &parsedDefinition); err != nil {
		return configTag{}, errors.Wrapf(err, "Parsing configTag from '%s'", configTagStr)
	}

	if len(parsedDefinition.Name) == 0 {
		return configTag{}, fmt.Errorf("Missing required config tag field 'name' on '%s'", configTagStr)
	}

	result := configTag{
		// update name to reflect the hierarchy
		Name:        fullFieldName(nameOfParent, parsedDefinition.Name),
		Description: parsedDefinition.Description,
		desiredType: typeOfEntry,
	}

	// only in case a default value is given
	if parsedDefinition.Def != nil {
		castedValue, err := castToTargetType(parsedDefinition.Def, typeOfEntry)
		if err != nil {
			return configTag{}, errors.Wrap(err, "Casting parsed default value to target type")
		}
		result.Def = castedValue
	}
	return result, nil
}

// extractConfigTags extracts recursively all configTags from the given type.
func extractConfigTags(tCfg reflect.Type, nameOfParentType string, parent configTag) ([]configTag, error) {

	entries := make([]configTag, 0)

	// use the element type if we have a pointer
	if tCfg.Kind() == reflect.Ptr {
		tCfg = tCfg.Elem()
	}
	debug("[Extract-(%s)] structure-type=%v definition=%v\n", nameOfParentType, tCfg, parent)

	for i := 0; i < tCfg.NumField(); i++ {
		field := tCfg.Field(i)
		fType := field.Type

		fieldName := fullFieldName(nameOfParentType, field.Name)
		logPrefix := fmt.Sprintf("[Extract-(%s)]", fieldName)
		debug("%s field-type=%s\n", logPrefix, fType)

		// find out if we already have a primitive type
		isPrimitive, err := isOfPrimitiveType(fType)
		if err != nil {
			return nil, errors.Wrapf(err, "Checking for primitive type failed for field '%s'", fieldName)
		}

		cfgSetting, hasCfgTag := getConfigTagDeclaration(field)
		// skip all fields without the cfg tag
		if !hasCfgTag {
			debug("%s no tag found entry will be skipped\n", logPrefix)
			continue
		}
		debug("%s tag found cfgSetting=%v\n", logPrefix, cfgSetting)

		eDef, err := parseConfigTag(cfgSetting, fType, parent.Name)
		if err != nil {
			return nil, errors.Wrapf(err, "Parsing the config definition failed for field '%s'", fieldName)
		}
		debug("%s parsed config entry=%v\n", logPrefix, eDef)

		debug("%s is primitive=%t\n", logPrefix, isPrimitive)
		if !isPrimitive {
			subEntries, err := extractConfigTags(fType, fieldName, eDef)
			if err != nil {
				return nil, errors.Wrap(err, "Extracting subentries")
			}
			entries = append(entries, subEntries...)
			debug("%s added entries %v\n", logPrefix, entries)
			continue
		}

		entries = append(entries, eDef)
		debug("%s created new entry=%v\n", logPrefix, eDef)
	}
	return entries, nil
}

// isOfPrimitiveType returns true if the given type is a primitive one (can be easily casted).
// This is also the case for slices.
func isOfPrimitiveType(fieldType reflect.Type) (bool, error) {
	kind := fieldType.Kind()
	switch kind {
	case reflect.Struct:
		return false, nil
	case reflect.String, reflect.Bool, reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128, reflect.Int, reflect.Int16,
		reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		return true, nil
	case reflect.Ptr:
		elementType := fieldType.Elem()
		return isOfPrimitiveType(elementType)
	case reflect.Slice:
		return true, nil
	default:
		return false, fmt.Errorf("Kind '%s' with type '%s' is not supported", kind, fieldType)
	}
}
