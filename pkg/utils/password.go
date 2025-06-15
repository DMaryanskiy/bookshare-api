package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(pswd string) (string, error) {
	hashedPswd, err := bcrypt.GenerateFromPassword([]byte(pswd), 14)
	return string(hashedPswd), err
}

func CheckPasswordHash(hashedPswd, pswd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPswd), []byte(pswd))
	return err == nil, err
}
