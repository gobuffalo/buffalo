package randx

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"
)

func String(length int) string {
	s := ""
	for len(s) < length {
		data := []byte(fmt.Sprintf("%x%x", time.Now().UnixNano(), rand.Int63()))
		s += fmt.Sprintf("%x", md5.Sum(data))
	}
	if len(s) > length {
		return s[:length]
	}
	return s
}
