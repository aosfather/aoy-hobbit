package main

import "github.com/aosfather/bingo_utils/lua"

//工作者
type LuaWorker struct {
	script *lua.LuaScript
}

func (this *LuaWorker) Init() {
	this.script = &lua.LuaScript{}
	this.script.SetPool(lua.NewLuaPool(10, nil))
	this.script.Log = this.log
}

func (this *LuaWorker) RunScript(filename string) {
	this.script.Loadfile(filename)
	this.script.Call()
}

func (this *LuaWorker) log(s string) {
	debug(s)
}
