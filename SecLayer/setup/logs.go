package setup

import (
	"encoding/json"
	"fmt"

	"github.com/beego/beego/v2/core/logs"
)

func InitLogger() (err error) {
	config := make(map[string]interface{})
	config["filename"] = AppConfig.LogPath
	config["level"] = convertLogLevel(AppConfig.LogLevel)

	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("marshal failed, err: ", err)
		return
	}

	logs.SetLogger(logs.AdapterFile, string(configStr))

	return
}

func convertLogLevel(level string) int {
	switch level {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	}

	return logs.LevelDebug
}
