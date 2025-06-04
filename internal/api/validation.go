package api

// Validator is an interface for types that can validate themselves.
// The Valid method should return a map of validation errors, where the key
// is the field name and the value is the error message.
// If the object is valid, Valid() should return nil or an empty map.
type Validator interface {
	Valid() map[string]string
}
