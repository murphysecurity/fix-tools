package fix

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v47/github"
	"github.com/rs/xid"
	"golang.org/x/oauth2"
)

func (t *FixParams) GithubFix() (prUrl string, preview []Preview, err error) {
	var (
		client *github.Client

		branch              = "oscs_fix_" + xid.New().String()
		defaultBranch       string
		targetDefaultBranch string
	)

	ctx, cancel := context.WithTimeout(context.Background(), t.TimeOut)
	defer cancel()
	path := "./"
	gitRemote := fmt.Sprintf("https://github.com/%s/%s.git", t.UserName, t.Repo)

	// 获取github连接
	client, err = GetClient(ctx, t.Token, t.ProxyUrl)
	if err != nil {
		err = errors.New("获取github连接失败: " + err.Error())
		return
	}
	// fork2次 第一次fork的默认分支会出问题
	_, _, _, _, _ = CreateFork(t.Token, t.TargetOwner, t.Repo, t.ProxyUrl)
	time.Sleep(time.Second)
	defaultBranch, targetDefaultBranch, _, _, err = CreateFork(t.Token, t.TargetOwner, t.Repo, t.ProxyUrl)
	if err != nil && (defaultBranch == "" || targetDefaultBranch == "") {
		err = errors.New("fork到本地失败")
		return
	}
	// 创建自己仓库新的分支
	if err = CreateBranch(ctx, client, t.UserName, t.Repo, defaultBranch, branch); err != nil {
		return
	}

	// git配置 克隆文件
	respoName := t.UserName + "_" + t.Repo + "_" + strconv.FormatInt(time.Now().Unix(), 10)
	repoPath := filepath.Join(path, respoName)
	defer func() {
		// 删除文件夹
		DelDir(repoPath)
	}()
	_, err = GitConfig(ctx, path, repoPath, branch, gitRemote, t.CommitHash, t.ProxyUrl, t.UserName, t.Password, t.Token)
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
	// 切换分支
	_, err = RunGitCommand(ctx, repoPath, "git", "checkout", branch)
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
	// 7. 提交文件
	err = GitPushCode(ctx, repoPath, t.UserName, t.Token)
	if err != nil {
		return
	}
	// 提交pr

	prUrl, _, _, _, err = CreatePr(ctx, client, t.TargetOwner, t.Repo, t.UserName+":"+branch, targetDefaultBranch, t.Title, t.Body)
	return
}

func GetClient(ctx context.Context, token, proxyUrl string) (*github.Client, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	if proxyUrl != "" {
		proxy := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(proxyUrl)
		}
		ctx = context.WithValue(ctx, oauth2.HTTPClient, &http.Client{
			Transport: &http.Transport{Proxy: proxy},
		})
	}
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	_, _, err := client.Users.Get(ctx, "")
	return client, err
}

func GithubPostHttpByToken(httpUrl, token, proxyUrl string, body interface{}) (error, *http.Response) {

	var (
		client = http.Client{}
		resp   *http.Response
	)
	if proxyUrl != "" {
		proxy := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(proxyUrl)
		}
		client.Transport = &http.Transport{Proxy: proxy}
	}
	marshal, err := json.Marshal(body)
	if err != nil {
		return err, nil
	}
	req, err := http.NewRequest(http.MethodPost, httpUrl, bytes.NewBuffer(marshal))

	if err != nil {
		return err, nil
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = client.Do(req)
	return err, resp

}

// 从GitHub源仓库fork自己的仓库
func CreateFork(token, TargetOwner, Repo, proxyUrl string) (defaultBranch, targetDefaultBranch string, private, forkBool bool, err error) {
	type forkBody struct {
		Organization      string `json:"organization,omitempty"`
		Name              string `json:"name,omitempty"`
		DefaultBranchOnly bool   `json:"default_branch_only"`
	}
	var (
		res  []byte
		resp *http.Response
	)
	body := forkBody{
		DefaultBranchOnly: true,
	}

	url := fmt.Sprintf("https://api.github.com/repos/%v/%v/forks", TargetOwner, Repo)

	err, resp = GithubPostHttpByToken(url, token, proxyUrl, body)
	if err != nil {
		return
	}
	res, err = io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	var githubForkResp GithubForkResp
	err = json.Unmarshal(res, &githubForkResp)
	if err != nil {
		return
	}
	defaultBranch = githubForkResp.DefaultBranch

	targetDefaultBranch = githubForkResp.Parent.DefaultBranch
	if len(targetDefaultBranch) == 0 {
		targetDefaultBranch = githubForkResp.DefaultBranch
	}
	private = githubForkResp.Private
	forkBool = githubForkResp.Fork

	return
}

