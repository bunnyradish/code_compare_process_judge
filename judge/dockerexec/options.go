package dockerexec

import "fmt"

var (
	DockerRunPath = fmt.Sprintf("/var/www/html/zt-code-evaluation/dockerF")//docker运行路径
	DockerWorkPath = fmt.Sprintf("/work")//容器工作路径
	LimitTime = fmt.Sprintf("1500")//默认最大运行时长
	LimitMemory = fmt.Sprintf("65535")//默认最大运行空间
)