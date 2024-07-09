package version

import (
	"fmt"
)

const (
	IkubeopsDescription = "这是一个 gin 模板项目"
)

var (
	IkubeopsTag         string
	IkubeopsCommit      string
	IkubeopsBranch      string
	IkubeopsBuildTime   string
	IkubeopsGoVersion   string
	IkubeopsUrl         string
	IkubeopsProjectName = "ikubeops-gin-dem"
	IkubeopsConfigType  string
	IkubeopsConfigFile  string
	IkubeopsConfigEtcd  string
)

// FullVersion show the version info
func FullTagVersion() string {
	version := fmt.Sprintf("Version   : %s\nBuild Time: %s\nGit URL: %s\nGit Branch: %s\nGo Version: %s\n", IkubeopsCommit, IkubeopsBuildTime, IkubeopsUrl, IkubeopsBranch, IkubeopsGoVersion)
	return version
}

// Short 版本缩写
func ShortTagVersion() string {
	return fmt.Sprintf("%s[%s %s]", IkubeopsTag, IkubeopsBuildTime, IkubeopsCommit)
}

func GetWebUrl(httpAddr string, tls bool) string {
	addr := fmt.Sprintf("%s", httpAddr)
	if tls {
		return fmt.Sprintf("https://%s", addr)
	} else {
		return fmt.Sprintf("http://%s", addr)
	}
}

func GetConfig() string {
	if IkubeopsConfigType == "file" {
		return fmt.Sprintf("File: %s", IkubeopsConfigFile)
	} else if IkubeopsConfigType == "env" {
		return "ENV"
	} else {
		return fmt.Sprintf("ETCD: %s", IkubeopsConfigEtcd)
	}
}
