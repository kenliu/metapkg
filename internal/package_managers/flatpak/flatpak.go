package flatpak

type FlatpakPackageState struct {
}

func (f *FlatpakPackageState) IsInstalled(name string, arguments []string) (bool, error) {
	return false, nil
}

func (f *FlatpakPackageState) Install(name string) error {
	println("Installing flatpak package: ", name)
	println("Flatpak package manager is not implemented yet")
	return nil
}
