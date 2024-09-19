package models

import "crypto/rsa"

type Contact struct {
	Name      string         `json:"name"`
	IP        string         `json:"ip"`
	Port      string         `json:"port"`
	PublicKey []byte         `json:"public_key"`
	KeyObject *rsa.PublicKey `json:"-"`
}

func (c Contact) Identifier() string {
	return c.IP + ":" + c.Port
}
