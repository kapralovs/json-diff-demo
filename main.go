package main

import (
	"encoding/json"
	"fmt"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/wI2L/jsondiff"
)

type User struct {
	Name     string   `json:"name,omitempty"`
	Age      int      `json:"age,omitempty"`
	IsAdult  bool     `json:"is_adult"`
	Bp       Backpack `json:"bp,omitempty"`
	Items    []string `json:"items,omitempty"`
	Vehicles []Car    `json:"vehicles,omitempty"`
}

type Backpack struct {
	Notebook string `json:"notebook,omitempty"`
	IsFull   bool   `json:"is_full,omitempty"`
}

type Car struct {
	Vendor string `json:"vendor,omitempty"`
	Model  string `json:"model,omitempty"`
}

type Event struct {
	Rollback interface{}
	Update   interface{}
}

func main() {
	u1 := &User{Name: "John", Age: 27, IsAdult: true,
		Bp:       Backpack{IsFull: false},
		Items:    []string{"drink", "smartphone", "bubblegum"},
		Vehicles: []Car{{Vendor: "VW", Model: "Polo"}},
	}
	u2 := &User{Name: "John", Age: 28, IsAdult: false,
		Bp:    Backpack{Notebook: "Macbook", IsFull: true},
		Items: []string{"smartphone", "bubblegum"},
		Vehicles: []Car{
			{Vendor: "VW", Model: "Polo"},
			{Vendor: "Ford", Model: "Transit"},
		},
	}

	u1Serialized, _ := json.Marshal(u1)
	u2Serialized, _ := json.Marshal(u2)

	updatePatch, err := jsondiff.CompareJSONOpts(u1Serialized, u2Serialized, jsondiff.Invertible())
	if err != nil {
		panic(err)
	}
	rollbackPatch, err := jsondiff.CompareJSONOpts(u2Serialized, u1Serialized, jsondiff.Invertible())
	if err != nil {
		panic(err)
	}

	updated := applyPatch(u1Serialized, updatePatch)
	rollbacked := applyPatch(updated, rollbackPatch)

	fmt.Println(updatePatch.String())
	fmt.Println()
	fmt.Println(updatePatch.String())
	fmt.Printf("Original document: %s\n", u1Serialized)
	fmt.Printf("Original document2: %s\n", u2Serialized)
	fmt.Printf("After Update: %s\n", updated)
	fmt.Printf("After Rollback: %s\n", rollbacked)
}

func applyPatch(entity []byte, patch jsondiff.Patch) []byte {
	patchSerialized, _ := json.Marshal(patch)
	// fmt.Printf("PATCH DATA: %v\n", string(patchSerialized))
	p, _ := jsonpatch.DecodePatch(patchSerialized)
	patched, _ := p.Apply(entity)
	return patched
}
