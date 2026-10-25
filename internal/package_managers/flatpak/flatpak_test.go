package flatpak

import (
	"strings"
	"testing"
)

func TestFlatpakPackageState_IsInstalled(t *testing.T) {

}

func TestFlatpakPackageState_Install(t *testing.T) {
}

func TestInstallCommand(t *testing.T) {
	f := &FlatpakPackageState{}

	got := strings.Join(f.InstallCommand("org.kde.filelight"), " ")
	expected := "flatpak install -y org.kde.filelight"
	if got != expected {
		t.Errorf("InstallCommand(\"org.kde.filelight\") = %q, expected %q", got, expected)
	}

	got = strings.Join(flathubInstallCommand("org.kde.filelight"), " ")
	expected = "flatpak install -y flathub org.kde.filelight"
	if got != expected {
		t.Errorf("flathubInstallCommand(\"org.kde.filelight\") = %q, expected %q", got, expected)
	}

	description := f.InstallDescription("org.kde.filelight")
	if len(description) != 1 || !strings.Contains(description[0], expected) || !strings.Contains(description[0], "may retry") {
		t.Errorf("InstallDescription(\"org.kde.filelight\") = %q, expected one line mentioning a retry with %q", description, expected)
	}
}
