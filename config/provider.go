package config

import (
	"fmt"

	gconf "github.com/ThomasObenaus/go-conf"
	gconfIf "github.com/ThomasObenaus/go-conf/interfaces"
)

// Provider is a structure containing the parsed configuration
type Provider struct {
	gconfIf.Provider
}

type providerCfg struct {
	parameterName      string
	shortParameterName string
}

// ProviderOption represents an option for the Provider
type ProviderOption func(cfg *providerCfg)

// CfgFile specifies a default value
func CfgFile(parameterName, shortParameterName string) ProviderOption {
	return func(cfg *providerCfg) {
		cfg.parameterName = parameterName
		cfg.shortParameterName = shortParameterName
	}
}

// NewProvider creates a new config provider that is able to parse the command line, env vars and config file based
// on the given entries
func NewProvider(configEntries []Entry, configName, envPrefix string, options ...ProviderOption) Provider {
	configEntriesT := make([]gconf.Entry, 0, len(configEntries))
	for _, entry := range configEntries {
		configEntriesT = append(configEntriesT, entryToGConfEntry(entry))
	}

	optionsT := pOptsToGConfPOpts(options)
	provider, err := gconf.NewProvider(configEntriesT, configName, envPrefix, optionsT...)
	if err != nil {
		panic(fmt.Sprintf("Error creating config provider: %s", err))
	}

	return Provider{
		Provider: provider,
	}
}

func (p Provider) String() string {
	return p.Provider.String()
}

func entryToGConfEntry(entry Entry) gconf.Entry {
	opts := []gconf.EntryOption{
		gconf.Bind(entry.Bind()),
		gconf.Default(entry.DefaultValue()),
		gconf.ShortName(entry.ShortName()),
	}

	return gconf.NewEntry(entry.Name(), entry.Usage(), opts...)
}

func pOptsToGConfPOpts(opts []ProviderOption) []gconf.ProviderOption {
	pCfg := &providerCfg{}
	for _, opt := range opts {
		opt(pCfg)
	}
	pOpts := []gconf.ProviderOption{}

	if len(pCfg.parameterName) > 0 && len(pCfg.shortParameterName) > 0 {
		pOpts = append(pOpts, gconf.CfgFile(pCfg.parameterName, pCfg.shortParameterName))
	}
	return pOpts
}
