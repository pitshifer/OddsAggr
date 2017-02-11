package jsonodds

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/pitshifer/oddsaggr/entity"
)

type Config struct {
	Key, Url string
}

type client struct {
	apiKey 		string
	url		string
	client 		http.Client
}

func New(cfg Config) *client {
	cli := client{
		apiKey: cfg.Key,
		url: cfg.Url,
	}

	return &cli
}

func (cli client) GetSports() (entity.Sports, error) {
	var sports entity.Sports
	var data map[int]string

	sportsByte, err := cli.request("sports")
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

	otByte, err := cli.request("oddtype")
	if err != nil {
		return ot, err
	}

	if err := json.Unmarshal(otByte, &data); err != nil {
		return ot, err
	}
	ot.SetData(data)

	return ot, nil
}

func (cli client) request(path string) ([]byte, error) {
	client := &http.Client{}

	url := cli.url + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("JsonOdds-API-Key", cli.apiKey)

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)

	return b, nil
}