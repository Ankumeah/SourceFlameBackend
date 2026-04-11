package pat

import (
  "crypto/rand"
  "log"
  "encoding/base64"
)

func Get_PAT_Handler(prefix string, length uint8) *PAT_Handler {
  return &PAT_Handler{
    prefix: prefix,
    length: length,
  }
}

func (h *PAT_Handler) Genrate_PAT() (string, error) {
  _pat := make([]byte, h.length)
  if _, err := rand.Read(_pat); err != nil {
    log.Printf("Error while genrating PAT: %v\n", err.Error())
    return "", err
  }

  pat := base64.RawURLEncoding.EncodeToString(_pat)

  return h.prefix + pat, nil
}
