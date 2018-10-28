package Utils

import "os"

type Configs struct {
	ConfigDefault string
	ConfigEnvKey string
}

var configMap = map [string]Configs{
	DATABASE_HOST_KEY : Configs{DATABASE_HOST,DATABASE_HOST_KEY},
	DATABASE_PORT_KEY : Configs{DATABASE_PORT, DATABASE_PORT_KEY},
	DATABASE_USERNAME_KEY : Configs{DATABASE_USERNAME, DATABASE_USERNAME_KEY},
	DATABASE_PASSWORD_KEY : Configs{DATABASE_PASSWORD, DATABASE_PASSWORD_KEY},

	WEBSERVICE_PORT_KEY : Configs{WEBSERVICE_PORT, WEBSERVICE_PORT_KEY},

	KEYSPACE_KEY : Configs{KEYSPACE, KEYSPACE_KEY},
}

func GetConfig(configKey string) string  {
	config := configMap[configKey]
	configVal := config.ConfigDefault

	envConfigVal := os.Getenv(config.ConfigEnvKey)
	if envConfigVal != "" {
		configVal = envConfigVal
	}
	return configVal
}
