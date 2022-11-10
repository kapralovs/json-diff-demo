package main

import (
	"encoding/json"
	"fmt"

	"github.com/wI2L/jsondiff"
)

type User struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Weight  float32  `json:"weight"`
	IsAdult bool     `json:"is_adult"`
	Items   []string `json:"items"`
}

type Event struct {
	Initiator string      `json:"initiator,omitempty"`
	Subject   string      `json:"subject,omitempty"`
	Action    string      `json:"action,omitempty"`
	Rollback  interface{} `json:"rollback,omitempty"`
	Update    interface{} `json:"update,omitempty"`
}

func main() {
	user1 := &User{
		ID:      1,
		Name:    "Some person",
		Weight:  80.1,
		IsAdult: false,
		Items: []string{
			"pencil",
			"sandwitch",
			"money",
		},
	}
	user2 := &User{
		ID:      1,
		Name:    "Some person",
		Weight:  80.1,
		IsAdult: true,
		Items: []string{
			"pencil",
			"sandwitch",
			"money",
			"smartphone",
			"gun",
		},
	}

	u1Serialized, err := json.Marshal(user1)
	if err != nil {
		fmt.Println(err)
	}
	u2Serialized, err := json.Marshal(user2)
	if err != nil {
		fmt.Println(err)
	}

	before := make(map[string]interface{})
	after := make(map[string]interface{})

	update, err := jsondiff.CompareJSONOpts(u1Serialized, u2Serialized, jsondiff.Invertible())
	if err != nil {
		fmt.Println(err)
	}

	rollback, err := jsondiff.CompareJSONOpts(u2Serialized, u1Serialized, jsondiff.Equivalent())
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(rollback.String())
	fmt.Println()
	fmt.Println(update.String())

	e := &Event{
		Initiator: user1.Name,
		Subject:   user1.Name,
		Action:    "update_user",
		Rollback:  before,
		Update:    after,
	}
	evnt, err := json.Marshal(e)
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println(string(u1Serialized))
	// fmt.Println(string(u2Serialized))
	fmt.Println(string(evnt))
}
