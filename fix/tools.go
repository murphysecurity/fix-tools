package fix

import (
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
			err = errors.New(" clone 失败  " + err.Error())
			return "", err
		}
	} else {

		_, err = RunGitCommand(ctx, path, "git", "config", "http.proxy", proxyUrl)

		if err != nil {
			err = errors.New(" 设置代理 失败  " + err.Error())
			return "", err
		}

		_, err = RunGitCommand(ctx, path, "git", "config", "https.proxy", proxyUrl)

		if err != nil {
			err = errors.New(" 设置代理 失败  " + err.Error())
			return "", err
		}

		_, err = RunGitCommand(ctx, path, "git", "clone", gitRemote, repoPath)

		if err != nil {
			err = errors.New(" clone 失败  " + err.Error())
			return "", err
		}

	}

	cmd := exec.CommandContext(ctx, "git", "branch")
	cmd.Dir = repoPath
	out, err := cmd.Output()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			newErr := errors.New(string(out) + " 获取默认分支 失败   ==  " + exitError.Error() + string(exitError.Stderr))
			return string(out), newErr
		}

		err = errors.New(" 获取默认分支 失败  " + string(out) + err.Error())
		return "", err
	}

	if len(string(out)) == 0 {
		return "", errors.New("无法获得默认分支")
	}
	defBranch := strings.TrimSpace(strings.ReplaceAll(string(out), "*", ""))

	_, err = RunGitCommand(ctx, repoPath, "git", "checkout", "-b", branch, commitHash)

	if err != nil {
		err = errors.New(" checkout 失败  " + err.Error())
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

		if exitError, ok := err.(*exec.ExitError); ok {
			newErr := errors.New(string(out) + "  ==  " + exitError.Error() + string(exitError.Stderr))
			return string(out), newErr
		}
		newErr := errors.New(string(out) + "  ==  " + err.Error())
		return string(out), newErr
	} else {
		return string(out), nil
	}
}

// 执行任意cmd命令的封装
func RunGiteePush(ctx context.Context, path, name string, arg ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, arg...)
	cmd.Dir = path
	if out, err := cmd.Output(); err != nil {
		//检测报错是否是因为超时引起的
		if ctx.Err() != nil && ctx.Err() == context.DeadlineExceeded {
			return "command timeout", errors.New("command timeout")
		}

		if exitError, ok := err.(*exec.ExitError); ok {
			newErr := errors.New(string(out) + "  ==  " + exitError.Error() + string(exitError.Stderr))
			return string(out), newErr
		}
		newErr := errors.New(string(out) + "  ==  " + err.Error())
		return string(out), newErr
	} else {
		return string(out), nil
	}
}
