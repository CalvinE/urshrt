package keygenerator

// KeyGenerator is designed to be the mechanism that generates the URL short keys.
type KeyGenerator interface {
	GenerateKey(length int) (string, error)
}
