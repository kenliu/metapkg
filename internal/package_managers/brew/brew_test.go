package brew

import (
	"strings"
	"testing"
)

func TestInstallCommand(t *testing.T) {
	got := strings.Join((&BrewPackageState{}).InstallCommand("bat"), " ")
	expected := "brew install bat"
	if got != expected {
		t.Errorf("InstallCommand(\"bat\") = %q, expected %q", got, expected)
	}
}
