package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/ThomasObenaus/go-base/config"
	"github.com/pkg/errors"
)

// configTag represents the definition for a config read from the type tag.
// A config tag on a type is expected to be defined as:
//
// `cfg:"{'name':'<name of the config>','desc':'<description>','default':<default value>}"`
//
type configTag struct {
	Name        string      `json:"name,omitempty"`
	Description string      `json:"desc,omitempty"`
	Def         interface{} `json:"default,omitempty"`
}

func (e configTag) String() string {
	return fmt.Sprintf(`n:"%s",d:"%s",df:%v`, e.Name, e.Description, e.Def)
}

func (e configTag) IsRequired() bool {
	return e.Def == nil
}

// parseConfigTag parses a definition like
// `cfg:"{'name':'<name of the config>','desc':'<description>','default':<default value>}"`
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

	result := configTag{
		// update name to reflect the hierarchy
		Name:        fullFieldName(nameOfParent, parsedDefinition.Name),
		Description: parsedDefinition.Description,
	}

	// only in case a default value is given
	if parsedDefinition.Def != nil {

		// TODO: Enable defaults on struct level, see how it is done for slices of stucts
		if typeOfEntry.Kind() == reflect.Struct {
			return configTag{}, fmt.Errorf("Default values on struct level are not allowed")
		}

		switch typedDefaultValue := parsedDefinition.Def.(type) {
		case []interface{}:
			// obtain the element type
			elementType := typeOfEntry.Elem()
			sliceInTargetType := reflect.MakeSlice(typeOfEntry, 0, len(typedDefaultValue))

			for _, rawDefaultValueElement := range typedDefaultValue {

				switch castedRawElement := rawDefaultValueElement.(type) {
				case map[string]interface{}:
					// handles structs
					castedToTargetType, err := createAndMapStruct(elementType, castedRawElement)
					if err != nil {
						return configTag{}, errors.Wrap(err, "Handling default value for element in a slice of structs")
					}
					sliceInTargetType = reflect.Append(sliceInTargetType, castedToTargetType)
				default:
					// handles primitive elements (int, string, ...)
					castedToTargetType := reflect.ValueOf(rawDefaultValueElement).Convert(elementType)
					sliceInTargetType = reflect.Append(sliceInTargetType, castedToTargetType)
				}

			}

			result.Def = sliceInTargetType.Interface()
		default:
			// cast the parsed default value to the target type
			castedToTargetType := reflect.ValueOf(parsedDefinition.Def).Convert(typeOfEntry)
			result.Def = castedToTargetType.Interface()
		}
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

		cfgSetting, hasCfgTag := field.Tag.Lookup("cfg")
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

func extractConfigDefinition(tCfg reflect.Type, nameOfParentType string, parent configTag) ([]config.Entry, error) {

	entries := make([]config.Entry, 0)

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

		cfgSetting, hasCfgTag := field.Tag.Lookup("cfg")
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
			subEntries, err := extractConfigDefinition(fType, fieldName, eDef)
			if err != nil {
				return nil, errors.Wrap(err, "Extracting subentries")
			}
			entries = append(entries, subEntries...)
			debug("%s added entries %v\n", logPrefix, entries)
			continue
		}

		// create and append the new config entry
		entry := config.NewEntry(eDef.Name, eDef.Description, config.Default(eDef.Def))
		entries = append(entries, entry)
		debug("%s created new entry=%v\n", logPrefix, entry)
	}
	return entries, nil
}
