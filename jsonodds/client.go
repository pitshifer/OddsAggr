package jsonodds

import (
	"encoding/json"
	"github.com/pitshifer/oddsaggr/entity"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Config struct {
	Key, Url, OddFormat string
}

type client struct {
	apiKey     string
	url        string
	oddFormat  string
	client     http.Client
}

func New(cfg Config) *client {
	if cfg.Key == "" || cfg.Url == "" {
		log.Fatalln("Api-key and Url are required")
	}
	cli := client{
		apiKey:     cfg.Key,
		url:        cfg.Url,
		oddFormat: cfg.OddFormat,
	}

	return &cli
}

func (cli client) GetSports() (*entity.Sports, error) {
	var sports entity.Sports
	var data map[int8]string

	sportsByte, err := cli.request("sports", "0")
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(sportsByte, &data); err != nil {
		return nil, err
	}
	for i, n := range data {
		sport := entity.Sport{
			Id:	i,
			Name:	n,
		}
		sports.Sports = append(sports.Sports, sport)
	}

	return &sports, nil
}

func (cli client) GetOddTypes() (*entity.OddTypes, error) {
	var ot entity.OddTypes

	otByte, err := cli.request("oddtype", "0")
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(otByte, &ot); err != nil {
		return nil, err
	}

	return &ot, nil
}

func (cli client) GetOddsBySport(sport, source string) (*[]entity.EventOdds, error) {
	var eo []entity.EventOdds

	eoByte, err := cli.request("odds/" + sport, source)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(eoByte, &eo); err != nil {
		return nil, err
	}

	return &eo, nil
}

func (cli client) request(path, source string) ([]byte, error) {
	client := &http.Client{}

	u, err := url.Parse(cli.url + path)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("source", source)
	q.Set("oddFormat", cli.oddFormat)
	u.RawQuery = q.Encode()
	req := &http.Request{
		Method: "GET",
		URL:    u,
		Header: http.Header{
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