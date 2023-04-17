package fix

import (
	"context"
	"github.com/rs/xid"
	"github.com/xanzy/go-gitlab"
	"path/filepath"
	"strconv"
	"time"
)

func (t *FixParams) GitlabFix() (preview []Preview, err error) {
	var (
		branch    = "fix_" + xid.New().String()
		defBranch string
	)

	ctx, cancel := context.WithTimeout(context.Background(), t.TimeOut)
	defer cancel()
	git, err := gitlab.NewClient(t.Token, gitlab.WithBaseURL(t.GitlabUrl))

	// git配置 克隆文件
	respoName := t.UserName + "_" + strconv.FormatInt(time.Now().Unix(), 10)
	repoPath := filepath.Join("./", respoName)
	defer func() {
		// 删除文件夹
		DelDir(repoPath)
	}()

	defBranch, err = GitConfig(ctx, "./", repoPath, branch, t.GitlabUrl+"/"+t.TargetOwner+"/"+t.Repo+".git", t.CommitHash, t.proxyUrl, t.UserName, t.Password)
	if err != nil {
		return
	}

	// 设置git邮箱和用户名
	_, err = RunGitCommand(ctx, repoPath, "git", "config", "user.email", t.UserEmail)
	if err != nil {
		return
	}
	_, err = RunGitCommand(ctx, repoPath, "git", "config", "user.name", t.UserName)
	if err != nil {
		return
	}

	t.Dir = repoPath
	preview, err = t.LocalFix()
	if err != nil {
		return
	}

	// 查看是否有修改
	_, err = RunGitCommand(ctx, repoPath, "git", "status", "--short")
	if err != nil {
		return
	}

	// commit代码,要执行的参数 commit msg
	_, err = RunGitCommand(ctx, repoPath, "git", "commit", "-am", "update "+t.CompName+" "+t.CompVersion+" to "+t.MinFixVersion)
	if err != nil {
		return
	}
	//  提交文件
	_, err = RunGitCommand(ctx, repoPath, "git", "push", "--set-upstream origin", branch)
	if err != nil {
		return
	}
	g := gitlab.CreateMergeRequestOptions{
		Title:        gitlab.String(t.Title),
		Description:  gitlab.String(t.Body),
		SourceBranch: gitlab.String(branch),
		TargetBranch: gitlab.String(defBranch),
	}
	rule, response, err := git.MergeRequests.CreateMergeRequest(t.TargetOwner+"/"+t.Repo, &g)
	println(rule)
	println(response)
	return
}
