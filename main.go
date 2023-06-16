package main

import "github.com/murphysecurity/fix-tools/fix"

func main() {
	// todo
	println(666)
	//params := fix.FixParams{
	//	ShowOnly: true,
	//	CompList: []fix.Comp{{
	//		CompName:    "org.springframework.boot:spring-boot-configuration-processor",
	//		CompVersion: "2.0.5.RELEASE",
	//	}},
	//	PackageManager: "maven",
	//	RepoType:       "local",
	//	Dir:            "E:\\project\\java_project\\test\\java",
	//}
	params := fix.FixParams{
		ShowOnly: false,
		TimeOut:  60,
		RepoType: "github",
		CompList: []fix.Comp{{
			CompName:      "com.fasterxml.jackson.core:jackson-databind",
			CompVersion:   "2.8.7",
			MinFixVersion: "2.9.10.8",
		}},
		PackageManager: "maven",
		Dir:            "",
		ProxyUrl:       "http://127.0.0.1:7890",
		GitlabUrl:      "gitee.com",
		CommitHash:     "244babb436c3896b253495f6012d4ae36c9ed331",
		TargetOwner:    "645775992",
		Owner:          "",
		Repo:           "light-4j",
		UserName:       "645775992",
		Password:       "",
		UserEmail:      "645775992@qq.com",
		Token:          "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",
		Title:          "提交pr title",
		Body:           "提交pr body",
	}
	params.Fix()
	//if err != nil {
	//	println(err.Error())
	//}
	//for _, p := range preview {
	//	print(p.Line)
	//	print("--------")
	//	println(p.Path)
	//	for _, content := range p.Content {
	//		print(content.Line)
	//		print("    ")
	//		println(content.Text)
	//
	//	}
	//
	//}
}
