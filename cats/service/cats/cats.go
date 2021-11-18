package cats

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	. "moussaud.org/cats/internal"
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

var calls = 0

var shift = 0

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func db_authentication(r *http.Request) {
	span := NewServerSpan(r, "db_authentication")
	defer span.Finish()
	
	if GlobalConfig.Service.Delay.Period > 0 {
		y := float64(calls+shift ) * math.Pi / float64(2*GlobalConfig.Service.Delay.Period)
		sin_y := math.Sin(y)
		abs_y := math.Abs(sin_y)
		sleep := int(abs_y * GlobalConfig.Service.Delay.Amplitude * 1000.0)
		fmt.Printf("waitTime %d - %f - %f - %f  -> sleep %d ms  \n", calls, y, math.Sin(y), abs_y, sleep)
		start := time.Now()
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		elapsed := time.Since(start)
		fmt.Printf("Current Unix Time: %s\n", elapsed)
	}
}

func index(w http.ResponseWriter, r *http.Request) {

	span := NewServerSpan(r, "index")
	defer span.Finish()

	setupResponse(&w, r)
	time.Sleep(time.Duration(10) * time.Millisecond)

	
	db_authentication(r)

	//fmt.Printf("Handling %+v\n", r)
	//fmt.Printf("MODE %s\n", mode)
	cat1 := Cat{"Orphee", "Persan", 12, "https://cdn.pixabay.com/photo/2020/02/29/13/51/cat-4890133_960_720.jpg"}
	cat2 := Cat{"Pirouette", "Bengal", 1, "https://upload.wikimedia.org/wikipedia/commons/thumb/b/ba/Paintedcats_Red_Star_standing.jpg/934px-Paintedcats_Red_Star_standing.jpg"}
	cat3 := Cat{"Pamina", "Angora", 120, "https://upload.wikimedia.org/wikipedia/commons/thumb/a/a5/Turkish_Angora_Odd-Eyed.jpg/440px-Turkish_Angora_Odd-Eyed.jpg"}
	cat4 := Cat{"Clochette", "Siamois", 120, "https://www.woopets.fr/assets/races/home/siamois-124x153.jpg"}
	cats := Cats{4, "Unknown", []Cat{cat1, cat2, cat3, cat4}}

	calls = calls + 1
	if GlobalConfig.Service.Mode == "RANDOM_NUMBER" {
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

	

	if GlobalConfig.Service.FrequencyError > 0 && calls%GlobalConfig.Service.FrequencyError == 0 {
		fmt.Printf("Fails this call (%d)", calls)
		http.Error(w, "Unexpected Error when querying the cats repository", http.StatusServiceUnavailable)
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

func readiness_and_liveness(w http.ResponseWriter, r *http.Request) {
	span := NewServerSpan(r, "readiness_and_liveness")
	defer span.Finish()

	w.WriteHeader(200)
	w.Write([]byte("ok"))
}

func Start() {
	config := LoadConfiguration()

	http.HandleFunc("/cats/v1/data", index)

	http.HandleFunc("/cats/liveness", readiness_and_liveness)
	http.HandleFunc("/cats/readiness", readiness_and_liveness)

	http.HandleFunc("/liveness", readiness_and_liveness)
	http.HandleFunc("/readiness", readiness_and_liveness)

	//http.HandleFunc("/", fallback)

	rand.Seed(time.Now().UnixNano())
	shift = rand.Intn(100)

	fmt.Printf("******* Starting to the cats service on port %s, mode %s\n", config.Service.Port, config.Service.Mode)
	fmt.Printf("******* Delay Period %d Amplitude %f shift %d \n", config.Service.Delay.Period, config.Service.Delay.Amplitude, shift)
	fmt.Printf("******* Frequency Error %d\n", config.Service.FrequencyError)

	

	log.Fatal(http.ListenAndServe(config.Service.Port, nil))
}
