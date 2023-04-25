package fix

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"github.com/murphysecurity/fix-tools/fix/xml/parser"
	"io"
	"net/http"
	"strings"
)

type SimpleXMLListener struct {
	parser.BaseXMLParserListener
	targetElement string
	value         string
	line          int
}

func (l *SimpleXMLListener) ExitElement(ctx *parser.ElementContext) {
	if ctx.Content() != nil {
		if ctx.Name(0).GetText() == l.targetElement && ctx.Content().GetText() == l.value {
			l.line = ctx.Name(0).GetSymbol().GetLine()
		}
	}

}

func FindPropertiesLine(pomPath, targetElement, value string) int {

	input, _ := antlr.NewFileStream(pomPath)
	lexer := parser.NewXMLLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewXMLParser(stream)
	listener := &SimpleXMLListener{
		targetElement: targetElement,
		value:         value,
	}
	antlr.ParseTreeWalkerDefault.Walk(listener, p.Document())
	return listener.line
}

type ChildXMLListener struct {
	parser.BaseXMLParserListener
	pomPath         string
	relativePomPath string
	compName        string
	newVersion      string
	compVersion     string
	modelMap        map[string][]PropertyModel
	fixModelList    []FixModel
}

func (l *ChildXMLListener) EnterElement(ctx *parser.ElementContext) {
	name := ctx.Name(0).GetText()
	if name == "dependency" {

		model := FixModel{
			PomPath:         l.pomPath,
			relativePomPath: l.relativePomPath,
		}
		if ctx.Content() != nil && len(ctx.Content().AllElement()) > 0 {
			elements := ctx.Content().AllElement()
			for _, element := range elements {
				text := element.Name(0).GetText()
				value := element.Content().GetText()
				if text == "groupId" {
					model.GroupId = value
				}
				if text == "artifactId" {
					model.ArtifactId = value
				}
				if text == "version" {
					model.OldVersion = value
					model.Line = element.Name(0).GetSymbol().GetLine()
				}
			}
		}
		if model.Line != 0 && model.GroupId != "" && model.ArtifactId != "" {
			if model.GroupId+":"+model.ArtifactId == l.compName {
				if strings.Contains(model.OldVersion, "${") && strings.Contains(model.OldVersion, "}") {
					model.OldVersion = strings.ReplaceAll(model.OldVersion, "${", "")
					model.OldVersion = strings.ReplaceAll(model.OldVersion, "}", "")

					if propertyModel, ok := l.modelMap[model.OldVersion]; ok {
						for _, m := range propertyModel {
							newModel := FixModel{
								Line:            m.Line,
								OldVersion:      model.OldVersion,
								NewVersion:      l.newVersion,
								CompName:        l.compName,
								PomPath:         l.pomPath,
								relativePomPath: l.relativePomPath,
							}
							l.fixModelList = append(l.fixModelList, newModel)
						}
					}
				} else {
					model.CompName = l.compName
					model.NewVersion = l.newVersion
					l.fixModelList = append(l.fixModelList, model)
				}

			}
		}
	}

}

func GetFixModelList(pomPath, relativePomPath, compName, compVersion, newVersion string, model map[string][]PropertyModel) []FixModel {

	input, _ := antlr.NewFileStream(pomPath)
	lexer := parser.NewXMLLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewXMLParser(stream)
	listener := &ChildXMLListener{
		pomPath:         pomPath,
		relativePomPath: relativePomPath,
		compName:        compName,
		compVersion:     compVersion,
		newVersion:      newVersion,
		modelMap:        model,
		fixModelList:    nil,
	}
	antlr.ParseTreeWalkerDefault.Walk(listener, p.Document())
	return listener.fixModelList
}

type ParentXMLListener struct {
	parser.BaseXMLParserListener
	pomPath         string
	relativePomPath string
	compName        string
	newVersion      string
	compVersion     string
	modelMap        map[string][]PropertyModel
	fixModelList    []FixModel
}

