package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ThomasObenaus/go-base/config"
	"github.com/pkg/errors"
)

// TODO: Fail in case there are duplicate settings configured
// TODO: SubStructs

type Cfg struct {
	//Setting1 time.Duration `cfg:"name:duration;;desc:A duration;;default:23h10m5s"`
	//Setting2 bool          `cfg:"name:really;;desc:A bool;;default:true"`
	//Setting3 string        `cfg:"name:name;;desc:A string;;default:Hans"`
	//Setting4 int           `cfg:"name:how-many;;desc:A int;;default:-19"`
	//Setting5 uint          `cfg:"name:max;;desc:A uint;;default:256"`
	//Setting6 float64       `cfg:"name:temp;;desc:A float;;default:-256.12302"`
	LevelA LevelA `cfg:"name:a;;desc:desc"`
	Port   int
	DryRun bool
	//Setting3 int    `cfg:"name:bla.setting2;;desc:This is;;default:sdfsdf"`
	//Setting4 int    `cfg:"name:bla.setting4;;desc:This is;;default:989"`
	//Setting1  string `cfg:"name:bla.setting-one111;;desc:This is;;default:bla_default"`
	//Setting2  string `cfg:"name:bla.setting-two;;desc:desc"`
	//Setting3  CfgSub `cfg:"name:bla.setting-three;;desc:desc"`
}

type LevelA struct {
	LevelA1 LevelB `cfg:"name:a.b;;desc:This is;;default:bla_default"`
	LevelA2 string `cfg:"name:a.c;;desc:desc"`
}

type LevelB struct {
	LevelB1 string `cfg:"name:a.b.1;;desc:This is;;default:bla_default"`
	LevelB2 string `cfg:"name:a.b.2;;desc:desc"`
}

func main() {

	args := []string{"--port=1234", "--dry-run"} //, "--duration=15m", "--really", "--name=Harry"}

	parsedConfig, err := New(args, "ABCDE")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Success %v\n", parsedConfig)
}

func unmarshal(provider config.Provider, target interface{}) error {

	apply(provider, target)
	return nil
}

func apply(provider config.Provider, target interface{}) {
	tCfg := reflect.TypeOf(target)
	vCfg := reflect.ValueOf(target)

	// TODO move this outside to the unmarshal func
	if vCfg.Kind() != reflect.Ptr || vCfg.IsNil() {
		panic("skfskfj")
	}

	// use the element type if we have a pointer
	if tCfg.Kind() == reflect.Ptr {
		tCfg = tCfg.Elem()
		vCfg = vCfg.Elem()
	}

	for i := 0; i < tCfg.NumField(); i++ {
		field := tCfg.Field(i)
		fType := field.Type
		v := vCfg.Field(i)
		cfgSetting, ok := field.Tag.Lookup("cfg")
		if !ok {
			continue
		}

		// find out if we already have a primitive type
		isPrimitive, err := isOfPrimitiveType(fType)
		if err != nil {
			fmt.Printf("Error ignoring '%s' because %s\n", cfgSetting, err.Error())
			continue
		}

		if !isPrimitive {
			apply(provider, v)
			continue
		}

		eDef, err := parseCfgEntry(cfgSetting, fType)
		if err != nil {
			fmt.Printf("Error ignoring '%s' because %s\n", cfgSetting, err.Error())
			continue
		}

		// apply the value
		if provider.IsSet(eDef.name) {
			val := provider.Get(eDef.name)
			v.Set(reflect.ValueOf(val))
		}
	}
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

	for i := 0; i < tCfg.NumField(); i++ {
		field := tCfg.Field(i)
		fType := field.Type

		printableName := fmt.Sprintf("%s.%s", nameOfParent, field.Name)
		debug("Extracting %s\n", printableName)

		// find out if we already have a primitive type
		isPrimitive, err := isOfPrimitiveType(fType)
		if err != nil {
			return nil, errors.Wrapf(err, "Checking for primitive type failed for field '%s'", printableName)
		}

		if !isPrimitive {
			subEntries, err := extractConfigDefinition(fType, printableName)
			if err != nil {
				return nil, err
			}
			entries = append(entries, subEntries...)
			continue
		}

		cfgSetting, hasCfgTag := field.Tag.Lookup("cfg")
		// skip all fields without the cfg tag
		if !hasCfgTag {
			continue
		}

		eDef, err := parseCfgEntry(cfgSetting, fType)
		if err != nil {
			return nil, errors.Wrapf(err, "Parsing the config definition failed for field '%s'", printableName)
		}

		// create and append the new config entry
		entry := config.NewEntry(eDef.name, eDef.description, config.Default(eDef.def))
		entries = append(entries, entry)

		debug("Extracted %v \n", eDef)
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
	default:
		return errors.New("Unsupported kind: " + v.Kind().String())
	}
	return nil
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
			return entryDefinition{}, err
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
	cfg.DryRun = provider.GetBool(dryRun.Name())
	cfg.Port = provider.GetInt(port.Name())

	//	cfg.Setting3 = "Thomas (OVERWRITTEN)"
	return nil
}
