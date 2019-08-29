package main

import (
	"encoding/json"
	"fmt"
	"github.com/kardianos/service"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Config struct {
	ListenPort string `json:"port"`
	Backends   []struct {
		Url     string `json:"url"`
		PetType string `json:"type"`
	} `json:"backends"`
}

type Pet struct {
	Name string
	Type string
	Kind string
	Age  int
	Url  string
}

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
		fmt.Printf("* Accessing %d\t %s\t %s\n", i, backend.PetType, backend.Url)
		pets := queryPets(backend.Url)
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

var logger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func GetLocation(file string) string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return filepath.Join(exPath, file)
}

func LoadConfiguration(file string) Config {
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
	return config
}

func (p *program) run() {

	fmt.Printf("******* %s\n", configLocation)
	config := LoadConfiguration(configLocation)
	port := config.ListenPort
	http.HandleFunc("/", index)
	fmt.Printf("******* Starting to the Pets service on port %s\n", port)
	for i, backend := range config.Backends {
		fmt.Printf("* Managing %d\t %s\t %s\n", i, backend.PetType, backend.Url)
	}
	log.Fatal(http.ListenAndServe(port, nil))
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	fmt.Printf("******* Pet Service 1.0.4 \n")
	svcConfig := &service.Config{
		Name:        "PetService",
		DisplayName: "Core Pet Service",
		Description: "The core cat service",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
