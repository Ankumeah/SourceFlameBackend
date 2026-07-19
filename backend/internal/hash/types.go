package hash

type Hasher struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
	saltLen uint32
}

type Hash struct {
	Hash []byte
	Salt []byte
}
