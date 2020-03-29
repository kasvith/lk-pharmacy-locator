package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type DataCache struct {
	Data      DistrictData
	Districts []string
	sync.RWMutex
}

var cachedData DataCache

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func CacheDataStore() error {
	data, err := FetchData()
	if err != nil {
		return err
	}

	cachedData.Lock()
	cachedData.Data = data
	cachedData.Unlock()

	return nil
}

func main() {
	host := getEnv("HOST", "127.0.0.1")
	port := getEnv("PORT", "8080")

	// lets pre fetch data
	log.Println("Fetching data from upstream")
	err := CacheDataStore()
	if err != nil {
		log.Fatalf("error fetching data from upstream %s", err)
	}
	log.Println("Data fetched")

	// start a go routine which will fetch latest data for 1hr
	log.Println("Starting background fetcher")
	ticker := time.NewTicker(time.Hour * 1)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("fetching data from upstream")
				err := CacheDataStore()
				if err != nil {
					log.Printf("error: fetching data failed, cache intacted %s", err)
				}
			}
		}
	}()

	r := mux.NewRouter()
	r.HandleFunc("/api/districts", DistrictHandler).Methods("GET")
	r.HandleFunc("/api/districts/{district}/areas", AreasHandler).Methods("GET")
	r.HandleFunc("/api/districts/{district}/areas/{area}/pharmacies", PharmacyHandler).Methods("GET")

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%s", host, port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Running on %s\n", fmt.Sprintf("%s:%s", host, port))
	log.Fatal(srv.ListenAndServe())
}

func PharmacyHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	district := strings.ToLower(vars["district"])
	area := strings.ToLower(vars["area"])

	var results []*PharmacyEntry
	if d, has := cachedData.Data[district]; has {
		if a, found := d[area]; found {
			results = a
		}
	}

	writer.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(PharmacyResponse{Pharmacies: results})
}

func AreasHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	district := strings.ToLower(vars["district"])

	var results []string
	if areas, has := cachedData.Data[district]; has {
		results = make([]string, len(areas))
		i := 0
		for a := range areas {
			results[i] = a
			i++
		}
	}
	sort.Strings(results)
	writer.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(AreasResponse{Areas: results})
}

func DistrictHandler(writer http.ResponseWriter, request *http.Request) {
	districts := make([]string, len(cachedData.Data))
	i := 0
	for k := range cachedData.Data {
		districts[i] = k
		i++
	}

	sort.Strings(districts)
	writer.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(DistrictsResponse{Districts: districts})
}
