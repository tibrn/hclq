package actions

import (
	"io"
	"sync"

	"github.com/tibrn/hclq/hclq"
)

var (
	cache = map[string]*hclq.HclDocument{}
	lock  = sync.Mutex{}
)

func GetDocument(filename string, reader io.Reader) (*hclq.HclDocument, error) {
	lock.Lock()
	defer lock.Unlock()

	if doc, isDoc := cache[filename]; isDoc {
		return doc, nil
	}

	doc, err := hclq.FromReader(reader)

	if err == nil {
		cache[filename] = doc
	}

	return doc, err
}
