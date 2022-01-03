package dogs

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	. "moussaud.org/dogs/internal"
)

//Dog type
type Dog struct {
	Name string
	Kind string
	Age  int
	URL  string
	From string
}

//Dogs type
type Dogs struct {
	Total    int
	Hostname string
	Dogs     []Dog `json:"Pets"`
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
		time.Sleep(time.Duration(2) * time.Millisecond)
	}
}

func db() Dogs {
	pet1 := Dog{"Medor", "BullDog", 18, "https://www.petmd.com/sites/default/files/10New_Bulldog_0.jpeg", GlobalConfig.Service.From}
	pet2 := Dog{"Bil", "Bull Terrier", 12, "https://www.petmd.com/sites/default/files/07New_Collie.jpeg", GlobalConfig.Service.From}
	pet3 := Dog{"Rantaplan", "Labrador Retriever", 24, "https://www.petmd.com/sites/default/files/01New_GoldenRetriever.jpeg", GlobalConfig.Service.From}
	pet4 := Dog{"Lassie", "Golden Retriever", 20, "https://www.petmd.com/sites/default/files/11New_MixedBreed.jpeg", GlobalConfig.Service.From}
	pet5 := Dog{"Beethoven", "Great St Bernard", 30, "https://upload.wikimedia.org/wikipedia/commons/6/64/Hummel_Vedor_vd_Robandahoeve.jpg", GlobalConfig.Service.From}
	pets := Dogs{5, "UKN", []Dog{pet1, pet2, pet3, pet4, pet5}}

	host, err := os.Hostname()

	if err != nil {
		pets.Hostname = "Unknown"
	} else {
		pets.Hostname = host
	}

	return pets
}

func index(w http.ResponseWriter, r *http.Request) {

	span := NewServerSpan(r, "index")
	defer span.Finish()

	setupResponse(&w, r)

	pets := db()

	time.Sleep(time.Duration(len(pets.Dogs)) * time.Millisecond)
	db_authentication(r)

	calls = calls + 1
	if GlobalConfig.Service.Mode == "RANDOM_NUMBER" {
		total := rand.Intn(pets.Total) + 1
		fmt.Printf("reduce results to total %d/%d\n", total, pets.Total)
		for i := 1; i < total; i++ {
			pets.Dogs = pets.Dogs[:len(pets.Dogs)-1]
			pets.Total = pets.Total - 1
		}
	}

	js, err := json.Marshal(pets)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if GlobalConfig.Service.Delay.Period > 0 {
		y := float64(calls+shift) * math.Pi / float64(2*GlobalConfig.Service.Delay.Period)
		sin_y := math.Sin(y)
		abs_y := math.Abs(sin_y)
		sleep := int(abs_y * GlobalConfig.Service.Delay.Amplitude * 1000.0)
		fmt.Printf("waitTime %d - %f - %f - %f  -> sleep %d ms  \n", calls, y, math.Sin(y), abs_y, sleep)
		start := time.Now()
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		elapsed := time.Since(start)
		fmt.Printf("Current Unix Time: %s\n", elapsed)
	}

	if GlobalConfig.Service.FrequencyError > 0 && calls%GlobalConfig.Service.FrequencyError == 0 {
		fmt.Printf("Fails this call (%d)", calls)
		http.Error(w, "Unexpected Error when querying the dogs repository", http.StatusServiceUnavailable)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func single(w http.ResponseWriter, r *http.Request) {

	span := NewServerSpan(r, "single")
	defer span.Finish()

	setupResponse(&w, r)
	time.Sleep(time.Duration(10) * time.Millisecond)

	db_authentication(r)
	dogs := db()

	re := regexp.MustCompile(`/`)
	submatchall := re.Split(r.URL.Path, -1)
	id, _ := strconv.Atoi(submatchall[4])

	if id >= len(dogs.Dogs) {
		http.Error(w, fmt.Sprintf("invalid index %d", id), http.StatusInternalServerError)
	} else {
		element := dogs.Dogs[id]
		fmt.Println(element)
		w.Header().Set("Content-Type", "application/json")
		js, err := json.Marshal(element)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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

	http.HandleFunc("/dogs/v1/data", index)
	http.HandleFunc("/dogs/v1/data/", single)

	http.HandleFunc("/dogs/liveness", readiness_and_liveness)
	http.HandleFunc("/dogs/readiness", readiness_and_liveness)

	http.HandleFunc("/liveness", readiness_and_liveness)
	http.HandleFunc("/readiness", readiness_and_liveness)

	//http.HandleFunc("/", fallback)

	rand.Seed(time.Now().UnixNano())
	shift = rand.Intn(100)

	fmt.Printf("******* Starting to the dogs service on port %s, mode %s\n", config.Service.Port, config.Service.Mode)
	fmt.Printf("******* Delay Period %d Amplitude %f shift %d \n", config.Service.Delay.Period, config.Service.Delay.Amplitude, shift)
	fmt.Printf("******* Frequency Error %d\n", config.Service.FrequencyError)

	log.Fatal(http.ListenAndServe(config.Service.Port, nil))
}
