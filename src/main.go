package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type SourceList struct {
	name string
	url  string
}

var sourceList []SourceList = []SourceList{
	{"Tuna mirror", "mirrors.tuna.tsinghua.edu.cn"},
	{"Aliyun mirror", "mirrors.aliyun.com"},
}

func getProcessOwner() string {
	stdout, err := exec.Command("ps", "-o", "user=", "-p", strconv.Itoa(os.Getpid())).Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	data := strings.ReplaceAll(string(stdout), " ", "")
	data = strings.ReplaceAll(data, "\n", "")
	return data
}

func getCmdResult(name string, arg ...string) string {

	cmd := exec.Command(name, arg...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("combined out:\n%s\n", string(out))
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	return strings.Trim(string(out), "\n")
}

func check() {
	support_os := []string{"Debian", "Ubuntu"}
	system_dis := getCmdResult("cat", "/etc/issue")
	system_dis = strings.Split(system_dis, " ")[0]
	for index, val := range support_os {
		if strings.Contains(system_dis, val) {
			break
		}

		if len(support_os)-1 == index {
			fmt.Printf("Does not support current operating system: %s\n", system_dis)
			os.Exit(0)
		}

	}

	user := getProcessOwner()
	fmt.Println(user)
	if user != "root" {
		fmt.Println(user != "root")
		fmt.Println("Need to run as root user: " + user)
		os.Exit(0)
	}

}

func readFile2Lines(filepath string) []string {

	fi, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return nil
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	lines := make([]string, 0)
	for {
		line, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		lineStr := string(line)
		lineStr = strings.Trim(lineStr, " ")
		if lineStr == "" {
			continue
		}

		lines = append(lines, string(lineStr))
	}
	return lines
}

func replaceUrl(old string, new string) string {
	// fmt.Println(*old)

	replaceStr := strings.Split(old, `://`)[1]
	replaceStr = strings.Split(replaceStr, `/`)[0]
	newStr := strings.Replace(old, replaceStr, new, -1)
	old = newStr
	return old

}

func updateNewMirrorsUrl(data []string, url string) string {
	for index, val := range data {
		if val[0:1] == "#" || val == "" {
			continue
		}
		data[index] = replaceUrl(val, url)
	}
	newData := strings.Join(data, "\n")
	return newData
}

func writeFile(filepath string, data string) {
	ioutil.WriteFile(filepath, []byte(data), 0664)
}

func init() {
	check()
	for index, val := range sourceList {
		fmt.Printf("# %d %s\n", index+1, val.name)
	}

}

func main() {
	var selectID int
	var confirm string
	fmt.Print("Please select mirror id, enter 0 exit\n>:")
	fmt.Scan(&selectID)
	if selectID == 0 {
		os.Exit(0)
	}
	filepath := "/etc/apt/sources.list"
	lines := readFile2Lines(filepath)
	data := updateNewMirrorsUrl(lines, sourceList[selectID-1].url)
	fmt.Println("/etc/apt/sources.list will be updated with the following text content, please check.")
	fmt.Println("---------------------------")
	fmt.Println(data)
	fmt.Println("---------------------------")
	fmt.Print("Please enter \"yes\" to confirm the update, enter \"no\" exit\n>:")
	fmt.Scan(&confirm)
	if strings.ToLower(confirm) == "yes" {
		writeFile(filepath, data)
	}

}
