# fix-tool-go

修复模块包含
拉取项目
fork到本地
修复项目
提交pr


// 配置GoMod私有仓库
go env -w GOPRIVATE="git@git.murphy-int.com"
// 配置不加密访问
go env -w GOINSECURE="git.murphy-int.com"
// 配置不使用代理
go env -w GONOPROXY="git.murphy-int.com"
// 配置不验证包
go env -w GONOSUMDB="git.murphy-int.com"