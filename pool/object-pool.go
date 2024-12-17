package pool

import (
	"fmt"
	"sync"
)

type objectPool struct {
}

func (o *objectPool) New(newObject func() any) any {
	var pool = sync.Pool{
		New: newObject,
	}
	fmt.Print(pool)
	return nil
}
