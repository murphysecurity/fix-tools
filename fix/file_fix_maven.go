package fix

import (
	"bufio"
	"crypto/md5"

	"github.com/antchfx/xmlquery"
	"io"
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
			pomPathList = append(pomPathList, d.Name())
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

	if len(preview) == 0 {
		params.getExtensionFix(t, pomPathList)
		params.modifyPom(t)
	}
	if len(preview) == 0 {
		params.getInheritFix(t, pomPathList)
		params.modifyPom(t)
	}
	preview = params.preview
	return
}

func (p *mavenParams) modifyPom(params FixParams) (err error) {

	for _, comp := range params.CompList {
		for _, model := range p.fixModelList {
			var (
				file *os.File
			)
			contentList := make([]Content, 0)
			line := 1
			fixText := ""
			file, err = os.Open(model.PomPath)
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
			file.Close()

			if !params.ShowOnly {
				var writeFile *os.File

				//读写方式打开文件
				writeFile, err = os.OpenFile(model.PomPath, os.O_RDWR, 0666)
				if err != nil {
					return
				}
				//beforeMD5 := md5.New()
				//if _, err = io.Copy(beforeMD5, writeFile); err != nil {
				//	return
				//}
				//befroeHash := beforeMD5.Sum(nil)

				//读取文件内容到io中
				reader := bufio.NewReader(writeFile)
				lineNum := 1
				pos := int64(0)
				for {
					//读取每一行内容
					lineText, readerErr := reader.ReadString('\n')
					if readerErr != nil {
						//读到末尾
						if readerErr == io.EOF {
							break
						} else {
							err = readerErr
							return
						}
					}
					if lineNum == model.Line {
						bytes := make([]byte, 0)
						//if strings.Contains(fixText, "\n") {
						//	bytes = []byte(fixText)
						//} else {
						//	bytes = []byte(fixText + "\n")
						//}
						bytes = []byte(fixText)
						_, readerErr = writeFile.WriteAt(bytes, pos)
						if readerErr != nil {
							err = readerErr
							return
						}
						break
					}

					//每一行读取完后记录位置
					pos += int64(len(lineText))
					lineNum++
				}

				afterMD5 := md5.New()
				if _, err = io.Copy(afterMD5, file); err != nil {
					return
				}
				//afterHash := afterMD5.Sum(nil)
				//if string(befroeHash) == string(afterHash) {
				//
				//}

				writeFile.Close()
			}

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
