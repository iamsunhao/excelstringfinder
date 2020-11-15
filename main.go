package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	dir := "/home/sun/Code"
	flist := listFiles(dir)
	for _, j := range flist {
		fmt.Println(j)
	}

}

// listFiles 通过递归的方式将指定目录下的所有文件写入一个切片中
func listFiles(dir string) []string {
	// 判断输入的路径是否合法
	dirInfo, err := os.Stat(dir)
	if err != nil {
		log.Fatal(err)
	}
	if !dirInfo.IsDir() {
		log.Fatal("not a path")
	}
	// not end ....
	ans := make([]string, 0)
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
