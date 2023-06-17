package fix

import (
	"bufio"
	"fmt"
	"github.com/antchfx/xmlquery"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func (t FixParams) MavenFix() (preview []Preview, dmPreview []Preview, haveDMList map[string]int, err error) {

	fileSystem := os.DirFS(t.Dir)

	err = fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if d.Name() == "pom.xml" {
			t.pomPathList = append(t.pomPathList, path)
		}
		return nil
	})
	if err != nil {
		return
	}
	params := mavenParams{
		propertyMap:      make(map[string][]PropertyModel, 0),
		fixModelList:     make([]FixModel, 0),
		preview:          make([]Preview, 0),
		haveDmMap:        make(map[string]int, 0),
		dependenciesLine: make(map[string]int, 0),
	}
	params.parsePropertyNode(t, t.pomPathList)
	params.getFixModelList(t, t.pomPathList)
	params.modifyPom(t)

	if len(params.preview) == 0 && !t.DmFix {
		params.getExtensionFix(t, t.pomPathList)
		params.modifyPom(t)
	}
	if len(params.preview) == 0 && !t.DmFix {
		params.getInheritFix(t, t.pomPathList)
		params.modifyPom(t)
	}
	preview = params.preview
	dmPreview = params.dmPreview
	haveDMList = params.haveDmMap
	return
}

