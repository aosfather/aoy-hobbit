package main

import (
	"flag"
	"fmt"
	"os"
)

/**
  执行逻辑
   1、根据配置文件，连接到指定的服务端
   2、获取任务
   3、执行任务
   4、上传任务进度
*/
func main() {
	//
	var configfile string
	var scriptfile string

	flag.StringVar(&configfile, "c", "", "the config file")
	flag.StringVar(&scriptfile, "s", "", "script file name")
	flag.Usage = usage
	flag.Parse()
	if configfile == "" && scriptfile == "" {
		flag.Usage()
	} else if scriptfile != "" {
		fmt.Println(scriptfile)
		w := LuaWorker{}
		w.Init()
		w.RunScript(scriptfile)
	}
}

/**
  输出使用文档
*/
func usage() {
	fmt.Fprintf(os.Stderr, `hobbit version: hobbit/1.0.0
Usage: hobbit [-h] [-c filename] |[-s filename]

Options:
`)
	flag.PrintDefaults()
}

//任务连接配置
type Config struct {
	Name     string //执行者名称
	User     string //所属用户
	PoolName string //用户下任务池名称
	Token    string //连接用的Token
}
