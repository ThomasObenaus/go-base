package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ThomasObenaus/go-base/config"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
)

/*
{
	"name":"conf1",
	"prio":1,
	"immutable":true,
	"config_store": {
	  "file-path": "cfgs",
	  "target-secrets": [
		{
		  "name": "secret1",
		  "key": "23948239842kdsdkfj",
		  "count": 12
		}
	  ]
	},
	"tasks": {
	  "table-name": "tasks",
	  "use-db": true
	}
}

// according struct
type cfg struct {
	Name string `cfg:"name:name;;desc:the name of the config"`
	Prio int `cfg:"name:prio;;desc:the prio;;default:0"`
	Immutable bool `cfg:"name:immutable;;desc:can be modified or not;;default:false"`
	ConfigStore configStore `cfg:"name:config-store;;desc:the config store"`
}

type configStore struct {
	FilePath string `cfg:"name:file-path;;desc:the path;;default:configs/"`
	TargetSecrets []targetSecret `cfg:"name:target-secrets;;desc:list of target secrets;;"`
	//TargetSecrets []targetSecret `cfg:"name:target-secrets;;desc:list of target secrets;;default:[{}]"`
}

type targetSecret struct {
	Name string `cfg:"name:name;;desc:the name of the config"`
	Key string `cfg:"name:name;;desc:the name of the config"`
	Count int `cfg:"name:name;;desc:the name of the config;;default:0"`
}
*/

// TODO: Fail in case there are duplicate settings configured
// TODO: Custom function hooks for complex parsing
// TODO: Handle missing default as required

// HINT: Desired schema:
// cfg:"name:<name>;;desc:<description>;;default:<default value>"
// ';;' is the separator
// if no default value is given then the config field is treated as required

type Cfg struct {
	ShouldBeSkipped string // this should be ignored since its not annotated
	Name            string `cfg:"name:name;;desc:the name of the config"`
	Prio            int    `cfg:"name:prio;;desc:the prio;;default:0"`
	Immutable       bool   `cfg:"name:immutable;;desc:can be modified or not;;default:false"`
	//ConfigStore configStore `cfg:"name:config-store;;desc:the config store"`
}

type configStore struct {
	FilePath      string         `cfg:"name:file-path;;desc:the path;;default:configs/"`
	TargetSecrets []targetSecret `cfg:"name:target-secrets;;desc:list of target secrets;;"`
	//TargetSecrets []targetSecret `cfg:"name:target-secrets;;desc:list of target secrets;;default:[{}]"`
}

type targetSecret struct {
	Name  string `cfg:"name:name;;desc:the name of the config"`
	Key   string `cfg:"name:name;;desc:the name of the config"`
	Count int    `cfg:"name:name;;desc:the name of the config;;default:0"`
}

func main() {

	args := []string{
		"--prio=10",
	} //, "--a.c=Hello World", "--a.b.duration=51m42s"} //, "--duration=15m", "--really", "--name=Harry"}

	parsedConfig, err := New(args, "ABCDE")
	if err != nil {
		panic(err)
	}
	fmt.Println("")
	fmt.Println("SUCCESS")
	spew.Dump(parsedConfig)
}

func unmarshal(provider config.Provider, target interface{}) error {
	return apply(provider, target, "")
}

