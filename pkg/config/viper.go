package config

import (
	"doan/pkg/constants"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"gopkg.in/yaml.v2"
	"strings"
)

type Viper struct {
	ConfigType    string
	FilePath      string
	RemoteSchema  string
	RemoteAddress string
	RemoteKeys    string
	ConfigFile    string
}

func (v *Viper) InitConfig() error {
	if v.ConfigType == constants.ConfigTypeFile {
		return v.LoadConfigFromFile()
	} else {
		return v.LoadConfigFromConsul()
	}
}

func (v *Viper) LoadConfigFromFile() error {
	viper.AddConfigPath(v.FilePath)
	viper.SetConfigName(v.ConfigFile)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Errorf("Failed while read config from file, file-path: %s, err: %v", v.FilePath, err)
		return err
	}
	log.Infof("Config loaded from file %s/%s.%s", v.FilePath, v.ConfigFile, "yaml")
	return nil
}

func (v *Viper) loadConsulConfigMap(configRemoteKeys []string) map[string]interface{} {
	result := make(map[string]interface{})
	client, err := api.NewClient(
		&api.Config{
			Address:    v.RemoteAddress,
			Scheme:     v.RemoteSchema,
			Datacenter: "dc1",
		},
	)
	if err != nil {
		panic(err)
	}

	// Get a handle to the KV API
	kv := client.KV()

	for _, remoteConfigKey := range configRemoteKeys {
		pairs, _, err := kv.List(remoteConfigKey, nil)
		if err != nil {
			panic(err)
		}
		kvPair := pairs[0]
		mapValue := make(map[string]interface{})
		err = yaml.Unmarshal(kvPair.Value, mapValue)
		if err != nil {
			panic(err)
		}
		for key, value := range mapValue {
			result[key] = value
		}
	}
	return result
}

func (v *Viper) LoadConfigFromConsul() error {
	viper.SetEnvPrefix("app")
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()
	configRemoteKeys := strings.Split(v.RemoteKeys, ",")
	err := viper.MergeConfigMap(v.loadConsulConfigMap(configRemoteKeys))
	if err != nil {
		log.Errorf("Failed while read remote config, err: %v", err)
		return err
	}
	log.Infof("Config loaded from remote")
	return nil
}
