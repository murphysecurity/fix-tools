package fix

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var reMavenProperties = regexp.MustCompile(`<([^>]+)>([^<]+)</[^>]+>`)

/*
name:version  =  properties的 name:version = path lines content
*/
func (t FixParams) MavenFix() (preview []Preview, err error) {
	params := mavenParams{
		PropertiesMP: make(map[string]ModPosition),
		Preview:      make([]Preview, 0),
		DmPreview:    make([]Preview, 0),
		HaveDMList:   map[string]int{},
	}

	return params.CheckMavenLine(t.Dir, t.CompList)
}

func (c mavenParams) CheckMavenLine(dir string, compList []Comp) (previews []Preview, err error) {
	var repeat = make(map[string]bool)
	modInfoMap, err := c.scanMavenMod(dir)
	for k, _ := range modInfoMap {
		fmt.Println(k)
	}
	if err != nil {
		panic(err)
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

// 遍历所有的go.mod 文件
func (c mavenParams) scanMavenMod(dir string) (fileInfo map[string]any, err error) {
	/*
		modName+modVersion:
			0:filepath
			1:line
	*/
	fileInfo = make(map[string]any)
	fileSystem := os.DirFS(dir)
	var filePath []string
	err = fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if d.Name() == "pom.xml" {
			filePath = append(filePath, filepath.Join(dir, path))
			c.initPropertiesMP(filepath.Join(dir, path))
			return nil
		}
		return nil
	})
	if err != nil {
		return
	}

	for _, j := range filePath {
		infoOne := c.UpdateVersion(j)
		for k, v := range infoOne {
			fileInfo[k] = v
		}
	}
	return
}

func (c mavenParams) initPropertiesMP(filePath string) {
	var lines []string
	haveProject := false
	isProperties := false
	// 打开pom.xml文件
	file, err := os.Open(filePath)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	//需要额外的获取peoject.Version
	projectVersion := struct {
		XMLName xml.Name `xml:"project"`
		Version string   `xml:"version"`
	}{}
	file.Seek(0, 0)
	by, _ := io.ReadAll(file)
	xml.Unmarshal(by, &projectVersion)
	if projectVersion.Version != "" {
		haveProject = true
	}
	for i, j := range lines {
		if !isProperties && haveProject {
			if j = strings.ReplaceAll(j, " ", ""); j == "<version>"+projectVersion.Version+"</version>" {
				var mod ModPosition
				mod.Line = i + 1
				mod.Path = filePath
				mod.Content = scanContent(lines, mod.Line)
				c.PropertiesMP["project.version:"+projectVersion.Version] = mod

			}
		}
		if isProperties {
			matches := reMavenProperties.FindAllStringSubmatch(j, -1)
			for _, match := range matches {
				if len(match) >= 3 {
					versionName := match[1]
					versionNumber := match[2]
					var mod ModPosition
					mod.Line = i + 1
					mod.Path = filePath
					mod.Content = scanContent(lines, mod.Line)
					c.PropertiesMP[versionName+":"+versionNumber] = mod
				}
			}
		}
		if strings.Contains(j, "</properties>") {
			break
		}
		if strings.Contains(j, "<properties>") {
			isProperties = true
		}
	}

}

func (c mavenParams) getLine(filePath string) (map[string]any, *projectXML) {
	var dependency bool // dependency==true进入包扫描区间
	var startLine int
	var infoString [3]string
	var lines []string
	project := UnmarshalXML(filePath)

	file, _ := os.Open(filePath)

	scanner := bufio.NewScanner(file)
	var resultMap = make(map[string]any)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	for i, j := range lines {

		text := j

		//当一个模块扫描开始了
		if strings.Contains(text, "<dependency>") {
			dependency = true
			continue
		}

		//当一个模块寻找结束了重置标志位
		if strings.Contains(text, " </dependency>") {
			info := infoString[0] + infoString[1] + infoString[2]
			info = strings.ReplaceAll(info, " ", "")
			for _, j := range project.DependencyManagement.Dependencies.Dependency {
				key := "<groupId>" + j.GroupId + "</groupId>" + "<artifactId>" + j.ArtifactId + "</artifactId>" + "<version>" + j.Version + "</version>"
				if info == key {
					mapKey := j.GroupId + ":" + j.ArtifactId + j.Version
					var mod ModPosition
					mod.Line = startLine
					mod.Path = filePath
					mod.Content = scanContent(lines, startLine)
					modSLI := checkNilModPosition(resultMap[mapKey])
					modSLI = append(modSLI, mod)
					resultMap[mapKey] = modSLI
				}
			}
			for _, j := range project.Dependencies.Dependency {
				key := "<groupId>" + j.GroupId + "</groupId>" + "<artifactId>" + j.ArtifactId + "</artifactId>" + "<version>" + j.Version + "</version>"
				if info == key {
					mapKey := j.GroupId + ":" + j.ArtifactId + j.Version
					var mod ModPosition
					mod.Line = startLine
					mod.Path = filePath
					mod.Content = scanContent(lines, i+1)
					modSLI := checkNilModPosition(resultMap[mapKey])
					modSLI = append(modSLI, mod)
					resultMap[mapKey] = modSLI
				}
			}
			dependency = false
			startLine = 0
		}
		//如果正在扫描一个模块
		if dependency {

			if strings.Contains(text, "<groupId>") && !strings.Contains(text, "<!--") {
				infoString[0] = text
			}
			if strings.Contains(text, "<artifactId>") && !strings.Contains(text, "<!--") {
				infoString[1] = text
			}
			if strings.Contains(text, "<version>") && !strings.Contains(text, "<!--") {
				startLine = i + 1
				infoString[2] = text
			}
		}
	}
	return resultMap, project
}
func (c mavenParams) UpdateVersion(dir string) map[string]any {
	getLineEndMp, _ := c.getLine(dir)

	for k, _ := range getLineEndMp {
		for i, j := range c.PropertiesMP {
			var version string
			infos := strings.Split(i, ":")
			name := infos[0]
			if len(infos) >= 2 {
				version = infos[1]
			}

			//如果他是需要替换版本号的形式
			//那就直接用properties里面的行数据
			if strings.Contains(k, "${"+name+"}") {
				delete(getLineEndMp, k)
				k = strings.ReplaceAll(k, "${"+name+"}", version)
				getLineEndMp[k] = j
				var mods []ModPosition
				mods = checkNilModPosition(getLineEndMp)
				mods = append(mods, j)
				getLineEndMp[k] = mods

			}
		}
	}

	return getLineEndMp
}
