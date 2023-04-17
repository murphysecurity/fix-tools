package fix

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
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
				var writeFile *os.File

				//读写方式打开文件
				writeFile, err = os.OpenFile(filePath, os.O_RDWR, 0666)
				if err != nil {
					fmt.Println("open file filed.", err)
					return
				}
				beforeMD5 := md5.New()
				if _, err = io.Copy(beforeMD5, writeFile); err != nil {
					return
				}
				befroeHash := beforeMD5.Sum(nil)

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
					if lineNum == fixLine {
						bytes := []byte(fixText + "\n")
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
				afterHash := afterMD5.Sum(nil)
				if string(befroeHash) == string(afterHash) {

				}

				writeFile.Close()
			}

		}

	}

	return
}
