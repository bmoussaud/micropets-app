package internal

import (
	"fmt"

	"os"

	"github.com/spf13/viper"
)

//Config Structure
type Config struct {
	Service struct {
		Port   string
		Listen bool
		Mode string
		FrequencyError int
		Delay  struct {
			Period int
			Amplitude float64
		}
		From string
	}
	Observability struct {
		Application string
		Service     string
		Cluster     string
		Shard       string
		Server      string
		Token       string
		Source      string
		Enable      bool
	}
	//internal flag
	setup bool
}

var GlobalConfig Config
	
//LoadConfiguration method
func LoadConfiguration() Config {
	if !GlobalConfig.setup {
		viper.SetConfigType("json")
		viper.SetEnvPrefix("mp")           // will be uppercased automatically eg. MP_OBSERVABILITY.TOKEN=$(TO_TOKEN)
		viper.SetConfigName("pets_config") // name of config file (without extension)
		viper.AutomaticEnv()

		if serviceConfigDir := os.Getenv("SERVICE_CONFIG_DIR"); serviceConfigDir != "" {
			fmt.Printf("Load configuration from %s\n", serviceConfigDir)
			viper.AddConfigPath(serviceConfigDir)

		}
		//add default config path
		viper.AddConfigPath("/config/")  // path to look for the config file in
		viper.AddConfigPath("/etc/micropets/")  // path to look for the config file in
		viper.AddConfigPath("$HOME/.micropets") // call multiple times to add many search paths
		viper.AddConfigPath(".")                // optionally look for config in the working directory

		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {             // Handle errors reading the config file
			panic(fmt.Errorf("fatal error config file: %s ", err))
		}

		err = viper.Unmarshal(&GlobalConfig)
		if err != nil {
			panic(fmt.Errorf("unable to decode into struct, %v", err))
		}


		GlobalConfig.setup = true
		fmt.Printf("Resolved Configuration\n")
		fmt.Printf("%+v\n", GlobalConfig)

	}
	return GlobalConfig
}
