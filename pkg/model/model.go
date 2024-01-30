package model

type Person struct {
	UID       int    `json:"uid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Result struct {
	Result bool   `json:"result"`
	Info   string `json:"info"`
}
