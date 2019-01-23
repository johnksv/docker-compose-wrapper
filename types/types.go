package types

// EnvLookup A function returning a value for the given key, or nil to use the default
type EnvLookup func(key string, defaultValue string) string
