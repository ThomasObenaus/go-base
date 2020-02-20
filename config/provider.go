package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Provider is a structure containing the parsed configuration
type Provider struct {
	// config entries are all definitions of config entries that should be regarded
	configEntries []Entry

	// this special entry is used to specify the name + location of the config file
	configFileEntry Entry

	configName string

	// the environment prefix (will be added to all env vars) <envPrefix>_<name of config entry>
	// e.g. assuming the envPrefix is "myApp" and the name of the config entry is "my-entry"
	// then the env var is MYAPP_MY_ENTRY
	envPrefix string

	// instance of pflag, needed to parse command line parameters
	pFlagSet *pflag.FlagSet

	// instance of viper, needed to parse env vars and to read from cfg-file
	*viper.Viper
}

// ProviderOption represents an option for the Provider
type ProviderOption func(p *Provider)

// CfgFile specifies a default value
func CfgFile(parameterName, shortParameterName string) ProviderOption {
	return func(p *Provider) {
		p.configFileEntry = NewEntry(parameterName, "Specifies the full path and name of the configuration file", ShortName(shortParameterName))
	}
}

// NewProvider creates a new config provider that is able to parse the command line, env vars and config file based
// on the given entries
func NewProvider(configEntries []Entry, configName, envPrefix string, options ...ProviderOption) Provider {

	defaultConfigFileEntry := NewEntry("config-file", "Specifies the full path and name of the configuration file", Bind(true, true))
	provider := Provider{
		configEntries:   configEntries,
		configName:      configName,
		envPrefix:       envPrefix,
		pFlagSet:        pflag.NewFlagSet(configName, pflag.ContinueOnError),
		Viper:           viper.New(),
		configFileEntry: defaultConfigFileEntry,
	}

	// apply the options
	for _, opt := range options {
		opt(&provider)
	}

	return provider
}

func (p Provider) String() string {
	return fmt.Sprintf("%s: %v", p.configName, p.AllSettings())
}
