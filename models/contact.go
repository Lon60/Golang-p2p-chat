package models

type Contact struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	Port string `json:"port"`
}

func (c Contact) Identifier() string {
	return c.IP + ":" + c.Port
}