// 创建分支
func CreateBranch(ctx context.Context, client *github.Client, owner, repo, oldBranchName, newBranchName string) error {
	// 获取当前的ref信息
	oldBranch := "refs/heads/" + oldBranchName
	ref, response, err := client.Git.GetRef(ctx, owner, repo, oldBranch)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return errors.New("获取当前分支信息失败" + strconv.Itoa(response.StatusCode))
	}

	// 创建一个ref
	newRef := &github.Reference{Ref: github.String("refs/heads/" + newBranchName),
		Object: &github.GitObject{
			SHA: ref.Object.SHA,
		},
	}
	_, response, err = client.Git.CreateRef(ctx, owner, repo, newRef)
	if err != nil {
		return err
	}
	if response.StatusCode != 201 {
		return errors.New("创建分支请求失败, 状态码为" + strconv.Itoa(response.StatusCode))
	}
	return nil
}

// 创建pr
func CreatePr(ctx context.Context, client *github.Client, owner, repo, head, base, title, body string) (string, string, int, time.Time, error) {
	newPR := &github.NewPullRequest{
		Title: github.String(title),
		Head:  github.String(head), // 旧分支
		Base:  github.String(base), // 新分支
		Body:  github.String(body),
	}

	pr, response, err := client.PullRequests.Create(ctx, owner, repo, newPR)
	if err != nil {
		return "", "", 0, time.Time{}, err
	}
	if response.StatusCode == 403 {
		return "", "", 0, time.Time{}, errors.New("github create pr no authorization")
	} else if response.StatusCode != 201 {
		return "", "", 0, time.Time{}, errors.New("提交pr请求状态码返回错误, 状态码为" + strconv.Itoa(response.StatusCode))
	}

	return pr.GetHTMLURL(), pr.GetURL(), pr.GetNumber(), pr.GetCreatedAt(), nil
}

// 带用户名提交push
func GitPushCode(ctx context.Context, path, username, password string) error {

	// 克隆pom文件（执行git 命令）
	r, err := git.PlainOpen(path)
	if err != nil {
		return err
	}
	err = r.PushContext(ctx, &git.PushOptions{
		Auth: &githttp.BasicAuth{
			Username: username,
			Password: password,
		},
		Force: true,
	})
	return err
}

