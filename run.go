package main

/**
  执行逻辑
   1、根据配置文件，连接到指定的服务端
   2、获取任务
   3、执行任务
   4、上传任务进度
*/
func main() {
	//

}

//任务连接配置
type Config struct {
	Name     string //执行者名称
	User     string //所属用户
	PoolName string //用户下任务池名称
	Token    string //连接用的Token
}
