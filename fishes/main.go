package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/magiconair/properties"
)

//Fish Struct
type Fish struct {
	Name string
	Kind string
	Age  int
	URL  string
}

//fishes Struct
type Fishes struct {
	Total    int
	Hostname string
	Fishes   []Fish `json:"Pets"`
}

var mode = "ALL"
var frequencyError = -1
var calls = 0

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func index(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)

	fmt.Printf("Handling %+v\n", r)

	host, err := os.Hostname()
	if err != nil {
		host = "Unknown"
	}

	fishes := Fishes{4,
		host,
		[]Fish{
			{"Nemo", "Poisson Clown", 14,
				"https://www.sciencesetavenir.fr/assets/img/2019/07/10/cover-r4x3w1000-5d258790dd324-f96f05d4901fc6ce0ab038a685e4d5c99f6cdfe2-jpg.jpg"},
			{"Glumpy", "Neon Tetra", 11,
				"https://www.fishkeepingworld.com/wp-content/uploads/2018/02/Neon-Tetra-New.jpg"},
			{"Dory", "Pacific regal blue tang", 12,
				"http://www.oceanlight.com/stock-photo/palette-surgeonfish-image-07922-671143.jpg"},
			{"Argo", "Combattant", 27,
				"https://www.aquaportail.com/pictures1003/anemone-clown_1267799900_poisson-combattant.jpg"}}}

	if mode == "RANDOM_NUMBER" {
		total := rand.Intn(fishes.Total) + 1
		fmt.Printf("total %d\n", total)
		for i := 1; i < total; i++ {
			fishes.Fishes = fishes.Fishes[:len(fishes.Fishes)-1]
			fishes.Total = fishes.Total - 1
		}
	}

	js, err := json.Marshal(fishes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	calls = calls + 1
	if frequencyError > 0 && calls%frequencyError == 0 {
		fmt.Printf("Fails this call (%d)", calls)
		http.Error(w, "Unexpected Error when querying the fish repository", http.StatusServiceUnavailable)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}

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
func main() {
	var port = ":7007"

	configLocation := GetLocation("config.properties")
	fmt.Printf("******* %s\n", configLocation)
	properties, err := properties.LoadFile(configLocation, properties.UTF8)

	if err != nil {
		fmt.Printf("config file not found, use default values\n")
	} else {
		readPort := properties.GetString("listen.port", port)
		//fmt.Printf(readPort)
		if strings.HasPrefix(readPort, ":{{") {
			fmt.Printf("config file fount but it contains unreplaced values %s\n", readPort)
		} else {
			port = readPort
		}

		readMode := properties.GetString("mode", mode)
		if strings.HasPrefix(readPort, ":{{") {
			fmt.Printf("config file fount but it contains unreplaced values %s\n", readMode)
		} else {
			mode = readMode
		}

		readFrequencyError := properties.GetInt("frequencyError", frequencyError)
		frequencyError = readFrequencyError
	}

	http.HandleFunc("/", index)
	fmt.Printf("******* Starting to the Fishes service on port %s, mode %s, frequency Error %d\n", port, mode, frequencyError)
	log.Fatal(http.ListenAndServe(port, nil))
}
