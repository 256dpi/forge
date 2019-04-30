package forge

// Value is a general value that is used to transfer application specific data.
type Value = interface{}

// Signal is an empty value that is used to construct closable channels.
type Signal = struct{}
