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

type IDNameExtensions struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

func flattenManagedBy(managedBy ManagedBy) interface{} {
	return []map[string]interface{}{
		{
			"id":   managedBy.ID,
			"name": managedBy.Name,
		},
	}
}

func flattenLastModifiedBy(lastModifiedBy LastModifiedBy) interface{} {
	return []map[string]interface{}{
		{
			"id":   lastModifiedBy.ID,
			"name": lastModifiedBy.Name,
		},
	}
}

func flattenLocation(location Location) interface{} {
	return []map[string]interface{}{
		{
			"id":   location.ID,
			"name": location.Name,
		},
	}
}
