package main

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/aosfather/bingo_utils/lua"
	l "github.com/yuin/gopher-lua"
	"strconv"
	"strings"
)

//工作上下文
type Worker interface {
	SetInputs(inputs map[string]interface{})
	AddInput(key string, value interface{})
	LoadRecords(records map[string]interface{})
	GetRecords() map[string]interface{}
	RunScript(filename string)
}

//工作者
type LuaWorker struct {
	script  *lua.LuaScript
	outpath string
	inputs  map[string]interface{}
	records map[string]interface{}
	Log     lua.LuaLogFunction
}

func (this *LuaWorker) Init() {
	this.script = &lua.LuaScript{}
	debug("set log")
	if this.Log != nil {
		this.script.Log = this.Log
	} else {
		this.script.Log = this.log
	}

	this.script.Context = make(map[string]interface{})
	libs := make(map[string]l.LGFunction)
	//http访问接口
	libs["http_get"] = this.lua_http_get
	libs["http_post"] = this.lua_http_post
	////本地程序调用
	libs["os_run"] = this.lua_os_run
	libs["os_shell"] = this.lua_os_shell
	libs["os_cmd"] = this.lua_os_cmd
	//发送邮件
	libs["sendmail"] = this.lua_sendmail
	//html处理
	libs["html_doc"] = this.lua_html_doc
	//查找节点
	libs["html_find"] = this.lua_html_find
	//获取text内容
	libs["html_text"] = this.lua_html_text
	//获取子属性内容
	libs["html_attr"] = this.lua_html_attr
	//获取子属性内容
	libs["html_children_attr"] = this.lua_html_children_attr
	libs["pack"] = this.lua_pack
	this.script.SetPool(lua.NewLuaPool(10, "aoy", libs))
}
func (this *LuaWorker) SetInputs(inputs map[string]interface{}) {
	this.inputs = inputs
}
func (this *LuaWorker) AddInput(key string, value interface{}) {
	if key == "" || value == "" {
		return
	}
	if this.inputs == nil {
		this.inputs = make(map[string]interface{})
	}
	this.inputs[key] = value
}

func (this *LuaWorker) LoadRecords(records map[string]interface{}) {
	this.records = records
}

func (this *LuaWorker) GetRecords() map[string]interface{} {
	return this.records
}
func (this *LuaWorker) beforeRun(l *l.LState) {

	//压入输入参数，只读
	l.SetGlobal("_inputs", lua.SetReadOnly(l, lua.ToLuaTable2(l, this.inputs)))
	//压入记录
	l.SetGlobal("_records", lua.ToLuaTable2(l, this.records))
}

func (this *LuaWorker) afterRun(l *l.LState) {
	debug(this.records)
	//读取记录
	rec := l.GetGlobal("_records")
	this.records = lua.ToGoMap(rec)
}

func (this *LuaWorker) RunScript(filename string) {
	this.script.Loadfile(filename)
	this.script.Call(this.beforeRun, this.afterRun)
}

func (this *LuaWorker) log(s string) {
	debug("log----")
	debug(s)
}

//------系统功能调用--------//
func (this *LuaWorker) lua_os_run(l *l.LState) int {
	args := l.Get(-1)
	l.Pop(1)
	cmd := l.Get(-1).String()
	l.Pop(1)
	buffer := new(bytes.Buffer)
	RunCMD(cmd, buffer, lua.ToGoStringArray(args)...)
	l.Push(lua.ToLuaValue(buffer.String()))
	return 1
}

func (this *LuaWorker) lua_os_shell(l *l.LState) int {
	cmd := l.Get(-1).String()
	l.Pop(1)
	buffer := new(bytes.Buffer)
	RunCMDWithShell(cmd, buffer)
	l.Push(lua.ToLuaValue(buffer.String()))
	return 1
}

func (this *LuaWorker) lua_os_cmd(l *l.LState) int {
	args := l.Get(-1)
	l.Pop(1)
	cmd := l.Get(-1).String()
	l.Pop(1)
	buffer := new(bytes.Buffer)
	RunCMDWithDos(cmd, buffer, lua.ToGoStringArray(args)...)
	l.Push(lua.ToLuaValue(buffer.String()))
	return 1
}

//获取文本内容
func (this *LuaWorker) lua_http_get(l *l.LState) int {
	content := l.Get(-1).String()
	l.Pop(1)
	buffer := new(bytes.Buffer)
	err := DoGet(content, buffer, nil)
	if err != nil {
		l.Push(lua.ToLuaValue(""))
		l.Push(lua.ToLuaValue(err.Error()))
	} else {
		l.Push(lua.ToLuaValue(buffer.String()))
		l.Push(lua.ToLuaValue(""))
	}
	return 2
}

