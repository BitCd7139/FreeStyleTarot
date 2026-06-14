package config

import (
	"fmt"
	"time"
)

// ServerBootID 服务进程启动时生成，用于区分重启前后会话
var ServerBootID string

func InitBootID() {
	ServerBootID = fmt.Sprintf("%d", time.Now().UTC().UnixNano())
}
