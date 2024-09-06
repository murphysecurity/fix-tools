package fix

import (
	"bufio"
	"encoding/xml"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type nugetMod struct {
	XMLName     xml.Name `xml:"Project"`
	PackageRefs []struct {
		Include string `xml:"Include,attr"`
		Version string `xml:"Version,attr"`
	} `xml:"ItemGroup>PackageReference"`
	Reference []struct {
		Include string `xml:"Include,attr"`
		Version string `xml:"Version,attr"`
	} `xml:"ItemGroup>Reference"`
}

var versionPattern = regexp.MustCompile(`Version=([^,]+)(?:,|$)`)

func (t FixParams) NugetFix() (preview []Preview, err error) {
	return checkNugetLine(t.Dir, t.CompList)
}

func checkNugetLine(dir string, compList []Comp) (previews []Preview, err error) {
	var repeat = make(map[string]bool)
	modInfoMap, err := ScanNugetMod(dir)
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
func ScanNugetMod(dir string) (fileInfo map[string]any, err error) {
	/*
		modName+modVersion:
			0:filepath
			1:line
	*/
	fileInfo = make(map[string]any)
	filePath, err := findCsproj(dir)
	if err != nil {
		return nil, err
	}
	for _, j := range filePath {
		readNugetMod(j, fileInfo)
	}
	return
}

// 解析package.json 所有包名 包版本  所在行号
func readNugetMod(modPath string, fileInfo map[string]any) {

	var lines []string
	file, err := os.Open(modPath)
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	for i, j := range lines {
		var include, version string

		if strings.Contains(j, "<PackageReference") {
			include, version, err = packageReference(j)
			if err != nil {
				continue
			}
		}
		if strings.Contains(j, "<Reference") {
			include, version, err = reference(j)
			if err != nil {
				continue
			}
		}
		var modInfo ModPosition
		modInfo.Path = modPath
		modInfo.Line = i + 1
		modInfo.Content = scanContent(lines, i+1)
		mod := checkNilModPosition(fileInfo[include+version])
		mod = append(mod, modInfo)
		fileInfo[include+version] = mod
	}
}
func packageReference(xmlData string) (string, string, error) {
	var s struct {
		Include string `xml:"Include,attr"`
		Version string `xml:"Version,attr"`
	}
	err := xml.Unmarshal([]byte(xmlData), &s)
	if err != nil {
		return "", "", err
	}
	return s.Include, s.Version, nil
}
func reference(xmlData string) (string, string, error) {
	var includePackage string
	if index := strings.Index(xmlData, ","); index != -1 && !strings.Contains(xmlData, " ") {
		includePackage = xmlData[:index]
	} else {
		return "", "", errors.New("includePackage nil")
	}
	versionMatches := versionPattern.FindStringSubmatch(xmlData)
	if len(versionMatches) == 0 {
		return "", "", errors.New("version nil")
	}
	return includePackage, versionMatches[1], nil
}

func findCsproj(path string) ([]string, error) {
	var csprojPath []string
	return csprojPath, filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		matched, err := filepath.Match("*.csproj", d.Name())
		if err != nil {
			return err
		}
		// 检查当前路径是否是.csproj文件
		if !d.IsDir() && matched {
			csprojPath = append(csprojPath, path)
		}
		return nil
	})
}

func analysis(path string, fileInfo map[string]any) (result map[string]any, err error) {
	var proj nugetMod
	xmlData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := xml.Unmarshal(xmlData, &proj); err != nil {
		return nil, err
	}

	for _, j := range proj.Reference {

		var mod struct {
			Include string `xml:"Include,attr"`
			Version string `xml:"Version,attr"`
		}

		var includePackage string
		if index := strings.Index(j.Include, ","); index != -1 && !strings.Contains(j.Include[:1], " ") {
			includePackage = j.Include[:index]

		} else {
			continue
		}
		versionMatches := versionPattern.FindStringSubmatch(j.Include)
		if len(versionMatches) == 0 {
			continue
		}
		mod.Include = includePackage
		mod.Version = versionMatches[1]
		if mod.Version != "" {
			var modInfo ModPosition
			modInfo.Path = path
			mods := checkNilModPosition(fileInfo[mod.Include+mod.Version])
			mods = append(mods, modInfo)
			fileInfo[mod.Include+mod.Version] = mods
		}

	}

	return
}
