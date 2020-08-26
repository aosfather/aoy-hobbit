package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/smtp"
	"strings"
	"time"
)

// 发送邮件
func SendMail(host string, port int, from, to, subject, pwd string, body string, names []string, filename ...string) bool {
	sc := SmtpClient{}
	sc.Host = host
	sc.User = from
	sc.Pwd = pwd
	sc.Port = port

	msg := EmailMessage{}
	msg.From = sc.User
	msg.To = strings.Split(to, ",")
	msg.Subject = subject
	msg.Body = body

	for index, f := range filename {
		msg.AddAttachment(Attachment{CT_TEXT, names[index], f})
	}

	return sc.SendMessage(&msg)
}

const (
	CT_HTML ContentType = "html"
	CT_TEXT ContentType = "plain"
)

//内容类型：包括html\plain\
type ContentType string

func (this ContentType) ToString() string {
	return "Content-Type: text/" + fmt.Sprintf("%s", this) + "; charset=UTF-8"
}

//smtp客户端，完成登录及发送邮件
type SmtpClient struct {
	Host string
	User string
	Pwd  string
	Port int
	auth smtp.Auth
}

//发送邮件
func (this *SmtpClient) SendMessage(m *EmailMessage) bool {
	if this.auth == nil {
		this.auth = smtp.PlainAuth("", this.User, this.Pwd, this.Host)
	}

	e := smtp.SendMail(fmt.Sprintf("%s:%d", this.Host, this.Port), this.auth, m.From, m.To, m.ToByte())
	if e != nil {
		log.Println(e)
		return false
	}
	return true
}

//邮件
type EmailMessage struct {
	Type    ContentType
	Charset string       //字符编码
	From    string       //发件人
	To      []string     //接收者
	Cc      []string     //接收者
	Bcc     []string     //接收者
	Subject string       //标题
	Body    string       //邮件内容
	attachs []Attachment //附件
}

//构建邮件头
func (this *EmailMessage) toHeader(boundary string) map[string]string {
	Header := make(map[string]string)
	Header["From"] = this.From
	Header["To"] = strings.Join(this.To, ";")
	if len(this.Cc) > 0 {
		Header["Cc"] = strings.Join(this.Cc, ";") //抄送
	}

	if len(this.Bcc) > 0 {
		Header["Bcc"] = strings.Join(this.Bcc, ";") //密送
	}

	Header["Subject"] = this.Subject
	Header["Content-Type"] = "multipart/mixed;boundary=" + boundary
	Header["Date"] = time.Now().String()
	return Header

}

//增加附件
func (this *EmailMessage) AddAttachment(a Attachment) {
	if a.FileName != "" {
		this.attachs = append(this.attachs, a)
	}

}

//生成邮件发送格式体
func (this *EmailMessage) ToByte() []byte {
	buffer := bytes.NewBuffer(nil)
	//唯一标识
	boundary := fmt.Sprintf("send4kindle_%d", time.Now().UnixNano())

	//写header头 key：value的形式
	headers := this.toHeader(boundary)
	header := ""
	for key, value := range headers {
		header += key + ":" + value + "\r\n"
	}
	header += "\r\n"
	buffer.WriteString(header)

	//写正文内容
	body := "\r\n--" + boundary + "\r\n"
	body += "Content-Transfer-Encoding:base64\r\n"
	body += "Content-Type: text/html; charset=UTF-8 \r\n"

	buffer.WriteString(body)
	writeAsBase64([]byte(this.Body), buffer)

	//写附件
	for _, a := range this.attachs {
		a.Write(boundary, buffer)
	}

	return buffer.Bytes()
}

//附件
type Attachment struct {
	Type     ContentType //类型
	Name     string      //显示用的名称
	FileName string      //真实文件地址
}

//写入附件的头及文件内容(使用base64)
func (this *Attachment) Write(boundary string, buffer *bytes.Buffer) {
	filename := mime.BEncoding.Encode("utf-8", this.Name)
	attachment := "\r\n--" + boundary + "\r\n"
	attachment += "Content-Disposition: attachment;filename=\"" + filename + "\"\r\n"
	attachment += "Content-Transfer-Encoding:base64\r\n"
	attachment += this.Type.ToString() + ";name=\"" + filename + "\"\r\n"
	//内嵌用的语法
	//attachment += "Content-ID: <" + this.Name + "> \r\n\r\n"
	buffer.WriteString(attachment)
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf(err.(string))
		}
	}()

	this.writeFile(buffer)

}

//写入附件文件内容，使用base64编码
func (this *Attachment) writeFile(buffer *bytes.Buffer) {
	file, err := ioutil.ReadFile(this.FileName)
	if err != nil {
		panic(err.Error())
	}

	writeAsBase64(file, buffer)
}

//将内容使用base64格式写入到buffer中
func writeAsBase64(file []byte, buffer *bytes.Buffer) {
	payload := make([]byte, base64.StdEncoding.EncodedLen(len(file)))
	base64.StdEncoding.Encode(payload, file)
	buffer.WriteString("\r\n")
	for index, line := 0, len(payload); index < line; index++ {
		buffer.WriteByte(payload[index])
		if (index+1)%76 == 0 {
			buffer.WriteString("\r\n")
		}
	}

}
