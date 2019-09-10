package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kardianos/service"
	"github.com/magiconair/properties"
)

var count = 1

//Dog type
type Dog struct {
	Name string
	Kind string
	Age  int
	URL  string
}

//Dogs type
type Dogs struct {
	Total int
	Dogs  []Dog `json:"Pets"`
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func index(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	fmt.Printf("Handling %+v\n", r)
	pet1 := Dog{"Medor", "BullDog", 18, "https://www.petmd.com/sites/default/files/10New_Bulldog_0.jpeg"}
	pet2 := Dog{"Bil", "Bull Terrier", 12, "https://www.petmd.com/sites/default/files/07New_Collie.jpeg"}
	pet3 := Dog{"Rantaplan", "Labrador Retriever", 24, "https://www.petmd.com/sites/default/files/01New_GoldenRetriever.jpeg"}
	pet4 := Dog{"Lassie", "Golden Retriever", 20, "https://www.petmd.com/sites/default/files/11New_MixedBreed.jpeg"}
	pets := Dogs{4, []Dog{pet1, pet2, pet3, pet4}}

	js, err := json.Marshal(pets)
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

//GetLocation returns the full path of the config file based on the current executable location
func GetLocation(file string) string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return filepath.Join(exPath, file)
}

func (p *program) run() {

	configLocation := GetLocation("config.properties")
	fmt.Printf("******* %s\n", configLocation)
	properties, err := properties.LoadFile(configLocation, properties.UTF8)
	var port = ":7003"
	if err != nil {
		fmt.Printf("config file not found, use default values\n")
	} else {
		var readPort string
		readPort = properties.GetString("listen.port", port)
		//fmt.Printf(readPort)
		if strings.HasPrefix(readPort, ":{{") {
			fmt.Printf("config file fount but it contains unreplaced values %s\n", readPort)
		} else {
			port = readPort
		}

	}

	http.HandleFunc("/", index)
	fmt.Printf("******* Starting to the Dog service on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	fmt.Printf("******* Dog Service 1.0.4 \n")
	svcConfig := &service.Config{
		Name:        "DogService",
		DisplayName: "Core Dog Service",
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
