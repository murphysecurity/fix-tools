package fix

func (t *FixParams) Fix() (prUrl string, preview []Preview, dmPreview []Preview, haveDMList map[string]int, err error) {
	err = t.check()
	if err != nil {
		return
	}
	switch t.RepoType {
	case "github":
		prUrl, preview, err = t.GithubFix()
	case "gitee":
		prUrl, preview, err = t.GiteeFix()
	case "gitlab":
		preview, err = t.GitlabFix()
	case "local":
		preview, dmPreview, haveDMList, err = t.LocalFix()

	}
	return
}
