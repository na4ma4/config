package config_test

import (
	"sync"
	"testing"
	"time"

	"github.com/na4ma4/config"

	"github.com/spf13/viper"
)

func expectGetString(t *testing.T, vcfg config.Conf, key, expectValue string) {
	t.Helper()

	if v := vcfg.GetString(key); v != expectValue {
		t.Errorf("GetString(): got '%s', want '%s'", v, expectValue)
	}
}

func expectGetDuration(t *testing.T, vcfg config.Conf, key string, expectValue time.Duration) {
	t.Helper()

	if v := vcfg.GetDuration(key); v.String() != expectValue.String() {
		t.Errorf("GetDuration(): got '%s', want '%s'", v.String(), expectValue.String())
	}
}

func expectGetInt(t *testing.T, vcfg config.Conf, key string, expectValue int) {
	t.Helper()

	if v := vcfg.GetInt(key); v != expectValue {
		t.Errorf("GetInt(): got '%d', want '%d'", v, expectValue)
	}
}

//nolint:gochecknoglobals // viper is a singleton, so we need to lock it for tests that use it.
var viperLock sync.Mutex

func viperGetViper(t *testing.T) *viper.Viper {
	t.Helper()
	viperLock.Lock()
	defer viperLock.Unlock()

	v := viper.GetViper()

	o := viper.New()
	if err := o.MergeConfigMap(v.AllSettings()); err != nil {
		t.Fatalf("MergeConfigMap(): error, got '%s', want 'nil'", err)
	}

	return o
}
