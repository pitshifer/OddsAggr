package entity

import "encoding/json"

type Stringer interface {
	String() string
}

type Sport struct {
	Id	int8
	Name	string
}

type Sports struct {
	Stringer
	Sports	[]Sport
}

type OddTypes struct {
	Stringer
	data []string
}

type Event struct {
	Stringer
	Id       string `json:"ID"`
	HomeTeam string
	AwayTeam string
	SportId  int    `json:"Sport"`
	Time     string `json:"MatchTime"`
	Details  string
	League   League
	Region   string
	HomeROT  string
	AwayROT  string
}

type Odds struct {
	Id          string `json:"ID"`
	EventId     string `json:"EventID"`
	BmId        int    `json:"SiteID"`
	Home        string `json:"MoneyLineHome"`
	Away        string `json:"MoneyLineAway"`
	Draw        string `json:"DrawLine"`
	OddType     string
	LastUpdated string
}

type League struct {
	Name string
}

type EventOdds struct {
	Event
	Odds []Odds
}

func (s Sports) String() string {
	b, err := json.Marshal(s.Sports)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func (ot OddTypes) String() string {
	b, err := json.Marshal(ot.data)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func (ot *OddTypes) SetData(data []string) {
	ot.data = data
}

func (eo Event) String() {

}
