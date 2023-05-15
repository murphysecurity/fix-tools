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

func (t FixParams) MavenFix() (preview []Preview, err error) {

	var (
		pomPathList = make([]string, 0)
	)

	fileSystem := os.DirFS(t.Dir)

	err = fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if d.Name() == "pom.xml" {
			pomPathList = append(pomPathList, path)
		}
		return nil
	})
	if err != nil {
		return
	}
	params := mavenParams{
		propertyMap:  make(map[string][]PropertyModel, 0),
		fixModelList: make([]FixModel, 0),
		preview:      make([]Preview, 0),
	}
	params.parsePropertyNode(t, pomPathList)
	params.getFixModelList(t, pomPathList)
	params.modifyPom(t)

	if len(params.preview) == 0 {
		params.getExtensionFix(t, pomPathList)
		params.modifyPom(t)
	}
	if len(params.preview) == 0 {
		params.getInheritFix(t, pomPathList)
		params.modifyPom(t)
	}
	preview = params.preview
	return
}

func (p *mavenParams) modifyPom(params FixParams) (err error) {

	for _, comp := range params.CompList {
		for _, model := range p.fixModelList {
			err = fileUpdate(model, comp, p, params)
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

func (p *mavenParams) parsePropertyNode(params FixParams, pomPathList []string) {
	for _, pomPath := range pomPathList {
		joinPath := filepath.Join(params.Dir, pomPath)
		f, err := os.Open(joinPath)
		if err != nil {
			return
		}
		doc, err := xmlquery.Parse(f)
		properties := xmlquery.Find(doc, "//properties")
		if len(properties) > 0 {
			node0 := properties[0]
			propertiesChilds := xmlquery.Find(node0, "child::*")
			for _, item := range propertiesChilds {
				model := PropertyModel{
					Line:       FindPropertiesLine(joinPath, item.Data, item.InnerText()),
					OldVersion: item.InnerText(),
					TagName:    item.Data,
					PomPath:    pomPath,
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
			list := GetFixModelList(filepath.Join(params.Dir, pomPath), pomPath, comp.CompName, comp.CompVersion, comp.MinFixVersion, p.propertyMap)
			p.fixModelList = append(p.fixModelList, list...)

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
