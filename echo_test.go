package echo

import "testing"

func TestEcho(t *testing.T) {
	input := "Hello, world"
	expected := "Hello, world"
	output := Echo(input)
	if output != expected {
		t.Errorf("Echo(%q) = %q; want %q", input, output, expected)
	}
}
