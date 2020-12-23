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

// ViperConfD is a Conf compatible Viper configuration object.
type ViperConfD struct {
	viper    *viper.Viper
	lock     *sync.Mutex
	filename string
}

// NewViperConfDFromViper returns a Conf compatible ViperConfD object copied from the system viper.Viper.
func NewViperConfDFromViper(vcfg *viper.Viper, confdpath string, filename ...string) Conf {
	allset := vcfg.AllSettings()
	v := &ViperConfD{
		viper:    viper.New(),
		lock:     &sync.Mutex{},
		filename: vcfg.ConfigFileUsed(),
	}

	for key, val := range allset {
		v.Set(key, val)
	}

	if len(filename) > 0 {
		v.filename = filename[0]
	}

	_ = v.loadConfigPath(confdpath)

	return v
}

// NewViperConfD returns a Conf compatible ViperConfD object.
func NewViperConfD(project string, confdpath string, filename ...string) Conf {
	if len(filename) > 0 {
		for i, fname := range filename {
			v := &ViperConfD{
				viper:    viper.New(),
				lock:     &sync.Mutex{},
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

			_ = v.loadConfigPath(confdpath)

			return v
		}
	}

	fname := fmt.Sprintf("%s.toml", project)
	v := &ViperConfD{
		viper:    viper.New(),
		lock:     &sync.Mutex{},
		filename: fname,
	}
	v.initConfig(project)

	if !strings.EqualFold(v.viper.ConfigFileUsed(), "") {
		v.filename = v.viper.ConfigFileUsed()
	}

	_ = v.loadConfigPath(confdpath)

	return v
}

func (v *ViperConfD) readFromFile(project, filename string) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	v.viper.SetConfigName(project)
	v.viper.SetConfigType("toml")
	v.viper.SetConfigFile(filename)

	return v.viper.ReadInConfig()
}

func (v *ViperConfD) setFilename(filename string) {
	v.lock.Lock()
	v.filename = filename
	v.viper.SetConfigType("toml")
	v.viper.SetConfigFile(filename)
	v.lock.Unlock()
}

func (v *ViperConfD) loadConfigPath(confdpath string) error {
	if confdpath != "" {
		v.lock.Lock()
		v.viper.SetConfigType("toml")
		v.lock.Unlock()

		abspath, err := filepath.Abs(confdpath)
		if err != nil {
			abspath = confdpath
		}

		m, err := filepath.Glob(fmt.Sprintf("%s/*.toml", abspath))
		if err != nil || m == nil {
			return fmt.Errorf("unable to find config files \"%s\": %w", fmt.Sprintf("%s/*.toml", abspath), err)
		}

		for _, fn := range m {
			if err := v.mergeConfigFile(fn); err != nil {
				return err
			}
		}
	}

	return nil
}

func (v *ViperConfD) mergeConfigFile(filename string) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("unable to open config file \"%s\": %w", filename, err)
	}
	defer f.Close()

	if err = v.viper.MergeConfig(f); err != nil {
		return fmt.Errorf("unable to merge config file \"%s\": %w", filename, err)
	}

	return nil
}

func (v *ViperConfD) initConfig(project string) {
	v.lock.Lock()
	defer v.lock.Unlock()

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

	_ = v.viper.ReadInConfig()
}

// SetDefault sets the default value for this key.
// SetDefault is case-insensitive for a key.
// Default only used when no value is provided by the user via flag, config or ENV.
func (v *ViperConfD) SetDefault(key string, value interface{}) {
	v.lock.Lock()
	defer v.lock.Unlock()
	v.viper.SetDefault(key, value)
}

// AllSettings merges all settings and returns them as a map[string]interface{}.
func (v *ViperConfD) AllSettings() map[string]interface{} {
	return v.viper.AllSettings()
}

// Get can retrieve any value given the key to use.
// Get is case-insensitive for a key.
// Get has the behavior of returning the value associated with the first
// place from where it is set. Viper will check in the following order:
// override, flag, env, config file, key/value store, default
//
// Get returns an interface. For a specific value use one of the Get____ methods.
func (v *ViperConfD) Get(key string) interface{} {
	v.lock.Lock()
	defer v.lock.Unlock()
	val := v.viper.Get(key)

	return val
}

// GetBool returns the value associated with the key as a boolean.
func (v *ViperConfD) GetBool(key string) bool {
	v.lock.Lock()
	defer v.lock.Unlock()
	val := v.viper.GetBool(key)

	return val
}

