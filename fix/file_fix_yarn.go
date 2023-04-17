package fix

import (
	"bufio"
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (t FixParams) YarnFix() (preview []Preview, err error) {
	var (
		file                *os.File
		packagejsonPathList = make([]string, 0)
	)
	// 修改了的comp组件名称
	if t.Dir == "" || t.CompName == "" || t.CompVersion == "" || t.MinFixVersion == "" {
		err = errors.New("项目目录|组件名|组件版本|最小修复版本不能为空")
		return
	}

	fileSystem := os.DirFS(t.Dir)

	err = fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if d.Name() == "package.json" {
			packagejsonPathList = append(packagejsonPathList, path)
		}
		return nil
	})
	if err != nil {
		return
	}

	for _, filePath := range packagejsonPathList {
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
			if strings.Contains(text, t.CompName) && strings.Contains(text, t.CompVersion) {
				if strings.Contains(text, "*"+t.CompName) {
					continue
				}
				if strings.Contains(text, "^"+t.CompName) {
					if strings.Split(text, ".")[0] == strings.Split(t.MinFixVersion, ".")[0] {
						continue
					}
				}
				if strings.Contains(text, "~"+t.CompName) {
					if strings.Split(text, ".")[0] == strings.Split(t.MinFixVersion, ".")[0] &&
						strings.Split(text, ".")[1] == strings.Split(t.MinFixVersion, ".")[1] {
						continue
					}
				}
				fixText = strings.ReplaceAll(text, t.CompVersion, t.MinFixVersion)
				if fixText != text {
					fixLine = line
					isFix = true
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

			package_lockPath := strings.ReplaceAll(filePath, "package.json", "yarn.lock")
			exists := Exists(package_lockPath)
			if exists {
				ctx := context.Background()
				cmd := exec.CommandContext(ctx, "yarn", "install")
				dir, _ := filepath.Split(filePath)
				cmd.Dir = filepath.Join(t.Dir, dir)

				var stdout, stderr bytes.Buffer
				cmd.Stdout = &stdout
				cmd.Stderr = &stderr

				err = cmd.Run()
				if err != nil {
					err = errors.New(fmt.Sprintf("执行失败 err: %s stdout: %s stderr %s", err.Error(), stdout.String(), stderr.String()))
					return
				}
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

	return
}
