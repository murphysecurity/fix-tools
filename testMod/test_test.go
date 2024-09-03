package testMod

import (
	"encoding/json"
	"fix-tools/fix"
	"fmt"
	"os"
	"testing"
)

func Test_GO(t *testing.T) {
	by, _ := os.ReadFile("./go/goBefer.json")
	var s []fix.Comp
	json.Unmarshal(by, &s)
	params := fix.FixParams{
		CompList: s,
		Dir:      "C:\\Users\\陈浩轩\\Desktop\\mofei\\PackageDependency\\测试源码包\\go\\go-01\\go-master",
	}
	p, _ := params.GoFix()
	b, _ := json.Marshal(&p)
	os.WriteFile("./go/goEnd.json", b, 0777)
}
func Test_Npm(t *testing.T) {
	by, _ := os.ReadFile("./npm/npmBefer.json")
	var s []fix.Comp
	json.Unmarshal(by, &s)
	params := fix.FixParams{
		CompList: s,
		Dir:      "C:\\Users\\陈浩轩\\Desktop\\mofei\\PackageDependency\\测试源码包\\pnpm1",
	}
	p, _ := params.NpmFix()
	b, _ := json.Marshal(&p)
	os.WriteFile("./npm/npmEnd.json", b, 0777)
}

func Test_Python(t *testing.T) {
	by, _ := os.ReadFile("./python/pythonBefer.json")
	var s []fix.Comp
	json.Unmarshal(by, &s)
	params := fix.FixParams{
		CompList: s,
		Dir:      "C:\\Users\\陈浩轩\\Desktop\\mofei\\PackageDependency\\测试源码包\\python\\python-master",
	}
	p, _ := params.PythonFix()
	b, _ := json.Marshal(&p)
	os.WriteFile("./python/pythonEnd.json", b, 0777)
}
func Test_Nuget(t *testing.T) {
	by, _ := os.ReadFile("./nuget/nugetBefer.json")
	var s []fix.Comp
	json.Unmarshal(by, &s)
	params := fix.FixParams{
		CompList: s,
		Dir:      "C:\\Users\\陈浩轩\\Desktop\\mofei\\PackageDependency\\aerospike-client-csharp",
	}
	p, _ := params.NugetFix()
	b, _ := json.Marshal(&p)
	os.WriteFile("./nuget/nugetEnd.json", b, 0777)
}
func Test_Maven(t *testing.T) {
	by, _ := os.ReadFile("./maven/mavenBefer.json")
	var s []fix.Comp
	json.Unmarshal(by, &s)
	params := fix.FixParams{
		CompList: s,
		Dir:      "C:\\Users\\陈浩轩\\Desktop\\mofei\\PackageDependency\\测试源码包\\java\\java-1\\TLog-1.0.2",
	}
	p, _ := params.MavenFix()
	b, _ := json.Marshal(&p)
	fmt.Println(len(p))
	os.WriteFile("./maven/mavenEnd.json", b, 0777)
}
