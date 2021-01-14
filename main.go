package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

func main() {

	res, err := http.Get("https://sub.paasmi.com/subscribe/58648/pzODI9eUg3vO?mode=3")
	if err != nil {
		log.Fatal("不能访问订阅链接")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("读取订阅内容失败")
	}
	content, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		log.Fatal("base64解码订阅内容失败")
	}
	//fmt.Println(string(content))

	chanVmess := make(chan vmess)

	var wg sync.WaitGroup

	scanner := bufio.NewScanner(bytes.NewReader(content))

	for scanner.Scan() {
		wg.Add(1)
		line := scanner.Text()
		go func(l string) {
			defer wg.Done()
			protocol := line[:strings.Index(l, ":")]
			switch protocol {
			case "ss":
				//ParseSS(l)
			case "ssr":
				//ParseSSR(l)
			case "vmess":
				v, err := ParseVMESS(l)
				if err == nil {
					chanVmess <- v
				}
			default:
				log.Printf("解析不了的协议：%s", protocol)
			}

		}(line)
	}

	done := make(chan bool)
	go func() {
		for vms := range chanVmess {
			fmt.Println(vms)
		}
		done <- true
	}()

	wg.Wait()

	close(chanVmess)

	select {
	case <-done:
	}

}

func ParseSS(str string) {

}

func ParseSSR(str string) {

}

type vmess struct {
	V    string `json:"v"`
	Ps   string `json:"ps"`
	Add  string `json:"add"`
	Port string `json:"port"`
	Id   string `json:"id"`
	Aid  int `json:"aid"`
	Net  string `json:"net"`
	Type string `json:"type"`
	Host string `json:"host"`
	Path string `json:"path"`
	Tls  string `json:"tls"`
}

func ParseVMESS(str string) (vmess, error) {
	cont := StripProtocol(str)
	var v vmess

	decode, err := base64.StdEncoding.DecodeString(cont)
	if err != nil {
		log.Printf("base64解码失败，%s", cont)
		return v, err
	}

	//fmt.Println(string(decode))

	err = json.Unmarshal(decode, &v)
	if err != nil {
		log.Printf("反序列化json失败，%s", err)
		return v, err
	}
	fmt.Println(v)
	return v, nil
}

func StripProtocol(str string) string {
	return str[strings.Index(str, "//")+2:]
}
