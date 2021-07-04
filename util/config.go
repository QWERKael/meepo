package util

import (
	"flag"
	"fmt"
	"path/filepath"
	"utility-go/config"
)

var (
	Config     Conf
	configPath = flag.String("c", "config/config.yml", "the path of config file")
)

func init() {
	fmt.Println("初始化模块: config.go")
	if err := ParseConfig(); err != nil {
		fmt.Printf("init config error: %s\n", err.Error())
		SugarLogger.Fatal(err)
	}
	err := Config.Default()
	if err != nil {
		SugarLogger.Fatal(err)
	}
	fmt.Printf("配置文件信息：\n%#v", Config)
}

func ParseConfig() error {
	flag.Parse()
	Config = Conf{}
	err := config.ParserFromPath(*configPath, &Config)
	if err != nil {
		return err
	}
	err = Config.Default()
	if err != nil {
		return err
	}
	LogInit()
	return nil
}

type Conf struct {
	ListenHost          string `yaml:"host"`
	ListenPort          int    `yaml:"listen"`
	ActorName           string `yaml:"name"`
	WorkDir             string `yaml:"work dir"`
	WorkDirAbs          string
	LogPath             string `yaml:"log path"`
	LogPathAbs          string
	PluginDir           string `yaml:"plugin dir"`
	PluginDirAbs        string
	LuaPackagePath      string `yaml:"lua package path"`
	LuaPackagePathAbs   string
	LuaScriptsPath      string `yaml:"lua scripts path"`
	LuaScriptsPathAbs   string
	PromptConfigPath    string `yaml:"prompt config"`
	PromptConfigPathAbs string
}

func (c *Conf) Default() error {
	var err error

	//设置默认的工作目录，默认为当前目录
	addStringDefaultValue(&c.WorkDir, "./")
	if c.WorkDirAbs, err = filepath.Abs(c.WorkDir); err != nil {
		return err
	}

	//设置默认的日志地址，日志地址为空字符串时，日志输出到控制台，默认为空字符串
	addStringDefaultValue(&c.LogPath, "")
	if !filepath.IsAbs(c.LogPath) {
		if c.LogPath == "" {
			c.LogPathAbs = ""
		} else {
			c.LogPathAbs = filepath.Join(c.WorkDirAbs, c.LogPath)
		}
	}

	//设置默认的插件目录，默认为"plugin"目录
	addStringDefaultValue(&c.PluginDir, "plugin")
	if !filepath.IsAbs(c.PluginDir) {
		c.PluginDirAbs = filepath.Join(c.WorkDirAbs, c.PluginDir)
	}

	//设置默认的lua包地址，默认为"lua/module/?.lua"
	addStringDefaultValue(&c.LuaPackagePath, "lua/module/?.lua")
	if !filepath.IsAbs(c.LuaPackagePath) {
		c.LuaPackagePathAbs = filepath.Join(c.WorkDirAbs, c.LuaPackagePath)
	}

	//设置默认的lua脚本地址，默认为"lua"
	addStringDefaultValue(&c.LuaScriptsPath, "lua")
	if !filepath.IsAbs(c.LuaScriptsPath) {
		c.LuaScriptsPathAbs = filepath.Join(c.WorkDirAbs, c.LuaScriptsPath)
	}

	//设置默认的prompt脚本地址，默认为"config/prompt.yml"
	addStringDefaultValue(&c.PromptConfigPath, "config/prompt.yml")
	if !filepath.IsAbs(c.PromptConfigPath) {
		c.PromptConfigPathAbs = filepath.Join(c.WorkDirAbs, c.PromptConfigPath)
	}

	return nil
}

func addStringDefaultValue(key *string, defaultValue string) {
	if *key == "" {
		*key = defaultValue
	}
}
