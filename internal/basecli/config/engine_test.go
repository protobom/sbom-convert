package config

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

type TestConfigWithBase struct {
	BaseApplication `yaml:",omitempty,inline" json:",omitempty,inline" mapstructure:",squash"`
	Str             string `yaml:"str,omitempty" json:"str,omitempty" mapstructure:"str"`
	Int             int    `yaml:"int,omitempty" json:"int,omitempty" mapstructure:"int"`
}

func (d TestConfigWithBase) GetConfigPath() string {
	return d.ConfigPath
}

type TestConfig struct {
	ConfigPath string `yaml:"config_path,omitempty" json:"config_path,omitempty" mapstructure:"config_path,omitempty"`
	Str        string `yaml:"str,omitempty" json:"str,omitempty" mapstructure:"str"`
	Int        int    `yaml:"int,omitempty" json:"int,omitempty" mapstructure:"int"`
	NoOmit     string `yaml:"noomit" json:"noomit" mapstructure:"testint"`
}

func (d TestConfig) GetConfigPath() string {
	return d.ConfigPath
}

func chdir(t *testing.T, dir string) func() {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("chdir %s: %v", dir, err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	return func() {
		if err := os.Chdir(wd); err != nil {
			t.Fatalf("restoring working directory: %v", err)
		}
	}
}

func TestReadConfig(t *testing.T) {
	logger := logrus.New()
	testCases := []struct {
		name              string
		applicationConfig ApplicationConfig
		applicationName   string
		err               bool
		expectedErr       error
	}{
		{
			name:              "no config file found",
			applicationConfig: &TestConfig{},
			applicationName:   "no_config_application",
			err:               true,
			expectedErr:       ErrApplicationConfigNotFound,
		},
		{
			name: "error config file path",
			applicationConfig: &TestConfig{
				ConfigPath: "error_cfg.yml",
			},
			applicationName: "test",
			err:             true,
		},
		{
			name: "good config file path",
			applicationConfig: &TestConfig{
				ConfigPath: "good_cfg.yml",
			},
			applicationName: "test",
			err:             false,
		},
		{
			name: "good config file path with used with base config",
			applicationConfig: &TestConfigWithBase{
				BaseApplication: BaseApplication{
					ConfigPath: "good_cfg.yml",
				},
			},
			applicationName: "test", //Look for .goodapp/config.yaml in local dir
			err:             false,
		},
		{
			name:              "good config by application name",
			applicationConfig: &TestConfig{},
			applicationName:   "goodapp", //Look for .goodapp/config.yaml in local dir
			err:               false,
		},
		{
			name:              "good config by application name",
			applicationConfig: &TestConfig{},
			applicationName:   "goodapp", //Look for .goodapp.yaml in local dir
			err:               false,
		},
		{
			name:              "error config by application name",
			applicationConfig: &TestConfig{},
			applicationName:   "errorapp", //Look for .errorapp.yaml in local dir
			err:               true,
		},
	}
	defer chdir(t, "fixtures")()

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			conf := New(logger, test.applicationName)
			err := conf.readConfig(test.applicationConfig)
			if test.err {
				assert.Error(t, err)
				if test.expectedErr != nil {
					assert.ErrorIs(t, err, ErrApplicationConfigNotFound)
				}
				return
			}
			assert.NoError(t, err)

		})
	}

}

func TestLoadApplicationConfig(t *testing.T) {
	logger := logrus.New()
	testCases := []struct {
		name              string
		applicationConfig ApplicationConfig
		applicationName   string
		err               bool
		expectedConfig    string
	}{
		{
			name:              "no config file",
			applicationConfig: &TestConfig{},
			applicationName:   "no_config_application",
			err:               false,
			expectedConfig: `noomit: ""
`,
		},
		{
			name:              "good config by application name in subdir",
			applicationConfig: &TestConfig{},
			applicationName:   "goodapp", //Look for .goodapp/config.yaml in local dir
			err:               false,
			expectedConfig: `str: teststr_value
int: 10
noomit: ""
`,
		},
	}

	defer chdir(t, "fixtures")()

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			conf := New(logger, test.applicationName)
			err := conf.LoadApplicationConfig(test.applicationConfig)
			if test.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				v, err := yaml.Marshal(test.applicationConfig)
				assert.Nil(t, err)
				assert.Equal(t, string(v), test.expectedConfig)
			}
		})
	}

}
