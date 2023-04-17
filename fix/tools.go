package fix

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
)

// 删除文件夹
func DelDir(path string) error {
	if Exists(path) {
		return os.RemoveAll(path)
	}
	return nil
}

// 判断文件夹或文件是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// 设置git配置
func GitConfig(ctx context.Context, path, repoPath, branch, gitRemote, commitHash, proxyUrl, username, password string) (string, error) {
	var (
		err error
	)
	if len(username) != 0 && len(password) != 0 {
		index := strings.Index(gitRemote, "://")
		i := index + 3
		username = strings.ReplaceAll(username, "@", "%40")
		password = strings.ReplaceAll(password, "@", "%40")
		gitRemote = gitRemote[0:i] + username + ":" + password + "@" + gitRemote[i:]
	}

	// 克隆pom文件（执行git 命令）
	// 1. 初始化仓库
	if len(proxyUrl) == 0 {
		_, err = RunGitCommand(ctx, path, "git", "clone", gitRemote, repoPath)
		if err != nil {
			return "", err
		}
	} else {
		_, err = RunGitCommand(ctx, path, "git", "clone", gitRemote, repoPath, "http.proxy="+proxyUrl)

		if err != nil {
			return "", err
		}

	}

	cmd := exec.CommandContext(ctx, "git", "branch")
	cmd.Dir = path
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	if err != nil {
		return "", err
	}
	if len(stdout.String()) == 0 {
		return "", errors.New("无法获得默认分支")
	}
	defBranch := strings.TrimSpace(strings.ReplaceAll(stdout.String(), "*", ""))

	_, err = RunGitCommand(ctx, repoPath, "git", "checkout", "-b", branch, commitHash)

	if err != nil {
		return defBranch, err
	}

	return defBranch, nil
}

// 执行任意cmd命令的封装
func RunGitCommand(ctx context.Context, path, name string, arg ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, arg...)
	cmd.Dir = path
	if out, err := cmd.Output(); err != nil {
		//检测报错是否是因为超时引起的
		if ctx.Err() != nil && ctx.Err() == context.DeadlineExceeded {
			return "command timeout", errors.New("command timeout")
		}
		return string(out), err
	} else {
		return string(out), nil
	}
}
