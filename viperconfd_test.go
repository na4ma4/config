package config_test

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/na4ma4/config"
	"github.com/spf13/viper"
)

func TestViperConfD_ThreadSafe(t *testing.T) {
	v := config.NewViperConfD("test", "testdata/conf.d")

	v.SetString("test.string", "string")

	go func() {
		for range 1000 {
			v.SetString("test.string", "string")
		}
	}()

	go func() {
		for range 1000 {
			v.GetString("test.string")
		}
	}()

	expectGetString(t, v, "test.string", "string")
}

func TestViperConfD_LoadingFromFile(t *testing.T) {
	v := config.NewViperConfD("test-project", "testdata/conf.d")

	expectGetString(t, v, "category1.string", "foobar")
	expectGetInt(t, v, "category2.int", 8335)
}

func TestViperConfD_LoadingFromSpecifiedFile(t *testing.T) {
	v := config.NewViperConfD("test", "testdata/conf.d", "testdata/test-project.toml")

	expectGetString(t, v, "category1.string", "foobar")
	// should not find category2
	expectGetInt(t, v, "category2.int", 0)
}

func TestViperConfD_WriteToFile(t *testing.T) {
	tempfile, tempfileErr := os.CreateTemp(t.TempDir(), "*-dummy-file.toml")
	if tempfileErr != nil {
		t.Errorf("CreateTemp(): error, got '%s', want 'nil'", tempfileErr)
	}
	defer os.Remove(tempfile.Name())

	vcfg := config.NewViperConfD("test", "", tempfile.Name())

	vcfg.SetString("category.test", "barfoo")

	if err := vcfg.Save(); err != nil {
		t.Errorf("config.Save(): error, got '%s', want 'nil'", err)
	}

	b, bErr := os.ReadFile(tempfile.Name())
	if bErr != nil {
		t.Errorf("os.ReadFile(): error, got '%s', want 'nil'", bErr)
	}

	expectedOutput := "[category]\ntest = 'barfoo'\n"

	if diff := cmp.Diff(string(b), expectedOutput); diff != "" {
		t.Errorf("config.Save(): config file -got +want:\n%s", diff)
	}
}

func TestViperConfD_WriteToWriter(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	vcfg := config.NewViperConfD("test", "")

	vcfg.SetString("category.test", "barfoo")

	if v, ok := vcfg.(*config.ViperConfD); ok {
		if err := v.Write(buf); err != nil {
			t.Errorf("config.Write(): error, got '%s', want 'nil'", err)
		}

		expectedOutput := "\n[category]\n  test = \"barfoo\"\n"

		if diff := cmp.Diff(buf.String(), expectedOutput); diff != "" {
			t.Errorf("config.Write(): config file -got +want:\n%s", diff)
		}
	} else {
		t.Error("config.Write(): vcfg not config.ViperConfD")
	}
}

func TestViperConfD_ImportingSystemViper(t *testing.T) {
	viper.SetDefault("sesame.open", "open.sesame")
	viper.SetDefault("system.test.duration", "30s")
	viper.Set("fooman", "barwoman")

	viper.SetDefault("system.default", "default")

	vcfg := config.NewViperConfDFromViper(viper.GetViper(), "testdata/conf.d")

	vcfg.Set("system.default", "override")

	expectGetString(t, vcfg, "sesame.open", "open.sesame")
	expectGetString(t, vcfg, "fooman", "barwoman")
	expectGetString(t, vcfg, "system.default", "override")
	expectGetDuration(t, vcfg, "system.test.duration", 30*time.Second)
}

func TestViperConfD_CanSetDefaultValue(t *testing.T) {
	vcfg := config.NewViperConfDFromViper(viper.GetViper(), "test/conf.d")

	if vp, ok := vcfg.(*config.ViperConfD); ok {
		vp.SetDefault("some-key-with-default", "custom-default-value")
	}

	expectGetString(t, vcfg, "some-key-with-default", "custom-default-value")

	vcfg.Set("some-key-with-default", "new-value")

	expectGetString(t, vcfg, "some-key-with-default", "new-value")
}

func TestViperConfD_LoadTheConfDFiles(t *testing.T) {
	vcfg := config.NewViperConfDFromViper(viper.GetViper(), "testdata/conf.d")
	expectGetInt(t, vcfg, "category2.int", 8335)
}

func TestViperConfD_OrderLoad(t *testing.T) {
	vcfg := config.NewViperConfDFromViper(viper.GetViper(), "testdata/conf.d")

	expectGetString(t, vcfg, "category3.first", "foo")
	expectGetString(t, vcfg, "category3.second", "foobar")
}

func TestViperConfD_GetterSetter_Bool(t *testing.T) {
	vcfg := config.NewViperConfD("test-project", "")

	expect := true

	if v := vcfg.GetBool("typing.bool"); v != expect {
		t.Errorf("config.GetBool(): conf-val got '%t', want '%t'", v, expect)
	}

	vcfg.SetBool("typing.bool", false)
	expect = false
	if v := vcfg.GetBool("typing.bool"); v != expect {
		t.Errorf("config.GetBool(): set-val got '%t', want '%t'", v, expect)
	}

	expectString := "false"
	if v := vcfg.GetString("typing.bool"); v != expectString {
		t.Errorf("config.GetString(): stringify got '%s', want '%s'", v, expectString)
	}
}