//post请求
func (this *LuaWorker) lua_http_post(l *l.LState) int {
	body := l.Get(-1).String()
	l.Pop(1)

	url := l.Get(-1).String()
	l.Pop(1)

	buffer := new(bytes.Buffer)
	err := DoPost(url, body, buffer, nil)
	if err != nil {
		l.Push(lua.ToLuaValue(""))
		l.Push(lua.ToLuaValue(err.Error()))
	} else {
		l.Push(lua.ToLuaValue(buffer.String()))
		l.Push(lua.ToLuaValue(""))
	}
	return 2
}

//下载

type mail struct {
	Host      string
	Port      int
	Pwd       string
	From      string
	To        string
	Subject   string
	Body      string
	FileNames []string
	FilePaths []string
}

func (this *mail) InitByValue(value l.LValue) {
	v := value.(*l.LTable)
	v.ForEach(func(key, value l.LValue) {
		keystr := key.String()
		v := value.String()
		switch keystr {
		case "host":
			this.Host = v
		case "port":
			this.Port, _ = strconv.Atoi(v)
		case "pwd":
			this.Pwd = v
		case "from":
			this.From = v
		case "to":
			this.To = v
		case "subject":
			this.Subject = v
		case "body":
			this.Body = v
		case "att_names":
			this.FileNames = strings.Split(v, ";")
		case "att_files":
			this.FilePaths = strings.Split(v, ";")

		}

	})
}

//发送邮箱,通过table传入参数
func (this *LuaWorker) lua_sendmail(l *l.LState) int {
	mv := l.Get(-1)
	l.Pop(1)
	m := mail{}
	m.InitByValue(mv)
	debug(m)
	b := SendMail(m.Host, m.Port, m.From, m.To, m.Subject, m.Pwd, m.Body, m.FileNames, m.FilePaths...)
	l.Push(lua.ToLuaValue(b))
	return 1
}

//合并文件:文件列表，及合并的size，返回合并后的文件列表。 pack(数组,大小)
func (this *LuaWorker) lua_pack(l *l.LState) int {
	//打包的大小。参数获取的出栈顺序是先入后出，先出栈的是后面的参数。所以取参数的顺序与调用输入的顺序是相反的。
	size := l.Get(-1).String()
	l.Pop(1)
	_size, _ := strconv.Atoi(size)

	//文件列表
	files := l.Get(-1)
	l.Pop(1)
	_files := lua.ToGoValue(files, lua.NewLuaOption()).([]interface{})
	var target []string
	for _, v := range _files {
		target = append(target, v.(string))
	}
	//打包文件
	bag := TxtBagMan{BagMan{Size: _size, OutDir: this.outpath}}
	l.Push(lua.StringArrayToLuaTable(l, bag.MakeBag(target)))
	return 1
}

//--------------------------HTML 相关函数--------------------------//
//html 读取doc
func (this *LuaWorker) lua_html_doc(l *l.LState) int {
	url := l.Get(-1).String()
	l.Pop(1)
	doc := GetDocument(url)
	if doc != nil {
		this.script.Context["_doc"] = doc
		l.Push(lua.ToLuaValue(true))
	} else {
		l.Push(lua.ToLuaValue(false))
	}
	return 1
}

//查询node
func (this *LuaWorker) lua_html_find(l *l.LState) int {
	htmlpath := l.Get(-1).String()
	l.Pop(1)
	doc := this.script.Context["_doc"]
	if doc != nil {
		node := DocFind(doc.(*goquery.Document), htmlpath)
		if node != nil {
			this.script.Context["_node"] = node
			l.Push(lua.ToLuaValue(true))
			return 1
		}
	}

	l.Push(lua.ToLuaValue(false))

	return 1
}

//将当前节点转换成文本
func (this *LuaWorker) lua_html_text(l *l.LState) int {
	node := this.script.Context["_node"]
	if node != nil {
		l.Push(lua.ToLuaValue(NodeText(node.(*goquery.Selection))))
		return 1
	}
	l.Push(lua.ToLuaValue(""))
	return 1
}

func (this *LuaWorker) lua_html_attr(l *l.LState) int {
	attr := l.Get(-1).String()
	l.Pop(1)
	node := this.script.Context["_node"]
	if node != nil {
		list := NodeChildAttrToArray(node.(*goquery.Selection), attr)

		l.Push(lua.StringArrayToLuaTable(l, list))
		return 1
	}
	l.Push(lua.ToLuaValue(""))
	return 1
}

func (this *LuaWorker) lua_html_children_attr(l *l.LState) int {
	attr := l.Get(-1).String()
	l.Pop(1)
	node := this.script.Context["_node"]
	debug(node.(*goquery.Selection).Nodes)
	if node != nil {
		list := NodesChildrenAttrToArray(node.(*goquery.Selection), attr)
		l.Push(lua.StringArrayToLuaTable(l, list))
		return 1
	}
	l.Push(lua.ToLuaValue(""))
	return 1
}
