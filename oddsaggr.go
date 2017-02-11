package main

import (
	"flag"
	"net/http"
	"log"
	"fmt"
	"github.com/pitshifer/oddsaggr/jsonodds"
	"github.com/pitshifer/oddsaggr/entity"
)

var client interface{
	GetSports() (entity.Sports, error)
	GetOddTypes() (entity.OddTypes, error)
}

const DEFAULT_PORT = "5050"

func main() {
	var port string

	flag.StringVar(&port, "-p", DEFAULT_PORT, "port")
	flag.Parse()

	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(resp, "OddsAggregator")
	})

	http.HandleFunc("/sports/", showSports)
	http.HandleFunc("/oddtypes/", showOddTypes)

	client = jsonodds.New(jsonodds.Config{
		Url: 	"https://jsonodds.com/api/",
		Key:	"5f0ceb75-edfe-11e6-a75b-0667553df103",
	})

	log.Println("Server started on port: " + port)
	log.Fatal(http.ListenAndServe("localhost:" + port, nil))
}

func showSports(resp http.ResponseWriter, req *http.Request) {
	sports, err := client.GetSports()
	if err != nil {
		fmt.Fprintln(resp, err)
		return
	}

	fmt.Fprintf(resp, "%s", sports)
}

func showOddTypes(resp http.ResponseWriter, req *http.Request) {
	ot, err := client.GetOddTypes()
	if err != nil {
		fmt.Fprintln(resp, err)
		return
	}

	fmt.Fprintf(resp, "%s", ot)
}
