package main

import (
	"github.com/aosfather/bingo_utils/files"
	"os"
)

// 任务执行
type Task struct {
	Id           string      //执行的序列唯一id
	Code         string      //任务的标识符
	Label        string      //任务显示用名称
	Script       string      //脚本名称
	ScriptType   string      //脚本类型
	InitParamter []Parameter //执行用的参数
}

type Parameter struct {
	Code  string //参数名
	Label string //参数显示名称
	Value string //参数值
}

type TaskExecutor struct {
	WorkRootPath string //工作用的根目录
	t            *Task
	workDir      string //工作目录
}

func (this *TaskExecutor) Init(t *Task) {
	if t == nil || t.Code == "" || t.Script == "" {
		return
	}

	this.t = t
	this.workDir = this.WorkRootPath + "/" + t.Code

}
func (this *TaskExecutor) Run() {
	if this.t == nil {
		//检查任务
		return
	}
	//创建工作目录
	if !files.IsFileExist(this.workDir) {
		os.Mkdir(this.workDir, os.ModePerm)
	}
	var w Worker
	//根据脚本类型使用不同的worker
	switch this.t.ScriptType {
	case "lua":
		l := &LuaWorker{}
		l.Init()
		w = l
	case "javascript":
		w = nil

	}

	//加载任务参数
	if len(this.t.InitParamter) > 0 {
		for _, v := range this.t.InitParamter {
			w.AddInput(v.Code, v.Value)
		}
	}
	//加载任务持久化状态（例如工作进度等)

	w.LoadRecords(this.loadRecords())
	w.RunScript(this.t.Script)
	//写入记录

}

//读取记录
func (this *TaskExecutor) loadRecords() map[string]interface{} {
	//从文件中读取记录
	return make(map[string]interface{})
}
