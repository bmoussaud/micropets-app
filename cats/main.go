package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
	Total int
	Cats  []Cat `json:"Pets"`
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func index(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	fmt.Printf("Handling %+v\n", r)
	cat1 := Cat{"Orphee", "Persan", 12, "https://www.pets4homes.co.uk/images/breeds/21/db349a9afb9b6973fa3b40f684a37bb9.jpg"}
	cat2 := Cat{"Pirouette", "Bengal", 1, "https://upload.wikimedia.org/wikipedia/commons/thumb/b/ba/Paintedcats_Red_Star_standing.jpg/934px-Paintedcats_Red_Star_standing.jpg"}
	cat3 := Cat{"Pamina", "Angora", 120, "https://upload.wikimedia.org/wikipedia/commons/thumb/a/a5/Turkish_Angora_Odd-Eyed.jpg/440px-Turkish_Angora_Odd-Eyed.jpg"}
	cat4 := Cat{"Clochette", "Siamois", 120, "https://www.woopets.fr/assets/races/000/380/mobile/siamois.jpg"}
	cats := Cats{3, []Cat{cat1, cat2, cat3, cat4}}

	js, err := json.Marshal(cats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

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
			fmt.Printf("config file fount but it contains unreplaced values %s\n", readPort)
		} else {
			port = readPort
		}

	}

	http.HandleFunc("/", index)
	fmt.Printf("******* Starting to the Cats service on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
