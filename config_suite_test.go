package config_test

import (
	"testing"
	"time"

	"github.com/na4ma4/config"
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
