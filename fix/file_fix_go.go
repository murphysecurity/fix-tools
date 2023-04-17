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

func (t FixParams) GoFix() (preview []Preview, err error) {

	var (
		file        *os.File
		modPathList = make([]string, 0)
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
		if d.Name() == "go.mod" {
			modPathList = append(modPathList, path)
		}
		return nil
	})
	if err != nil {
		return
	}
	for _, filePath := range modPathList {
		contentList := make([]Content, 0)
		isFix := false
		line := 1
		fixLine := 0
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
				isFix = true
				fixLine = line
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
			if !t.ShowOnly {
				beforeMD5 := md5.New()
				if _, err = io.Copy(beforeMD5, file); err != nil {
					return
				}
				befroeHash := beforeMD5.Sum(nil)

				ctx := context.Background()
				cmd := exec.CommandContext(ctx, "go", "get", fmt.Sprintf("%s@v%s", t.CompName, t.MinFixVersion))
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

				afterMD5 := md5.New()
				if _, err = io.Copy(afterMD5, file); err != nil {
					return
				}
				afterHash := afterMD5.Sum(nil)
				if string(befroeHash) == string(afterHash) {

				}
			}

		}
		file.Close()

	}

	return
}
