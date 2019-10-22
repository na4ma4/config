package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/pelletier/go-toml"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// ViperConf is a Conf compatible Viper configuration object
type ViperConf struct {
	viper    *viper.Viper
	mutex    *sync.Mutex
	filename string
}

// NewViperConfigFromViper returns a Conf compatible ViperConf object copied from the system viper.Viper.
func NewViperConfigFromViper(vcfg *viper.Viper, filename ...string) Conf {
	allset := vcfg.AllSettings()
	v := &ViperConf{
		viper:    viper.New(),
		mutex:    &sync.Mutex{},
		filename: vcfg.ConfigFileUsed(),
	}
	for key, val := range allset {
		v.Set(key, val)
	}
	if len(filename) > 0 {
		v.filename = filename[0]
	}
	return v
}

// NewViperConfig returns a Conf compatible ViperConf object.
func NewViperConfig(project string, filename ...string) Conf {
	if len(filename) > 0 {
		for i, fname := range filename {
			v := &ViperConf{
				viper:    viper.New(),
				mutex:    &sync.Mutex{},
				filename: fname,
			}
			err := v.readFromFile(project, fname)
			if i == len(filename)-1 {
				// If filenames are specified, the last one is used as the fallback
				// and is then used for the `Save()` method.
				v.setFilename(fname)
				return v
			}

			// Error loading file, and not the last filename in the list
			if err != nil {
				continue
			}

			// No error, so the file was loaded successfully
			v.filename = v.viper.ConfigFileUsed()
			if v.viper.ConfigFileUsed() == "" {
				continue
			}
			return v
		}
	}

	fname := fmt.Sprintf("%s.toml", project)
	v := &ViperConf{
		viper:    viper.New(),
		mutex:    &sync.Mutex{},
		filename: fname,
	}
	v.initConfig(project)
	if !strings.EqualFold(v.viper.ConfigFileUsed(), "") {
		v.filename = v.viper.ConfigFileUsed()
	}
	return v
}

func (v *ViperConf) readFromFile(project, filename string) error {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	v.viper.SetConfigName(project)
	v.viper.SetConfigType("toml")
	v.viper.SetConfigFile(filename)
	return v.viper.ReadInConfig()
}

func (v *ViperConf) setFilename(filename string) {
	v.mutex.Lock()
	v.filename = filename
	v.viper.SetConfigType("toml")
	v.viper.SetConfigFile(filename)
	v.mutex.Unlock()
}

func (v *ViperConf) initConfig(project string) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	v.viper.SetConfigName(project)
	v.viper.SetConfigType("toml")
	v.viper.AddConfigPath("./artifacts")
	v.viper.AddConfigPath("./test")
	v.viper.AddConfigPath("$HOME/.config")
	v.viper.AddConfigPath("/etc")
	v.viper.AddConfigPath(fmt.Sprintf("/etc/%s", project))
	v.viper.AddConfigPath(fmt.Sprintf("/usr/local/%s/etc", project))
	v.viper.AddConfigPath("/run/secrets")
	v.viper.AddConfigPath(".")

	v.viper.ReadInConfig()
}

// Get can retrieve any value given the key to use.
// Get is case-insensitive for a key.
// Get has the behavior of returning the value associated with the first
// place from where it is set. Viper will check in the following order:
// override, flag, env, config file, key/value store, default
//
// Get returns an interface. For a specific value use one of the Get____ methods.
func (v *ViperConf) Get(key string) interface{} {
	v.mutex.Lock()
	val := v.viper.Get(key)
	v.mutex.Unlock()
	return val
}

// GetBool returns the value associated with the key as a boolean.
func (v *ViperConf) GetBool(key string) bool {
	v.mutex.Lock()
	val := v.viper.GetBool(key)
	v.mutex.Unlock()
	return val
}

// GetDuration returns the value associated with the key as a duration.
func (v *ViperConf) GetDuration(key string) time.Duration {
	v.mutex.Lock()
	val := v.viper.GetDuration(key)
	v.mutex.Unlock()
	return val
}