func apply(provider config.Provider, target interface{}, nameOfParent string) error {
	tCfg := reflect.TypeOf(target)
	vCfg := reflect.ValueOf(target)

	isNilPtr := vCfg.Kind() == reflect.Ptr && vCfg.IsNil()
	isNotSupportedField := vCfg.Kind() != reflect.Ptr
	if isNotSupportedField || isNilPtr {
		return fmt.Errorf("Can't handle %v (kind=%s,value=%v) (probably the type is not supported)", tCfg, tCfg.Kind(), vCfg)
	}

	// use the element type since we have a pointer
	tCfg = tCfg.Elem()
	vCfg = vCfg.Elem()

	debug("[Apply-(%s)] typeOfInputConfig=%v stateOfInputConfig=%v\n", nameOfParent, tCfg, vCfg)

	for i := 0; i < tCfg.NumField(); i++ {
		field := tCfg.Field(i)
		fType := field.Type
		v := vCfg.Field(i)
		fieldValue := v.Addr().Interface()
		printableName := fmt.Sprintf("%s.%s", nameOfParent, field.Name)
		logPrefix := fmt.Sprintf("[Apply-(%s)]", printableName)

		debug("%s type=%s\n", logPrefix, fType)

		// find out if we already have a primitive type
		isPrimitive, err := isOfPrimitiveType(fType)
		if err != nil {
			return errors.Wrapf(err, "Checking for primitive type failed for field '%s'", printableName)
		}
		debug("%s is primitive=%t\n", logPrefix, isPrimitive)

		// handling of non primitives (stucts)
		if !isPrimitive {
			if err := apply(provider, fieldValue, printableName); err != nil {
				return errors.Wrap(err, "Applying non primitive")
			}
			debug("%s applied non primitive %v\n", logPrefix, fieldValue)
			continue
		}

		cfgSetting, hasCfgTag := field.Tag.Lookup("cfg")
		// ignore fields without a config tag
		if !hasCfgTag {
			debug("%s no tag found entry will be skipped\n", logPrefix)
			continue
		}
		debug("%s tag found cfgSetting=%v\n", logPrefix, cfgSetting)

		eDef, err := parseCfgEntry(cfgSetting, fType)
		if err != nil {
			return errors.Wrapf(err, "Parsing the config definition failed for field '%s'", printableName)
		}

		if !provider.IsSet(eDef.name) {
			debug("%s parameter not provided, nothing will be applied\n", logPrefix)
			continue
		}

		// apply the value
		val := provider.Get(eDef.name)
		newValue := reflect.ValueOf(val)
		v.Set(newValue)
		debug("%s apply value '%v' to '%s' based on config '%s'\n", logPrefix, newValue, printableName, eDef.name)
	}
	return nil
}

var verbose = true

func debug(format string, a ...interface{}) {
	if verbose {
		fmt.Print("[DBG]")
		fmt.Printf(format, a...)
	}
}

func extractConfigDefinition(tCfg reflect.Type, nameOfParent string) ([]config.Entry, error) {

	entries := make([]config.Entry, 0)

	// use the element type if we have a pointer
	if tCfg.Kind() == reflect.Ptr {
		tCfg = tCfg.Elem()
	}
	debug("[Extract-(%s) cfgType=%v\n", nameOfParent, tCfg)

	for i := 0; i < tCfg.NumField(); i++ {
		field := tCfg.Field(i)
		fType := field.Type

		printableName := fmt.Sprintf("%s.%s", nameOfParent, field.Name)
		logPrefix := fmt.Sprintf("[Extract-(%s)]", printableName)
		debug("%s type=%s\n", logPrefix, fType)

		// find out if we already have a primitive type
		isPrimitive, err := isOfPrimitiveType(fType)
		if err != nil {
			return nil, errors.Wrapf(err, "Checking for primitive type failed for field '%s'", printableName)
		}

		debug("%s is primitive=%t\n", logPrefix, isPrimitive)
		if !isPrimitive {
			subEntries, err := extractConfigDefinition(fType, printableName)
			if err != nil {
				return nil, errors.Wrap(err, "Extracting subentries")
			}
			entries = append(entries, subEntries...)
			debug("%s added entries %v\n", logPrefix, entries)
			continue
		}

		cfgSetting, hasCfgTag := field.Tag.Lookup("cfg")
		// skip all fields without the cfg tag
		if !hasCfgTag {
			debug("%s no tag found entry will be skipped\n", logPrefix)
			continue
		}
		debug("%s tag found cfgSetting=%v\n", logPrefix, cfgSetting)

		eDef, err := parseCfgEntry(cfgSetting, fType)
		if err != nil {
			return nil, errors.Wrapf(err, "Parsing the config definition failed for field '%s'", printableName)
		}
		debug("%s parsed config entry=%v\n", logPrefix, eDef)

		// create and append the new config entry
		entry := config.NewEntry(eDef.name, eDef.description, config.Default(eDef.def))
		entries = append(entries, entry)

		debug("%s created new entry=%v\n", logPrefix, entry)
	}
	return entries, nil
}

type entryDefinition struct {
	name        string
	description string
	def         interface{}
}

