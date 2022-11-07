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
		ID:     1,
		Name:   "Serj",
		Weight: 80.1,
		Items: []string{
			"pencil",
			"sandwitch",
			"money",
			"smartphone",
		},
	}
	user2 := &User{
		ID:      1,
		Name:    "Serj",
		Weight:  85.9,
		IsAdult: true,
		Items: []string{
			"pen",
			"smartphone",
			"drink",
			"backpack",
			"eraser",
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
		fmt.Printf("Before operation: %v\n", before)
		fmt.Printf("Before operation: %v\n", after)
		pathParts := strings.Split(op.Path.String()[1:], "/")
		switch op.Type {
		case "test":
			fmt.Println("TEST")
			testTypeCase(pathParts, op, before, after)
			if patch[idx+1].Type == "remove" {
				fmt.Println("REMOVE")
				nextPathParts := strings.Split(patch[idx+1].String()[1:], "/")
				removeTypeCase(nextPathParts, patch[idx+1], after)
				fmt.Printf("After operation: %v\n", before)
				fmt.Printf("After operation: %v\n", after)
				continue
			}
			if patch[idx+1].Type == "replace" {
				fmt.Println("REPLACE")
				nextPathParts := strings.Split(patch[idx+1].Path.String()[1:], "/")
				replaceTypeCase(op.Value, nextPathParts, patch[idx+1], after)
				fmt.Printf("After operation: %v\n", before)
				fmt.Printf("After operation: %v\n", after)
				continue
			}
			after[patch[idx+1].Path.String()[1:]] = patch[idx+1].Value
			fmt.Printf("After operation: %v\n", before)
			fmt.Printf("After operation: %v\n", after)
		case "add":
			fmt.Println("ADD")
			addTypeCase(idx, pathParts, op, after)
			fmt.Printf("After operation: %v\n", before)
			fmt.Printf("After operation: %v\n", after)
		}
	}

	fmt.Println(patch.String())
	fmt.Println(before)
	fmt.Println(after)
}

func testTypeCase(pathParts []string, op jsondiff.Operation, before, after map[string]interface{}) {
	if len(pathParts) > 1 {
		if _, ok := before[pathParts[0]]; !ok {
			items := []interface{}{op.Value}
			before[pathParts[0]] = items
			after[pathParts[0]] = items
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

func removeTypeCase(pathParts []string, op jsondiff.Operation, after map[string]interface{}) {
	if len(pathParts) > 1 {
		if _, ok := after[pathParts[0]]; ok {
			if values, ok := after[pathParts[0]].([]interface{}); ok {
				for i := range values {
					if values[i] == op.Value {
						values[i] = nil
					}
				}
			}
		}
	}
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

func replaceTypeCase(oldValue interface{}, pathParts []string, op jsondiff.Operation, after map[string]interface{}) {
	if len(pathParts) > 1 {
		if _, ok := after[pathParts[0]]; !ok {
			after[pathParts[0]] = []interface{}{op.Value}
			return
		}
		if values, ok := after[pathParts[0]].([]interface{}); ok {
			fmt.Println("values", len(values))
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
