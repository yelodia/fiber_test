package models

import "fmt"

type RequestState string

const (
	RequeststateSuccess RequestState = "SUCCESS"
	RequestStateError   RequestState = "ERROR"
	RequestStateNew     RequestState = "NEW"
	RequestStatePending RequestState = "PENDING"
)

type Request struct {
	UUID    string       `json:"uuid"`
	Value   int          `json:"value"`
	Timeout int          `json:"timeout"`
	Result  int          `json:"result"`
	State   RequestState `json:"state"`
}

func (r *Request) init() {
	fmt.Println("init req")
}
