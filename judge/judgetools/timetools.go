package judgetools

import (
	"fmt"
	"time"
)

func GetMillisecond() string {
	t1 := time.Now().Unix()
	t2 := time.Now().UnixNano() / 1e6
	t3 := fmt.Sprintf("%.0f", float64(t1)+float64(t2)*1000)
	return t3
}