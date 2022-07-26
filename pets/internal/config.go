package internal

import (
	"context"
	"fmt"

	"os"
	"strconv"

	"github.com/spf13/viper"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

//Config Structure
type Config struct {
	Service struct {
		Port   string
		Listen bool
	}
	Backends []struct {
		Name    string `json:"name"`
		Host    string `json:"host"`
		Port    string `json:"port"`
		Context string `json:"context"`
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

		if len(GlobalConfig.Backends) == 0 {
			fmt.Printf("* No defined backends, use dynamic mode!\n")
			var dynamicConfig = QueryBackendService()
			GlobalConfig.Backends = dynamicConfig.Backends
			DumpBackendConfig(GlobalConfig)
		}

		GlobalConfig.setup = true
		fmt.Printf("Resolved Configuration\n")
		fmt.Printf("%+v\n", GlobalConfig)

	}
	return GlobalConfig
}

func DumpBackendConfig(config Config) {
	fmt.Printf("******* Backends are:\n")
	for i, backend := range config.Backends {
		fmt.Printf("* Managing %d\t %s\t %s:%s%s\n", i, backend.Name, backend.Host, backend.Port, backend.Context)
	}
}

func QueryBackendService() Config {
	var config Config
	ctx := context.Background()
	k8sconfig := ctrl.GetConfigOrDie()
	clientset := kubernetes.NewForConfigOrDie(k8sconfig)

	//TODO: manage namespace
	namespace := "dev-tap"
	items, err := GetK8SServices(clientset, ctx, namespace)

	if err != nil {
		fmt.Println(err)
		return config
	} else {
		for _, item := range items {
			var svcName = item.ObjectMeta.Name
			var svcPort int32 = item.Spec.Ports[0].Port
			config.Backends = append(config.Backends, struct {
				Name    string "json:\"name\""
				Host    string "json:\"host\""
				Port    string "json:\"port\""
				Context string "json:\"context\""
			}{svcName, fmt.Sprintf("%s.%s.svc.cluster.local", svcName, namespace), strconv.FormatUint(uint64(svcPort), 10), fmt.Sprintf("/%s/v1/data", svcName)})
		}
	}
	return config
}

func GetK8SServices(clientset *kubernetes.Clientset, ctx context.Context,
	namespace string) ([]v1.Service, error) {
	listOptions := metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/part-of=micro-pet, micropets/kind=backend",
		Limit:         100,
	}

	list, err := clientset.CoreV1().Services(namespace).
		List(ctx, listOptions)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}