// GetDuration returns the value associated with the key as a duration.
func (v *ViperConfD) GetDuration(key string) time.Duration {
	v.lock.Lock()
	defer v.lock.Unlock()
	val := v.viper.GetDuration(key)

	return val
}

// GetFloat64 returns the value associated with the key as a float64.
func (v *ViperConfD) GetFloat64(key string) float64 {
	v.lock.Lock()
	defer v.lock.Unlock()
	val := v.viper.GetFloat64(key)

	return val
}

// GetInt returns the value associated with the key as an int.
func (v *ViperConfD) GetInt(key string) int {
	v.lock.Lock()
	defer v.lock.Unlock()
	val := v.viper.GetInt(key)

	return val
}

// GetIntSlice returns the value associated with the key as a slice of ints.
func (v *ViperConfD) GetIntSlice(key string) []int {
	v.lock.Lock()
	defer v.lock.Unlock()
	val := cast.ToIntSlice(v.viper.Get(key))

	return val
}

// GetString returns the value associated with the key as a string.
func (v *ViperConfD) GetString(key string) string {
	v.lock.Lock()
	defer v.lock.Unlock()
	val := v.viper.GetString(key)

	return val
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (v *ViperConfD) GetStringSlice(key string) []string {
	v.lock.Lock()
	defer v.lock.Unlock()
	val := v.viper.GetStringSlice(key)

	return val
}

// Set sets the value for the key in the viper object.
func (v *ViperConfD) Set(key string, value interface{}) {
	v.lock.Lock()
	defer v.lock.Unlock()
	v.viper.Set(key, value)
}

// SetBool sets the value for the key in the viper object.
func (v *ViperConfD) SetBool(key string, value bool) {
	v.lock.Lock()
	defer v.lock.Unlock()
	v.viper.Set(key, value)
}

// SetDuration sets the value for the key in the viper object.
func (v *ViperConfD) SetDuration(key string, value time.Duration) {
	v.lock.Lock()
	defer v.lock.Unlock()
	v.viper.Set(key, value)
}

// SetFloat64 sets the value for the key in the viper object.
func (v *ViperConfD) SetFloat64(key string, value float64) {
	v.lock.Lock()
	defer v.lock.Unlock()
	v.viper.Set(key, value)
}

// SetInt sets the value for the key in the viper object.
func (v *ViperConfD) SetInt(key string, value int) {
	v.lock.Lock()
	defer v.lock.Unlock()
	v.viper.Set(key, value)
}

// SetIntSlice sets the value for the key in the viper object.
func (v *ViperConfD) SetIntSlice(key string, value []int) {
	v.lock.Lock()
	defer v.lock.Unlock()
	v.viper.Set(key, value)
}

// SetString sets the value for the key in the viper object.
func (v *ViperConfD) SetString(key string, value string) {
	v.lock.Lock()
	defer v.lock.Unlock()
	v.viper.Set(key, value)
}

// SetStringSlice sets the value for the key in the viper object.
func (v *ViperConfD) SetStringSlice(key string, value []string) {
	v.lock.Lock()
	defer v.lock.Unlock()
	v.viper.Set(key, value)
}

// Save writes the config to the file system.
func (v *ViperConfD) Save() error {
	v.lock.Lock()
	defer v.lock.Unlock()

	if err := os.MkdirAll(filepath.Dir(v.filename), os.ModePerm); err != nil {
		return fmt.Errorf("unable to create directory: %w", err)
	}

	if _, err := os.Create(v.filename); err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}

	return v.viper.WriteConfigAs(v.filename)
}

func (v *ViperConfD) Write(out io.Writer) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	c := v.viper.AllSettings()

	t, err := toml.TreeFromMap(c)
	if err != nil {
		return fmt.Errorf("unable to make tree from map: %w", err)
	}

	s := t.String()

	if _, err := io.WriteString(out, s); err != nil {
		return fmt.Errorf("unable to write config file: %w", err)
	}

	return nil
}

// ZapConfig returns a zap logger configuration derived from settings in the viper config.
func (v *ViperConfD) ZapConfig() zap.Config {
	v.lock.Lock()
	defer v.lock.Unlock()

	if v.viper.GetBool("debug") {
		return zap.NewDevelopmentConfig()
	}

	return zap.NewProductionConfig()
}