func (p *mavenParams) modifyPom(params FixParams) (err error) {

	for _, item := range params.CompList {
		if params.DmFix {
			newParams := mavenParams{
				propertyMap:      make(map[string][]PropertyModel, 0),
				fixModelList:     make([]FixModel, 0),
				preview:          make([]Preview, 0),
				haveDmMap:        make(map[string]int, 0),
				dependenciesLine: make(map[string]int, 0),
			}
			fixParams := params
			fixParams.CompList = params.DirectDependencyList
			compName := params.CompList[0].CompName
			compVersion := params.CompList[0].CompVersion
			newParams.getFixModelList(fixParams, fixParams.pomPathList)
			if len(newParams.dmModelList) > 0 {
				for _, dm := range newParams.dmModelList {
					var (
						file *os.File
					)
					line := 1
					fixText := ""
					file, err = os.Open(dm.PomPath)
					defer file.Close()
					if err != nil {
						return
					}
					scanner := bufio.NewScanner(file)
					for scanner.Scan() {
						text := scanner.Text()

						if dm.Line == line {
							fixText = strings.ReplaceAll(text, compVersion, compVersion)
							break
						}
						line++
					}
					// 检查是否有任何错误
					if err = scanner.Err(); err != nil {
						return
					}

					if !params.ShowOnly {
						// 打开文件并读取所有内容
						file, err = os.Open(dm.PomPath)
						if err != nil {
							panic(err)
						}
						defer file.Close()

						scanner := bufio.NewScanner(file)
						var lines []string
						for scanner.Scan() {
							lines = append(lines, scanner.Text())
						}

						lines[dm.Line-1] = fixText
						newContent := []byte(fmt.Sprintf("%s\n", lines[0]))
						for _, line := range lines[1:] {
							newContent = append(newContent, []byte(fmt.Sprintf("%s\n", line))...)
						}

						err = os.WriteFile(dm.PomPath, newContent, 0644)
						if err != nil {
							return
						}
					}
				}
				return
			}
			if len(newParams.haveDmMap) > 0 {
				for pomPath, _ := range newParams.haveDmMap {
					var (
						file *os.File
					)
					dependenciesLine, ok := newParams.dependenciesLine[pomPath]
					if !ok {
						break
					}
					// 打开文件并读取所有内容
					file, err = os.Open(filepath.Join(params.Dir, pomPath))
					if err != nil {
						panic(err)
					}
					defer file.Close()

					scanner := bufio.NewScanner(file)
					var lines []string
					for scanner.Scan() {
						lines = append(lines, scanner.Text())
					}

					dependencies := lines[dependenciesLine]
					num := len(dependencies) - len(strings.TrimLeft(dependencies, " "))
					spaces := ""
					for i := 0; i < num; i++ {
						spaces += " "
					}
					split := strings.Split(compName, ":")
					line1 := spaces + "<dependency>"
					line2 := spaces + "    <groupId>" + split[0] + "</groupId>"
					line3 := spaces + "    <artifactId>" + split[1] + "</artifactId>"
					line4 := spaces + "    <version>" + compVersion + "</version>"
					line5 := spaces + "</dependency>"

					lines2 := make([]string, len(lines[dependenciesLine:]))
					lines3 := make([]string, len(lines[:dependenciesLine]))
					lines4 := make([]string, 0)
					lines5 := make([]string, 0)

					copy(lines2, lines[dependenciesLine:])
					copy(lines3, lines[:dependenciesLine])
					lines4 = append(lines3, line1, line2, line3, line4, line5)
					lines5 = append(lines4, lines2...)

					newContent := []byte(fmt.Sprintf("%s\n", lines5[0]))
					for _, line := range lines5[1:] {
						newContent = append(newContent, []byte(fmt.Sprintf("%s\n", line))...)
					}

					err = os.WriteFile(pomPath, newContent, 0644)
					if err != nil {
						return
					}
					return
				}

			}
			for pomPath, line := range newParams.dependenciesLine {
				var (
					file *os.File
				)

				// 打开文件并读取所有内容
				file, err = os.Open(filepath.Join(params.Dir, pomPath))
				if err != nil {
					panic(err)
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				var lines []string
				for scanner.Scan() {
					lines = append(lines, scanner.Text())
				}

				dependencies := lines[line-1]
				num := len(dependencies) - len(strings.TrimLeft(dependencies, " "))
				spaces := ""
				for i := 0; i < num; i++ {
					spaces += " "
				}
				split := strings.Split(compName, ":")
				line1 := spaces + "    <dependencyManagement>"
				line2 := spaces + "        <dependencies>"
				line3 := spaces + "            <dependency>"
				line4 := spaces + "                <groupId>" + split[0] + "</groupId>"
				line5 := spaces + "                <artifactId>" + split[1] + "</artifactId>"
				line6 := spaces + "                <version>" + compVersion + "</version>"
				line7 := spaces + "            </dependency>"
				line8 := spaces + "        </dependencies>"
				line9 := spaces + "    </dependencyManagement>"

				lines2 := make([]string, len(lines[line-1:]))
				lines3 := make([]string, len(lines[:line]))
				lines4 := make([]string, 0)
				lines5 := make([]string, 0)

				copy(lines2, lines[line:])
				copy(lines3, lines[:line])
				lines4 = append(lines3, line1, line2, line3, line4, line5, line6, line7, line8, line9)
				lines5 = append(lines4, lines2...)

				newContent := []byte(fmt.Sprintf("%s\n", lines5[0]))
				for _, line := range lines5[1:] {
					newContent = append(newContent, []byte(fmt.Sprintf("%s\n", line))...)
				}

				err = os.WriteFile(pomPath, newContent, 0644)
				if err != nil {
					return
				}
				return
			}

		} else {
			for _, model := range p.fixModelList {
				err = fileUpdate(model, item, p, params)
				if err != nil {
					return
				}
			}
		}

		for _, model := range p.dmModelList {
			err = dmPreview(model, p)
			if err != nil {
				return
			}
		}

	}

	return
}

func fileUpdate(model FixModel, comp Comp, p *mavenParams, params FixParams) (err error) {
	var (
		file *os.File
	)
	contentList := make([]Content, 0)
	line := 1
	fixText := ""
	file, err = os.Open(model.PomPath)
	defer file.Close()
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		contentList = append(contentList, Content{
			Line: line,
			Text: text,
		})
		if model.Line == line {
			fixText = strings.ReplaceAll(text, comp.CompVersion, comp.MinFixVersion)
			for i := 0; i < 5; i++ {
				line++
				if scanner.Scan() {
					text2 := scanner.Text()
					contentList = append(contentList, Content{
						Line: line,
						Text: text2,
					})
				}
			}
			if len(contentList) > 11 {
				contentList = contentList[len(contentList)-11:]
			}
			break
		}
		line++
	}
	// 检查是否有任何错误
	if err = scanner.Err(); err != nil {
		return
	}
	p.preview = append(p.preview, Preview{
		Path:    model.relativePomPath,
		Line:    model.Line,
		Content: contentList,
	})

	if !params.ShowOnly {
		// 打开文件并读取所有内容
		file, err = os.Open(model.PomPath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		lines[model.Line-1] = fixText
		newContent := []byte(fmt.Sprintf("%s\n", lines[0]))
		for _, line := range lines[1:] {
			newContent = append(newContent, []byte(fmt.Sprintf("%s\n", line))...)
		}

		err = os.WriteFile(model.PomPath, newContent, 0644)
		if err != nil {
			return
		}

	}

	return

}

func dmPreview(model FixModel, p *mavenParams) (err error) {
	var (
		file *os.File
	)
	contentList := make([]Content, 0)
	line := 1
	file, err = os.Open(model.PomPath)
	defer file.Close()
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		contentList = append(contentList, Content{
			Line: line,
			Text: text,
		})
		if model.Line == line {
			for i := 0; i < 5; i++ {
				line++
				if scanner.Scan() {
					text2 := scanner.Text()
					contentList = append(contentList, Content{
						Line: line,
						Text: text2,
					})
				}
			}
			if len(contentList) > 11 {
				contentList = contentList[len(contentList)-11:]
			}
			break
		}
		line++
	}
	// 检查是否有任何错误
	if err = scanner.Err(); err != nil {
		return
	}
	p.dmPreview = append(p.dmPreview, Preview{
		Path:    model.relativePomPath,
		Line:    model.Line,
		Content: contentList,
	})
	return

}

