package environment

import "fmt"

var (
	CompareRootPath = fmt.Sprintf("/var/www/html/zt-code-evaluation/compareGo/")//对拍根路径
	CompareJudgePath = fmt.Sprintf("/var/www/html/zt-code-evaluation/compareGo/judgepro")//对拍c++启动器的路径
	CodeRunRootPath = fmt.Sprintf("/var/www/html/zt-code-evaluation/codeGo/")//代码运行根路径
	CodeJudgePath = fmt.Sprintf("/var/www/html/zt-code-evaluation/codeGo/judgepro")//代码运行c++启动器的路径
)
