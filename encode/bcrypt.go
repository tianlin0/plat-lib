package encode

import "golang.org/x/crypto/bcrypt"

func BCryptPasswordEncoder(password string, cost ...int) (string, error) {
	oneCost := bcrypt.DefaultCost
	if len(cost) > 0 {
		if cost[0] > bcrypt.MinCost && cost[0] < bcrypt.MaxCost {
			oneCost = cost[0]
		}
	}
	retByte, err := bcrypt.GenerateFromPassword([]byte(password), oneCost)
	if err != nil {
		return "", err
	}
	return string(retByte), nil
}
func BCryptCompareHashAndPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
