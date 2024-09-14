package models

import "crypto/rsa"

type Contact struct {
	Name      string         `json:"name"`
	IP        string         `json:"ip"`
	Port      string         `json:"port"`
	PublicKey []byte         `json:"public_key"` // Public Key als []byte für Übertragung
	KeyObject *rsa.PublicKey `json:"-"`          // Public Key als *rsa.PublicKey für Nutzung, nicht für JSON
}

func (c Contact) Identifier() string {
	return c.IP + ":" + c.Port
}
