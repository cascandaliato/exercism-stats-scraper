package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/apex/gateway"
)

type solutions struct {
	Total string `json:"total"`
}

func main() {
	port := flag.Int("port", -1, "specify a port to use http rather than AWS Lambda")
	flag.Parse()

	listener := gateway.ListenAndServe
	portStr := ":3030"
	if *port != -1 {
		listener = http.ListenAndServe
		portStr = fmt.Sprintf(":%d", *port)
	}

	http.HandleFunc("/", solutionsJSON)

	log.Fatal(listener(portStr, nil))
}

func solutionsJSON(w http.ResponseWriter, r *http.Request) {
	os.Exit(1)
	res, err := http.Get("https://exercism.io/profiles/casca")
	if err != nil {
		log.Panic(err)
	}

	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Panic(err)
	}

	re := regexp.MustCompile(`Showing (\d+) solutions`)
	matches := re.FindSubmatch(robots)

	w.Header().Set("Content-Type", "application/json")
	s := solutions{string(matches[1])}
	json.NewEncoder(w).Encode(s)
}
