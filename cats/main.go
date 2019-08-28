package main

import (
	"encoding/json"
	"fmt"
	"github.com/kardianos/service"
	"github.com/magiconair/properties"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var count = 1

type Cat struct {
	Name string
	Kind string
	Age  int
	Url  string
}

type Cats struct {
	Total int
	Cats  []Cat
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Handling %+v\n", r)
	cat1 := Cat{"Orphee", "Persan", 12, "https://www.pets4homes.co.uk/images/breeds/21/db349a9afb9b6973fa3b40f684a37bb9.jpg"}
	cat2 := Cat{"Pirouette", "Bengal", 1, "https://upload.wikimedia.org/wikipedia/commons/thumb/b/ba/Paintedcats_Red_Star_standing.jpg/934px-Paintedcats_Red_Star_standing.jpg"}
	cat3 := Cat{"Pamina", "Angora", 120, "https://upload.wikimedia.org/wikipedia/commons/thumb/a/a5/Turkish_Angora_Odd-Eyed.jpg/440px-Turkish_Angora_Odd-Eyed.jpg"}
	cats := Cats{3, []Cat{cat1, cat2, cat3}}

	js, err := json.Marshal(cats)
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

func (p *program) run() {

	configLocation := GetLocation("config.properties")
	fmt.Printf("******* %s\n", configLocation)
	properties, err := properties.LoadFile(configLocation, properties.UTF8)
	var port = ":7000"
	if err != nil {
		fmt.Printf("config file not found, use default values\n")
	} else {
		port = properties.GetString("listen.port", port)
	}

	http.HandleFunc("/", index)
	fmt.Printf("******* Starting to the Cat service on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	fmt.Printf("******* Cat Service 1.0.4 \n")
	svcConfig := &service.Config{
		Name:        "CatService",
		DisplayName: "Core Cat Service",
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
