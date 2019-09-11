package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

//Config Structure
type Config struct {
	ListenPort string `json:"port"`
	Backends   []struct {
		PetType string `json:"type"`
		URL     string `json:"url"`
	} `json:"backends"`
}

//EnvConfig Structure
type EnvConfig struct {
	//export PET_LISTENPORT=9999
	ListenPort string
	//export PETS_BACKENDS=cat_http://localhost:9000,dogs_http://localhost:9002
	Backends []string
}

//Pet Structure
type Pet struct {
	Name string
	Type string
	Kind string
	Age  int
	URL  string
}

//Pets Structure
type Pets struct {
	Total int
	Pets  []Pet `json:"Pets"`
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

var configLocation string = "config.json"

func queryPets(backend string) Pets {
	var pets Pets
	var emptyPets Pets

	fmt.Printf("Connecting backend %s\n", backend)
	resp, herr := http.Get(backend)

	if herr != nil {
		fmt.Sprintf("Error communicating with backend: %v", herr)
		return emptyPets
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Sprintf("Backend seems unhealthy: %v", resp)
		return emptyPets
	}

	body, ierr := ioutil.ReadAll(resp.Body)

	if ierr != nil {
		fmt.Sprintf("Backend producing garbage: %v", ierr)
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
	config := LoadConfiguration(configLocation)

	var all Pets

	for i, backend := range config.Backends {
		fmt.Printf("* Accessing %d\t %s\t %s\n", i, backend.PetType, backend.URL)
		pets := queryPets(backend.URL)
		all.Total = all.Total + pets.Total
		for _, pet := range pets.Pets {
			pet.Type = backend.PetType
			all.Pets = append(all.Pets, pet)
		}
	}

	js, err := json.Marshal(all)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

//GetLocation returns the fullpath
func GetLocation(file string) string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return filepath.Join(exPath, file)
}

//LoadConfiguration method
func LoadConfiguration(file string) Config {
	fmt.Printf("LoadConfiguration from File %s\n", file)
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fullPath := filepath.Join(exPath, file)

	var config Config
	configFile, err := os.Open(fullPath)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	//Overide using env variable
	fmt.Printf("LoadConfiguration from Env prefix PETS\n")
	var envConfig EnvConfig
	err = envconfig.Process("pets", &envConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
	if len(envConfig.ListenPort) > 0 && envConfig.ListenPort != "" {
		fmt.Printf("* override Listen Port with %s\n", envConfig.ListenPort)
		config.ListenPort = envConfig.ListenPort
	}

	for _, b := range envConfig.Backends {
		if len(b) > 0 && b != "" {
			fmt.Println("* override Backends:")
			result := strings.Split(b, "_")
			pet := result[0]
			url := result[1]
			//fmt.Printf(" from env (%s) (%s)\n", pet, url)
			for i, service := range config.Backends {
				//fmt.Printf("  from config file (%s) (%s)\n", service.PetType, service.URL)
				if strings.EqualFold(service.PetType, pet) {
					fmt.Printf("* %s  %s <- %s \n", pet, service.URL, url)
					config.Backends[i].URL = url
					break
				}
			}
		}
	}

	fmt.Printf("Resolved Configuration\n")
	fmt.Printf("%+v\n", config)
	return config
}

func main() {

	config := LoadConfiguration(configLocation)
	port := config.ListenPort
	http.HandleFunc("/", index)
	fmt.Printf("******* Starting to the Pets service on port %s\n", port)
	for i, backend := range config.Backends {
		fmt.Printf("* Managing %d\t %s\t %s\n", i, backend.PetType, backend.URL)
	}
	log.Fatal(http.ListenAndServe(port, nil))
}
