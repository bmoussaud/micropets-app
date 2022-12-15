package internal

import (
	"context"
	"fmt"

	"os"
	"strconv"

	"github.com/spf13/viper"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Config Structure
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

// LoadConfiguration method
func LoadConfiguration() Config {
	if !GlobalConfig.setup {
		var LocalConfig Config

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

		fmt.Printf("* Unmarshal config!\n")
		err = viper.Unmarshal(&LocalConfig)
		if err != nil {
			panic(fmt.Errorf("unable to decode into struct, %v", err))
		}

		if len(LocalConfig.Backends) == 0 {
			fmt.Printf("* No defined backends, use dynamic mode!\n")
			var dynamicConfig = QueryBackendService()
			LocalConfig.Backends = dynamicConfig.Backends
			DumpBackendConfig(LocalConfig)
		}

		fmt.Printf("Resolved Configuration\n")
		GlobalConfig = LocalConfig
		//re-read the configuration again & again
		GlobalConfig.setup = false
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
	fmt.Printf("* QueryBackendService....\n")
	var config Config

	//TODO: manage namespace
	namespace := "dev-tap"
	config, err := GetK8SKNativeServices(namespace)

	if err != nil {
		fmt.Println(err)
		return config
	} else {

		if len(config.Backends) == 0 {
			fmt.Printf("* No KNative services switch back to Svc....\n")
			config, _ := GetK8SServices(namespace)
			fmt.Printf("* QueryBackendService config %+v\n", config)
			return config
		}

	}
	fmt.Printf("* QueryBackendService config %+v\n", config)
	return config
}

func GetK8SServices(namespace string) (Config, error) {

	ctx := context.Background()
	k8sconfig := ctrl.GetConfigOrDie()
	clientset := kubernetes.NewForConfigOrDie(k8sconfig)

	listOptions := metav1.ListOptions{
		LabelSelector: "micropets/kind=backend",
		Limit:         100,
	}

	fmt.Printf("* GetK8SServices in %s: labelSelector is %s\n", namespace, listOptions.LabelSelector)
	var config Config
	items, err := clientset.CoreV1().Services(namespace).
		List(ctx, listOptions)
	if err != nil {
		return config, err
	} else {
		fmt.Printf("* GetK8SServices found size:%d\n", len(items.Items))
		for _, item := range items.Items {
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

	return config, nil
}

func GetK8SKNativeServices(namespace string) (Config, error) {
	var config Config
	ksvcRes := schema.GroupVersionResource{Group: "serving.knative.dev", Version: "v1", Resource: "services"}
	listOptions := metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/part-of=micropets-app, micropets/kind=backend",
		Limit:         100,
	}

	fmt.Printf("* GetK8SKNativeServices in %s: labelSelector is %s\n", namespace, listOptions.LabelSelector)
	ctx := context.Background()
	k8sconfig := ctrl.GetConfigOrDie()
	client := dynamic.NewForConfigOrDie(k8sconfig)

	list, err := client.Resource(ksvcRes).Namespace(namespace).List(ctx, listOptions)
	if err != nil {
		return config, err
	}
	fmt.Printf("* GetK8SKNativeServices found size:%d\n", len(list.Items))
	for _, d := range list.Items {
		address, found, err := unstructured.NestedMap(d.Object, "status", "address")
		url := address["url"]
		if err != nil || !found {
			fmt.Printf("ERR: URL not found for ksvc %s: error=%s\n", d.GetName(), err)
			continue
		}
		fmt.Printf(" * %s -> url : %s \n", d.GetName(), url)

		config.Backends = append(config.Backends, struct {
			Name    string "json:\"name\""
			Host    string "json:\"host\""
			Port    string "json:\"port\""
			Context string "json:\"context\""
		}{d.GetName(), fmt.Sprintf("%s", url), "80", fmt.Sprintf("/%s/v1/data", d.GetName())})
	}
	return config, nil
}
