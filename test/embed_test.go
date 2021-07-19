package test

import (
	_ "embed"
	"testing"
)

//go:embed examples/getuid-block.yml
var example1 string

func TestGoEmbed(t *testing.T) {
	if example1 == "" {
		t.Errorf("//go:embed should have worked")
	}
}
