package jsonodds

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/pitshifer/oddsaggr/entity"
	"log"
	"net/url"
)

type Config struct {
	Key, Url, OddsFormat	string
}

type client struct {
	apiKey 		string
	url		string
	oddsFormat	string
	client 		http.Client
}

func New(cfg Config) *client {
	if cfg.Key == "" && cfg.Url == "" {
		log.Fatalln("Api-key and Url are required")
	}
	cli := client{
		apiKey: 	cfg.Key,
		url:		cfg.Url,
		oddsFormat:	cfg.OddsFormat,
	}

	return &cli
}

func (cli client) GetSports() (entity.Sports, error) {
	var sports entity.Sports
	var data map[int]string

	sportsByte, err := cli.request("sports", 0)
	if err != nil {
		return sports, err
	}

	if err := json.Unmarshal(sportsByte, &data); err != nil {
		return sports, err
	}
	sports.SetData(data)

	return sports, nil
}

func (cli client) GetOddTypes() (entity.OddTypes, error) {
	var ot entity.OddTypes
	var data []string

	otByte, err := cli.request("oddtype", 0)
	if err != nil {
		return ot, err
	}

	if err := json.Unmarshal(otByte, &data); err != nil {
		return ot, err
	}
	ot.SetData(data)

	return ot, nil
}

func (cli client) GetOddsBySport(sport string, source int) ([]entity.EventOdds, error) {
	var eo []entity.EventOdds
	eoByte, err := cli.request("odds/" + sport, 0)
	if err != nil {
		return eo, err
	}
	if err := json.Unmarshal(eoByte, &eo); err != nil {
		return eo, err
	}

	return eo, nil
}

func (cli client) request(path string, source int) ([]byte, error) {
	client := &http.Client{}

	url, err := url.Parse(cli.url + path);
	if err != nil {
		return nil, err
	}
	q := url.Query()
	q.Set("source", string(source))
	q.Set("oddsFormat", cli.oddsFormat)
	url.RawQuery = q.Encode()
	req := &http.Request{
		Method:	"GET",
		URL:	url,
		Header:	http.Header{
			"JsonOdds-API-Key": {cli.apiKey},
		},
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)

	return b, nil
}