package main

import (
	"fmt"
	"git.murphy-int.com/wangxiaodong/fix_tool_go/gitub_fix"
)

func main() {
	params := gitub_fix.GithubFixParams{
		UserName:       "645775992",
		UserEmail:      "645775992@qq.com",
		Token:          "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",
		TargetUserName: "645775992",
		TargetRepo:     "gotest",
		CompName:       "github.com/ecnepsnai/web",
		CompVersion:    "v1.4.0",
		MinFixVersion:  "1.5.2",
		PackageManager: "gomod",
		Title:          "pr title",
		Body:           "pr body",
	}
	err := gitub_fix.GithubFix(params, "http://127.0.0.1:10809")
	if err != nil {
		fmt.Println("err ======" + err.Error())
	} else {
		fmt.Println("修改成功")

	}

}
