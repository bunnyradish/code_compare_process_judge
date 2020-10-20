package zapconf

import (
	"fmt"
	"time"
)

var (
	LogPath = fmt.Sprintf("/data/logs/golang/" + time.Now().Format("2006-01-02") + ".log")
	LogMaxSize = 32
	LogMaxBackups = 5
)