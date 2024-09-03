package fix

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type packageJSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func (t FixParams) NpmFix() (preview []Preview, err error) {
	return checkNpmLine(t.Dir, t.CompList)
}

func checkNpmLine(dir string, compList []Comp) (previews []Preview, err error) {
	var repeat = make(map[string]bool)
	modInfoMap, err := ScanNpmMod(dir)
	if err != nil {
		return
	}
	for _, comp := range compList {
		if v, ok := modInfoMap[comp.CompName+comp.CompVersion]; ok && !repeat[comp.CompName+comp.CompVersion] {
			ModPositions := checkNilModPosition(v)
			repeat[comp.CompName+comp.CompVersion] = true
			for _, j := range ModPositions {
				var preview Preview
				preview.Path = comp.CompName
				preview.Line = j.Line
				preview.Content = j.Content
				previews = append(previews, preview)
			}
		}
	}
	return
}

// 遍历所有的package.json 文件
func ScanNpmMod(dir string) (fileInfo map[string]any, err error) {
	/*
		modName:
			0:filepath
			1:modVersion
			2:line
	*/
	fileInfo = make(map[string]any)
	fileSystem := os.DirFS(dir)
	err = fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if d.Name() == "package.json" {
			readNpmMod(filepath.Join(dir, path), fileInfo)
			return nil
		}
		return nil
	})
	if err != nil {
		return
	}
	return
}

// 解析package.json 所有包名 包版本  所在行号
func readNpmMod(modPath string, fileInfo map[string]any) {
	var lines []string
	file, err := os.Open(modPath)
	if err != nil {
		return
	}
	defer file.Close()

	var pkg packageJSON
	if err := json.NewDecoder(file).Decode(&pkg); err != nil {
		return
	}

	for modName, version := range pkg.Dependencies {

		fileInfo[modName] = info{
			Version: version,
		}
	}
	for modName, version := range pkg.DevDependencies {
		fileInfo[modName] = info{
			Version: version,
		}
	}
	file.Seek(0, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for i, j := range lines {

		modName, version := parseDependency(j)
		if modName == "" || version == "" {
			continue
		}
		fmt.Printf("name:%s : version%s\n", modName, version)
		modInfo := ModPosition{
			Path:    modPath,
			Line:    i + 1,
			Content: scanContent(lines, i+1),
		}
		key := modName + version
		mod := checkNilModPosition(fileInfo[key])
		mod = append(mod, modInfo)
		fileInfo[key] = mod
	}
}
func parseDependency(dependency string) (string, string) {
	dependency = strings.ReplaceAll(dependency, ",", "")
	dependency = strings.ReplaceAll(dependency, `"`, "")
	dependency = strings.TrimSpace(dependency)
	dependency = strings.Trim(dependency, "\"")

	parts := strings.Split(dependency, ":")
	if len(parts) != 2 {
		return "", ""
	}

	name := strings.TrimSpace(parts[0])
	version := strings.TrimSpace(parts[1])

	return name, version
}
