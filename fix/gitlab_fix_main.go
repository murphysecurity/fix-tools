package fix

import (
	"context"
	"github.com/rs/xid"
	"github.com/xanzy/go-gitlab"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (t *FixParams) GitlabFix() (PrUrl string, preview []Preview, err error) {
	var (
		branch       = "fix_" + xid.New().String()
		defBranch    string
		mergeRequest *gitlab.MergeRequest
	)

	ctx, cancel := context.WithTimeout(context.Background(), t.TimeOut)
	defer cancel()
	client, err := gitlab.NewClient(t.Token, gitlab.WithBaseURL(t.GitlabUrl))

	// git配置 克隆文件
	respoName := t.UserName + "_" + strconv.FormatInt(time.Now().Unix(), 10)
	repoPath := filepath.Join("./", respoName)
	defer func() {
		// 删除文件夹
		DelDir(repoPath)
	}()

	defBranch, err = GitConfig(ctx, "./", repoPath, branch, t.GitlabUrl+"/"+t.TargetOwner+"/"+t.Repo+".git", t.CommitHash, t.ProxyUrl, t.UserName, t.Password)
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
	preview, _, _, err = t.LocalFix()
	if err != nil {
		return
	}

	// 查看是否有修改
	_, err = RunGitCommand(ctx, repoPath, "git", "status", "--short")
	if err != nil {
		return
	}

	// commit代码,要执行的参数 commit msg
	_, err = RunGitCommand(ctx, repoPath, "git", "commit", "-am", "fix vuln")
	if err != nil {
		return
	}
	//  提交文件
	// git push https://gitlab-ci-token:<access_token>@gitlab.com/myuser/myrepo.git <branch_name>
	//_, err = RunGitCommand(ctx, repoPath, "git", "push", "--set-upstream", "origin", branch)
	httpType := "http"
	if strings.Contains(t.GitlabUrl, "https") {
		httpType += "s"
	}
	httpType += "://"
	gitlabUrlEnd := strings.ReplaceAll(t.GitlabUrl, httpType, "")
	if gitlabUrlEnd[len(gitlabUrlEnd)-1] == '/' {
		gitlabUrlEnd = gitlabUrlEnd[0 : len(gitlabUrlEnd)-1]
	}
	_, err = RunGitCommand(ctx, repoPath, "git", "push", "--set-upstream", httpType+"gitlab-ci-token:"+t.Token+"@"+gitlabUrlEnd+"/"+t.TargetOwner+"/"+t.Repo+".git", branch)
	if err != nil {
		return
	}
	g := gitlab.CreateMergeRequestOptions{
		Title:        gitlab.String(t.Title),
		Description:  gitlab.String(t.Body),
		SourceBranch: gitlab.String(branch),
		TargetBranch: gitlab.String(defBranch),
	}
	mergeRequest, _, err = client.MergeRequests.CreateMergeRequest(t.TargetOwner+"/"+t.Repo, &g)
	PrUrl = mergeRequest.WebURL
	return
}
