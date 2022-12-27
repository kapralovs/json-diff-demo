package main

import (
	"encoding/json"
	"fmt"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/wI2L/jsondiff"
)

type User struct {
	Name    string `json:"name,omitempty"`
	Age     int    `json:"age,omitempty"`
	IsAdult bool   `json:"is_adult"`
	// Bp       Backpack `json:"bp,omitempty"`
	// Items    []string `json:"items,omitempty"`
	// Vehicles []Car    `json:"vehicles,omitempty"`
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
	testNames := []string{
		"Sam",
		"Serj",
		"Aaron",
		"Henry",
		"Steven",
		"Jackie",
		"Alex",
		"Johan",
		"Bjorn",
		"Anders",
	}
	updtPatches := make([]jsondiff.Patch, 0)
	rbPatches := make([]jsondiff.Patch, 0)
	u1 := &User{Name: "John", Age: 27, IsAdult: true} // Bp:       Backpack{IsFull: false},
	// Items:    []string{"drink", "smartphone", "bubblegum"},
	// Vehicles: []Car{{Vendor: "VW", Model: "Polo"}},

	u2 := *u1

	// u1Serialized, _ := json.Marshal(u1)

	for idx, name := range testNames {
		if idx == 3 {
			u2.Age += 3
		}
		if idx == 6 {
			u2.IsAdult = false
		}
		u2.Name = name
		u1Serialized, _ := json.Marshal(u1)
		u2Serialized, _ := json.Marshal(&u2)
		updatePatch, err := jsondiff.CompareJSONOpts(u1Serialized, u2Serialized, jsondiff.Invertible())
		if err != nil {
			panic(err)
		}
		rollbackPatch, err := jsondiff.CompareJSONOpts(u2Serialized, u1Serialized, jsondiff.Invertible())
		if err != nil {
			panic(err)
		}
		if idx == 3 {
			u1.Age += 3
		}
		if idx == 6 {
			u1.IsAdult = false
		}
		u1.Name = name
		fmt.Printf("ROLLBACK PATCH %d:\n", idx)
		rbPatches = append(rbPatches, rollbackPatch)
		fmt.Println(rollbackPatch.String())
		fmt.Printf("UPDATE PATCH %d:\n", idx)
		updtPatches = append(updtPatches, updatePatch)
		fmt.Println(updatePatch.String())
		fmt.Println()
	}

	u2Ptr := &u2

	fmt.Println("ROLLBACK:")
	for i := len(rbPatches) - 1; i >= 0; i-- {
		fmt.Println(i)
		u2Serialized, _ := json.Marshal(u2Ptr)
		// fmt.Println(string(u2Serialized))
		patched, err := applyPatch(u2Serialized, rbPatches[i])
		if err != nil {
			fmt.Println(err)
		}
		err = json.Unmarshal(patched, u2Ptr)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(patched))
	}
	fmt.Println()
	fmt.Println("UPDATE:")
	for i, patch := range updtPatches {
		fmt.Println(i)
		u2Serialized, _ := json.Marshal(u2Ptr)
		// fmt.Println(string(u2Serialized))
		patched, err := applyPatch(u2Serialized, patch)
		if err != nil {
			fmt.Println(err)
		}
		err = json.Unmarshal(patched, u2Ptr)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(patched))
	}
}

func applyPatch(entity []byte, patch jsondiff.Patch) ([]byte, error) {
	patchSerialized, _ := json.Marshal(patch)
	p, err := jsonpatch.DecodePatch(patchSerialized)
	if err != nil {
		return nil, err
	}
	patched, err := p.Apply(entity)
	if err != nil {
		return nil, err
	}
	return patched, err
}
