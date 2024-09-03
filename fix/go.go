package fix

import (
	"bufio"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

func (t FixParams) GoFix() (preview []Preview, err error) {
	return checkGoModLine(t.Dir, t.CompList)
}

func checkGoModLine(dir string, compList []Comp) (previews []Preview, err error) {
	var repeat = make(map[string]bool)
	modInfoMap, err := scanGoMod(dir)

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

// 遍历所有的go.mod 文件
func scanGoMod(dir string) (fileInfo map[string]any, err error) {
	/*
		modName+modVersion:
			0:filepath
			1:line
	*/
	fileInfo = make(map[string]any)
	fileSystem := os.DirFS(dir)
	err = fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if d.Name() == "go.mod" {
			readMod(filepath.Join(dir, path), fileInfo)
			return nil
		}
		return nil
	})
	if err != nil {
		return
	}
	return
}

// 解析mod 所有包名 包版本  所在行号
func readMod(modPath string, fileInfo map[string]any) {
	var lines []string
	file, err := os.Open(modPath)
	if err != nil {
		log.Fatalf("Failed to read go.mod: %v", err)
	}
	goModData, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read go.mod: %v", err)
	}

	fileData, err := modfile.Parse("go.mod", goModData, nil)
	if err != nil {
		log.Fatalf("Failed to parse go.mod: %v", err)
	}
	file.Seek(0, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	for _, req := range fileData.Require {
		if req.Syntax.Start.Line != 1 && req.Mod.Version != "" {
			modInfo := ModPosition{
				Path:    modPath,
				Line:    req.Syntax.Start.Line,
				Content: scanContent(lines, req.Syntax.Start.Line),
			}
			key := req.Mod.Path + req.Mod.Version
			mod := checkNilModPosition(fileInfo[key])
			mod = append(mod, modInfo)
			fileInfo[key] = mod
		}
	}
}
