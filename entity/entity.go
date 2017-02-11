package entity

import "encoding/json"

type Stringer interface {
	String() string
}

type Sports struct {
	Stringer
	data 	map[int]string
}

type OddTypes struct {
	Stringer
	data 	[]string
}

func (s Sports) String() string {
	b, err := json.Marshal(s.data)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func (s *Sports) SetData(data map[int]string) {
	s.data = data
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


