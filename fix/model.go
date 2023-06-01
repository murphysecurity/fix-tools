package fix

import (
	"errors"
	"strconv"
	"time"
)

type FixParams struct {
	// 必填
	ShowOnly       bool          // 仅展示pr预览 不进行真的pr提交
	TimeOut        time.Duration //超时时间 单位秒 默认60秒
	RepoType       string        // 仓库类型 github gitee gitlab local
	CompList       []Comp
	PackageManager string // 包管理器
	// local必填
	Dir string // 检测类型中需要指定文件夹

	// 非local必填
	ProxyUrl    string // 可能需要的代理
	GitlabUrl   string //   Gitlab地址
	CommitHash  string //   提交的hash
	TargetOwner string // 目标 仓库拥有者
	Owner       string //   当前用户 仓库拥有者
	Repo        string //   仓库路径 不变

	// 仓库类型 用户相关
	UserName  string //   设置git用户名
	Password  string
	UserEmail string //   设置git用户邮箱
	Token     string // GitHub用户token Gitee用户AccessToken  Gitlab用户token

	// pr提交相关
	Title string // pr 标题
	Body  string // pr 内容

	// 内部使用不暴露
	branch    string
	defBranch string
	delDir    string
}

type Comp struct {
	CompName      string `json:"comp_name"`       // 组件名称
	CompVersion   string `json:"comp_version"`    // 组件版本
	MinFixVersion string `json:"min_fix_version"` // 最小修复版本
}

func (t *FixParams) check() error {

	switch t.PackageManager {
	case "maven":
	case "go":
	case "npm":
	case "yarn":
	case "python":
	default:
		return errors.New("不支持的包管理器: " + t.PackageManager)
	}

	switch t.RepoType {
	case "github":
	case "gitee":
	case "gitlab":
		if len(t.GitlabUrl) == 0 {
			return errors.New("gitlab检测请指定路径GitlabUrl")
		}
	case "local":
		if len(t.Dir) == 0 {
			return errors.New("本地检测请指定路径")
		}

	default:
		return errors.New("不支持的仓库类型: " + t.RepoType)
	}
	if t.TimeOut == 0 {
		t.TimeOut = 60 * time.Second
	} else {
		t.TimeOut = t.TimeOut * time.Second
	}
	t.branch = "fix_" + strconv.FormatInt(time.Now().Unix(), 10)
	return nil
}

type mavenParams struct {
	propertyMap  map[string][]PropertyModel
	fixModelList []FixModel
	preview      []Preview
}

type Preview struct {
	Path    string    `json:"path"`
	Line    int       `json:"line"`
	Content []Content `json:"content"`
}

type Content struct {
	Line int    `json:"line"`
	Text string `json:"text"`
}

type FixModel struct {
	Line            int
	OldVersion      string
	NewVersion      string
	GroupId         string
	ArtifactId      string
	CompName        string
	PomPath         string
	relativePomPath string
}

type PropertyModel struct {
	Line       int
	OldVersion string
	TagName    string
	PomPath    string
	CompName []string
}
