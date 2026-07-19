package pat

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

func GetPATHandler(prefix string, length uint8) *PATHandler {
	return &PATHandler{
		prefix: prefix,
		length: length,
	}
}

func (h *PATHandler) GeneratePAT() (string, error) {
	_pat := make([]byte, h.length)
	if _, err := rand.Read(_pat); err != nil {
		log.Printf("Error while generating PAT: %v\n", err.Error())
		return "", err
	}

	pat := base64.RawURLEncoding.EncodeToString(_pat)

	return h.prefix + pat, nil
}
