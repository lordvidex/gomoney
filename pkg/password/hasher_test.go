package password

import "testing"

func TestPasswordHash(t *testing.T) {
	hash, err := CreatePasswordHash("password")
	if err != nil {
		t.Error(err)
	}

	if err := CheckPasswordHash(hash, "password"); err != nil {
		t.Error(err)
	}

	if err := CheckPasswordHash(hash, "123123123"); err == nil {
		t.Error("password should not match")
	}
}
