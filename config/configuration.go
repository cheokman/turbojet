package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"turbojet/cli"
)

const (
	configPath               = "/.turbojet"
	configFile               = "config.json"
	DefaultConfigProfileName = "default"
)

type Configuration struct {
	CurrentProfile string    `json:"current"`
	Profiles       []Profile `json:"profiles"`
	MetaPath       string    `json:"meta_path"`
}

func NewConfiguration() Configuration {
	return Configuration{
		CurrentProfile: DefaultConfigProfileName,
		Profiles: []Profile{
			NewProfile(DefaultConfigProfileName),
		},
	}
}

func (c *Configuration) NewProfile(pn string) Profile {
	p, ok := c.GetProfile(pn)
	if !ok {
		p = NewProfile(pn)
		c.PutProfile(p)
	}
	return p
}

func (c *Configuration) GetCurrentProfile(ctx *cli.Context) Profile {
	var profileName string
	if envPN := os.Getenv("TJ_PROVIDER"); envPN != "" {
		profileName = envPN
	} else {
		profileName = ProfileFlag(ctx.Flags()).GetStringOrDefault(c.CurrentProfile)
	}
	p, _ := c.GetProfile(profileName)
	p.OverwriteWithFlags(ctx)
	return p
}

func (c *Configuration) GetProfile(pn string) (Profile, bool) {
	for _, p := range c.Profiles {
		if p.Name == pn {
			return p, true
		}
	}
	return Profile{Name: pn}, false
}

func (c *Configuration) PutProfile(profile Profile) {
	for i, p := range c.Profiles {
		if p.Name == profile.Name {
			c.Profiles[i] = profile
			return
		}
	}
	c.Profiles = append(c.Profiles, profile)
}

func LoadConfiguration(path string, w io.Writer) (Configuration, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return NewConfiguration(), nil
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return NewConfiguration(), fmt.Errorf("reading config from '%s' failed %v", path, err)
	}

	return NewConfigFromBytes(bytes)
}

func SaveConfiguration(config Configuration) error {
	bytes, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}
	path := GetConfigPath() + string(os.PathSeparator) + configFile
	err = ioutil.WriteFile(path, bytes, 0600)
	if err != nil {
		return err
	}
	return nil
}

func NewConfigFromBytes(bytes []byte) (Configuration, error) {
	var conf Configuration
	err := json.Unmarshal(bytes, &conf)
	if err != nil {
		return conf, err
	}
	return conf, nil
}

func GetConfigPath() string {
	path := GetHomePath() + configPath
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			panic(err)
		}
	}
	return path
}

func GetHomePath() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}

	return os.Getenv("HOME")
}

func LoadProfile(path string, w io.Writer, name string) (Profile, error) {
	var p Profile
	config, err := LoadConfiguration(path, w)
	if err != nil {
		return p, fmt.Errorf("init config failed %v", err)
	}

	if name == "" {
		name = config.CurrentProfile
	}
	p, ok := config.GetProfile(name)
	if !ok {
		return p, fmt.Errorf("unknown profile %s, run configure to check", name)
	}
	return p, nil
}

func LoadProfileWithContext(ctx *cli.Context) (profile Profile, err error) {
	var currentPath string
	if envCP := os.Getenv("TJ_CONFIG_PATH"); envCP != "" {
		currentPath = envCP
	} else if path, ok := ConfigurePathFlag(ctx.Flags()).GetValue(); ok {
		currentPath = path
	} else {
		currentPath = GetConfigPath() + string(os.PathSeparator) + configFile
	}

	if name, ok := ProfileFlag(ctx.Flags()).GetValue(); ok {
		profile, err = LoadProfile(currentPath, ctx.Writer(), name)
	} else {
		profile, err = LoadProfile(currentPath, ctx.Writer(), "")
	}
	if err != nil {
		return
	}
	profile.OverwriteWithFlags(ctx)
	err = profile.Validate()
	return
}
