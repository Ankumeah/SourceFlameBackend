package jwt

type Handler func(string) (bool, error)
