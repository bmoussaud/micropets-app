package main

import (
	"fmt"
	"github.com/kardianos/service"
	"github.com/magiconair/properties"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

var count = 1

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Handling %+v\n", r)
	configLocation := GetLocation("config.properties")
	properties := properties.MustLoadFile(configLocation, properties.UTF8)
	backend := properties.GetString("backend.url", "http://127.0.0.1:8002")

	fmt.Printf("Connecting backend %s\n", backend)
	resp, herr := http.Get(backend)

	if herr != nil {
		http.Error(w, fmt.Sprintf("Error communicating with backend: %v", herr), 500)
		return
	}

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Backend seems unhealthy: %v", resp), 500)
		return
	}

	body, ierr := ioutil.ReadAll(resp.Body)

	if ierr != nil {
		http.Error(w, fmt.Sprintf("Backend producing garbage: %v", ierr), 500)
		return
	}

	host, err := os.Hostname()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving hostname: %v", err), 500)
		return
	}
	msg2 := fmt.Sprintf("%s", host)

	alternateColor := properties.GetString("ALTERNATE_COLOR", "gray")

	t, terr := template.ParseFiles(GetLocation("content/index.html"))

	if terr != nil {
		http.Error(w, fmt.Sprintf("Error loading template: %v", terr), 500)
		return
	}

	config := map[string]string{
		"Message":      string(body),
		"Feature":      os.Getenv("FEATURE"),
		"FrontendHost": string(msg2),
		"Config":       alternateColor,
	}

	t.Execute(w, config)
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
	properties := properties.MustLoadFile(configLocation, properties.UTF8)
	port := properties.GetString("listen.port", "8010")
	http.HandleFunc("/", index)
	fmt.Printf("******* Starting to frontend service on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	fmt.Printf("******* Frontend Service 1.0.4 \n")
	svcConfig := &service.Config{
		Name:        "FrontendService",
		DisplayName: "Core Frontend Service",
		Description: "The core frontend service",
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
