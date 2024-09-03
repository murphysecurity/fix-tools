package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/murphysecurity/fix-tools/fix"
)

func main() {
	// todo
	start := time.Now().Unix()
	println(666)
	by, _ := os.ReadFile("./a.json")
	var s []fix.Comp
	json.Unmarshal(by, &s)
	params := fix.FixParams{
		ShowOnly:       true,
		CompList:       s,
		PackageManager: "Npm",
		RepoType:       "local",
		Dir:            "C:\\Users\\陈浩轩\\Desktop\\mofei\\PackageDependency\\测试源码包\\java\\java-1\\TLog-1.0.2",
	}
	preview, _, _, _ := params.MavenFix()
	// dmPreviewBy, _ := json.Marshal(&dmPreview)
	// haveDMListBy, _ := json.Marshal(&haveDMList)
	previewby, _ := json.Marshal(&preview)
	// os.WriteFile("haveDMList.json", haveDMListBy, 0777)
	// os.WriteFile("dmPreviewBy.json", dmPreviewBy, 0777)
	os.WriteFile("npm.json", previewby, 0777)
	end := time.Now().Unix() - start
	fmt.Println(len(preview))
	fmt.Println(end)
	// by, _ = json.Marshal(&p)
	// fmt.Println(string(by))
	// fmt.Println(len(p))
	// params := fix.FixParams{
	// 	ShowOnly: false,
	// 	TimeOut:  60,
	// 	RepoType: "github",
	// 	CompList: []fix.Comp{{
	// 		CompName:      "com.fasterxml.jackson.core:jackson-databind",
	// 		CompVersion:   "2.8.7",
	// 		MinFixVersion: "2.9.10.8",
	// 	}},
	// 	PackageManager: "maven",
	// 	Dir:            "",
	// 	ProxyUrl:       "http://127.0.0.1:7890",
	// 	GitlabUrl:      "gitee.com",
	// 	CommitHash:     "244babb436c3896b253495f6012d4ae36c9ed331",
	// 	TargetOwner:    "645775992",
	// 	Owner:          "",
	// 	Repo:           "light-4j",
	// 	UserName:       "645775992",
	// 	Password:       "",
	// 	UserEmail:      "645775992@qq.com",
	// 	Token:          "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",
	// 	Title:          "提交pr title",
	// 	Body:           "提交pr body",
	// }
	// params.Fix()
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