func (l *ParentXMLListener) EnterElement(ctx *parser.ElementContext) {
	name := ctx.Name(0).GetText()
	if name == "parent" {

		model := FixModel{
			PomPath: l.relativePomPath,
		}
		if ctx.Content() != nil && len(ctx.Content().AllElement()) > 0 {
			elements := ctx.Content().AllElement()
			for _, element := range elements {
				text := element.Name(0).GetText()
				value := element.Content().GetText()
				if text == "groupId" {
					model.GroupId = value
				}
				if text == "artifactId" {
					model.ArtifactId = value
				}
				if text == "version" {
					model.OldVersion = value
					model.Line = element.Name(0).GetSymbol().GetLine()
				}
			}
		}
		if model.Line != 0 && model.GroupId != "" && model.ArtifactId != "" {
			if model.GroupId+":"+model.ArtifactId == l.compName {
				if strings.Contains(model.OldVersion, "${") && strings.Contains(model.OldVersion, "}") {
					model.OldVersion = strings.ReplaceAll(model.OldVersion, "${", "")
					model.OldVersion = strings.ReplaceAll(model.OldVersion, "}", "")

					if propertyModel, ok := l.modelMap[model.OldVersion]; ok {
						for _, m := range propertyModel {
							newModel := FixModel{
								Line:            m.Line,
								OldVersion:      model.OldVersion,
								NewVersion:      l.newVersion,
								CompName:        l.compName,
								PomPath:         l.pomPath,
								relativePomPath: l.relativePomPath,
							}
							l.fixModelList = append(l.fixModelList, newModel)
						}
					}
				} else {
					model.CompName = l.compName
					model.NewVersion = l.newVersion
					l.fixModelList = append(l.fixModelList, model)
				}

			}
		}
	}

}

func GetExtensionFixModelList(pomPath, relativePomPath, compName, compVersion, newVersion string, model map[string][]PropertyModel) []FixModel {

	input, _ := antlr.NewFileStream(pomPath)
	lexer := parser.NewXMLLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewXMLParser(stream)
	listener := &ChildXMLListener{
		pomPath:         pomPath,
		relativePomPath: relativePomPath,
		compName:        compName,
		compVersion:     compVersion,
		newVersion:      newVersion,
		modelMap:        model,
		fixModelList:    nil,
	}
	antlr.ParseTreeWalkerDefault.Walk(listener, p.Document())
	return listener.fixModelList
}

type InheritXMLListener struct {
	parser.BaseXMLParserListener
	relativePomPath string
	pomPath         string
	compName        string
	newVersion      string
	compVersion     string
	projectPathMap  map[string]string
}

func (l *InheritXMLListener) EnterElement(ctx *parser.ElementContext) {
	//name := ctx.Name(0).GetText()

	if ctx.GetParent() == nil {

		model := FixModel{
			PomPath: l.relativePomPath,
		}
		if ctx.Content() != nil && len(ctx.Content().AllElement()) > 0 {
			elements := ctx.Content().AllElement()
			for _, element := range elements {
				text := element.Name(0).GetText()
				if element.Content() != nil {
					value := element.Content().GetText()
					if text == "groupId" {
						model.GroupId = value
					}
					if text == "artifactId" {
						model.ArtifactId = value
					}
					if text == "version" {
						model.OldVersion = value
						model.Line = element.Name(0).GetSymbol().GetLine()
					}
				}

			}
		}
		if model.Line != 0 && model.GroupId != "" && model.ArtifactId != "" && model.OldVersion != "" {
			l.projectPathMap[fmt.Sprintf("%s:%s@%s", model.GroupId, model.ArtifactId, model.OldVersion)] = l.pomPath
		}
	}

}

type InheritParentXMLListener struct {
	parser.BaseXMLParserListener
	pomPath         string
	relativePomPath string
	compName        string
	newVersion      string
	compVersion     string
	compNameVarMap  map[string]string
	fixModelList    []FixModel
	projectPathMap  map[string]string
}