func setValueFromString(v reflect.Value, strVal string) error {

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		// handle duration
		if v.Type() == reflect.TypeOf(time.Duration(0)) {
			dur, err := time.ParseDuration(strVal)
			if err != nil {
				return err
			}
			v.SetInt(dur.Nanoseconds())
			return nil
		}

		// handle the usual int
		val, err := strconv.ParseInt(strVal, 0, 64)
		if err != nil {
			return err
		}
		if v.OverflowInt(val) {
			return errors.New("Int value too big: " + strVal)
		}
		v.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(strVal, 0, 64)
		if err != nil {
			return err
		}
		if v.OverflowUint(val) {
			return errors.New("UInt value too big: " + strVal)
		}
		v.SetUint(val)
	case reflect.Float32:
		val, err := strconv.ParseFloat(strVal, 32)
		if err != nil {
			return err
		}
		v.SetFloat(val)
	case reflect.Float64:
		val, err := strconv.ParseFloat(strVal, 64)
		if err != nil {
			return err
		}
		v.SetFloat(val)
	case reflect.String:
		v.SetString(strVal)
	case reflect.Bool:
		val, err := strconv.ParseBool(strVal)
		if err != nil {
			return err
		}
		v.SetBool(val)
	case reflect.Slice:
		arr, err := strToValueSlice(v.Type().Elem(), strVal)
		if err != nil {
			return errors.Wrap(err, "Setting a value from given string for a slice.")
		}
		v.Set(arr)
	default:
		return fmt.Errorf("Unsupported kind: %s", v.Kind())
	}
	return nil
}

// strToValueSlice takes the given string and tries to convert it to a slice of reflect.Type
// but using the according type of the given reflect.Type.
// It is expected that the string has the form "<element_1>,<element_n+1>, ...,<element_n>".
// It is expected that the values of the given array (encoded in the given string) can be converted to the given type.
// The parameter elementType is used as the target type of the slice to be generated.
// Only primitive types are supported.
func strToValueSlice(elementType reflect.Type, strVal string) (reflect.Value, error) {
	debug("strToValueSlice(%v,'%s')\n", elementType, strVal)

	splittedValues := strings.Split(strVal, ",")
	numSplittedValues := len(splittedValues)

	arr := reflect.MakeSlice(reflect.SliceOf(elementType), 0, numSplittedValues)
	for _, element := range splittedValues {
		value, err := strToValue(elementType, element)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "Converting string to slice of reflect.Value")
		}
		arr = reflect.Append(arr, value)
	}
	return arr, nil
}

