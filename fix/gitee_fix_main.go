package fix

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (t *FixParams) GiteeFix() (preview []Preview, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), t.TimeOut)
	defer cancel()

	resp, err := fork(t.Token, t.TargetOwner, t.Repo)
	if err != nil {
		return
	}
	respByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	forksResponse := new(forksResponse)
	err = json.Unmarshal(respByte, forksResponse)
	if err != nil {
		return
	}
	// git配置 克隆文件
	respoName := t.UserName + "_" + strconv.FormatInt(time.Now().Unix(), 10)
	repoPath := filepath.Join("./", respoName)
	defer func() {
		// 删除文件夹
		DelDir(repoPath)
	}()
	t.defBranch, err = GitConfig(ctx, "./", repoPath, t.branch, forksResponse.HtmlUrl, t.CommitHash, t.proxyUrl, t.UserName, t.Password)
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
	_, err = RunGitCommand(ctx, repoPath, "git", "commit", "-am", "fix vuln")
	if err != nil {
		return
	}
	//  提交文件
	_, err = RunGitCommand(ctx, repoPath, "git", "push", "--set-upstream origin", t.branch)
	if err != nil {
		return
	}

	_, err = CreatePullRequest(t.Token, t.TargetOwner, t.Repo, t.Title, t.Body, t.branch, t.defBranch, forksResponse.Namespace.Path)
	return
}

func fork(accessToken, owner, repo string) (resp *http.Response, err error) {
	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/%s/forks", owner, repo)
	sprintf := fmt.Sprintf("{\n    \"access_token\":\"%s\"\n}", accessToken)
	reader := strings.NewReader(sprintf)
	return http.Post(url, "application/json", reader)
}

func CreatePullRequest(accessToken, owner, repo, title, body, currentBranch, targetBranch, currentOwner string) (resp *http.Response, err error) {
	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/%s/pulls", owner, repo)
	request := CreatePullRequestRequest{
		AccessToken: accessToken,
		Title:       title,
		Head:        currentOwner + ":" + currentBranch,
		Base:        targetBranch,
		Body:        body,
	}
	marshal, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	reader := strings.NewReader(string(marshal))
	return http.Post(url, "application/json", reader)
}
