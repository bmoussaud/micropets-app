package main

import (
	"fmt"
	"github.com/kardianos/service"
	"github.com/magiconair/properties"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var count = 1

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Handling %+v\n", r)

	host, err := os.Hostname()

	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving hostname: %v", err), 500)
		return
	}

	msg := fmt.Sprintf("**** Host: %s\n* Successful Requests: %d\n* %s \n", host, count, GetLocation("config.properties"))
	count += 1

	io.WriteString(w, msg)
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
	port := properties.GetString("listen.port", "8001")
	http.HandleFunc("/", index)
	fmt.Printf("******* Starting to service on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
	// Do work here
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	fmt.Printf("******* Backend Service 1.0.4 \n")
	svcConfig := &service.Config{
		Name:        "BackendService",
		DisplayName: "Core Backend Service",
		Description: "The core backend service",
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
