package global

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var JiSiLu = &JiSiLuSetting{}

type JiSiLuSetting struct {
	UserName string
	Password string
}

func init() {
	vp := viper.New()
	// 优先读取app.local.yaml
	for _, f := range []string{"app.local.yaml", "app.yaml"} {
		if IsFileExist(filepath.Join("config", f)) {
			vp.SetConfigName(f)
			break
		}
	}

	vp.AddConfigPath("config/")
	vp.SetConfigType("yaml")
	if err := vp.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := vp.UnmarshalKey("jisilu", JiSiLu); err != nil {
		panic(err)
	}
}

// 检查文件是否存在
func IsFileExist(path string) bool {
	f, err := os.Stat(path)
	if err == nil {
		return !f.IsDir()
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}
