package hash

type Hasher struct {
	time uint32
	memory uint32
	threads uint8
	key_len uint32
	salt_len uint32
}

type Hash struct {
  Hash []byte
  Salt []byte
}
