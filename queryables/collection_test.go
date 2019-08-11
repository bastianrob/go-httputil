package queryables

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/bastianrob/go-httputil/queryables/qtype"
)

func TestCollection_ToQuery(t *testing.T) {
	createRequest := func(rawQuery string) *http.Request {
		r, _ := http.NewRequest("GET", "", nil)
		r.URL.RawQuery = rawQuery

		return r
	}
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		c    Collection
		args args
		want map[string]interface{}
	}{{
		name: "OK#1",
		c: Collection{
			{QueryKey: "name", DBKey: "fullname", Type: qtype.String, Transform: func(k string, v []interface{}) (string, interface{}, error) { return k, v[0], nil }},
			{QueryKey: "status", DBKey: "status", Type: qtype.String, Transform: func(k string, v []interface{}) (string, interface{}, error) { return k, "in (alive, deceased)", nil }},
			{QueryKey: "age", DBKey: "age", Type: qtype.String, Transform: func(k string, v []interface{}) (string, interface{}, error) { return k, "between (0, 100)", nil }},
			{QueryKey: "number", DBKey: "number", Type: qtype.Integer, Transform: nil},
		},
		args: args{createRequest("name=john&status=alive,deceased&age=0,100&not_handled=something")},
		want: map[string]interface{}{
			"fullname": "john",
			"status":   "in (alive, deceased)",
			"age":      "between (0, 100)",
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.ToQuery(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collection.ToQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
