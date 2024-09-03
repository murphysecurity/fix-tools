package fix

func (t *FixParams) LocalFix() (preview []Preview, err error) {

	switch t.PackageManager {
	case "maven":
		preview, err = t.MavenFix()

	case "go":
		preview, err = t.GoFix()

	case "npm":
		preview, err = t.NpmFix()

	case "yarn":
		preview, err = t.YarnFix()

	case "python":
		preview, err = t.PythonFix()

	case "nuget":
		preview, err = t.NugetFix()

	}
	return
}
