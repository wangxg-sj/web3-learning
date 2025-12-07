package utils

import "golang.org/x/crypto/bcrypt"

func CreatePassword(password string) string {
	// 实际应用中，这里应该使用 bcrypt 等密码哈希算法
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hashedPassword)
}

func CheckPassword(password, hashedPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false
	}
	return true
}
