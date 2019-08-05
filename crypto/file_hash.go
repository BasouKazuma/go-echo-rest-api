package crypto

import (
	"encoding/json"
	"fmt"
	"crypto/sha256"
)

type HashConfig struct {
	Salt		string		`json:"salt"`
	MinLength	int			`json:"min_length"`
}

type FileHashData struct {
	UserId	int64	`json:"userId" db:"userId"`
	Name	string	`json:"name" db:"name"`
}

func CreateFileHash(data FileHashData) (string) {
	// Create Hash
	dataBytes, _ := json.Marshal(data)
	encoded := sha256.Sum256(dataBytes)
	return fmt.Sprintf("%x", encoded)
}
