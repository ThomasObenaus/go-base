package main

import (
	"fmt"

	"github.com/ThomasObenaus/go-base/config"
	"github.com/ThomasObenaus/go-base/config/interfaces"
	"github.com/davecgh/go-spew/spew"
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
// TODO: Check if pointer fields are supported
// TODO: Add support for shorthand flags

// HINT: Desired schema:
// cfg:"name:<name>;;desc:<description>;;default:<default value>"
// ';;' is the separator
// if no default value is given then the config field is treated as required

type Cfg struct {
	DryRun        bool           // this should be ignored since its not annotated, but it can be still read using on the usual way
	Name          string         `cfg:"{'name':'name','desc':'the name of the config'}"`
	Prio          int            `cfg:"{'name':'prio','desc':'the prio','default':0}"`
	Immutable     bool           `cfg:"{'name':'immutable','desc':'can be modified or not','default':false}"`
	NumericLevels []int          `cfg:"{'name':'numeric-levels','desc':'allowed levels','default':[1,2]}"`
	Levels        []string       `cfg:"{'name':'levels','desc':'allowed levels','default':['a','b']}"`
	ConfigStore   configStore    `cfg:"{'name':'config-store','desc':'the config store'}"`
	TargetSecrets []targetSecret `cfg:"{'name':'target-secrets','desc':'list of target secrets','default':[{'name':'1mysecret1','key':'sdlfks','count':231},{'name':'mysecret2','key':'sdlfks','count':231}]}"`
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
		"--dry-run",
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

var dryRun = config.NewEntry("dry-run", "If true, then sokar won't execute the planned scaling action. Only scaling\n"+
	"actions triggered via ScaleBy end-point will be executed.", config.Default(false))
var configEntries = []config.Entry{
	dryRun,
}

func New(args []string, serviceAbbreviation string) (Cfg, error) {
	cfg := Cfg{}

	provider, err := config.NewConfigProvider(
		&cfg,
		serviceAbbreviation,
		serviceAbbreviation,
		config.CustomConfigEntries(configEntries),
		config.Logger(interfaces.InfoLogger),
	)
	if err != nil {
		return Cfg{}, err
	}
	err = provider.ReadConfig(args)
	if err != nil {
		return Cfg{}, err
	}

	if err := cfg.fillCfgValues(provider); err != nil {
		return Cfg{}, err
	}

	return cfg, nil
}

func (cfg *Cfg) fillCfgValues(provider interfaces.Provider) error {
	cfg.DryRun = provider.GetBool(dryRun.Name())
	cfg.Name = "Thomas (OVERWRITTEN)"
	return nil
}
