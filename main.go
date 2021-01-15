package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

func main() {
	const url = "https://sub.paasmi.com/subscribe/58648/pzODI9eUg3vO?mode=3"

	content := subsribeContent(url)

	// 解析订阅
	scanner := bufio.NewScanner(bytes.NewReader(content))
	chanVmess := make(chan Vmess)
	var wg sync.WaitGroup
	for scanner.Scan() {
		wg.Add(1)
		line := scanner.Text()
		go func(l string) {
			defer wg.Done()
			protocol := line[:strings.Index(l, ":")]
			switch protocol {
			case "ss":
				//parseSS(l)
			case "ssr":
				//parseSSR(l)
			case "vmess":
				v, err := parseVMESS(l)
				if err == nil {
					//fmt.Println(v)
					chanVmess <- v
				}
			default:
				log.Printf("解析不了的协议：%s", protocol)
			}

		}(line)
	}

	// 将订阅转为 v2ray 格式
	done := make(chan bool)
	go func() {
		// vmess
		v2rayConfig := NewV2rayConfig(withDefaultLog(), withDefaultInbound(), vmessBound(chanVmess))
		//config, _ := json.MarshalIndent(v2rayConfig, "", "  ")
		//fmt.Println(string(config))
		PrintV2rayOutbounds(v2rayConfig)

		done <- true
	}()

	wg.Wait()
	close(chanVmess)

	select {
	case <-done:
	}

}

func vmessBound(ch <-chan Vmess) OutboundV2ray {
	var v2rays []V2ray
	for v := range ch {
		v2ray, err := v.toV2ray()
		if err != nil {
			log.Println(err)
		}
		v2rays = v2ray.Join(v2rays)
	}
	// 构建outbound
	outbound := NewV2rayOutBound("vmess", v2rays)

	return outbound
}

func subsribeContent(url string) []byte {
	res, err := http.Get(url)
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
	return content
}
