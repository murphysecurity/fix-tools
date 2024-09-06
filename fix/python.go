package fix

import (
	"bufio"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var pattern = `(\w+)(?:[><=!~=]+)([\d\.]+)`
var re = regexp.MustCompile(pattern)

func (t FixParams) PythonFix() (preview []Preview, err error) {
	return checkRequirementsLine(t.Dir, t.CompList)
}
func checkRequirementsLine(dir string, compList []Comp) (previews []Preview, err error) {
	var repeat = make(map[string]bool)
	modInfoMap, err := scanRequirements(dir)
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

// 遍历所有的requirements.txt 文件
func scanRequirements(dir string) (fileInfo map[string]any, err error) {
	fileInfo = make(map[string]any)
	fileSystem := os.DirFS(dir)
	err = fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if d.Name() == "requirements.txt" {
			readRequirementsMod(filepath.Join(dir, path), fileInfo)
			return nil
		}
		return nil
	})
	if err != nil {
		return
	}
	return
}

func readRequirementsMod(modPath string, fileInfo map[string]any) {
	var lines []string
	file, err := os.Open(modPath)
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for i, j := range lines {
		var modInfo ModPosition
		packageName, version, err := extractVersionInfo(j)
		if err != nil {
			continue
		}
		modInfo.Path = modPath
		modInfo.Line = i + 1
		modInfo.Content = scanContent(lines, i+1)
		modInfos := checkNilModPosition(fileInfo[packageName+version])
		modInfos = append(modInfos, modInfo)
		fileInfo[packageName+version] = modInfos
	}
}

// 去除注释
func spli(text string) string {
	index := strings.Index(text, "#")
	if index != -1 {
		text = text[0:index]
	}
	return text
}
func extractVersionInfo(line string) (packageName string, version string, err error) {
	line = spli(line)
	if strings.Contains(line, "=") {
		matches := re.FindStringSubmatch(line)
		if len(matches) == 3 {
			return matches[1], matches[2], nil
		} else {
			return "", "", errors.New("")
		}
	} else {
		return strings.ReplaceAll(line, " ", ""), "", nil
	}

}
