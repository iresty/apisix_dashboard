package conf

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_mergeSchema(t *testing.T) {
	type args struct {
		apisixSchema    []byte
		customizeSchema []byte
	}
	tests := []struct {
		name           string
		args           args
		wantRes        []byte
		wantErr        bool
		wantErrMessage string
	}{
		{
			name: "should failed when have duplicates key",
			args: args{
				apisixSchema:    []byte(`{"main":{"a":1,"b":2},"plugins":{"a":1}}`),
				customizeSchema: []byte(`{"main":{"b":1}}`),
			},
			wantErr:        true,
			wantErrMessage: "duplicates key: main.b between schema.json and customize_schema.json",
		},
		{
			name: "should success",
			args: args{
				apisixSchema:    []byte(`{"main":{"a":1,"b":2},"plugins":{"a":1}}`),
				customizeSchema: []byte(`{"main":{"c":3}}`),
			},
			wantErr: false,
			wantRes: []byte(`{"main":{"a":1,"b":2,"c":3},"plugins":{"a":1}}`),
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			var (
				wantMap map[string]interface{}
				gotMap  map[string]interface{}
			)

			got, err := mergeSchema(tt.args.apisixSchema, tt.args.customizeSchema)
			if tt.wantErr {
				assert.Equal(t, tt.wantErrMessage, err.Error())
				return
			}

			assert.NoError(t, err)
			err = json.Unmarshal(got, &gotMap)
			assert.NoError(t, err)
			err = json.Unmarshal(tt.wantRes, &wantMap)
			assert.NoError(t, err)
			assert.Equal(t, wantMap, gotMap)
		})
	}
}

func Test_unmarshalConfig(t *testing.T) {
	tests := []struct {
		name   string
		init   func()
		config []byte
		assert func(*Config, *testing.T)
	}{
		{
			name:   "should correctly parse config without environment variables",
			init:   func() {},
			config: []byte("conf:\n  listen:\n    port: \"9000\""),
			assert: func(config *Config, t *testing.T) {
				assert.Equal(t, 9000, config.Conf.Listen.Port)
			},
		},
		{
			name: "should correctly substitute int port from environment variables",
			init: func() {
				os.Setenv("PORT", "8080")
			},
			config: []byte("conf:\n  listen:\n    port: \"$PORT\""),
			assert: func(config *Config, t *testing.T) {
				assert.Equal(t, 8080, config.Conf.Listen.Port)
			},
		},
		{
			name: "should correctly substitute string etcd endpoint from environment variables",
			init: func() {
				os.Setenv("ETCD_ENDPOINT", "127.0.0.1:2379")
			},
			config: []byte("conf:\n  etcd:\n    endpoints:\n      - $ETCD_ENDPOINT"),
			assert: func(config *Config, t *testing.T) {
				assert.Equal(t, "127.0.0.1:2379", config.Conf.Etcd.Endpoints[0])
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			viper.SetConfigType("yaml")
			viper.AutomaticEnv()

			tt.init()

			err := viper.ReadConfig(bytes.NewBuffer(tt.config))
			if err != nil {
				t.Errorf("unable to read config: %v", err)
				return
			}
			config, err := unmarshalConfig()
			if err != nil {
				t.Errorf("unable to unmarshall config: %v", err)
				return
			}

			tt.assert(config, t)
		})
	}
}
