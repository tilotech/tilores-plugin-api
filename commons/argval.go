package commons

import "fmt"

// Value returns the value as an interface for the provided key from args or error if not found
func Value(args map[string]interface{}, key string) (interface{}, error) {
	if val, ok := args[key]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("key %v not found in args", key)
}
