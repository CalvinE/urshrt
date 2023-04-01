package keygenerator

import (
	"fmt"
	"math/rand"
	"time"
)

type randomKeyGenerator struct {
	randNumGen *rand.Rand
}

func (r *randomKeyGenerator) GenerateKey(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be greater than 0: %d", length)
	}
	data := make([]rune, length)
	for i := 0; i < length; i++ {
		randIndex := r.randNumGen.Intn(num_characters)
		data[i] = rune(valid_characters[randIndex])
	}
	return string(data), nil
}

func NewRandomKeyGenerator() (KeyGenerator, error) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	return &randomKeyGenerator{
		randNumGen: r,
	}, nil
}
