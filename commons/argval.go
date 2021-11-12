package commons

import "fmt"

func Value(args map[string]interface{}, key string) (interface{}, error) {
	if val, ok := args[key]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("key %v not found in args", key)
}

func StringValue(args map[string]interface{}, key string) (string, error) {
	val, ok := args[key]
	if !ok {
		return "", fmt.Errorf("key %v not found in args", key)
	}
	if s, ok := val.(string); ok {
		return s, nil
	}
	return "", fmt.Errorf("key %v is not a string but a %T", key, val)
}
