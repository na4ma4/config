package config

import (
	"time"

	"go.uber.org/zap"
)

// Conf is the default configuration object interface.
type Conf interface {
	Get(key string) interface{}
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	GetFloat64(key string) float64
	GetInt(key string) int
	GetIntSlice(key string) []int
	GetString(key string) string
	GetStringSlice(key string) []string

	Set(key string, value interface{})
	SetBool(key string, value bool)
	SetDuration(key string, value time.Duration)
	SetFloat64(key string, value float64)
	SetInt(key string, value int)
	SetIntSlice(key string, value []int)
	SetString(key string, value string)
	SetStringSlice(key string, value []string)

	ZapConfig() zap.Config
	Save() error
	// Write(out io.Writer) error
}
