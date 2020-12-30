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

// TODO: Fail in case there are duplicate settings (names) configured
// TODO: Custom function hooks for complex parsing
// TODO: Custom logger hook function
// TODO: Check if pointer fields are supported

// HINT: Desired schema:
// cfg:"name:<name>;;desc:<description>;;default:<default value>"
// ';;' is the separator
// if no default value is given then the config field is treated as required

type Cfg struct {
	ShouldBeSkipped string         // this should be ignored since its not annotated
	Name            string         `cfg:"{'name':'name','desc':'the name of the config'}"`
	Prio            int            `cfg:"{'name':'prio','desc':'the prio','default':0}"`
	Immutable       bool           `cfg:"{'name':'immutable','desc':'can be modified or not','default':false}"`
	NumericLevels   []int          `cfg:"{'name':'numeric-levels','desc':'allowed levels','default':[1,2]}"`
	Levels          []string       `cfg:"{'name':'levels','desc':'allowed levels','default':['a','b']}"`
	ConfigStore     configStore    `cfg:"{'name':'config-store','desc':'the config store'}"`
	TargetSecrets   []targetSecret `cfg:"{'name':'target-secrets','desc':'list of target secrets','default':[{'name':'1mysecret1','key':'sdlfks','count':231},{'name':'mysecret2','key':'sdlfks','count':231}]}"`
}

type configStore struct {
	FilePath     string       `cfg:"{'name':'file-path','desc':'the path','default':'configs/'}"`
	TargetSecret targetSecret `cfg:"{'name':'target-secret','desc':'the secret'}"`
}

type targetSecret struct {
	Name  string `cfg:"{'name':'name','desc':'the name of the secret'}"`
	Key   string `cfg:"{'name':'key','desc':'the key'}"`
	Count int    `cfg:"{'name':'count','desc':'the count','default':0}"`
}

func main() {

	args := []string{
		"--name=hello",
		"--prio=23",
		"--immutable=true",
		"--config-store.file-path=/devops",
		"--config-store.target-secret.key=#lsdpo93",
		"--config-store.target-secret.name=mysecret",
		"--config-store.target-secret.count=2323",
		"--numeric-levels=1,2,3",
		"--target-secrets=[{'name':'mysecret1','key':'sdlfks','count':231},{'name':'mysecret2','key':'sdlfks','count':231}]",
	}

	parsedConfig, err := New(args, "ABCDE")
	if err != nil {
		panic(err)
	}
	fmt.Println("")
	fmt.Println("SUCCESS")
	spew.Dump(parsedConfig)
}

func unmarshal(provider config.Provider, target interface{}) error {
	return config.Apply(provider, target)
}

func isSliceOfStructs(t reflect.Type) bool {
	if t.Kind() != reflect.Slice {
		return false
	}
	elementType := t.Elem()
	return elementType.Kind() == reflect.Struct
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

	extractedConfigEntries, err := config.Extract(&cfg)
	if err != nil {
		return Cfg{}, errors.Wrapf(err, "Extracting config tags from %v", cfgType)
	}
	configEntries = append(configEntries, extractedConfigEntries...)

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
