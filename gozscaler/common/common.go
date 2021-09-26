package common

// These are common struct fields used in other traffic forwarding resources
// This is an immutable reference to an entity. which mainly consists of id and name
// SD-WAN Partner that manages the location. If a partner does not manage the locaton, this is set to Self.
type ManagedBy struct {
	ID         string                 `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}
type LastModifiedBy struct {
	ID         string                 `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}
type Location struct {
	ID         string                 `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}
