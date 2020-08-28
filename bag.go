package main

import (
	"container/list"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

/**
  打包
*/
type BagMan struct {
	Size   int    //包的大小
	OutDir string //输出目录
}

type TxtBagMan struct {
	BagMan
}

func (this *TxtBagMan) MakeBag(filename []string) []string {
	now := time.Now().Unix()
	var files *list.List
	var bags []*list.List
	for index, name := range filename {
		if index%this.Size == 0 {
			files = list.New()
			bags = append(bags, files)
		}
		files.PushBack(name)
	}
	log.Println(len(bags))
	var result []string
	for index, files := range bags {
		filename := fmt.Sprintf("%s/%d_%d.txt", this.OutDir, now, index)
		result = append(result, filename)
		file, err := os.Create(filename)
		if err != nil {
			return nil
		}
		for e := files.Front(); e != nil; e = e.Next() {
			f, err := os.Open(e.Value.(string))
			if err != nil {
				log.Println("open file error:", err.Error())
			}
			io.Copy(file, f)
			defer f.Close()
		}
		defer file.Close()

	}

	return result

}
