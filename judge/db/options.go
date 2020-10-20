package db

import "fmt"

var (
	DbMysql = fmt.Sprintf("user:pwd@tcp(127.0.0.1:3306)/code_evaluation")//mysql配置
)

//代码对拍表
type EvaCompare struct {
	CompareId     int    `json:"CompareId"`
	CompareName   string `json:"CompareName"`
	UserId        int    `json:"UserId"`
	FirstCodeId   int    `json:"FirstCodeId"`
	SecondCodeId  int    `json:"SecondCodeId"`
	InputDataPath string `json:"InputDataPath"`
	MaxInputGroup int    `json:"MaxInputGroup"`
	CreateTime    string `json:"CreateTime"`
	UpdateTime    string `json:"UpdateTime"`
	Remarks       string `json:"Remarks"`
}

//运行代码结果
type CodeRunData struct {
	Output     string
	Status     string
	TimeUsed   string
	MemoryUsed string
}

//对拍运行表中结果的json
type CodeMsg struct {
	First_runTime    string `json:"First_runTime"`
	First_runMemory  string `json:"First_runMemory"`
	Second_runTime   string `json:"Second_runTime"`
	Second_runMemory string `json:"Second_runMemory"`
	Input_path       string
}

//代码运行状态信息
type CodeStatusMsg struct {
	Status     string `json:"status"`
	TimeUsed   string `json:"timeUsed"`
	MemoryUsed string `json:"memoryUsed"`
}


//代码信息表
type EvaCode struct {
	CodeId 			int		`json:"CodeId"`
	CodeName		string	`json:"CodeName"`
	UserId 			int		`json:"UserId"`
	CodeText		string  `json:"Code_text"`
	Path 			string	`json:"Path"`
	CreateTime 		string 	`json:"CreateTime"`
	UpdateTime 		string	`json:"UpdateTime"`
}

//运行代码结构
type runCode struct {
	CodeId int
	UserId int
	RunData string
	InputPath string
	Version string
}
