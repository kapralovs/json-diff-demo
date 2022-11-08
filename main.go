package main

import (
	"encoding/json"
	"fmt"
	"strings"

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
	OldData   interface{} `json:"old_data,omitempty"`
	NewData   interface{} `json:"new_data,omitempty"`
}

func main() {
	user1 := &User{
		ID:      1,
		Name:    "Some person",
		Weight:  80.1,
		IsAdult: true,
		Items: []string{
			"pencil",
			"sandwitch",
			"money",
			"smartphone",
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

	patch, err := jsondiff.CompareJSONOpts(u1Serialized, u2Serialized, jsondiff.Invertible())
	if err != nil {
		fmt.Println(err)
	}

	for idx, op := range patch {
		if op.Type == "replace" || op.Type == "remove" {
			continue
		}
		pathParts := strings.Split(op.Path.String()[1:], "/")
		switch op.Type {
		case "test":
			fmt.Printf("BEFORE REMOVE: old: %v\n", before)
			fmt.Printf("BEFORE REMOVE: new: %v\n", after)
			testTypeCase(pathParts, op, before, after)
			fmt.Printf("AFTER TEST: old: %v\n", before)
			fmt.Printf("AFTER TEST: new: %v\n", after)
			if patch[idx+1].Type == "remove" {
				nextPathParts := strings.Split(patch[idx+1].Path.String()[1:], "/")
				removeTypeCase(pathParts, nextPathParts, after)
				fmt.Printf("AFTER REMOVE: old: %v\n", before)
				fmt.Printf("AFTER REMOVE: new: %v\n", after)
				continue
			}
			if patch[idx+1].Type == "replace" {
				nextPathParts := strings.Split(patch[idx+1].Path.String()[1:], "/")
				replaceTypeCase(op.Value, nextPathParts, patch[idx+1], before, after)
				continue
			}
			after[patch[idx+1].Path.String()[1:]] = patch[idx+1].Value
		case "add":
			addTypeCase(idx, pathParts, op, after)
		}
	}

	fmt.Println(string(u1Serialized))
	fmt.Println(string(u2Serialized))
	fmt.Println(patch.String())
	fmt.Println(before)
	fmt.Println(after)

	e := &Event{
		Initiator: user1.Name,
		Subject:   user1.Name,
		Action:    "update_user",
		OldData:   before,
		NewData:   after,
	}
	evnt, err := json.Marshal(e)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(evnt))
}

func testTypeCase(pathParts []string, op jsondiff.Operation, before, after map[string]interface{}) {
	if len(pathParts) > 1 {
		if _, ok := before[pathParts[0]]; !ok {
			itemsBefore := []interface{}{op.Value}
			itemsAfter := []interface{}{op.Value}
			before[pathParts[0]] = itemsBefore
			after[pathParts[0]] = itemsAfter
			return
		}
		if values, ok := before[pathParts[0]].([]interface{}); ok {
			values = append(values, op.Value)
			before[pathParts[0]] = values
		}
		if values, ok := after[pathParts[0]].([]interface{}); ok {
			values = append(values, op.Value)
			after[pathParts[0]] = values
		}
		return
	}

	before[pathParts[0]] = op.Value
}

func removeTypeCase(pathParts, nextPathParts []string, after map[string]interface{}) {
	if len(pathParts) > 1 {
		if _, ok := after[pathParts[0]]; ok {
			if values, ok := after[pathParts[0]].([]interface{}); ok {
				if len(values) > 0 {
					for i := range values {
						if pathParts[1] == nextPathParts[1] {
							after[pathParts[0]] = values[:i]
						}
						return
					}
				}
			}
		}
	}
	fmt.Println("FAIL")
	after[pathParts[0]] = nil
}

func addTypeCase(idx int, pathParts []string, op jsondiff.Operation, after map[string]interface{}) {
	if len(pathParts) > 1 {
		if _, ok := after[pathParts[0]]; !ok {
			after[pathParts[0]] = []interface{}{op.Value}
		}
		if values, ok := after[pathParts[0]].([]interface{}); ok {
			values = append(values, op.Value)
			after[pathParts[0]] = values
		}
		return
	}
	after[pathParts[0]] = op.Value
}

func replaceTypeCase(oldValue interface{}, pathParts []string, op jsondiff.Operation, before, after map[string]interface{}) {
	if len(pathParts) > 1 {
		if _, ok := after[pathParts[0]]; !ok {
			after[pathParts[0]] = []interface{}{op.Value}
			return
		}
		if values, ok := after[pathParts[0]].([]interface{}); ok {
			for i := range values {
				if values[i] == oldValue {
					values[i] = op.Value
					after[pathParts[0]] = values
					return
				}
			}
		}

	}
	after[pathParts[0]] = op.Value
}
