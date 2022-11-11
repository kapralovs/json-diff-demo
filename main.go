package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/wI2L/jsondiff"
)

type User struct {
	ID      int      `json:"id,omitempty"`
	Name    string   `json:"name,omitempty"`
	Weight  float32  `json:"weight,omitempty"`
	IsAdult bool     `json:"is_adult,omitempty"`
	Items   []string `json:"items,omitempty"`
	Vehicle []Car    `json:"vehicle,omitempty"`
}

type Car struct {
	Vendor     string   `json:"vendor,omitempty"`
	Model      string   `json:"model,omitempty"`
	Additional []string `json:"additional,omitempty"`
}

type Event struct {
	Initiator string      `json:"initiator,omitempty"`
	Subject   string      `json:"subject,omitempty"`
	Action    string      `json:"action,omitempty"`
	Rollback  interface{} `json:"rollback,omitempty"`
	Update    interface{} `json:"update,omitempty"`
}

type Diff struct {
	Op    string      `json:"op,omitempty"`
	Path  string      `json:"path,omitempty"`
	Value interface{} `json:"value,omitempty"`
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
			"smartphone",
		},
		Vehicle: []Car{
			{
				Vendor: "BMW",
				Model:  "M5",
				Additional: []string{
					"airbag",
					"music station",
					"4WD",
				},
			},
			{
				Vendor: "VW",
				Model:  "Polo",
				Additional: []string{
					"airbag",
					"music station",
				},
			},
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
		Vehicle: []Car{
			{
				Vendor: "Skoda",
				Model:  "Rapid",
			},
			{
				Vendor: "Chevrolet",
				Model:  "Tahoe",
				Additional: []string{
					"airbag",
					"music station",
					"offroad",
					"4WD",
				},
			},
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

	update, err := jsondiff.CompareJSONOpts(u1Serialized, u2Serialized, jsondiff.Invertible())
	if err != nil {
		fmt.Println(err)
	}
	rollback, err := jsondiff.CompareJSONOpts(u2Serialized, u1Serialized, jsondiff.Invertible())
	if err != nil {
		fmt.Println(err)
	}

	rollbackPatch := makePatch(rollback)
	updatePatch := makePatch(update)

	fmt.Println(update.String())
	fmt.Println()
	fmt.Println(rollback.String())

	// fmt.Println("ROLLBACK:")
	// for _, v := range rollbackPatch {
	// 	fmt.Printf("op: %v\n", v.Op)
	// 	fmt.Printf("path: %v\n", v.Path)
	// 	fmt.Printf("value: %v\n", v.Value)
	// }
	// fmt.Println("UPDATE:")
	// for _, v := range updatePatch {
	// 	fmt.Printf("op: %v\n", v.Op)
	// 	fmt.Printf("path: %v\n", v.Path)
	// 	fmt.Printf("value: %v\n", v.Value)
	// }

	rbSerialized, _ := json.Marshal(rollbackPatch)
	udSerialized, _ := json.Marshal(updatePatch)
	fmt.Printf("%v\n%v\n", string(rbSerialized), string(udSerialized))
	// e := &Event{
	// 	Initiator: user1.Name,
	// 	Subject:   user1.Name,
	// 	Action:    "update_user",
	// 	Rollback:  rollbackPatch,
	// 	Update:    updatePatch,
	// }
	// evnt, err := json.Marshal(e)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println(string(evnt))
}

func removeCase(pathParts []string, currentOp, previousOp jsondiff.Operation, diffs []*Diff) []*Diff {
	diff := &Diff{
		Op:   currentOp.Type,
		Path: pathParts[0],
	}

	if len(pathParts) > 1 {
		if diff.Value == nil {
			diff.Value = []interface{}{previousOp.Value}
			return append(diffs, diff)
		}
		if values, ok := diff.Value.([]interface{}); ok {
			if len(values) > 0 {
				diff.Value = append(values, previousOp.Value)
				return append(diffs, diff)
			}
		}
	}

	diff.Value = previousOp.Value
	return append(diffs, diff)
}

func addCase(pathParts []string, op jsondiff.Operation, diffs []*Diff) []*Diff {
	diff := &Diff{
		Op:   op.Type,
		Path: pathParts[0],
	}

	if len(pathParts) > 1 {
		if diff.Value == nil {
			diff.Value = []interface{}{op.Value}
			return append(diffs, diff)
		}
		if values, ok := diff.Value.([]interface{}); ok {
			diff.Value = append(values, op.Value)
			return append(diffs, diff)
		}
	}

	diff.Value = op.Value
	return append(diffs, diff)
}

func replaceCase(pathParts []string, op jsondiff.Operation, diffs []*Diff) []*Diff {
	diff := &Diff{
		Op:   op.Type,
		Path: pathParts[0],
	}
	if len(pathParts) > 1 {
		if diff.Value == nil {
			diff.Value = []interface{}{op.Value}
			return append(diffs, diff)
		}
		if values, ok := diff.Value.([]interface{}); ok {
			diff.Value = append(values, op.Value)
			return append(diffs, diff)
		}
	}
	diff.Value = op.Value
	return append(diffs, diff)
}

func makePatch(patch jsondiff.Patch) []*Diff {
	diffs := []*Diff{}
	for idx, op := range patch {
		// fmt.Println(op.Type)
		pathParts := strings.Split(op.Path.String()[1:], "/")
		switch op.Type {
		case "add":
			diffs = addCase(pathParts, op, diffs)
		case "replace":
			diffs = replaceCase(pathParts, op, diffs)
		case "remove":
			diffs = removeCase(pathParts, op, patch[idx-1], diffs)
		}
	}
	return diffs
}
