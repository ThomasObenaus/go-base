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
	return fmt.Sprintf(`name:"%s",desc:"%s",default:%v (%T),required=%t`, e.Name, e.Description, e.Def, e.Def, e.IsRequired())
}

func (e configTag) IsRequired() bool {
	return e.Def == nil
}

// parseConfigTagDefinition parses a definition like
// 	`cfg:"{'name':'<name of the config>','desc':'<description>','default':<default value>}"`
// to a configTag
func parseConfigTagDefinition(configTagStr string, typeOfEntry reflect.Type, nameOfParent string) (configTag, error) {
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

// extractConfigTagFromStructField extracts the configTag from the given StructField.
// Beside the extracted configTag a bool value indicating if the given type is a primitive type is returned.
func extractConfigTagFromStructField(field reflect.StructField, parent configTag) (isPrimitive bool, tag *configTag, err error) {
	fType := field.Type

	// find out if we have a primitive type
	isPrimitive, err = isOfPrimitiveType(fType)
	if err != nil {
		return false, nil, errors.Wrapf(err, "Checking for primitive type failed for field '%v'", field)
	}

	configTagDefinition, hasCfgTag := getConfigTagDefinition(field)
	if !hasCfgTag {
		return isPrimitive, nil, nil
	}

	cfgTag, err := parseConfigTagDefinition(configTagDefinition, fType, parent.Name)
	if err != nil {
		return isPrimitive, nil, errors.Wrapf(err, "Parsing the config definition ('%s') failed for field '%v'", configTagDefinition, field)
	}

	return isPrimitive, &cfgTag, nil
}

// extractConfigTagsOfStruct extracts recursively all configTags from the given struct.
// Fields of the target struct that are not annotated with a configTag are ignored.
//
// target - the target that should be processed (has to be a pointer to a struct)
// nameOfParentField - the name of the targets parent field. This is needed since this function runs recursively through the given target struct.
// parent - the configTag of the targets parent field. This is needed since this function runs recursively through the given target struct.
func extractConfigTagsOfStruct(target interface{}, nameOfParentField string, parent configTag) ([]configTag, error) {

	entries := make([]configTag, 0)

	targetType := reflect.TypeOf(target)

	debug("[Extract-(%s)] structure-type=%v definition=%v\n", nameOfParentField, targetType, parent)

	err := processAllConfigTagsOfStruct(target, nameOfParentField, parent, func(fieldName string, isPrimitive bool, fieldType reflect.Type, fieldValue reflect.Value, cfgTag configTag) error {
		logPrefix := fmt.Sprintf("[Extract-(%s)]", fieldName)

		if !isPrimitive {
			fieldValueIf := fieldValue.Addr().Interface()
			subEntries, err := extractConfigTagsOfStruct(fieldValueIf, fieldName, cfgTag)
			if err != nil {
				return errors.Wrap(err, "Extracting subentries")
			}
			entries = append(entries, subEntries...)
			debug("%s added %d configTags.\n", logPrefix, len(entries))
			return nil
		}

		entries = append(entries, cfgTag)
		debug("%s added configTag entry=%v.\n", logPrefix, cfgTag)

		return nil
	})

	if err != nil {
		return nil, errors.Wrapf(err, "Extracting config tags for %v", targetType)
	}
	return entries, nil
}

// handleConfigTagFunc function type for handling an extracted configTag of a given field
type handleConfigTagFunc func(fieldName string, isPrimitive bool, fieldType reflect.Type, fieldValue reflect.Value, cfgTag configTag) error

// processAllConfigTagsOfStruct finds the configTag on each field of the given struct. Each of this configTags will handled by the given handleConfigTagFunc.
// Fields of the target struct that are not annotated with a configTag are ignored (handleConfigTagFunc won't be called).
//
// target - the target that should be processed (has to be a pointer to a struct)
// nameOfParentField - the name of the targets parent field. This is needed since this function runs recursively through the given target struct.
// parent - the configTag of the targets parent field. This is needed since this function runs recursively through the given target struct.
// handleConfigTagFun - a function that should be used to handle each of the targets struct fields.
func processAllConfigTagsOfStruct(target interface{}, nameOfParentField string, parent configTag, handleConfigTagFun handleConfigTagFunc) error {
	if target == nil {
		return fmt.Errorf("The target must not be nil")
	}

	targetType, targetValue, err := getTargetTypeAndValue(target)
	if err != nil {
		return errors.Wrapf(err, "Obtaining target type and -value for target='%v',nameOfParentField='%s',parent='%s'", target, nameOfParentField, parent)
	}

	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		fieldValue := targetValue.Field(i)
		fType := field.Type

		fieldName := fullFieldName(nameOfParentField, field.Name)
		logPrefix := fmt.Sprintf("[Process-(%s)]", fieldName)
		debug("%s field-type=%s\n", logPrefix, fType)

		isPrimitive, cfgTag, err := extractConfigTagFromStructField(field, parent)
		if err != nil {
			return errors.Wrap(err, "Extracting config tag")
		}

		// skip the field in case there is no config tag
		if cfgTag == nil {
			debug("%s no tag found entry will be skipped.\n", logPrefix)
			continue
		}

		debug("%s parsed config entry=%v. Is primitive=%t.\n", logPrefix, cfgTag, isPrimitive)

		err = handleConfigTagFun(fieldName, isPrimitive, fType, fieldValue, *cfgTag)
		if err != nil {
			return errors.Wrapf(err, "Handling configTag %s for field '%s'", *cfgTag, fieldName)
		}
	}
	return nil
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