func (p *mavenParams) parsePropertyNode(params FixParams, pomPathList []string) {
	for _, pomPath := range pomPathList {
		joinPath := filepath.Join(params.Dir, pomPath)
		f, err := os.Open(joinPath)
		if err != nil {
			return
		}
		doc, err := xmlquery.Parse(f)
		if err != nil {
			return
		}
		properties := xmlquery.Find(doc, "//properties")
		versions := xmlquery.Find(doc, "//version")

		if len(properties) > 0 {
			node0 := properties[0]
			propertiesChilds := xmlquery.Find(node0, "child::*")
		ok:
			for _, item := range propertiesChilds {
				compName := make([]string, 0)
				for _, version := range versions {
					if "${"+item.Data+"}" != version.InnerText() {
						continue
					}
					parent := version.Parent
					groupId := parent.SelectElement("groupId")
					if groupId == nil {
						break ok
					}
					artifactId := parent.SelectElement("artifactId")
					if artifactId == nil {
						break ok
					}
					compName = append(compName, groupId.InnerText()+":"+artifactId.InnerText())
				}
				if len(compName) == 0 {
					continue
				}

				//parent := propertieComp.Parent
				model := PropertyModel{
					Line:       FindPropertiesLine(joinPath, item.Data, item.InnerText()),
					OldVersion: item.InnerText(),
					TagName:    item.Data,
					PomPath:    pomPath,
					CompName:   compName,
					//CompName: nil,
				}
				if models, ok := p.propertyMap[item.Data]; ok {
					models = append(models, model)
					p.propertyMap[item.Data] = models
				} else {

					models := make([]PropertyModel, 0)
					models = append(models, model)
					p.propertyMap[item.Data] = models
				}

			}

		}

	}

}

func (p *mavenParams) parseDM(params FixParams, pomPathList []string) {
	for _, pomPath := range pomPathList {
		joinPath := filepath.Join(params.Dir, pomPath)
		f, err := os.Open(joinPath)
		if err != nil {
			return
		}
		doc, err := xmlquery.Parse(f)
		if err != nil {
			return
		}
		dependencyManagement := xmlquery.Find(doc, "//dependencyManagement")
		versions := xmlquery.Find(doc, "//version")

		if len(dependencyManagement) > 0 {
			node0 := dependencyManagement[0]
			propertiesChilds := xmlquery.Find(node0, "child::*")
		ok:
			for _, item := range propertiesChilds {
				compName := make([]string, 0)
				for _, version := range versions {
					if "${"+item.Data+"}" != version.InnerText() {
						continue
					}
					parent := version.Parent
					groupId := parent.SelectElement("groupId")
					if groupId == nil {
						break ok
					}
					artifactId := parent.SelectElement("artifactId")
					if artifactId == nil {
						break ok
					}
					compName = append(compName, groupId.InnerText()+":"+artifactId.InnerText())
				}
				if len(compName) == 0 {
					continue
				}

				//parent := propertieComp.Parent
				model := PropertyModel{
					Line:       FindPropertiesLine(joinPath, item.Data, item.InnerText()),
					OldVersion: item.InnerText(),
					TagName:    item.Data,
					PomPath:    pomPath,
					CompName:   compName,
					//CompName: nil,
				}
				if models, ok := p.propertyMap[item.Data]; ok {
					models = append(models, model)
					p.propertyMap[item.Data] = models
				} else {

					models := make([]PropertyModel, 0)
					models = append(models, model)
					p.propertyMap[item.Data] = models
				}

			}

		}

	}

}

func (p *mavenParams) getFixModelList(params FixParams, pomPathList []string) {
	for _, comp := range params.CompList {
		for _, pomPath := range pomPathList {
			list1, list2, haveDM, dependenciesLine := GetFixModelList(filepath.Join(params.Dir, pomPath), pomPath, comp.CompName, comp.CompVersion, comp.MinFixVersion, p.propertyMap)
			if len(list1) > 0 {
				p.fixModelList = append(p.fixModelList, list1...)
			}
			if len(list2) > 0 {
				p.dmModelList = append(p.dmModelList, list2...)
			}
			if haveDM > 0 {
				p.haveDmMap[pomPath] = haveDM
			}
			if dependenciesLine > 0 {
				p.dependenciesLine[pomPath] = dependenciesLine
			}

		}
	}

}

func (p *mavenParams) getExtensionFix(params FixParams, pomPathList []string) {
	p.fixModelList = make([]FixModel, 0)
	for _, comp := range params.CompList {
		split := strings.Split(comp.CompName, ":")
		if len(split) != 2 {
			return
		}
		if split[0] != "org.springframework.boot" {
			return

		}
		for _, pomPath := range pomPathList {
			list := GetExtensionFixModelList(filepath.Join(params.Dir, pomPath), pomPath, comp.CompName, comp.CompVersion, comp.MinFixVersion, p.propertyMap)
			p.fixModelList = append(p.fixModelList, list...)

		}
	}

}

func (p *mavenParams) getInheritFix(params FixParams, pomPathList []string) {
	p.fixModelList = make([]FixModel, 0)
	for _, comp := range params.CompList {
		split := strings.Split(comp.CompName, ":")
		if len(split) != 2 {
			return
		}
		if split[0] != "org.springframework.boot" {
			return

		}
		for _, pomPath := range pomPathList {
			list := GetInheritFixModelList(filepath.Join(params.Dir, pomPath), pomPath, comp.CompName, comp.CompVersion, comp.MinFixVersion, p.propertyMap)
			p.fixModelList = append(p.fixModelList, list...)
		}
	}

}
