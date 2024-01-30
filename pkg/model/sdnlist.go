package model

type SDNList struct {
	SDNs []SDNEntry `json:"sdns"`
}

type SDNEntry struct {
	SDNType   string `json:"sdnType"`
	UID       string `json:"uid"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
