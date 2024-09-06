package fix

import (
	"bufio"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func (t FixParams) YarnFix() (preview []Preview, err error) {
	return checkYarnLine(t.Dir, t.CompList)
}

func checkYarnLine(dir string, compList []Comp) (previews []Preview, err error) {
	var repeat = make(map[string]bool)
	modInfoMap, err := ScanYarnMod(dir)
	if err != nil {
		return
	}
	for _, comp := range compList {
		if v, ok := modInfoMap[comp.CompName+comp.CompVersion]; ok && !repeat[comp.CompName+comp.CompVersion] {
			ModPositions := checkNilModPosition(v)
			repeat[comp.CompName+comp.CompVersion] = true
			for _, j := range ModPositions {
				var preview Preview
				preview.Path = strings.ReplaceAll(j.Path, dir, "")
				preview.Line = j.Line
				preview.Content = j.Content
				previews = append(previews, preview)
			}
		}
	}
	return
}

// 遍历所有的package.json 文件
func ScanYarnMod(dir string) (fileInfo map[string]any, err error) {
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
			readYarnMod(filepath.Join(dir, path), fileInfo)
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
func readYarnMod(modPath string, fileInfo map[string]any) {
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
