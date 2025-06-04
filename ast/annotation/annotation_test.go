package annotation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnnotation(t *testing.T) {
	tests := []struct {
		input    string
		expected Annotation
	}{
		{"@format", Annotation{Name: "format"}},
		{"@format()", Annotation{Name: "format"}},
		{"@format(a=b)", Annotation{Name: "format", Values: map[string]string{"a": "b"}}},
		{`@format(a="b")`, Annotation{Name: "format", Values: map[string]string{"a": "b"}}},
		{"@format(a = b)", Annotation{Name: "format", Values: map[string]string{"a": "b"}}},
		{"@format(a = b, c = d)", Annotation{Name: "format", Values: map[string]string{"a": "b", "c": "d"}}},
		{`@format(a = "b", c = "d")`, Annotation{Name: "format", Values: map[string]string{"a": "b", "c": "d"}}},
		{`@format(a = "b", c = 123)`, Annotation{Name: "format", Values: map[string]string{"a": "b", "c": "123"}}},
	}

	for _, test := range tests {
		result := ParseAnnotation(test.input)
		require.Equal(t, test.expected, result.Output)
	}
}

func TestAnnotations(t *testing.T) {
	tests := []struct {
		input    string
		expected Annotations
	}{
		{"@format", Annotations{{Name: "format"}}},
		{"@format @check", Annotations{{Name: "format"}, {Name: "check"}}},
		{"@format(a=b) @check(c=d)", Annotations{{Name: "format", Values: map[string]string{"a": "b"}}, {Name: "check", Values: map[string]string{"c": "d"}}}},
	}

	for _, test := range tests {
		result := ParseAnnotations(test.input)
		require.Equal(t, test.expected, result.Output)
	}
}
