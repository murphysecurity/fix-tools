package fix

func (t *FixParams) Fix() (response Response) {
	response.Err = t.check()
	if response.Err != nil {
		return
	}
	switch t.RepoType {
	case "github":
		response.PrUrl, response.Preview, response.Err = t.GithubFix()
	case "gitee":
		response.PrUrl, response.Preview, response.Err = t.GiteeFix()
	case "gitlab":
		response.PrUrl, response.Preview, response.Err = t.GitlabFix()
	case "local":
		response.Preview, response.Err = t.LocalFix()
	}
	return
}
