package api

import (
	"encoding/gob"
	"time"
)

func init() {
	gob.Register(&SearchParameters{})
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})
	gob.Register(time.Time{})
}
