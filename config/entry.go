package config

import (
	"fmt"
	"reflect"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Entry is one item to define a configuration
type Entry struct {
	name         string
	usage        string
	defaultValue interface{}

	bindFlag      bool
	flagShortName string

	bindEnv bool
}

// EntryOption represents an option for the Entry
type EntryOption func(e *Entry)

// Default specifies a default value
func Default(value interface{}) EntryOption {
	return func(e *Entry) {
		e.defaultValue = value
	}
}

// ShortName specifies the shorthand (one-letter) flag name
func ShortName(fShort string) EntryOption {
	return func(e *Entry) {
		e.flagShortName = fShort
	}
}

// Bind enables/ disables binding of flag and env var
func Bind(flag, env bool) EntryOption {
	return func(e *Entry) {
		e.bindFlag = flag
		e.bindEnv = env
	}
}

// NewEntry creates a new Entry that is available as flag, config file entry and environment variable
func NewEntry(name, usage string, options ...EntryOption) Entry {
	entry := Entry{
		name:          name,
		usage:         usage,
		flagShortName: "",
		defaultValue:  nil,
		bindFlag:      true,
		bindEnv:       true,
	}

	// apply the options
	for _, opt := range options {
		opt(&entry)
	}

	return entry
}

func (e Entry) String() string {
	return fmt.Sprintf("--%s (-%s) [default:%v (%T)]\t- %s", e.name, e.flagShortName, e.defaultValue, e.defaultValue, e.usage)
}

// Name provides the specified name for this entry
func (e Entry) Name() string {
	return e.name
}

func checkViper(vp *viper.Viper, entry Entry) error {
	if vp == nil {
		return fmt.Errorf("Viper is nil")
	}

	if len(entry.name) == 0 {
		return fmt.Errorf("Name is missing")
	}

	return nil
}

func registerFlag(flagSet *pflag.FlagSet, entry Entry) error {
	if !entry.bindFlag {
		return nil
	}
	if flagSet == nil {
		return fmt.Errorf("FlagSet is nil")
	}
	if len(entry.name) == 0 {
		return fmt.Errorf("Name is missing")
	}

	// no default value availabl -> we can't deduce the type
	if entry.defaultValue == nil {
		flagSet.StringP(entry.name, entry.flagShortName, "", entry.usage)
		return nil
	}

	switch entry.defaultValue.(type) {
	case string:
		flagSet.StringP(entry.name, entry.flagShortName, entry.defaultValue.(string), entry.usage)
	case uint:
		flagSet.UintP(entry.name, entry.flagShortName, entry.defaultValue.(uint), entry.usage)
	case int:
		flagSet.IntP(entry.name, entry.flagShortName, entry.defaultValue.(int), entry.usage)
	case bool:
		flagSet.BoolP(entry.name, entry.flagShortName, entry.defaultValue.(bool), entry.usage)
	case time.Duration:
		flagSet.DurationP(entry.name, entry.flagShortName, entry.defaultValue.(time.Duration), entry.usage)
	case float64:
		flagSet.Float64P(entry.name, entry.flagShortName, entry.defaultValue.(float64), entry.usage)
	case []bool:
		flagSet.BoolSliceP(entry.name, entry.flagShortName, entry.defaultValue.([]bool), entry.usage)
	case []string:
		flagSet.StringArrayP(entry.name, entry.flagShortName, entry.defaultValue.([]string), entry.usage)
	case []time.Duration:
		flagSet.DurationSliceP(entry.name, entry.flagShortName, entry.defaultValue.([]time.Duration), entry.usage)
	case []int:
		flagSet.IntSliceP(entry.name, entry.flagShortName, entry.defaultValue.([]int), entry.usage)
	case []uint:
		flagSet.UintSliceP(entry.name, entry.flagShortName, entry.defaultValue.([]uint), entry.usage)
	case []float64:
		flagSet.Float64SliceP(entry.name, entry.flagShortName, entry.defaultValue.([]float64), entry.usage)
	case []float32:
		flagSet.Float32SliceP(entry.name, entry.flagShortName, entry.defaultValue.([]float32), entry.usage)
	default:
		return fmt.Errorf("Type %s is not yet supported", reflect.TypeOf(entry.defaultValue))
	}

	return nil
}

func setDefault(vp *viper.Viper, entry Entry) error {
	if err := checkViper(vp, entry); err != nil {
		return err
	}

	if entry.defaultValue != nil {
		vp.SetDefault(entry.name, entry.defaultValue)
	}

	return nil
}

func registerEnv(vp *viper.Viper, envPrefix string, entry Entry) error {
	if !entry.bindEnv {
		return nil
	}
	if err := checkViper(vp, entry); err != nil {
		return err
	}

	if len(envPrefix) > 0 {
		vp.SetEnvPrefix(envPrefix)
	}
	return vp.BindEnv(entry.name)
}
