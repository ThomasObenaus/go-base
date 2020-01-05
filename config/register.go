package config

import (
	"os"
	"strings"

	"github.com/spf13/pflag"
)

func (cfg *Provider) registerEnvParams() error {
	replacer := strings.NewReplacer("-", "_", ".", "_")
	cfg.Viper.SetEnvKeyReplacer(replacer)

	// register also the config file parameter
	if err := registerEnv(cfg.Viper, cfg.envPrefix, cfg.configFileEntry); err != nil {
		return err
	}

	for _, entry := range cfg.configEntries {
		if err := registerEnv(cfg.Viper, cfg.envPrefix, entry); err != nil {
			return err
		}
	}

	return nil
}

func (cfg *Provider) registerAndParseFlags(args []string) error {

	// register also the config file parameter
	if err := registerFlag(cfg.pFlagSet, cfg.configFileEntry); err != nil {
		return err
	}

	for _, entry := range cfg.configEntries {
		if err := registerFlag(cfg.pFlagSet, entry); err != nil {
			return err
		}
	}

	if err := cfg.pFlagSet.Parse(args); err != nil {

		if err == pflag.ErrHelp {
			os.Exit(0)
		}
		return err
	}
	cfg.Viper.BindPFlags(cfg.pFlagSet)

	return nil
}

func (cfg *Provider) setDefaults() error {

	// regard also the config file parameter
	if err := setDefault(cfg.Viper, cfg.configFileEntry); err != nil {
		return err
	}

	for _, entry := range cfg.configEntries {
		if err := setDefault(cfg.Viper, entry); err != nil {
			return err
		}
	}
	return nil
}
