package encryption

import "testing"

func TestPasswordHash(t *testing.T) {
	hasher := NewBcryptPasswordHasher()

	hash, err := hasher.CreatePasswordHash("password")
	if err != nil {
		t.Error(err)
	}

	if err := hasher.CheckPasswordHash(hash, "password"); err != nil {
		t.Error(err)
	}

	if err := hasher.CheckPasswordHash(hash, "123123123"); err == nil {
		t.Error("password should not match")
	}
}
