package fix

import (
	"bufio"
	"io/fs"
	"os"
	"strings"
)

type propertyData struct {
	Path             string
	AbsoluteFilePath string
	Tag              string
	Value            string
	Line             int
}

type dependencyManagementData struct {
	Path             string
	AbsoluteFilePath string
	GroupId          string
	ArtifactId       string
	Version          string
	Line             int
}

func (t FixParams) MavenFixNew() (preview []Preview, dmPreview []Preview, haveDMList map[string]int, err error) {

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
	result := Cli(t.Dir, t.pomPathList)

	propertyMap := make(map[string][]propertyData)
	dependencyManagementMap := make(map[string][]dependencyManagementData)
	for _, data := range result {
		for _, property := range data.Properties {
			val, ok := propertyMap[property.Tag]
			if !ok {
				val = make([]propertyData, 0)
			}
			val = append(val, propertyData{
				Path:             data.FilePath,
				AbsoluteFilePath: data.AbsoluteFilePath,
				Tag:              property.Tag,
				Value:            property.Value,
				Line:             property.ValueLine,
			})
			propertyMap[property.Tag] = val
		}
		for _, da := range data.DependencyManagement {
			val, ok := dependencyManagementMap[da.GroupId+":"+da.ArtifactId]
			if !ok {
				val = make([]dependencyManagementData, 0)
			}
			val = append(val, dependencyManagementData{
				Path:             data.FilePath,
				AbsoluteFilePath: data.AbsoluteFilePath,
				GroupId:          da.GroupId,
				ArtifactId:       da.ArtifactId,
				Version:          da.Version,
				Line:             da.VersionLine,
			})
			dependencyManagementMap[da.GroupId+":"+da.ArtifactId] = val
		}
	}

	for _, comp := range t.CompList {
		val, ok := dependencyManagementMap[comp.CompName]
		if ok {
			for _, x := range val {
				if comp.CompVersion == x.Version {
					preview = append(preview, Preview{
						CompName:    comp.CompName,
						CompVersion: comp.CompVersion,
						Path:        x.Path,
						Line:        x.Line,
						Content:     findContent(x.AbsoluteFilePath, x.Line),
					})
				} else {

					val2, ok := propertyMap[x.Version]
					if !ok {
						version := strings.ReplaceAll(x.Version, "${", "")
						version = strings.ReplaceAll(version, "}", "")
						val2, ok = propertyMap[version]
					}
					if ok {
						for _, item := range val2 {
							preview = append(preview, Preview{
								CompName:    comp.CompName,
								CompVersion: comp.CompVersion,
								Path:        item.Path,
								Line:        item.Line,
								Content:     findContent(item.AbsoluteFilePath, item.Line),
							})

						}

					}
				}

			}
			continue
		}
		for _, data := range result {
			for _, dependency := range data.Dependencies {
				if comp.CompName == dependency.GroupId+":"+dependency.ArtifactId {
					if comp.CompVersion == dependency.Version {
						preview = append(preview, Preview{
							CompName:    comp.CompName,
							CompVersion: comp.CompVersion,
							Path:        data.FilePath,
							Line:        dependency.VersionLine,
							Content:     findContent(data.AbsoluteFilePath, dependency.VersionLine),
						})
					} else {

						val, ok := propertyMap[dependency.Version]
						if !ok {
							version := strings.ReplaceAll(dependency.Version, "${", "")
							version = strings.ReplaceAll(version, "}", "")
							val, ok = propertyMap[version]
						}
						if ok {
							for _, x := range val {
								if x.Tag == comp.CompVersion {
									preview = append(preview, Preview{
										CompName:    comp.CompName,
										CompVersion: comp.CompVersion,
										Path:        x.Path,
										Line:        x.Line,
										Content:     findContent(x.AbsoluteFilePath, x.Line),
									})
								}

							}

						}
					}

				}
			}

		}

	}

	return
}

func findContent(path string, fileLine int) (contentList []Content) {
	var (
		file *os.File
		err  error
	)
	file, err = os.Open(path)
	defer file.Close()
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(file)
	line := 1
	for scanner.Scan() {
		text := scanner.Text()
		contentList = append(contentList, Content{
			Line: line,
			Text: text,
		})
		if fileLine == line {
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
	return
}
