package main

import (
	"encoding/json"
	"fmt"
	"github.com/aosfather/bingo_utils"
	"github.com/aosfather/bingo_utils/files"
	"io/ioutil"
	"os"
	"time"
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
	log          *Logger
	recordName   string //记录文件名
}

func (this *TaskExecutor) Init(t *Task) {
	if t == nil || t.Code == "" || t.Script == "" {
		return
	}

	this.t = t
	this.workDir = this.WorkRootPath + "/" + t.Code
	//向日志文件中写入日志
	this.log = &Logger{}
	this.log.Init(this.workDir+"/runtime.log", this.t.Code+"-"+this.t.Id)
	this.recordName = this.workDir + "/" + this.t.Code + ".rec"

}
func (this *TaskExecutor) Run() {
	if this.t == nil {
		//检查任务
		return
	}
	debug("检查工作目录:", this.workDir)
	//创建工作目录
	if !files.IsFileExist(this.workDir) {
		debug("工作目录不存在，创建目录:", this.workDir)
		err := os.Mkdir(this.workDir, os.ModePerm)
		if err != nil {
			errs("创建工作目录错误:", err.Error())
		}
	}
	var w Worker
	//根据脚本类型使用不同的worker
	switch this.t.ScriptType {
	case "lua":
		l := &LuaWorker{Log: this.log.Write}
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
	r := w.GetRecords()

	rf, err := os.OpenFile(this.recordName, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		errs("打开记录文件失败:", err.Error())
		return
	}
	defer rf.Close()
	debug(r)
	data, err := json.Marshal(r)
	if err != nil {
		errs("写入记录错误:", err.Error())
	} else {
		rf.Truncate(0)
		rf.Write(data)
	}

}

//关闭执行器
func (this *TaskExecutor) ShutDown() {
	this.log.Close()
}

//
//读取记录
func (this *TaskExecutor) loadRecords() map[string]interface{} {
	records := make(map[string]interface{})
	if files.IsFileExist(this.recordName) {
		data, _ := ioutil.ReadFile(this.recordName)
		debug(string(data))
		json.Unmarshal(data, &records)
	}
	//从文件中读取记录
	return records
}

type Logger struct {
	file   *os.File
	Prefix string
}

func (this *Logger) Init(filepath string, prefix string) {
	this.Prefix = prefix
	if filepath != "" {
		this.file, _ = os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	}
}
func (this *Logger) Write(s string) {
	if this.file != nil {
		this.file.WriteString(fmt.Sprintf("[%s][%s] %s\n", time.Now().Format(bingo_utils.FORMAT_DATETIME_LOG), this.Prefix, s))
	}
}

func (this *Logger) Close() {
	if this.file != nil {
		this.file.Close()
	}
}
