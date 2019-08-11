package queryables

import "net/http"

//Collection of queryables metadata
type Collection []Meta

//ToQuery convert collections of queryables metadata into actual query
func (c *Collection) ToQuery(r *http.Request) map[string]interface{} {
	arr := ([]Meta)(*c)
	query := make(map[string]interface{})
	for _, entry := range arr {
		key, val, err := entry.ToQuery(r)
		if err != nil {
			continue
		}

		query[key] = val
	}

	return query
}
