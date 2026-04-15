package hash

import "golang.org/x/crypto/bcrypt"

const defaultCost = bcrypt.DefaultCost

// HashPassword hashes a plaintext password using bcrypt.
func HashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), defaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// CheckPassword compares a plaintext password against a bcrypt hash.
// Returns true if they match.
func CheckPassword(plain, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}
