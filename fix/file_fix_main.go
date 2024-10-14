package fix

func (t *FixParams) LocalFix() (preview []Preview, dmPreview []Preview, haveDMList map[string]int, err error) {

	switch t.PackageManager {
	case "maven":
		preview, dmPreview, haveDMList, err = t.MavenFix()
	case "go":
		preview, err = t.GoFix()
	case "npm":
		preview, err = t.NpmFix()
	case "yarn":
		preview, err = t.YarnFix()
	case "python":
		preview, err = t.PythonFix()
	case "maven_cli":
		preview, dmPreview, haveDMList, err = t.MavenFixNew()

	}
	return
}
