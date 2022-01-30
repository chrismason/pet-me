package shared

import "github.com/spf13/cobra"

type EnvConfig struct {
	Name          string
	BaseUrl       string
	IsDevelopment bool
}

var envConfigMap = map[string]EnvConfig{
	"development": {
		Name:          "development",
		BaseUrl:       "http://localhost:8080",
		IsDevelopment: true,
	},
	"staging": {
		Name:    "staging",
		BaseUrl: "https://pet-me-staging.azurewebsites.net",
	},
	"production": {
		Name:    "production",
		BaseUrl: "https://pet-me.azurewebsites.net",
	},
}

type ConnectionInfo struct {
	Environment string
}

func RegisterEnvironmentFlags(info *ConnectionInfo, cmd *cobra.Command) {
	cmd.Flags().StringVarP(&info.Environment, "env", "e", "", "An environment name to connect to")
}

func GetEnvConfig(envName string) (EnvConfig, bool) {
	if envName == "" {
		return envConfigMap["production"], true
	}

	if config, ok := envConfigMap[envName]; ok {
		return config, true
	}

	return EnvConfig{}, false
}
