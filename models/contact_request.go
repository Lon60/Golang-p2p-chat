package models

import "crypto/rsa"

type ContactRequest struct {
	Name      string         `json:"name"`
	IP        string         `json:"ip"`
	Port      string         `json:"port"`
	PublicKey []byte         `json:"public_key"`
	KeyObject *rsa.PublicKey `json:"-"`
}
