m=require("aoy")
-- 下载文件
function downloadDoc(url,name)
  file = io.open(name, "w")
  io.output(file)
  m.html_doc(url.."/"..name)
  -- 写标题
  m.html_find("div.bookname h1")
  io.write(m.html_text())
  -- 写内容
  m.html_find("div#content")
  io.write(m.html_text())

  io.flush()
  io.close(file)

end

-- 获取章节列表
function getCaptions(url)
   if m.html_doc(url) then
     b=m.html_find("dl dd")
     if b then
       print(b)
       return m.html_attr("href")
     end
   end
   return {}
end
-- 合并文件

-- 发送到指定邮箱
function synckindle(att_name,att_file)
  mymail={}
  mymail["host"]=_inputs.host
  mymail["port"]=_inputs.port
  mymail["from"]=_inputs.from
  mymail["pwd"]=_inputs.pwd
  mymail["to"]=_inputs.to
  mymail["subject"]="书"
  mymail["body"]="同步的书籍"
  mymail["att_names"]=att_name
  mymail["att_files"]=att_file
  return m.sendmail(mymail)
end
-- 小说网下载同步--
url="http://www.wuxianxs.com/ls10-10050/"
--l=getCaptions(url)
--print(l)
--print(l[1])
--print(l[2])
--l={"12","234"}
--m.pack(l,10)
--downloadDoc(url,l[1])
--downloadDoc(url,l[2])
--nl={l[1],l[2]}
--m.pack(nl,2)
print(synckindle("2026485.txt","3026485.html"))
--print(table.getn(l))







