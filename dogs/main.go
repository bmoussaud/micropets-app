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

//Dog type
type Dog struct {
	Name string
	Kind string
	Age  int
	URL  string
}

//Dogs type
type Dogs struct {
	Total    int
	Hostname string
	Dogs     []Dog `json:"Pets"`
}

var mode = "ALL"

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func index(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	fmt.Printf("Handling %+v\n", r)
	pet1 := Dog{"Medor", "BullDog", 18, "https://www.petmd.com/sites/default/files/10New_Bulldog_0.jpeg"}
  pet2 := Dog{"BBil", "Bull Terrier", 12, "https://www.petmd.com/sites/default/files/07New_Collie.jpeg"}
	pet3 := Dog{"Rantaplan", "Labrador Retriever", 24, "https://www.petmd.com/sites/default/files/01New_GoldenRetriever.jpeg"}
	pet4 := Dog{"Lassie", "Golden Retriever", 20, "https://www.petmd.com/sites/default/files/11New_MixedBreed.jpeg"}
	pet5 := Dog{"Beethoven", "Great St Bernard", 30, "https://upload.wikimedia.org/wikipedia/commons/6/64/Hummel_Vedor_vd_Robandahoeve.jpg"}
	pets := Dogs{5, "UKN", []Dog{pet1, pet2, pet3, pet4, pet5}}

	if mode == "RANDOM_NUMBER" {
		total := rand.Intn(pets.Total) + 1
		fmt.Printf("total %d\n", total)
		for i := 1; i < total; i++ {
			pets.Dogs = pets.Dogs[:len(pets.Dogs)-1]
			pets.Total = pets.Total - 1
		}
	}

	host, err := os.Hostname()

	if err != nil {
		pets.Hostname = "Unknown"
	} else {
		pets.Hostname = host
	}

	js, err := json.Marshal(pets)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

func main() {
	var port = ":7003"

	configLocation := GetLocation("config.properties")
	fmt.Printf("******* %s\n", configLocation)
	properties, err := properties.LoadFile(configLocation, properties.UTF8)

	if err != nil {
		fmt.Printf("config file not found, use default values\n")
	} else {
		//var readPort string
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
	}

	http.HandleFunc("/", index)
	fmt.Printf("******* Starting to the Dogs service on port %s, mode %s\n", port, mode)
	log.Fatal(http.ListenAndServe(port, nil))
}