// GetFloat64 returns the value associated with the key as a float64.
func (v *ViperConf) GetFloat64(key string) float64 {
	v.mutex.Lock()
	val := v.viper.GetFloat64(key)
	v.mutex.Unlock()
	return val
}

// GetInt returns the value associated with the key as an int.
func (v *ViperConf) GetInt(key string) int {
	v.mutex.Lock()
	val := v.viper.GetInt(key)
	v.mutex.Unlock()
	return val
}

// GetIntSlice returns the value associated with the key as a slice of ints.
func (v *ViperConf) GetIntSlice(key string) []int {
	v.mutex.Lock()
	val := cast.ToIntSlice(v.viper.Get(key))
	v.mutex.Unlock()
	return val
}

// GetString returns the value associated with the key as a string.
func (v *ViperConf) GetString(key string) string {
	v.mutex.Lock()
	val := v.viper.GetString(key)
	v.mutex.Unlock()
	return val
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (v *ViperConf) GetStringSlice(key string) []string {
	v.mutex.Lock()
	val := v.viper.GetStringSlice(key)
	v.mutex.Unlock()
	return val
}

// Set sets the value for the key in the viper object.
func (v *ViperConf) Set(key string, value interface{}) {
	v.mutex.Lock()
	v.viper.Set(key, value)
	v.mutex.Unlock()
}

// SetBool sets the value for the key in the viper object.
func (v *ViperConf) SetBool(key string, value bool) {
	v.mutex.Lock()
	v.viper.Set(key, value)
	v.mutex.Unlock()
}

// SetDuration sets the value for the key in the viper object.
func (v *ViperConf) SetDuration(key string, value time.Duration) {
	v.mutex.Lock()
	v.viper.Set(key, value)
	v.mutex.Unlock()
}

// SetFloat64 sets the value for the key in the viper object.
func (v *ViperConf) SetFloat64(key string, value float64) {
	v.mutex.Lock()
	v.viper.Set(key, value)
	v.mutex.Unlock()
}

// SetInt sets the value for the key in the viper object.
func (v *ViperConf) SetInt(key string, value int) {
	v.mutex.Lock()
	v.viper.Set(key, value)
	v.mutex.Unlock()
}

// SetIntSlice sets the value for the key in the viper object.
func (v *ViperConf) SetIntSlice(key string, value []int) {
	v.mutex.Lock()
	v.viper.Set(key, value)
	v.mutex.Unlock()
}

// SetString sets the value for the key in the viper object.
func (v *ViperConf) SetString(key string, value string) {
	v.mutex.Lock()
	v.viper.Set(key, value)
	v.mutex.Unlock()
}

// SetStringSlice sets the value for the key in the viper object.
func (v *ViperConf) SetStringSlice(key string, value []string) {
	v.mutex.Lock()
	v.viper.Set(key, value)
	v.mutex.Unlock()
}

// Save writes the config to the file system.
func (v *ViperConf) Save() error {
	v.mutex.Lock()
	if err := os.MkdirAll(filepath.Dir(v.filename), os.ModePerm); err != nil {
		return err
	}
	if _, err := os.Create(v.filename); err != nil {
		return err
	}
	err := v.viper.WriteConfigAs(v.filename)
	v.mutex.Unlock()
	return err
}

func (v *ViperConf) Write(out io.Writer) error {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	c := v.viper.AllSettings()
	t, err := toml.TreeFromMap(c)
	if err != nil {
		return err
	}
	s := t.String()
	if _, err := io.WriteString(out, s); err != nil {
		return err
	}
	return nil
}

// ZapConfig returns a zap logger configuration derived from settings in the viper config.
func (v *ViperConf) ZapConfig() zap.Config {
	var cfg zap.Config
	v.mutex.Lock()
	if v.viper.GetBool("debug") {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}
	v.mutex.Unlock()
	return cfg
}
