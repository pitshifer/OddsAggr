package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/pitshifer/oddsaggr/entity"
	"github.com/pitshifer/oddsaggr/jsonodds"
	"log"
	"net/http"
	"os"
	"github.com/kr/pretty"
)

var client interface {
	GetSports() (*entity.Sports, error)
	GetOddTypes() (*entity.OddTypes, error)
	GetOddsBySport(sport, source string) (*[]entity.EventOdds, error)
}

type Config struct {
	Jsonodds JOConfig
}

type JOConfig struct {
	Key        string `toml:"api_key"`
	Url        string `toml:"url"`
	OddFormat  string `toml:"odd_format"`
}

const (
	DEFAULT_PORT        = "5050"
	DEFAULT_CONFIG_FILE = "./config.toml"
)

func main() {
	var port string
	var configFile string
	var config Config

	flag.StringVar(&port, "-p", DEFAULT_PORT, "port")
	flag.StringVar(&configFile, "-c", DEFAULT_CONFIG_FILE, "config file")
	flag.Parse()

	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Fatalln(err)
	}
	client = jsonodds.New(jsonodds.Config{
		Url:        config.Jsonodds.Url,
		Key:        config.Jsonodds.Key,
		OddFormat:  config.Jsonodds.OddFormat,
	})

	if result, err := client.GetOddsBySport("khl", "0"); err == nil{
		pretty.Println(result)
	} else {
		log.Fatal(err)
	}
	os.Exit(0)

	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(resp, "OddsAggregator")
	})

	http.HandleFunc("/sports/", showSports)
	http.HandleFunc("/events/soccer", showEvents)

	log.Println("Server started on port: " + port)
	log.Fatal(http.ListenAndServe("localhost:"+port, nil))
}

func showSports(resp http.ResponseWriter, req *http.Request) {
	sports, err := client.GetSports()
	if err != nil {
		fmt.Fprintln(resp, err)
		return
	}

	fmt.Fprintf(resp, "%s", sports)
}

func showEvents(resp http.ResponseWriter, req *http.Request) {
	client.GetOddsBySport("soccer", "0")

}
