package queryables

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"testing"

	"github.com/bastianrob/go-httputil/queryables/qtype"
)

func TestMeta_ToQuery(t *testing.T) {
	createRequest := func(rawQuery string) *http.Request {
		r, _ := http.NewRequest("GET", "", nil)
		r.URL.RawQuery = rawQuery

		return r
	}

	parsefloat32 := func(num string) float32 {
		f, _ := strconv.ParseFloat(num, 32)
		return float32(f)
	}

	type fields struct {
		QueryKey  string
		DBKey     string
		Type      Type
		Transform Transform
	}
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantKey   string
		wantValue interface{}
		wantErr   bool
	}{{
		name: "OK#1",
		fields: fields{
			QueryKey: "this",
			DBKey:    "this",
			Type:     qtype.String,
		},
		args:      args{createRequest("this=this&that=that")},
		wantKey:   "this",
		wantValue: []interface{}{"this"},
		wantErr:   false,
	}, {
		name: "OK#2",
		fields: fields{
			QueryKey: "notexists",
			DBKey:    "notexists",
			Type:     qtype.String,
		},
		args:      args{createRequest("this=this&that=that")},
		wantKey:   "",
		wantValue: nil,
		wantErr:   true,
	}, {
		name: "OK#3",
		fields: fields{
			QueryKey: "that",
			DBKey:    "db_that",
			Type:     qtype.Integer,
		},
		args:      args{createRequest("this=this&that=3")},
		wantKey:   "db_that",
		wantValue: []interface{}{int64(3)},
		wantErr:   false,
	}, {
		name: "NOK#3",
		fields: fields{
			QueryKey: "that",
			DBKey:    "db_that",
			Type:     qtype.Integer,
		},
		args:      args{createRequest("this=this&that=notnumber")},
		wantKey:   "",
		wantValue: nil,
		wantErr:   true,
	}, {
		name: "OK#4",
		fields: fields{
			QueryKey: "that",
			DBKey:    "db_that",
			Type:     qtype.Float32,
		},
		args:      args{createRequest("this=this&that=3.10")},
		wantKey:   "db_that",
		wantValue: []interface{}{parsefloat32("3.10")},
		wantErr:   false,
	}, {
		name: "NOK#4",
		fields: fields{
			QueryKey: "that",
			DBKey:    "db_that",
			Type:     qtype.Float32,
		},
		args:      args{createRequest("this=this&that=notnumber")},
		wantKey:   "",
		wantValue: nil,
		wantErr:   true,
	}, {
		name: "OK#5",
		fields: fields{
			QueryKey: "that",
			DBKey:    "db_that",
			Type:     qtype.Float64,
		},
		args:      args{createRequest("this=this&that=3.105")},
		wantKey:   "db_that",
		wantValue: []interface{}{float64(3.105)},
		wantErr:   false,
	}, {
		name: "NOK#5",
		fields: fields{
			QueryKey: "that",
			DBKey:    "db_that",
			Type:     qtype.Float64,
		},
		args:      args{createRequest("this=this&that=notnumber")},
		wantKey:   "",
		wantValue: nil,
		wantErr:   true,
	}, {
		name: "NOT#IMPLEMENTED",
		fields: fields{
			QueryKey: "that",
			DBKey:    "db_that",
			Type:     qtype.Invalid,
		},
		args:      args{createRequest("this=this&that=notnumber")},
		wantKey:   "",
		wantValue: nil,
		wantErr:   true,
	}, {
		name: "OK#6",
		fields: fields{
			QueryKey: "that",
			DBKey:    "db_that",
			Type:     qtype.String,
			Transform: func(k string, v []interface{}) (string, interface{}, error) {
				return fmt.Sprintf("%s transformed", k),
					fmt.Sprintf("%s transformed", v[0]), nil
			},
		},
		args:      args{createRequest("this=this&that=that value")},
		wantKey:   "db_that transformed",
		wantValue: "that value transformed",
		wantErr:   false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Meta{
				QueryKey:  tt.fields.QueryKey,
				DBKey:     tt.fields.DBKey,
				Type:      tt.fields.Type,
				Transform: tt.fields.Transform,
			}
			gotKey, gotValue, err := m.ToQuery(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Meta.ToQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotKey != tt.wantKey {
				t.Errorf("Meta.ToQuery() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
			if !reflect.DeepEqual(gotValue, tt.wantValue) {
				t.Errorf("Meta.ToQuery() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}
