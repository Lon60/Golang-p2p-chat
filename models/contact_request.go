package models

type ContactRequest struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	Port string `json:"port"`
}
