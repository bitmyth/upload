package registry

import (
	"crypto/md5"
	"hash"
)

type Registry struct {
}

func (r Registry) Hash() hash.Hash {
	return md5.New()
}
