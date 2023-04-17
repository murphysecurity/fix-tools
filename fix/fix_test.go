package fix

import "testing"

func TestName(t *testing.T) {
	params := FixParams{
		Dir:            "",
		proxyUrl:       "",
		GitlabUrl:      "https://git.murphy-int.com",
		CommitHash:     "199df55ec2c81924c2b79b9535c08c589b30cec5",
		TargetOwner:    "wangxiaodong",
		Owner:          "",
		Repo:           "maven-test",
		TimeOut:        30,
		UserName:       "wangxiaodong@murphysec.com",
		Password:       "a645775992",
		UserEmail:      "645775992@qq.com",
		Token:          "YnTrmdyU-8FC-2N9eHzR",
		RepoType:       "gitlab",
		CompName:       "com.github.pagehelper:pagehelper",
		CompVersion:    "4.2.1",
		MinFixVersion:  "4.2.9",
		PackageManager: "maven",
		Title:          "提交pr title",
		Body:           "提交pr body",
		ShowOnly:       true,
	}
	preview, err := params.Fix()
	if err != nil {
		println(err.Error())
	}
	println(preview)
}
