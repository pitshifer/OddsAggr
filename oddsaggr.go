package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/pitshifer/oddsaggr/entity"
	"github.com/pitshifer/oddsaggr/jsonodds"
	log "github.com/sirupsen/logrus"
	"net/http"
	"encoding/json"
)

var client interface {
	GetSources() (*[]entity.Source, error)
	GetSports() (*entity.Sports, error)
	GetOddTypes() (*entity.OddTypes, error)
	GetFinalType() (*[]string, error)
	GetOddsBySport(sport, source string) (*[]entity.EventOdds, error)
	GetOdds(oddType, source string) (*[]entity.EventOdds, error)
}

var config Config

type Config struct {
	Environment	string		`toml:"environment"`
	Jsonodds 	JOConfig
}

type JOConfig struct {
	Key        	string	 	`toml:"api_key"`
	Url        	string 		`toml:"url"`
	OddFormat  	string 		`toml:"odd_format"`
}

const (
	DEFAULT_PORT        = "5050"
	DEFAULT_CONFIG_FILE = "./config.toml"
)

func main() {
	var port string
	var configFile string

	flag.StringVar(&port, "-p", DEFAULT_PORT, "port")
	flag.StringVar(&configFile, "-c", DEFAULT_CONFIG_FILE, "config file")
	flag.Parse()

	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Fatalln(err)
	}

	loggerInit(config)

	client = jsonodds.New(jsonodds.Config{
		Url:        config.Jsonodds.Url,
		Key:        config.Jsonodds.Key,
		OddFormat:  config.Jsonodds.OddFormat,
	})

	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(resp, "OddsAggregator")
	})

	http.HandleFunc("/sports/", showSports)
	http.HandleFunc("/events/", showEvents)

	log.Info("Server started on port: " + port)
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
	var events []entity.Event

	sport := req.URL.Query().Get("sport")
	source := req.URL.Query().Get("source")
	if sport == "" {
		log.Debug("No specify kind of sport")
		fmt.Fprintln(resp, "Url must be like 'event/?sport=name'")
		return
	}
	if source == "" {
		source = "0"
	}

	eventOdds, err := client.GetOddsBySport(sport, source);
	if err != nil {
		log.Error(err)
		fmt.Fprintln(resp, err)
	}

	for _, eo := range *eventOdds {
		events = append(events, eo.Event)
	}

	resp.Header().Set("Content-type", "text/json")
	b, err := json.Marshal(events)
	if err != nil {
		log.Error("Error during marshal data: ", err)
		fmt.Fprint(resp, "Error during marshal data.")
	} else  {
		fmt.Fprintf(resp, "%s", b)
	}
}