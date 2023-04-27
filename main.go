package main

import "github.com/murphysecurity/fix-tools/fix"

func main() {
	// todo
	println(666)
	params := fix.FixParams{
		ShowOnly: true,
		CompList: []fix.Comp{{
			CompName:    "org.springframework.boot:spring-boot-configuration-processor",
			CompVersion: "2.0.5.RELEASE",
		}},
		PackageManager: "maven",
		RepoType:       "local",
		Dir:            "E:\\project\\java_project\\test\\java",
	}
	preview, err := params.Fix()
	if err != nil {
		println(err.Error())
	}
	for _, p := range preview {
		print(p.Line)
		print("--------")
		println(p.Path)
		for _, content := range p.Content {
			print(content.Line)
			print("    ")
			println(content.Text)

		}

	}
}
