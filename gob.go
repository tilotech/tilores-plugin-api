package api

import (
	"encoding/gob"
)

func init() {
	gob.Register(&SearchParameters{})
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})
}
