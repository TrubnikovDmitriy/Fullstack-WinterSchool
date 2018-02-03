package serv

import "crypto/sha256"

func PasswordHashing(password string) []byte {
	sha := sha256.New()
	sha.Write([]byte(password))
	return sha.Sum(nil)
}
