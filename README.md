# aoy-hobbit
hobbit worker。the lightway worker client implements。the remote agent programe，can run shell,send mail,let your set as a robot,by your command!
# 奥义-霍比特人
一个远程执行任务的工作软件，属于自动化控制的远程客户端。它接收任务指令，并执行，与以往的远程控制客户端不同的是，其主动查询任务服务端获取执行的指令，而不是对外开放端口作为控制的服务端来完成工作。
计划实现功能：
 * 执行shell命令（可以执行指定的程序）
 * 下载文件
 * 文件合并
 * 对指定的网络端口进行调用(GET\POST）
 * 发送邮件
 * 获取任务执行参数[每次执行都可以不同]
 * 任务执行结果保存[任务全局]
 
## 任务模型定义
一个任务由以下几个部分构成
 * 任务code。作为唯一标识
 * 任务名称。用于显示
 * 任务对应的脚本。每次执行都使用最新的版本
 * 任务执行的参数。该次任务加载进去的参数和上下文。对于脚本而言只读。
 
