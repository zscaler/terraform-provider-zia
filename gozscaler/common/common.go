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

func flattenManagedBy(managedBy []ManagedBy) []interface{} {
	managed := make([]interface{}, len(managedBy))
	for i, managedItem := range managedBy {
		managed[i] = map[string]interface{}{

			"id":         managedItem.ID,
			"name":       managedItem.Name,
			"extensions": managedItem.Extensions,
		}
	}

	return managed
}

func flattenLastModifiedBy(lastModifiedBy []LastModifiedBy) []interface{} {
	lastModified := make([]interface{}, len(lastModifiedBy))
	for i, lastModifiedByItem := range lastModifiedBy {
		lastModified[i] = map[string]interface{}{

			"id":         lastModifiedByItem.ID,
			"name":       lastModifiedByItem.Name,
			"extensions": lastModifiedByItem.Extensions,
		}
	}

	return lastModified
}

func flattenLocation(location []Location) []interface{} {
	locations := make([]interface{}, len(location))
	for i, locationItem := range location {
		locations[i] = map[string]interface{}{

			"id":         locationItem.ID,
			"name":       locationItem.Name,
			"extensions": locationItem.Extensions,
		}
	}

	return locations
}
