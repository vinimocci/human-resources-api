package utils

import(
	"human-resources-api/commons"
	"github.com/pelletier/go-toml"
)

func GetCurrentEnvironment (environment string) (*toml.Tree, error){
	var tomlErr error
	var config *toml.Tree

	if environment == commons.DevelopmentEnvironment {
		config, tomlErr = toml.LoadFile("config.development.toml")
	}

	if environment == commons.StagingEnvironment {
		config, tomlErr = toml.LoadFile("config.staging.toml")
	}

	if environment == commons.ProductionEnvironment {
		config, tomlErr = toml.LoadFile("config.production.toml")
	}

	if environment == commons.EmptyResult {
		panic("missing enviroment variable on current system")
	}

	return config, tomlErr
}