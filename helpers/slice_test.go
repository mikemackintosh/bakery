package helpers

import (
	"reflect"
	"testing"
)

var appendTest = []struct {
	Have   []string
	Append []string
	Want   []string
}{
	{
		Have:   []string{"a", "quick"},
		Append: []string{"brown", "fox"},
		Want:   []string{"a", "quick", "brown", "fox"},
	},
}

func TestAppend(t *testing.T) {
	for _, test := range appendTest {
		result := Append(test.Have, test.Append)
		if !reflect.DeepEqual(result, test.Want) {
			t.Errorf("failed to append. Have %#v, want %#v", result, test.Want)
		}
	}
}

var prependTest = []struct {
	Have    []string
	Prepend []string
	Want    []string
}{
	{
		Have:    []string{"brown", "fox"},
		Prepend: []string{"a", "quick"},
		Want:    []string{"a", "quick", "brown", "fox"},
	},
}

func TestPrepend(t *testing.T) {
	for _, test := range prependTest {
		result := Prepend(test.Prepend, test.Have)
		if !reflect.DeepEqual(result, test.Want) {
			t.Errorf("failed to prepend. Have %#v, want %#v", result, test.Want)
		}
	}
}
