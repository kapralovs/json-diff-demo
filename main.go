package main

import (
	"fmt"
)

type User struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Weight  float32  `json:"weight"`
	IsAdult bool     `json:"is_adult"`
	Items   []string `json:"items"`
}

type Car struct {
	ID     int    `json:"id,omitempty"`
	Vendor string `json:"vendor,omitempty"`
	Model  string `json:"model,omitempty"`
	Color  string `json:"color,omitempty"`
	Is4WD  bool   `json:"is_4_wd,omitempty"`
}

type Event struct {
	Initiator string      `json:"initiator,omitempty"`
	Subject   string      `json:"subject,omitempty"`
	Action    string      `json:"action,omitempty"`
	Context   string      `json:"context,omitempty"`
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

	car1 := &Car{
		ID:     2,
		Vendor: "BMW",
		Model:  "X5",
		Color:  "black",
		Is4WD:  true,
	}
	car2 := &Car{
		ID:     2,
		Vendor: "BMW",
		Model:  "X5",
		Color:  "blue",
		Is4WD:  true,
	}

	err := AddEvent("some_user", "another_user", "user_update", "user", user1, user2)
	if err != nil {
		fmt.Println(err)
	}
	err = AddEvent("another_user", "third_user", "car_update", "car", car1, car2)
	if err != nil {
		fmt.Println(err)
	}
}

func AddEvent(initiator, subject, action, ctx string, oldData, newData interface{}) error {
	// if ctx == "user" {
	// 	first, ok := oldData.(*User)
	// 	if !ok {
	// 		return errors.New("can't make type assertion")
	// 	}
	// 	second, ok := newData.(*User)
	// 	if !ok {
	// 		return errors.New("can't make type assertion")
	// 	}
	// 	UserDiff(first, second)
	// }
	// if ctx == "car" {
	// 	first, ok := oldData.(*Car)
	// 	if !ok {
	// 		return errors.New("can't make type assertion")
	// 	}
	// 	second, ok := newData.(*Car)
	// 	if !ok {
	// 		return errors.New("can't make type assertion")
	// 	}
	// 	CarDiff(first, second)
	// }

	return nil
}

// func UserDiff(first, second *User) error {

// 	return nil
// }

// func CarDiff(first, second *Car) error {
// 	return nil
// }
