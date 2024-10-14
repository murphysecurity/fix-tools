package fix

import (
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"github.com/murphysecurity/fix-tools/fix/xml/parser"
	"path/filepath"
)

type CliXMLListener struct {
	parser.BaseXMLParserListener
	FilePath             string
	AbsoluteFilePath     string
	DependencyManagement []Dependency
	Dependencies         []Dependency
	Properties           []Properties
}
type Dependency struct {
	GroupId     string
	ArtifactId  string
	Version     string
	VersionLine int
}
type Properties struct {
	Tag       string
	Value     string
	ValueLine int
}

func (l *CliXMLListener) EnterElement(ctx *parser.ElementContext) {
	name := ctx.Name(0).GetText()
	dependency := Dependency{}

	if name == "properties" {
		elements := ctx.Content().AllElement()
		for _, element := range elements {
			text := element.Name(0).GetText()
			content := element.Content()
			if content == nil {
				continue
			}
			value := content.GetText()
			ValueLine := element.Name(0).GetSymbol().GetLine()

			l.Properties = append(l.Properties, Properties{
				Tag:       text,
				Value:     value,
				ValueLine: ValueLine,
			})

		}
	}

	if name == "dependency" {
		dependencyParent1 := ctx.GetParent()
		if dependencyParent1 == nil {
			return
		}
		dependencyParent2 := dependencyParent1.GetParent()
		if dependencyParent2 == nil {
			return
		}
		dependencyParent2Ctx, ok := dependencyParent2.(*parser.ElementContext)
		if !ok {
			return
		}
		dependencyParent2CtxName := dependencyParent2Ctx.Name(0).GetText()
		if dependencyParent2CtxName != "dependencies" {
			return
		}
		dependencyParent3 := dependencyParent2Ctx.GetParent()
		if dependencyParent3 == nil {
			return
		}
		dependencyParent4 := dependencyParent3.GetParent()
		if dependencyParent4 == nil {
			return
		}
		dependencyParent4Ctx, ok := dependencyParent4.(*parser.ElementContext)
		if !ok {
			return
		}
		dependencyParent4CtxName := dependencyParent4Ctx.Name(0).GetText()

		if ctx.Content() != nil && len(ctx.Content().AllElement()) > 0 {
			elements := ctx.Content().AllElement()
			for _, element := range elements {
				text := element.Name(0).GetText()
				content := element.Content()
				if content == nil {
					continue
				}
				value := content.GetText()
				if text == "groupId" {
					dependency.GroupId = value
				}
				if text == "artifactId" {
					dependency.ArtifactId = value
				}
				if text == "version" {
					dependency.Version = value
					dependency.VersionLine = element.Name(0).GetSymbol().GetLine()
				}
			}
		}
		if dependency.VersionLine != 0 && dependency.GroupId != "" && dependency.ArtifactId != "" {
			if dependencyParent4CtxName == "dependencyManagement" {
				l.DependencyManagement = append(l.DependencyManagement, dependency)
			} else {
				l.Dependencies = append(l.Dependencies, dependency)
			}
		}
	}

}

func Cli(dir string, pomPathList []string) (result []CliXMLListener) {
	for _, pomPath := range pomPathList {
		absoluteFilePath := filepath.Join(dir, pomPath)
		input, _ := antlr.NewFileStream(absoluteFilePath)
		lexer := parser.NewXMLLexer(input)
		liste := &MyErrorListener{suppress: true}
		lexer.RemoveErrorListeners()
		lexer.AddErrorListener(liste)
		stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
		p := parser.NewXMLParser(stream)
		p.RemoveErrorListeners()
		p.AddErrorListener(liste)
		listener := &CliXMLListener{
			FilePath:         pomPath,
			AbsoluteFilePath: absoluteFilePath,
		}
		antlr.ParseTreeWalkerDefault.Walk(listener, p.Document())
		result = append(result, *listener)
	}
	return

}
