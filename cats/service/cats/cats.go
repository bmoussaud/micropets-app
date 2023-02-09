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
	"regexp"
	"strconv"
	"time"

	otrext "github.com/opentracing/opentracing-go/ext"
	otrlog "github.com/opentracing/opentracing-go/log"

	. "moussaud.org/cats/internal"
)

// Cat Struct
type Cat struct {
	Index int
	Name  string
	Kind  string
	Age   int
	URL   string
	From  string
	URI   string
}

// Cats Struct
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

func db() Cats {
	cat1 := Cat{20, "Orphée", "Persan", 12, "https://cdn.pixabay.com/photo/2020/02/29/13/51/cat-4890133_960_720.jpg", GlobalConfig.Service.From, "/cats/v1/data/0"}
	cat2 := Cat{21, "Pirouette", "Bengal", 1, "https://upload.wikimedia.org/wikipedia/commons/thumb/b/ba/Paintedcats_Red_Star_standing.jpg/934px-Paintedcats_Red_Star_standing.jpg", GlobalConfig.Service.From, "/cats/v1/data/1"}
	cat3 := Cat{22, "Pamina", "Angora", 120, "https://upload.wikimedia.org/wikipedia/commons/thumb/a/a5/Turkish_Angora_Odd-Eyed.jpg/440px-Turkish_Angora_Odd-Eyed.jpg", GlobalConfig.Service.From, "/cats/v1/data/2"}
	cat4 := Cat{23, "Tommy Lee", "Siamois", 120, "https://www.woopets.fr/assets/races/home/siamois-124x153.jpg", GlobalConfig.Service.From, "/cats/v1/data/3"}
	cats := Cats{4, "Unknown", []Cat{cat1, cat2, cat3, cat4}}
	host, err := os.Hostname()

	if err != nil {
		cats.Hostname = "Unknown"
	} else {
		cats.Hostname = host
	}

	return cats
}

func db_authentication(r *http.Request) {
	span := NewServerSpan(r, "db_authentication")
	defer span.Finish()

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
}

func single(w http.ResponseWriter, r *http.Request) {

	span := NewServerSpan(r, "single")
	defer span.Finish()

	setupResponse(&w, r)
	time.Sleep(time.Duration(10) * time.Millisecond)

	db_authentication(r)
	cats := db()

	re := regexp.MustCompile(`/`)
	submatchall := re.Split(r.URL.Path, -1)
	id, _ := strconv.Atoi(submatchall[4])

	if id >= len(cats.Cats) {
		http.Error(w, fmt.Sprintf("invalid index %d", id), http.StatusInternalServerError)
	} else {
		element := cats.Cats[id]
		element.From = GlobalConfig.Service.From
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

func index(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("%s", r.Method)
	span := NewServerSpan(r, "index")
	defer span.Finish()

	setupResponse(&w, r)
	time.Sleep(time.Duration(10) * time.Millisecond)

	db_authentication(r)

	cats := db()

	for i := 1; i < cats.Total; i++ {
		cats.Cats[i].From = GlobalConfig.Service.From
	}

	calls = calls + 1
	if GlobalConfig.Service.Mode == "RANDOM_NUMBER" {
		total := rand.Intn(cats.Total) + 1
		//fmt.Printf("reduce results to total %d/%d\n", total, cats.Total)
		for i := 1; i < total; i++ {
			cats.Cats = cats.Cats[:len(cats.Cats)-1]
			cats.Total = cats.Total - 1
		}
	}

	if GlobalConfig.Service.FrequencyError > 0 && calls%GlobalConfig.Service.FrequencyError == 0 {
		fmt.Printf("Fails this call (%d)", calls)
		otrext.Error.Set(span, true)
		span.LogFields(otrlog.String("error.kind", "Unexpected Error when querying the cats repository"))
		http.Error(w, "Unexpected Error when querying the cats repository", http.StatusServiceUnavailable)
	} else {
		w.Header().Set("Content-Type", "application/json")
		js, err := json.Marshal(cats)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(js)
	}
}

// GetLocation returns the full path of the config file based on the current executable location or using SERVICE_CONFIG_DIR env
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

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func Start() {
	config := LoadConfiguration()

	http.HandleFunc("/cats/v1/data", index)
	http.HandleFunc("/cats/v1/data/", single)

	http.HandleFunc("/cats/liveness", readiness_and_liveness)
	http.HandleFunc("/cats/readiness", readiness_and_liveness)

	http.HandleFunc("/liveness", readiness_and_liveness)
	http.HandleFunc("/readiness", readiness_and_liveness)

	http.HandleFunc("/", index)

	//http.HandleFunc("/", fallback)

	rand.Seed(time.Now().UnixNano())
	shift = rand.Intn(100)

	fmt.Printf("******* Starting to the cats service on port %s, mode %s\n", config.Service.Port, config.Service.Mode)
	fmt.Printf("******* Delay Period %d Amplitude %f shift %d \n", config.Service.Delay.Period, config.Service.Delay.Amplitude, shift)
	fmt.Printf("******* Frequency Error %d\n", config.Service.FrequencyError)

	log.Fatal(http.ListenAndServe(config.Service.Port, logRequest(http.DefaultServeMux)))
}
