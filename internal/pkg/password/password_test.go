package password

import "testing"

func TestAdminSeedPasswordHash(t *testing.T) {
	seedHash := "$2a$10$0LCwL/15uA7zMUXusqGU6OofPjnYqvm3.jOBzYBETOrXPALCz4m9q"
	if !Check("123456", seedHash) {
		newHash, err := Hash("123456")
		if err != nil {
			t.Fatal(err)
		}
		t.Fatalf("seed hash does not match password 123456; use hash %s", newHash)
	}
}
