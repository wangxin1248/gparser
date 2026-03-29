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
		{
			name: "test_case11 max",
			expr: "max(a, b) == 3",
			data: map[string]interface{}{
				"a": 1,
				"b": 3,
			},
			want: true,
		},
		{
			name: "test_case12 min",
			expr: "min(a, b) == 1",
			data: map[string]interface{}{
				"a": 1,
				"b": 3,
			},
			want: true,
		},
		{
			name: "test_case13 max multi",
			expr: "max(a, b, 5) == 5",
			data: map[string]interface{}{
				"a": 1,
				"b": 3,
			},
			want: true,
		},
		{
			name: "test_case14 max float",
			expr: "max(a, b) == 2.5",
			data: map[string]interface{}{
				"a": 1.2,
				"b": 2.5,
			},
			want: true,
		},
		{
			name: "test_case15 min float",
			expr: "min(a, b) == 1.2",
			data: map[string]interface{}{
				"a": 1.2,
				"b": 2.5,
			},
			want: true,
		},
		{
			name: "test_case16 int arithmetic",
			expr: "a + b * c == 11",
			data: map[string]interface{}{
				"a": 2,
				"b": 3,
				"c": 3,
			},
			want: true,
		},
		{
			name: "test_case17 division",
			expr: "a / b == 2",
			data: map[string]interface{}{
				"a": 10,
				"b": 5,
			},
			want: true,
		},
		{
			name: "test_case18 bool and or",
			expr: "(a == 1 && b == 2) || c == 1",
			data: map[string]interface{}{
				"a": 1,
				"b": 3,
				"c": 1,
			},
			want: true,
		},
		{
			name: "test_case19 non function fallback",
			expr: "in_array(c, []int{1,2,3})",
			data: map[string]interface{}{
				"c": 2,
			},
			want: true,
		},
		{
			name: "test_case20 string compare",
			expr: "a == \"hello\"",
			data: map[string]interface{}{
				"a": "hello",
			},
			want: true,
		},
		{
			name: "test_case21 comparison mix float",
			expr: "a < b && b <= 5.2",
			data: map[string]interface{}{
				"a": 3,
				"b": 5.2,
			},
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

func TestGoParser_Evaluate(t *testing.T) {
	tests := []struct {
		name string
		expr string
		data map[string]interface{}
		want interface{}
	}{
		{
			name: "arithmetic int",
			expr: "a + b",
			data: map[string]interface{}{"a": 1, "b": 2},
			want: int64(3),
		},
		{
			name: "arithmetic float",
			expr: "a + b",
			data: map[string]interface{}{"a": 1.2, "b": 2.3},
			want: 3.5,
		},
		{
			name: "bool expression",
			expr: "a == 1 && b == 2",
			data: map[string]interface{}{"a": 1, "b": 2},
			want: true,
		},
		{
			name: "function max",
			expr: "max(a,b,5)",
			data: map[string]interface{}{"a": 1, "b": 3},
			want: 5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Evaluate(tt.expr, tt.data)
			if err != nil {
				t.Fatalf("Evaluate failed: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Evaluate(%q) = %#v, want %#v", tt.expr, got, tt.want)
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
