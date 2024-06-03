package data

import (
	"fmt"
	"os"
)

type CertPath struct {
	ID       string `json:"id"`
	CertPath string `json:"cert_path"`
	KeyPath  string `json:"key_path"`
}

func NewCertPath(id, certPath, keyPath string) *CertPath {
	return &CertPath{
		ID:       id,
		CertPath: certPath,
		KeyPath:  keyPath,
	}
}

func (c *CertPath) GetID() string {
	return c.ID
}

func (c *CertPath) GetCertStr() string {
	certStr, err := os.ReadFile(c.CertPath)
	if err != nil {
		fmt.Println("Error reading cert file:", err)
		return ""
	}
	return string(certStr)
}

func (c *CertPath) GetKeyStr() string {
	certKey, err := os.ReadFile(c.KeyPath)
	if err != nil {
		fmt.Println("Error reading key file:", err)
		return ""
	}
	return string(certKey)
}
