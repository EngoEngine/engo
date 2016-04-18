package engo

import "testing"

func TestNewResourceWithoutExtension(t *testing.T) {
	url := "pineapple"

	r := NewResource(url)
	if r != (Resource{}) {
		t.Error("Expected empty resource.")
	}
}
