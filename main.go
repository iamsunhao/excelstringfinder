package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func main() {
	s := "地"
	dir := "/home/sun/Code/"
	flist := listFiles(dir)
	wg := sync.WaitGroup{}
	for _, f := range flist {
		wg.Add(1)
		go func(s, f string) {
			defer wg.Done()
			FindStringInXLSX(s, f)
		}(s, f)
	}
	wg.Wait()

}

// listFiles 通过递归的方式将指定目录下的所有文件写入到一个切片中
func listFiles(dir string) []string {
	// 判断输入的路径是否合法
	dirInfo, err := os.Stat(dir)
	if err != nil {
		log.Fatal(err)
	}
	if !dirInfo.IsDir() {
		log.Fatal("not a path")
	}
	// 新建一个长度为0的切片，用于返回文件列表
	ans := make([]string, 0)
	// 创建一个匿名函数，用于递归将文件添加到之前新建的切片中
	var addFileToList func(string)
	addFileToList = func(dir string) {
		items, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Fatal(err)
		}
		for _, item := range items {
			itemName := dir + "/" + item.Name()
			if !item.IsDir() {
				ans = append(ans, itemName)
			} else {
				addFileToList(itemName)
			}
		}
	}
	addFileToList(dir)
	return ans
}

// FindStringInXLSX 在一个 xlsx 文件中查找一个字符串
func FindStringInXLSX(s string, filename string) {
	// 对文件进行初步的检查，确定是否是 xlsx 文件
	if len(filename) < 4 {
		return
	}
	if filename[len(filename)-4:] != "xlsx" {
		return
	}
	f, err := excelize.OpenFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	rows := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	for i, row := range rows {
		for j, cell := range row {
			if strings.Contains(cell, s) {
				fmt.Printf("在 %s 的第 %d 行 %d 列中找到了 %s 。\n", filename, i+1, j+1, s)
			}
		}
	}
}