func (l *InheritParentXMLListener) EnterElement(ctx *parser.ElementContext) {
	name := ctx.Name(0).GetText()
	if name == "parent" {

		model := FixModel{
			PomPath: l.relativePomPath,
		}
		if ctx.Content() != nil && len(ctx.Content().AllElement()) > 0 {
			elements := ctx.Content().AllElement()
			for _, element := range elements {
				text := element.Name(0).GetText()
				if element.Content() != nil {
					value := element.Content().GetText()
					if text == "groupId" {
						model.GroupId = value
					}
					if text == "artifactId" {
						model.ArtifactId = value
					}
					if text == "version" {
						model.OldVersion = value
						model.Line = element.Name(0).GetSymbol().GetLine()
					}
				}

			}
		}
		if model.Line != 0 && model.GroupId != "" && model.ArtifactId != "" && model.OldVersion != "" {
			sprintf := fmt.Sprintf("%s:%s@%s", model.GroupId, model.ArtifactId, model.OldVersion)
			if _, ok := l.projectPathMap[sprintf]; ok {

				compNameFormat := strings.ReplaceAll(model.GroupId+"/"+model.ArtifactId, ".", "/")
				resp, err := http.Get(fmt.Sprintf("https://repo1.maven.org/maven2/%s/%s/%s-%s.pom", compNameFormat, model.OldVersion, model.ArtifactId, model.OldVersion))
				if err != nil {
					return
				}
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return
				}
				input := antlr.NewInputStream(string(body))
				lexer := parser.NewXMLLexer(input)
				stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
				p := parser.NewXMLParser(stream)
				listener := &InheritParentBodyXMLListener{}
				antlr.ParseTreeWalkerDefault.Walk(listener, p.Document())
				for key, val := range listener.compNameVarMap {
					l.compNameVarMap[key] = val
				}
			}
		}
	}

	if name == "dependency" {

		model := FixModel{
			PomPath: l.relativePomPath,
		}
		if ctx.Content() != nil && len(ctx.Content().AllElement()) > 0 {
			elements := ctx.Content().AllElement()
			for _, element := range elements {
				text := element.Name(0).GetText()
				if element.Content() != nil {
					value := element.Content().GetText()
					if text == "groupId" {
						model.GroupId = value
					}
					if text == "artifactId" {
						model.ArtifactId = value
					}
					if text == "version" {
						model.OldVersion = value
						model.Line = element.Name(0).GetSymbol().GetLine()
					}
				}

			}
		}

		if model.OldVersion != "" && model.GroupId != "" && model.ArtifactId != "" {
			if strings.Contains(model.OldVersion, "${") {
				model.OldVersion = strings.ReplaceAll(model.OldVersion, "${", "")
				model.OldVersion = strings.ReplaceAll(model.OldVersion, "}", "")
				l.compNameVarMap[model.GroupId+":"+model.ArtifactId] = model.OldVersion
			}
		}
	}

}

type InheritParentBodyXMLListener struct {
	parser.BaseXMLParserListener
	compNameVarMap map[string]string
}

func (l *InheritParentBodyXMLListener) EnterElement(ctx *parser.ElementContext) {
	name := ctx.Name(0).GetText()

	if name == "dependency" {

		model := FixModel{}
		if ctx.Content() != nil && len(ctx.Content().AllElement()) > 0 {
			elements := ctx.Content().AllElement()
			for _, element := range elements {
				text := element.Name(0).GetText()
				if element.Content() != nil {
					value := element.Content().GetText()
					if text == "groupId" {
						model.GroupId = value
					}
					if text == "artifactId" {
						model.ArtifactId = value
					}
					if text == "version" {
						model.OldVersion = value
						model.Line = element.Name(0).GetSymbol().GetLine()
					}
				}

			}
		}

		if model.OldVersion != "" && model.GroupId != "" && model.ArtifactId != "" {
			if strings.Contains(model.OldVersion, "${") {
				model.OldVersion = strings.ReplaceAll(model.OldVersion, "${", "")
				model.OldVersion = strings.ReplaceAll(model.OldVersion, "}", "")
				l.compNameVarMap[model.GroupId+":"+model.ArtifactId] = model.OldVersion
			}
		}
	}

}

func GetInheritFixModelList(pomPath, relativePomPath, compName, compVersion, newVersion string, model map[string][]PropertyModel) []FixModel {
	models := make([]FixModel, 0)
	input, _ := antlr.NewFileStream(pomPath)
	lexer := parser.NewXMLLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewXMLParser(stream)
	listener := &InheritParentXMLListener{
		pomPath:         pomPath,
		relativePomPath: relativePomPath,
		compName:        compName,
		compVersion:     compVersion,
		newVersion:      newVersion,
		fixModelList:    nil,
		projectPathMap:  make(map[string]string, 0),
		compNameVarMap:  make(map[string]string, 0),
	}
	antlr.ParseTreeWalkerDefault.Walk(listener, p.Document())
	if varName, ok := listener.compNameVarMap[compName]; ok {
		if len(varName) != 0 {
			if propertyModels, ok := model[varName]; ok {
				for _, propertyModel := range propertyModels {
					models = append(models, FixModel{
						Line:       propertyModel.Line,
						OldVersion: compVersion,
						NewVersion: newVersion,
						CompName:   compName,
						PomPath:    propertyModel.PomPath,
					})
				}

			}
		}
	}

	return models
}
