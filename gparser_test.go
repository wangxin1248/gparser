package gparser

import (
	"reflect"
	"testing"
)

func TestGoParser_Match(t *testing.T) {
	tests := []struct {
		name string
		expr string
		data map[string]interface{}
		want bool
	}{
		{
			name: "test_case1",
			expr: "a == 1 && b == 2",
			data: map[string]interface{}{
				"a": 1,
				"b": 2,
			},
			want: true,
		},
		{
			name: "test_case2",
			expr: "a == 1 && b == 2",
			data: map[string]interface{}{
				"a": 1,
				"b": 3,
			},
			want: false,
		},
		{
			name: "test_case3",
			expr: "a == 1 && b == 2 || c == \"test\"",
			data: map[string]interface{}{
				"a": 1,
				"b": 3,
				"c": "test",
			},
			want: true,
		},
		{
			name: "test_case4",
			expr: "a == 1 && b == 2 && c == \"test\"",
			data: map[string]interface{}{
				"a": 1,
				"b": 3,
				"c": "test",
			},
			want: false,
		},
		{
			name: "test_case5",
			expr: "a == 1 && b == 2 && c == \"test\" && d == true",
			data: map[string]interface{}{
				"a": 1,
				"b": 2,
				"c": "test",
				"d": true,
			},
			want: true,
		},
		{
			name: "test_case6",
			expr: "a == 1 && b == 2 && c == \"test\" && d == false",
			data: map[string]interface{}{
				"a": 1,
				"b": 2,
				"c": "test",
				"d": true,
			},
			want: false,
		},
		{
			name: "test_case7",
			expr: "!(a == 1 && b == 2 && c == \"test\" && d == false)",
			data: map[string]interface{}{
				"a": 1,
				"b": 2,
				"c": "test",
				"d": true,
			},
			want: true,
		},
		{
			name: "test_case8",
			expr: "!(a == 1 && b == 2) || (c == \"test\" && d == false)",
			data: map[string]interface{}{
				"a": 1,
				"b": 2,
				"c": "test",
				"d": false,
			},
			want: true,
		},
		{
			name: "test_case9",
			expr: "a == 1 && b == 2",
			data: nil,
			want: false,
		},
		{
			name: "test_case10",
			expr: "",
			data: nil,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Match(tt.expr, tt.data); !reflect.DeepEqual(got, tt.want) || err != nil {
				t.Errorf("goParser match failed, want=%v, got=%v, err=%v", tt.want, got, err)
			}
		})
	}
}

func BenchmarkGoParser_Match(b *testing.B) {
	// 规则表达式
	expr := `(a == 1 && b == "b" && in_array(c, []int{100,99,98,97})) || (d == false)`
	// 映射数据
	data := map[string]interface{}{
		"a": 1,
		"b": "b",
		"c": 100,
		"d": true,
	}
	for i := 0; i < b.N; i++ {
		Match(expr, data)
	}
}
