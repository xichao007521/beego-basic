package b_logger

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type FileConfig struct {
	FileName string
	MaxLines int  // 每个文件保存的最大行数，默认值 1000000
	MaxSize  int  // 每个文件保存的最大尺寸, 默认值是 1 << 28
	Daily    bool // 是否按照每天 logrotate，默认是 true
	MaxDays  int  // 文件最多保存多少天，默认保存 7 天
	Hourly   bool
	MaxHours int64
	Rotate   bool   // 默认是 true
	Level    int    // 默认是 Trace 级别
	Perm     string // 日志文件权限
}

type AllConfig struct {
	FileLoggers map[string]FileConfig `toml:"file"`
}

var loggerConfig AllConfig
var lock sync.Mutex

var AppConfig FileConfig

func init() {
	readConfig()
	AppConfig = loggerConfig.FileLoggers["app"]
	content, _ := json.Marshal(AppConfig)
	beego.SetLogger("file", string(content))

	beego.BConfig.Log.AccessLogs = true

	// 内置一个access日志
	AccessLogger = buildCustomLogger("access")
	AccessLogger.EnableFuncCallDepth(false)
}

// 自定义access日志
var AccessLogger *logs.BeeLogger

func buildCustomLogger(loggerName string) *logs.BeeLogger {
	loggerFileConfig := loggerConfig.FileLoggers[loggerName]
	configContent, _ := json.Marshal(loggerFileConfig)
	result := logs.NewLogger()
	result.SetLogger("file", string(configContent))
	result.SetLogger(logs.AdapterConsole)
	return result
}

// read config
func readConfig() {
	lock.Lock()
	defer lock.Unlock()
	filename := "logger.toml"
	runmode := beego.AppConfig.String("runmode")
	_, err := os.Stat("./conf/" + runmode + "." + filename)
	if err == nil {
		filename = runmode + "." + filename
	}

	if len(loggerConfig.FileLoggers) == 0 {
		data, err := ioutil.ReadFile("./conf/" + filename)
		if err != nil {
			log.Fatal(err)
		}
		var loggerToml AllConfig
		if _, err := toml.Decode(string(data), &loggerToml); err != nil {
			log.Fatal(err)
		}
		loggerConfig = loggerToml
		for _, fileLogger := range loggerConfig.FileLoggers {
			if fileLogger.MaxDays == 0 {
				fileLogger.MaxDays = 7
			}
			if fileLogger.MaxLines == 0 {
				fileLogger.MaxLines = 100000
			}
			if fileLogger.MaxSize == 0 {
				fileLogger.MaxSize = 1 << 28
			}
		}
	}
}
