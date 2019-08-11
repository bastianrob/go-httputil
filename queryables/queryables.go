package queryables

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/bastianrob/go-httputil/queryables/qtype"
)

//Error collection
var (
	ErrNotImplemented = errors.New("Not yet implemented")
	ErrNotExists      = errors.New("Query does not exists")
	ErrValueType      = errors.New("Value type conversion failed")
)

//Type of query data type implementation is in qtype package
type Type uint

//Transform function to transform http query into actual DB query
type Transform func(k string, v []interface{}) (rk string, rv interface{}, err error)

//Meta HTTP queryables
type Meta struct {
	QueryKey  string
	DBKey     string
	Type      Type
	Transform Transform
}

//ToQuery transform queryable meta into actual query
func (m *Meta) ToQuery(r *http.Request) (key string, value interface{}, err error) {
	q := r.URL.Query().Get(m.QueryKey)
	if q == "" {
		return "", nil, ErrNotExists
	}

	qvl := strings.Split(q, ",")
	val := make([]interface{}, len(qvl))
	for i, entry := range qvl {
		switch m.Type {
		case qtype.Integer:
			converted, err := strconv.ParseInt(entry, 10, 64)
			if err != nil {
				return "", nil, ErrValueType
			}

			val[i] = converted
		case qtype.Float32:
			converted, err := strconv.ParseFloat(entry, 32)
			if err != nil {
				return "", nil, ErrValueType
			}

			val[i] = float32(converted)
		case qtype.Float64:
			converted, err := strconv.ParseFloat(entry, 64)
			if err != nil {
				return "", nil, ErrValueType
			}

			val[i] = converted
		case qtype.String:
			val[i] = entry
		default:
			return "", nil, ErrNotImplemented
			//TODO
		}
	}

	if m.Transform != nil {
		return m.Transform(m.DBKey, val)
	}

	return m.DBKey, val, nil
}
