package fix

func (t *FixParams) Fix() (preview []Preview, err error) {
	err = t.check()
	if err != nil {
		return
	}
	switch t.RepoType {
	case "github":
		preview, err = t.GithubFix()
	case "gitee":
		preview, err = t.GiteeFix()
	case "gitlab":
		preview, err = t.GitlabFix()
	case "local":
		preview, err = t.LocalFix()

	}
	return
}
