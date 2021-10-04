package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/magiconair/properties"
)

//Cat Struct
type Cat struct {
	Name string
	Kind string
	Age  int
	URL  string
}

//Cats Struct
type Cats struct {
	Total    int
	Hostname string
	Cats     []Cat `json:"Pets"`
}

var mode = "ALL"

var configLocation string = "config.properties"

var delayPeriod = 0.0

var delayAmplitude = 0.0

var calls = 0.0

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func index(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	//fmt.Printf("Handling %+v\n", r)
	//fmt.Printf("MODE %s\n", mode)
	cat1 := Cat{"Orphee", "Persan", 12, "https://cdn.pixabay.com/photo/2020/02/29/13/51/cat-4890133_960_720.jpg"}
	cat2 := Cat{"Pirouette", "Bengal", 1, "https://upload.wikimedia.org/wikipedia/commons/thumb/b/ba/Paintedcats_Red_Star_standing.jpg/934px-Paintedcats_Red_Star_standing.jpg"}
	cat3 := Cat{"Pamina", "Angora", 120, "https://upload.wikimedia.org/wikipedia/commons/thumb/a/a5/Turkish_Angora_Odd-Eyed.jpg/440px-Turkish_Angora_Odd-Eyed.jpg"}
	cat4 := Cat{"Clochette", "Siamois", 120, "https://www.woopets.fr/assets/races/home/siamois-124x153.jpg"}
	cats := Cats{4, "Unknown", []Cat{cat1, cat2, cat3, cat4}}

	calls = calls + 1
	if mode == "RANDOM_NUMBER" {
		total := rand.Intn(cats.Total) + 1
		fmt.Printf("reduce results to total %d/%d\n", total, cats.Total)
		for i := 1; i < total; i++ {
			cats.Cats = cats.Cats[:len(cats.Cats)-1]
			cats.Total = cats.Total - 1
		}
	}

	host, err := os.Hostname()

	if err != nil {
		cats.Hostname = "Unknown"
	} else {
		cats.Hostname = host
	}

	js, err := json.Marshal(cats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if delayPeriod > 0 {
		y := (math.Pi) / (2 * delayPeriod) * float64(calls)
		sin_y := math.Sin(y)
		abs_y := math.Abs(sin_y)
		sleep := int(abs_y * delayAmplitude * 1000.0)
		//fmt.Printf("waitTime %f - %f - %f - %f  -> sleep %d seconds  \n", calls, y, math.Sin(y), abs_y, sleep)
		start := time.Now()
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		elapsed := time.Since(start)
		fmt.Printf("Current Unix Time: %s\n", elapsed)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

//GetLocation returns the full path of the config file based on the current executable location or using SERVICE_CONFIG_DIR env
func GetLocation(file string) string {
	if serviceConfigDirectory := os.Getenv("SERVICE_CONFIG_DIR"); serviceConfigDirectory != "" {
		fmt.Printf("Load configuration from %s\n", serviceConfigDirectory)
		return filepath.Join(serviceConfigDirectory, file)
	} else {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)
		return filepath.Join(exPath, file)
	}
}

func readiness_and_liveness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}

func main() {
	var port = ":7002"

	configLocation := GetLocation("config.properties")
	fmt.Printf("******* %s\n", configLocation)
	properties, err := properties.LoadFile(configLocation, properties.UTF8)

	if err != nil {
		fmt.Printf("config file not found, use default values\n")
	} else {
		var readPort string
		readPort = properties.GetString("listen.port", port)
		//fmt.Printf(readPort)
		if strings.HasPrefix(readPort, ":{{") {
			fmt.Printf("config file found but it contains unreplaced values %s\n", readPort)
		} else {
			port = readPort
		}

		var readMode string
		readMode = properties.GetString("mode", mode)
		if strings.HasPrefix(readPort, ":{{") {
			fmt.Printf("config file found but it contains unreplaced values %s\n", readMode)
		} else {
			mode = readMode
		}

		delayPeriod = properties.GetFloat64("delay.period", delayPeriod)
		delayAmplitude = properties.GetFloat64("delay.amplitude", delayAmplitude)
	}

	http.HandleFunc("/cats/v1/data", index)
	http.HandleFunc("/liveness", readiness_and_liveness)
	http.HandleFunc("/readiness", readiness_and_liveness)
	fmt.Printf("******* Starting to the Cats service on port %s, mode %s\n", port, mode)
	fmt.Printf("******* Delay Period %f Amplitude %f\n", delayPeriod, delayAmplitude)
	log.Fatal(http.ListenAndServe(port, nil))
}
