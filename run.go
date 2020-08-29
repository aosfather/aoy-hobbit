package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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
	var taskfile string

	flag.StringVar(&configfile, "c", "", "run with the config file")
	flag.StringVar(&taskfile, "t", "", "run task file")
	flag.StringVar(&scriptfile, "s", "", " run script file")
	flag.Usage = usage
	flag.Parse()
	if len(os.Args) != 3 || (configfile == "" && scriptfile == "" && taskfile == "") {
		flag.Usage()
	} else if scriptfile != "" { //执行指定的脚本文件，用于本地脚本测试
		fmt.Println(scriptfile)
		w := LuaWorker{outpath: "work"}
		w.Init()
		w.RunScript(scriptfile)
	} else if taskfile != "" { //执行指定的task文件，用于本地task测试之用
		fmt.Println(taskfile)
		taskcontent, err := ioutil.ReadFile(taskfile)
		if err != nil {
			debug(err.Error())
			return
		}
		t := &Task{}
		json.Unmarshal(taskcontent, t)
		debug(t)
		te := TaskExecutor{}
		te.WorkRootPath = ""
		te.Init(t)
		te.Run()

	} else {
		//根据configfile设置，连接到任务服务器ring上获取任务进行执行

	}
}

/**
  输出使用文档
*/
func usage() {
	fmt.Fprintf(os.Stderr, `hobbit version: hobbit/1.0.0
Usage: hobbit [-c filename] |[-s filename] | [-t filename]

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