type GithubForkResp struct {
	Id       int    `json:"id"`
	NodeId   string `json:"node_id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
	Owner    struct {
		Login             string `json:"login"`
		Id                int    `json:"id"`
		NodeId            string `json:"node_id"`
		AvatarUrl         string `json:"avatar_url"`
		GravatarId        string `json:"gravatar_id"`
		Url               string `json:"url"`
		HtmlUrl           string `json:"html_url"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReposUrl          string `json:"repos_url"`
		EventsUrl         string `json:"events_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"owner"`
	HtmlUrl          string      `json:"html_url"`
	Description      string      `json:"description"`
	Fork             bool        `json:"fork"`
	Url              string      `json:"url"`
	ForksUrl         string      `json:"forks_url"`
	KeysUrl          string      `json:"keys_url"`
	CollaboratorsUrl string      `json:"collaborators_url"`
	TeamsUrl         string      `json:"teams_url"`
	HooksUrl         string      `json:"hooks_url"`
	IssueEventsUrl   string      `json:"issue_events_url"`
	EventsUrl        string      `json:"events_url"`
	AssigneesUrl     string      `json:"assignees_url"`
	BranchesUrl      string      `json:"branches_url"`
	TagsUrl          string      `json:"tags_url"`
	BlobsUrl         string      `json:"blobs_url"`
	GitTagsUrl       string      `json:"git_tags_url"`
	GitRefsUrl       string      `json:"git_refs_url"`
	TreesUrl         string      `json:"trees_url"`
	StatusesUrl      string      `json:"statuses_url"`
	LanguagesUrl     string      `json:"languages_url"`
	StargazersUrl    string      `json:"stargazers_url"`
	ContributorsUrl  string      `json:"contributors_url"`
	SubscribersUrl   string      `json:"subscribers_url"`
	SubscriptionUrl  string      `json:"subscription_url"`
	CommitsUrl       string      `json:"commits_url"`
	GitCommitsUrl    string      `json:"git_commits_url"`
	CommentsUrl      string      `json:"comments_url"`
	IssueCommentUrl  string      `json:"issue_comment_url"`
	ContentsUrl      string      `json:"contents_url"`
	CompareUrl       string      `json:"compare_url"`
	MergesUrl        string      `json:"merges_url"`
	ArchiveUrl       string      `json:"archive_url"`
	DownloadsUrl     string      `json:"downloads_url"`
	IssuesUrl        string      `json:"issues_url"`
	PullsUrl         string      `json:"pulls_url"`
	MilestonesUrl    string      `json:"milestones_url"`
	NotificationsUrl string      `json:"notifications_url"`
	LabelsUrl        string      `json:"labels_url"`
	ReleasesUrl      string      `json:"releases_url"`
	DeploymentsUrl   string      `json:"deployments_url"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	PushedAt         time.Time   `json:"pushed_at"`
	GitUrl           string      `json:"git_url"`
	SshUrl           string      `json:"ssh_url"`
	CloneUrl         string      `json:"clone_url"`
	SvnUrl           string      `json:"svn_url"`
	Homepage         string      `json:"homepage"`
	Size             int         `json:"size"`
	StargazersCount  int         `json:"stargazers_count"`
	WatchersCount    int         `json:"watchers_count"`
	Language         interface{} `json:"language"`
	HasIssues        bool        `json:"has_issues"`
	HasProjects      bool        `json:"has_projects"`
	HasDownloads     bool        `json:"has_downloads"`
	HasWiki          bool        `json:"has_wiki"`
	HasPages         bool        `json:"has_pages"`
	ForksCount       int         `json:"forks_count"`
	MirrorUrl        interface{} `json:"mirror_url"`
	Archived         bool        `json:"archived"`
	Disabled         bool        `json:"disabled"`
	OpenIssuesCount  int         `json:"open_issues_count"`
	License          struct {
		Key    string `json:"key"`
		Name   string `json:"name"`
		SpdxId string `json:"spdx_id"`
		Url    string `json:"url"`
		NodeId string `json:"node_id"`
	} `json:"license"`
	AllowForking             bool          `json:"allow_forking"`
	IsTemplate               bool          `json:"is_template"`
	WebCommitSignoffRequired bool          `json:"web_commit_signoff_required"`
	Topics                   []interface{} `json:"topics"`
	Visibility               string        `json:"visibility"`
	Forks                    int           `json:"forks"`
	OpenIssues               int           `json:"open_issues"`
	Watchers                 int           `json:"watchers"`
	DefaultBranch            string        `json:"default_branch"`
	Permissions              struct {
		Admin    bool `json:"admin"`
		Maintain bool `json:"maintain"`
		Push     bool `json:"push"`
		Triage   bool `json:"triage"`
		Pull     bool `json:"pull"`
	} `json:"permissions"`
	Parent struct {
		Id       int    `json:"id"`
		NodeId   string `json:"node_id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Private  bool   `json:"private"`
		Owner    struct {
			Login             string `json:"login"`
			Id                int    `json:"id"`
			NodeId            string `json:"node_id"`
			AvatarUrl         string `json:"avatar_url"`
			GravatarId        string `json:"gravatar_id"`
			Url               string `json:"url"`
			HtmlUrl           string `json:"html_url"`
			FollowersUrl      string `json:"followers_url"`
			FollowingUrl      string `json:"following_url"`
			GistsUrl          string `json:"gists_url"`
			StarredUrl        string `json:"starred_url"`
			SubscriptionsUrl  string `json:"subscriptions_url"`
			OrganizationsUrl  string `json:"organizations_url"`
			ReposUrl          string `json:"repos_url"`
			EventsUrl         string `json:"events_url"`
			ReceivedEventsUrl string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"owner"`
		HtmlUrl          string      `json:"html_url"`
		Description      string      `json:"description"`
		Fork             bool        `json:"fork"`
		Url              string      `json:"url"`
		ForksUrl         string      `json:"forks_url"`
		KeysUrl          string      `json:"keys_url"`
		CollaboratorsUrl string      `json:"collaborators_url"`
		TeamsUrl         string      `json:"teams_url"`
		HooksUrl         string      `json:"hooks_url"`
		IssueEventsUrl   string      `json:"issue_events_url"`
		EventsUrl        string      `json:"events_url"`
		AssigneesUrl     string      `json:"assignees_url"`
		BranchesUrl      string      `json:"branches_url"`
		TagsUrl          string      `json:"tags_url"`
		BlobsUrl         string      `json:"blobs_url"`
		GitTagsUrl       string      `json:"git_tags_url"`
		GitRefsUrl       string      `json:"git_refs_url"`
		TreesUrl         string      `json:"trees_url"`
		StatusesUrl      string      `json:"statuses_url"`
		LanguagesUrl     string      `json:"languages_url"`
		StargazersUrl    string      `json:"stargazers_url"`
		ContributorsUrl  string      `json:"contributors_url"`
		SubscribersUrl   string      `json:"subscribers_url"`
		SubscriptionUrl  string      `json:"subscription_url"`
		CommitsUrl       string      `json:"commits_url"`
		GitCommitsUrl    string      `json:"git_commits_url"`
		CommentsUrl      string      `json:"comments_url"`
		IssueCommentUrl  string      `json:"issue_comment_url"`
		ContentsUrl      string      `json:"contents_url"`
		CompareUrl       string      `json:"compare_url"`
		MergesUrl        string      `json:"merges_url"`
		ArchiveUrl       string      `json:"archive_url"`
		DownloadsUrl     string      `json:"downloads_url"`
		IssuesUrl        string      `json:"issues_url"`
		PullsUrl         string      `json:"pulls_url"`
		MilestonesUrl    string      `json:"milestones_url"`
		NotificationsUrl string      `json:"notifications_url"`
		LabelsUrl        string      `json:"labels_url"`
		ReleasesUrl      string      `json:"releases_url"`
		DeploymentsUrl   string      `json:"deployments_url"`
		CreatedAt        time.Time   `json:"created_at"`
		UpdatedAt        time.Time   `json:"updated_at"`
		PushedAt         time.Time   `json:"pushed_at"`
		GitUrl           string      `json:"git_url"`
		SshUrl           string      `json:"ssh_url"`
		CloneUrl         string      `json:"clone_url"`
		SvnUrl           string      `json:"svn_url"`
		Homepage         string      `json:"homepage"`
		Size             int         `json:"size"`
		StargazersCount  int         `json:"stargazers_count"`
		WatchersCount    int         `json:"watchers_count"`
		Language         string      `json:"language"`
		HasIssues        bool        `json:"has_issues"`
		HasProjects      bool        `json:"has_projects"`
		HasDownloads     bool        `json:"has_downloads"`
		HasWiki          bool        `json:"has_wiki"`
		HasPages         bool        `json:"has_pages"`
		ForksCount       int         `json:"forks_count"`
		MirrorUrl        interface{} `json:"mirror_url"`
		Archived         bool        `json:"archived"`
		Disabled         bool        `json:"disabled"`
		OpenIssuesCount  int         `json:"open_issues_count"`
		License          struct {
			Key    string `json:"key"`
			Name   string `json:"name"`
			SpdxId string `json:"spdx_id"`
			Url    string `json:"url"`
			NodeId string `json:"node_id"`
		} `json:"license"`
		AllowForking             bool     `json:"allow_forking"`
		IsTemplate               bool     `json:"is_template"`
		WebCommitSignoffRequired bool     `json:"web_commit_signoff_required"`
		Topics                   []string `json:"topics"`
		Visibility               string   `json:"visibility"`
		Forks                    int      `json:"forks"`
		OpenIssues               int      `json:"open_issues"`
		Watchers                 int      `json:"watchers"`
		DefaultBranch            string   `json:"default_branch"`
	} `json:"parent"`
	Source struct {
		Id       int    `json:"id"`
		NodeId   string `json:"node_id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Private  bool   `json:"private"`
		Owner    struct {
			Login             string `json:"login"`
			Id                int    `json:"id"`
			NodeId            string `json:"node_id"`
			AvatarUrl         string `json:"avatar_url"`
			GravatarId        string `json:"gravatar_id"`
			Url               string `json:"url"`
			HtmlUrl           string `json:"html_url"`
			FollowersUrl      string `json:"followers_url"`
			FollowingUrl      string `json:"following_url"`
			GistsUrl          string `json:"gists_url"`
			StarredUrl        string `json:"starred_url"`
			SubscriptionsUrl  string `json:"subscriptions_url"`
			OrganizationsUrl  string `json:"organizations_url"`
			ReposUrl          string `json:"repos_url"`
			EventsUrl         string `json:"events_url"`
			ReceivedEventsUrl string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"owner"`
		HtmlUrl          string      `json:"html_url"`
		Description      string      `json:"description"`
		Fork             bool        `json:"fork"`
		Url              string      `json:"url"`
		ForksUrl         string      `json:"forks_url"`
		KeysUrl          string      `json:"keys_url"`
		CollaboratorsUrl string      `json:"collaborators_url"`
		TeamsUrl         string      `json:"teams_url"`
		HooksUrl         string      `json:"hooks_url"`
		IssueEventsUrl   string      `json:"issue_events_url"`
		EventsUrl        string      `json:"events_url"`
		AssigneesUrl     string      `json:"assignees_url"`
		BranchesUrl      string      `json:"branches_url"`
		TagsUrl          string      `json:"tags_url"`
		BlobsUrl         string      `json:"blobs_url"`
		GitTagsUrl       string      `json:"git_tags_url"`
		GitRefsUrl       string      `json:"git_refs_url"`
		TreesUrl         string      `json:"trees_url"`
		StatusesUrl      string      `json:"statuses_url"`
		LanguagesUrl     string      `json:"languages_url"`
		StargazersUrl    string      `json:"stargazers_url"`
		ContributorsUrl  string      `json:"contributors_url"`
		SubscribersUrl   string      `json:"subscribers_url"`
		SubscriptionUrl  string      `json:"subscription_url"`
		CommitsUrl       string      `json:"commits_url"`
		GitCommitsUrl    string      `json:"git_commits_url"`
		CommentsUrl      string      `json:"comments_url"`
		IssueCommentUrl  string      `json:"issue_comment_url"`
		ContentsUrl      string      `json:"contents_url"`
		CompareUrl       string      `json:"compare_url"`
		MergesUrl        string      `json:"merges_url"`
		ArchiveUrl       string      `json:"archive_url"`
		DownloadsUrl     string      `json:"downloads_url"`
		IssuesUrl        string      `json:"issues_url"`
		PullsUrl         string      `json:"pulls_url"`
		MilestonesUrl    string      `json:"milestones_url"`
		NotificationsUrl string      `json:"notifications_url"`
		LabelsUrl        string      `json:"labels_url"`
		ReleasesUrl      string      `json:"releases_url"`
		DeploymentsUrl   string      `json:"deployments_url"`
		CreatedAt        time.Time   `json:"created_at"`
		UpdatedAt        time.Time   `json:"updated_at"`
		PushedAt         time.Time   `json:"pushed_at"`
		GitUrl           string      `json:"git_url"`
		SshUrl           string      `json:"ssh_url"`
		CloneUrl         string      `json:"clone_url"`
		SvnUrl           string      `json:"svn_url"`
		Homepage         string      `json:"homepage"`
		Size             int         `json:"size"`
		StargazersCount  int         `json:"stargazers_count"`
		WatchersCount    int         `json:"watchers_count"`
		Language         string      `json:"language"`
		HasIssues        bool        `json:"has_issues"`
		HasProjects      bool        `json:"has_projects"`
		HasDownloads     bool        `json:"has_downloads"`
		HasWiki          bool        `json:"has_wiki"`
		HasPages         bool        `json:"has_pages"`
		ForksCount       int         `json:"forks_count"`
		MirrorUrl        interface{} `json:"mirror_url"`
		Archived         bool        `json:"archived"`
		Disabled         bool        `json:"disabled"`
		OpenIssuesCount  int         `json:"open_issues_count"`
		License          struct {
			Key    string `json:"key"`
			Name   string `json:"name"`
			SpdxId string `json:"spdx_id"`
			Url    string `json:"url"`
			NodeId string `json:"node_id"`
		} `json:"license"`
		AllowForking             bool     `json:"allow_forking"`
		IsTemplate               bool     `json:"is_template"`
		WebCommitSignoffRequired bool     `json:"web_commit_signoff_required"`
		Topics                   []string `json:"topics"`
		Visibility               string   `json:"visibility"`
		Forks                    int      `json:"forks"`
		OpenIssues               int      `json:"open_issues"`
		Watchers                 int      `json:"watchers"`
		DefaultBranch            string   `json:"default_branch"`
	} `json:"source"`
	NetworkCount     int `json:"network_count"`
	SubscribersCount int `json:"subscribers_count"`
}
