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
	if len(os.Args) <= 2 {
		fmt.Printf("Usage: %s string directory\nExample: %s hello .\n", os.Args[0], os.Args[0])
		return
	}
	s := os.Args[1]
	flist := listFiles(os.Args[2])
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

// listFiles 通过递归的方式将指定目录下的所有xlsx文件写入到一个切片中
func listFiles(dir string) []string {
	// 判断输入的路径是否合法
	dirInfo, err := os.Stat(dir)
	if err != nil {
		log.Fatal(err)
	}
	if !dirInfo.IsDir() {
		log.Fatal("not a path")
	}
	// 新建一个长度为0，容量为500的切片，用于返回文件列表
	ans := make([]string, 0, 500)
	// 创建一个匿名函数，用于递归将xlsx文件添加到之前新建的切片中
	var addFileToList func(string)
	addFileToList = func(dir string) {
		items, err := ioutil.ReadDir(dir)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, item := range items {
			itemName := dir + "/" + item.Name()
			// 如果是一个 xlsx 文件，则加入到切片中
			if !item.IsDir() && len(itemName) >= 4 && itemName[len(itemName)-4:] == "xlsx" {
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
	f, err := excelize.OpenFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 新建一个空字符串，用于一次性返回所有结果
	ans := ""
	// 遍历每一个 sheet
	for _, sheet := range f.GetSheetMap() {
		rows := f.GetRows(sheet)
		if err != nil {
			fmt.Println(err)
			return
		}
		for i, row := range rows {
			for j, cell := range row {
				if strings.Contains(cell, s) {
					ans = ans + fmt.Sprintf("在 %s 中 %s 的第 %d 行 %d 列找到了 %s 。\n", filename, sheet, i+1, j+1, s)
				}
			}
		}
	}
	if len(ans) != 0 {
		fmt.Print(ans)
	}
}
