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

type appHandler func(http.ResponseWriter, *http.Request) *appError

type appError struct {
	message string
	code    int
}

const localPort = ":3030"

func main() {
	local := flag.Bool("local", false, "specify during development, the endpoint will listen on localhost"+localPort)
	flag.Parse()

	listener := gateway.ListenAndServe
	if *local {
		fmt.Println("Listening on localhost" + localPort)
		listener = http.ListenAndServe
	}

	http.Handle("/", appHandler(solutionsJSON))

	log.Fatal(listener(localPort, nil))
}

func solutionsJSON(w http.ResponseWriter, r *http.Request) *appError {
	users, ok := r.URL.Query()["user"]
	if !ok || len(users) < 1 {
		return &appError{"parameter 'user' is missing", http.StatusBadRequest}
	}

	user := users[0]

	if err := validateUser(user); err != nil {
		return &appError{err.Error(), http.StatusBadRequest}
	}

	res, err := http.Get("https://exercism.io/profiles/" + user)
	if err != nil {
		return &appError{err.Error(), http.StatusInternalServerError}
	}

	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return &appError{err.Error(), http.StatusInternalServerError}
	}

	re := regexp.MustCompile(`Showing (\d+) solutions`)
	matches := re.FindSubmatch(robots)
	if len(matches) < 2 {
		return &appError{
			message: fmt.Sprintf("no match found for regular expression: '%s' - please check if this profile page exists: https://exercism.io/profiles/%s", re.String(), user),
			code:    http.StatusInternalServerError,
		}
	}

	total, err := strconv.Atoi(string(matches[1]))
	if err != nil {
		return &appError{err.Error(), http.StatusInternalServerError}
	}

	s := solutions{total}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
	return nil
}

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		log.Printf("Message: %s - Code: %d\n", err.message, err.code)
		w.WriteHeader(err.code)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse{err.message})
	}
}

func validateUser(user string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
	match := re.FindString(user)
	if len(match) == 0 {
		return errors.New("parameter 'user' should contain only letters, numbers and hyphens")
	}
	return nil
}
