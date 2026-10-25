package packages

import (
	"bufio"
	"os"
	"runtime"
	"strings"
)

// DefaultPackageManager returns the package manager to use for packages that
// do not name one explicitly. It is chosen based on the operating system:
//
//   - macOS: brew
//   - Debian, Ubuntu, and derivatives: apt
//   - Fedora, RHEL, and derivatives: dnf
//
// If the operating system cannot be identified, DEFAULT_PACKAGE_MANAGER is used.
func DefaultPackageManager() string {
	if runtime.GOOS == "darwin" {
		return "brew"
	}

	f, err := os.Open("/etc/os-release")
	if err != nil {
		return DEFAULT_PACKAGE_MANAGER
	}
	defer f.Close()

	return defaultPackageManagerForOSRelease(parseOSRelease(f))
}

// defaultPackageManagerForOSRelease picks a package manager from the ID and
// ID_LIKE fields of an os-release file.
// See https://www.freedesktop.org/software/systemd/man/latest/os-release.html
func defaultPackageManagerForOSRelease(osRelease map[string]string) string {
	ids := []string{osRelease["ID"]}
	ids = append(ids, strings.Fields(osRelease["ID_LIKE"])...)

	for _, id := range ids {
		switch strings.ToLower(id) {
		case "debian", "ubuntu":
			return "apt"
		case "fedora", "rhel", "centos":
			return "dnf"
		}
	}

	return DEFAULT_PACKAGE_MANAGER
}

// parseOSRelease reads the KEY=value pairs of an os-release file. Values may be
// optionally enclosed in single or double quotes.
func parseOSRelease(r interface{ Read([]byte) (int, error) }) map[string]string {
	result := make(map[string]string)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}

		value = strings.TrimSpace(value)
		if len(value) >= 2 && (value[0] == '"' || value[0] == '\'') && value[len(value)-1] == value[0] {
			value = value[1 : len(value)-1]
		}
		result[strings.TrimSpace(key)] = value
	}

	return result
}
