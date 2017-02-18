package jsonodds

import (
	"encoding/json"
	"github.com/pitshifer/oddsaggr/entity"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
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
	log.Debug("Created client config")

	return &cli
}

func (cli client) GetSports() (*entity.Sports, error) {
	var sports entity.Sports
	var data map[int8]string

	sportsByte, err := cli.request("sports", map[string]string{})
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

	otByte, err := cli.request("oddtype", map[string]string{})
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(otByte, &ot); err != nil {
		return nil, err
	}

	return &ot, nil
}

func (cli client) GetFinalType() (*[]string, error)  {
	var ft []string

	ftB, err := cli.request("finaltype", map[string]string{})
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(ftB, &ft); err != nil {
		return nil, err
	}

	return &ft, nil;
}

func (cli client) GetOddsBySport(sport, source string) (*[]entity.EventOdds, error) {
	var eo []entity.EventOdds

	eoByte, err := cli.request("odds/" + sport, map[string]string{"source": source})
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(eoByte, &eo); err != nil {
		return nil, err
	}

	return &eo, nil
}

func (cli client) GetOdds(oddType, source string) (*[]entity.EventOdds, error) {
	var eo []entity.EventOdds

	eoByte, err := cli.request("odds", map[string]string{"source": source, "oddType": oddType})
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(eoByte, &eo); err != nil {
		return nil, err
	}

	return &eo, nil
}

func (cli client) request(path string, params map[string]string) ([]byte, error) {
	client := &http.Client{}

	u, err := url.Parse(cli.url + path)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	for p, n := range params {
		q.Set(p, n)
	}
	q.Set("oddFormat", cli.oddFormat)
	u.RawQuery = q.Encode()
	req := &http.Request{
		Method: "GET",
		URL:    u,
		Header: http.Header{
			"JsonOdds-API-Key": {cli.apiKey},
		},
	}
	log.WithFields(log.Fields{
		"url": u.String(),
	}).Info("Do request")

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{
		"statusCode": resp.StatusCode,
	}).Info("Response")

	b, err := ioutil.ReadAll(resp.Body)

	return b, nil
}