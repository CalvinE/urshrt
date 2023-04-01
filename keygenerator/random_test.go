package keygenerator

import "testing"

func TestRandomGenerateKey(t *testing.T) {
	keyLength := 7
	iter := 1_000_000
	instances := make(map[string]bool)
	ranKeyGen, _ := NewRandomKeyGenerator()
	for i := 0; i < iter; i++ {
		randKey, err := ranKeyGen.GenerateKey(keyLength)
		if err != nil {
			t.Errorf("failed to generate key: %v", err)
		}
		_, ok := instances[randKey]
		if ok {
			t.Errorf("duplicate found after %d iterations", i)
		}
		instances[randKey] = true
	}
}
