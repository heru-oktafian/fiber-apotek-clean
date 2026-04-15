package idgen

import (
	"fmt"
	"math/rand"
	"time"
)

type Generator struct{}

func (Generator) New(prefix string) string {
	return fmt.Sprintf("%s%s%04d", prefix, time.Now().Format("060102150405"), rand.Intn(10000))
}
