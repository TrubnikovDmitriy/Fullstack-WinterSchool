package serv

import "crypto/sha256"

func PasswordHashing(password string) []byte {
	sha := sha256.New()
	sha.Write([]byte(password))
	return sha.Sum(nil)
}

func PasswordEqual(password string, hash []byte) bool {

	newHash := PasswordHashing(password)
	for i, symbol := range hash {
		if newHash[i] != symbol {
			return false
		}
	}
	return true
}
