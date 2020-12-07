package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"sort"

	"github.com/imroc/req"
	"github.com/spf13/viper"
)

//Config Structure
type Config struct {
	Service struct {
		Port   string
		Listen bool
	}
	Backends []struct {
		Name string `json:"name"`
		Host string `json:"host"`
		Port string `json:"port"`
	}
}

//Pet Structure
type Pet struct {
	Name     string
	Type     string
	Kind     string
	Age      int
	URL      string
	Hostname string
}

//Path Structure
type Path struct {
	Service  string
	Hostname string
}

//Pets Structure
type Pets struct {
	Total     int
	Hostname  string
	Hostnames []Path
	Pets      []Pet `json:"Pets"`
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

var configLocation string = "config.json"

func lookupService(service string) string {

	fmt.Fprintf(os.Stderr, "-- Service %v\n", service)
	ips, err := net.LookupIP(service)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
		return service
	}

	for _, ip := range ips {
		fmt.Printf("%s. IN A %s\n", service, ip.String())
	}

	return service
}

func queryPets2(backend string) Pets {

	header := req.Header{
		"Accept":  "application/json",
		"Expires": "10ms",
	}
	var pets Pets
	req.Debug = true
	fmt.Printf("##########################@ 2 Connecting backend %s\n", backend)
	r, _ := req.Get(backend, header)
	r.ToJSON(&pets)

	return pets
}

func queryPets(backend string) Pets {
	var pets Pets
	var emptyPets Pets

	fmt.Printf("Connecting backend %s\n", backend)
	resp, herr := http.Get(backend)

	if herr != nil {
		fmt.Printf("Error communicating with backend: %v", herr)
		return emptyPets
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Backend seems unhealthy: %v", resp)
		return emptyPets
	}

	body, ierr := ioutil.ReadAll(resp.Body)

	if ierr != nil {
		fmt.Printf("Backend producing garbage: %v", ierr)
		return emptyPets
	}

	err := json.Unmarshal([]byte(body), &pets)
	if err != nil {
		panic(err)
	}
	return pets
}

func index(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	fmt.Printf("Handling %+v\n", r)
	config := LoadConfiguration()

	var all Pets
	host, err := os.Hostname()
	if err != nil {
		host = "Unknown"
	}
	path := Path{"pets", host}
	all.Hostnames = []Path{path}

	for i, backend := range config.Backends {
		URL := fmt.Sprintf("http://%s:%s", backend.Host, backend.Port)
		fmt.Printf("* Accessing %d\t %s\t %s\n", i, backend.Name, URL)
		pets := queryPets2(URL)
		all.Total = all.Total + pets.Total
		all.Hostnames = append(all.Hostnames, Path{backend.Name, pets.Hostname})
		fmt.Printf("* Hostnames %s\n", all.Hostname)
		for _, pet := range pets.Pets {
			pet.Type = backend.Name
			all.Pets = append(all.Pets, pet)
		}
	}

	sort.SliceStable(all.Pets, func(i, j int) bool {
		return all.Pets[i].Name < all.Pets[j].Name
	})

	js, err := json.Marshal(all)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func index_stdout() {

	config := LoadConfiguration()

	var all Pets
	host, err := os.Hostname()
	if err != nil {
		host = "Unknown"
	}
	path := Path{"pets", host}
	all.Hostnames = []Path{path}

	for i, backend := range config.Backends {
		URL := fmt.Sprintf("http://%s:%s", backend.Host, backend.Port)
		lookupService(backend.Host)
		fmt.Printf("* Accessing %d\t %s\t %s\n", i, backend.Name, URL)
		pets := queryPets2(URL)
		all.Total = all.Total + pets.Total
		all.Hostnames = append(all.Hostnames, Path{backend.Name, pets.Hostname})
		fmt.Printf("* Hostnames %s\n", all.Hostname)
		for _, pet := range pets.Pets {
			pet.Type = backend.Name
			all.Pets = append(all.Pets, pet)
		}
	}

	sort.SliceStable(all.Pets, func(i, j int) bool {
		return all.Pets[i].Name < all.Pets[j].Name
	})

	js, err := json.Marshal(all)
	if err != nil {
		panic(fmt.Errorf("Fatal marshall json: %s ", err))
	}
	fmt.Printf("%s", js)

}

//LoadConfiguration method
func LoadConfiguration() Config {
	viper.SetConfigType("json")
	viper.SetConfigName("pets_config") // name of config file (without extension)
	if envCfgFile := os.Getenv("SERVICE_CONFIG"); envCfgFile != "" {
		fmt.Printf("Load configuration from %s\n", envCfgFile)
		viper.SetConfigFile(envCfgFile)
	} else {
		viper.AddConfigPath("/etc/micropets/")  // path to look for the config file in
		viper.AddConfigPath("$HOME/.micropets") // call multiple times to add many search paths
		viper.AddConfigPath(".")                // optionally look for config in the working directory
	}
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s ", err))
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}

	fmt.Printf("Resolved Configuration\n")
	fmt.Printf("%+v\n", config)
	return config
}

func main() {

	config := LoadConfiguration()
	if config.Service.Listen {
		port := config.Service.Port
		http.HandleFunc("/", index)
		fmt.Printf("******* Starting to the Pets service on port %s\n", port)
		for i, backend := range config.Backends {
			fmt.Printf("* Managing %d\t %s\t %s:%s\n", i, backend.Name, backend.Host, backend.Port)
		}
		log.Fatal(http.ListenAndServe(port, nil))
	} else {
		fmt.Printf("******* Execute Pets service and exit \n")
		index_stdout()

	}
}
