package fix

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func (t FixParams) PythonFix() (preview []Preview, err error) {
	var (
		file                 *os.File
		requirementsPathList = make([]string, 0)
	)

	fileSystem := os.DirFS(t.Dir)

	err = fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if d.Name() == "requirements.txt" {
			requirementsPathList = append(requirementsPathList, path)
		}
		return nil
	})
	if err != nil {
		return
	}

	for _, comp := range t.CompList {
		for _, filePath := range requirementsPathList {
			contentList := make([]Content, 0)
			isFix := false
			line := 1
			fixLine := 0
			fixText := ""
			file, err = os.Open(filepath.Join(t.Dir, filePath))
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
				if strings.Contains(text, comp.CompName) && strings.Contains(text, comp.CompVersion) {
					fixLine = line
					isFix = true
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

			if isFix {
				preview = append(preview, Preview{
					Path:    filePath,
					Line:    fixLine,
					Content: contentList,
				})

			}
			file.Close()

			if !t.ShowOnly {
				// 打开文件并读取所有内容
				file, err = os.Open(filePath)
				if err != nil {
					panic(err)
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				var lines []string
				for scanner.Scan() {
					lines = append(lines, scanner.Text())
				}

				lines[fixLine-1] = fixText
				newContent := []byte(fmt.Sprintf("%s\n", lines[0]))
				for _, line := range lines[1:] {
					newContent = append(newContent, []byte(fmt.Sprintf("%s\n", line))...)
				}

				err = os.WriteFile(filePath, newContent, 0644)
				if err != nil {
					return
				}

			}

		}

	}

	return
}