// strToValue converts the given string into a reflect.Value of the given elementType
func strToValue(elementType reflect.Type, strVal string) (reflect.Value, error) {

	// special treatment for time.Duration
	if elementType == reflect.TypeOf(time.Duration(0)) {
		dur, err := strToDuration(strVal)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(dur), nil
	}

	switch elementType.Kind() {
	case reflect.String:
		return reflect.ValueOf(strVal), nil
	case reflect.Int:
		val, err := strToInt64(elementType, strVal)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(int(val)), nil
	case reflect.Int16:
		val, err := strToInt64(elementType, strVal)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(int16(val)), nil
	case reflect.Int32:
		val, err := strToInt64(elementType, strVal)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(int32(val)), nil
	case reflect.Int64:
		val, err := strToInt64(elementType, strVal)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(int64(val)), nil
	case reflect.Uint:
		val, err := strToUInt64(elementType, strVal)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(uint(val)), nil
	case reflect.Uint16:
		val, err := strToUInt64(elementType, strVal)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(uint16(val)), nil
	case reflect.Uint32:
		val, err := strToUInt64(elementType, strVal)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(uint32(val)), nil
	case reflect.Uint64:
		val, err := strToUInt64(elementType, strVal)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(uint64(val)), nil
	case reflect.Float32:
		val, err := strconv.ParseFloat(strVal, 32)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(float32(val)), nil
	case reflect.Float64:
		val, err := strconv.ParseFloat(strVal, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(val), nil
	case reflect.Bool:
		val, err := strconv.ParseBool(strVal)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(val), nil
	default:
		return reflect.Value{}, fmt.Errorf("Unable to convert '%s' to type '%s' (type not supported)", strVal, elementType)
	}
}

func strToDuration(strVal string) (time.Duration, error) {
	dur, err := time.ParseDuration(strVal)
	if err != nil {
		return 0, errors.Wrap(err, "Parsing str to time.Duration")
	}
	return dur, nil
}

func strToInt64(elementType reflect.Type, strVal string) (int64, error) {
	val, err := strconv.ParseInt(strVal, 0, 64)
	if err != nil {
		return 0, errors.Wrap(err, "Parsing str to int")
	}

	// check for overflow
	v := reflect.New(elementType).Elem()
	if v.OverflowInt(int64(val)) {
		return 0, fmt.Errorf("Int value too big: %s for %s", strVal, elementType)
	}
	return val, nil
}

func strToUInt64(elementType reflect.Type, strVal string) (uint64, error) {
	val, err := strconv.ParseUint(strVal, 0, 64)
	if err != nil {
		return 0, errors.Wrap(err, "Parsing str to int")
	}

	// check for overflow
	v := reflect.New(elementType).Elem()
	if v.OverflowUint(uint64(val)) {
		return 0, fmt.Errorf("Uint value too big: %s for %s", strVal, elementType)
	}
	return val, nil
}

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

func (e entryDefinition) String() string {
	return fmt.Sprintf(`n:"%s",d:"%s",df:%v`, e.name, e.description, e.def)
}

func parseCfgEntry(setting string, cfgType reflect.Type) (entryDefinition, error) {
	setting = strings.TrimSpace(setting)
	parts := strings.Split(setting, ";;")

	elements := make(map[string]string)
	result := entryDefinition{}
	for _, part := range parts {
		kvp := strings.Split(part, ":")

		if len(kvp) != 2 {
			return entryDefinition{}, fmt.Errorf("unexpected len kvp (2!=%d)", len(kvp))
		}

		key := strings.ToLower(strings.TrimSpace(kvp[0]))
		value := kvp[1]
		elements[key] = value
	}

	name, ok := elements["name"]
	if !ok {
		return entryDefinition{}, fmt.Errorf("Config key 'name' is missing but must be set (e.g. cfg:\"name:setting-one\")")
	}
	result.name = name

	desc, ok := elements["desc"]
	if !ok {
		return entryDefinition{}, fmt.Errorf("Config key 'desc' is missing but must be set (e.g. cfg:\"desc:this setting does that\")")
	}
	result.description = desc

	defaultValue, ok := elements["default"]
	if ok {
		value := reflect.New(cfgType)
		if err := setValueFromString(value.Elem(), defaultValue); err != nil {
			return entryDefinition{}, errors.Wrap(err, "Parsing value from string")
		}
		result.def = value.Elem().Interface()
	}

	return result, nil
}

var port = config.NewEntry("port", "Port where sokar is listening.", config.Default(11000))
var dryRun = config.NewEntry("dry-run", "If true, then sokar won't execute the planned scaling action. Only scaling\n"+
	"actions triggered via ScaleBy end-point will be executed.", config.Default(false))
var configEntries = []config.Entry{
	port,
	dryRun,
}

func New(args []string, serviceAbbreviation string) (Cfg, error) {
	cfg := Cfg{}
	cfgType := reflect.TypeOf(cfg)

	newConfigEntries, err := extractConfigDefinition(cfgType, "")
	if err != nil {
		return Cfg{}, err
	}
	configEntries = append(configEntries, newConfigEntries...)

	provider := config.NewProvider(configEntries, serviceAbbreviation, serviceAbbreviation)
	err = provider.ReadConfig(args)
	if err != nil {
		return Cfg{}, err
	}

	if err := unmarshal(provider, &cfg); err != nil {
		return Cfg{}, err
	}

	if err := cfg.fillCfgValues(provider); err != nil {
		return Cfg{}, err
	}

	return cfg, nil
}

func (cfg *Cfg) fillCfgValues(provider config.Provider) error {
	//cfg.DryRun = provider.GetBool(dryRun.Name())
	//cfg.Port = provider.GetInt(port.Name())
	//
	//	cfg.Setting3 = "Thomas (OVERWRITTEN)"
	return nil
}
