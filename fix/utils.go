package fix

import (
	"encoding/xml"
	"io"
	"os"
)

func checkNilModPosition(t any) []ModPosition {
	if t != nil {
		if v, ok := t.([]ModPosition); ok {
			return v
		}
	}
	var mod []ModPosition
	return mod

}
func scanContent(lines []string, line int) []Content {

	var contents []Content
	var befer, end int

	if line-5 >= 1 {
		befer = line - 5
	} else {
		befer = 1
	}
	if line+5 <= len(lines) {
		end = line + 5
	} else {
		end = len(lines)
	}
	for i := befer; i <= end; i++ {
		var content Content
		content.Line = i
		content.Text = lines[i-1]
		contents = append(contents, content)
	}

	return contents
}

func UnmarshalXML(filePath string) (project *projectXML) {
	project = new(projectXML)
	// 打开pom.xml文件
	file, err := os.Open(filePath)
	if err != nil {
		return
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return
	}

	err = xml.Unmarshal(data, &project)
	if err != nil {
		return
	}
	return
}
