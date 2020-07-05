package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/apex/gateway"
)

type solutions struct {
	Total int `json:"total"`
}

type errorResponse struct {
	Error string `json:"error"`
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
	users, ok := r.URL.Query()["user"]
	if !ok || len(users) < 1 {
		replyError(w, errors.New("parameter 'user' is missing"))
		return
	}

	user := users[0]
	res, err := http.Get("https://exercism.io/profiles/" + user)
	if err != nil {
		replyError(w, err)
		return
	}

	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		replyError(w, err)
		return
	}

	re := regexp.MustCompile(`Showing (\d+) solutions`)
	matches := re.FindSubmatch(robots)
	if len(matches) < 2 {
		replyError(w, fmt.Errorf("no match found for regular expression: %s", re.String()))
		return
	}

	total, err := strconv.Atoi(string(matches[1]))
	if err != nil {

	}
	s := solutions{total}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func replyError(w http.ResponseWriter, err error) {
	log.Print(err)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(errorResponse{err.Error()})
}