func TestViperConfD_GetterSetter_Duration(t *testing.T) {
	vcfg := config.NewViperConfD("test-project", "")

	expect := 10 * time.Second

	if v := vcfg.GetDuration("typing.duration"); v != expect {
		t.Errorf("config.GetDuration(): conf-val got '%s', want '%s'", v, expect)
	}

	expect = 15 * time.Second
	vcfg.SetDuration("typing.duration", expect)
	if v := vcfg.GetDuration("typing.duration"); v != expect {
		t.Errorf("config.GetDuration(): set-val got '%s', want '%s'", v, expect)
	}

	expectString := "15s"
	if v := vcfg.GetString("typing.duration"); v != expectString {
		t.Errorf("config.GetString(): stringify got '%s', want '%s'", v, expectString)
	}
}

func TestViperConfD_GetterSetter_Float64(t *testing.T) {
	vcfg := config.NewViperConfD("test-project", "")

	expect := 3.1415

	if v := vcfg.GetFloat64("typing.float64"); v != expect {
		t.Errorf("config.GetFloat64(): conf-val got '%f', want '%f'", v, expect)
	}

	expect = 6.2
	vcfg.SetFloat64("typing.float64", expect)
	if v := vcfg.GetFloat64("typing.float64"); v != expect {
		t.Errorf("config.GetFloat64(): set-val got '%f', want '%f'", v, expect)
	}

	expectString := "6.2"
	if v := vcfg.GetString("typing.float64"); v != expectString {
		t.Errorf("config.GetString(): stringify got '%s', want '%s'", v, expectString)
	}
}

func TestViperConfD_GetterSetter_Int(t *testing.T) {
	vcfg := config.NewViperConfD("test-project", "")

	expect := 1337

	if v := vcfg.GetInt("typing.int"); v != expect {
		t.Errorf("config.GetInt(): conf-val got '%d', want '%d'", v, expect)
	}

	expect = 2600
	vcfg.SetInt("typing.int", expect)
	if v := vcfg.GetInt("typing.int"); v != expect {
		t.Errorf("config.GetInt(): set-val got '%d', want '%d'", v, expect)
	}

	expectString := "2600"
	if v := vcfg.GetString("typing.int"); v != expectString {
		t.Errorf("config.GetString(): stringify got '%s', want '%s'", v, expectString)
	}
}

func TestViperConfD_GetterSetter_IntSlice(t *testing.T) {
	vcfg := config.NewViperConfD("test-project", "")

	expect := []int{100, 200, 50}

	v := vcfg.GetIntSlice("typing.intslice")
	if diff := cmp.Diff(v, expect); diff != "" {
		t.Errorf("config.GetIntSlice(): conf-val -got +want:\n%s", diff)
	}

	expect = []int{1000, 2000, 500}
	vcfg.SetIntSlice("typing.intslice", expect)
	v = vcfg.GetIntSlice("typing.intslice")
	if diff := cmp.Diff(v, expect); diff != "" {
		t.Errorf("config.GetIntSlice(): set-val -got +want:\n%s", diff)
	}

	// expectString := ""
	// vs := vcfg.GetString("typing.intslice")
	// if diff := cmp.Diff(vs, expectString); diff != "" {
	// 	t.Errorf("config.GetIntSlice(): stringify -got +want:\n%s", diff)
	// }
}

func TestViperConfD_GetterSetter_String(t *testing.T) {
	vcfg := config.NewViperConfD("test-project", "")

	expect := "foobarmoo"

	if v := vcfg.GetString("typing.string"); v != expect {
		t.Errorf("config.GetString(): conf-val got '%s', want '%s'", v, expect)
	}

	expect = "moocowbar"
	vcfg.SetString("typing.string", expect)
	if v := vcfg.GetString("typing.string"); v != expect {
		t.Errorf("config.GetString(): set-val got '%s', want '%s'", v, expect)
	}
}

func TestViperConfD_GetterSetter_StringSlice(t *testing.T) {
	vcfg := config.NewViperConfD("test-project", "")

	expect := []string{"one", "two", "three"}

	v := vcfg.GetStringSlice("typing.stringslice")
	if diff := cmp.Diff(v, expect); diff != "" {
		t.Errorf("config.GetStringSlice(): conf-val -got +want:\n%s", diff)
	}

	expect = []string{"four", "five", "six"}
	vcfg.SetStringSlice("typing.stringslice", expect)
	v = vcfg.GetStringSlice("typing.stringslice")
	if diff := cmp.Diff(v, expect); diff != "" {
		t.Errorf("config.GetStringSlice(): set-val -got +want:\n%s", diff)
	}

	// expectString := ""
	// vs := vcfg.GetString("typing.stringslice")
	// if diff := cmp.Diff(vs, expectString); diff != "" {
	// 	t.Errorf("config.GetStringSlice(): stringify -got +want:\n%s", diff)
	// }
}

func TestViperConfD_GetterSetter_Interface(t *testing.T) {
	vcfg := config.NewViperConfD("test-project", "")

	expect := "foobarmoo"

	v := vcfg.Get("typing.string")
	if diff := cmp.Diff(v, expect); diff != "" {
		t.Errorf("config.Get(): conf-val -got +want:\n%s", diff)
	}

	expectNew := []string{"four", "five", "six"}
	vcfg.Set("typing.string", expectNew)
	v = vcfg.Get("typing.string")
	if diff := cmp.Diff(v, expectNew); diff != "" {
		t.Errorf("config.Get(): set-val -got +want:\n%s", diff)
	}
}

func TestViperConfD_ZapConfig(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		debug bool
	}{
		{"WithDebug", true},
		{"WithoutDebug", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			vcfg := config.NewViperConfD("test-project", "")

			vcfg.SetBool("debug", tt.debug)

			loggerConfig := vcfg.ZapConfig()

			if loggerConfig.Development != tt.debug {
				t.Errorf("config.ZapConfig(): development mode, got '%t', want '%t'", loggerConfig.Development, tt.debug)
			}
		})
	}
}
