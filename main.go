package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var logger *log.Logger

func main() {
	var url string
	var v2ray string
	var logPath string

	flag.StringVar(&url, "url", "", "订阅链接")
	flag.StringVar(&v2ray, "v2ray", "", "v2ray可执行文件路径")
	flag.StringVar(&logPath, "logPath", "", "日志路径")
	flag.Parse()

	if url == "" {
		fmt.Print("订阅链接不能为空，请输入订阅链接：")
		_, _ = fmt.Scanln(&url)
	}

	if v2ray == "" {
		fmt.Print("v2ray路径不能为空，请输入v2ray可执行文件路径：")
		_, _ = fmt.Scanln(&v2ray)
	}

	logWriter, err := LogWriter(logPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 输出
	multiWriter := io.MultiWriter(os.Stdout, logWriter)
	logger = log.New(multiWriter, "v3ray", log.LUTC)

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
				logger.Printf("解析不了的协议：%s", protocol)
			}

		}(line)
	}

	// 将订阅转为 v2ray 格式
	done := make(chan bool)
	go func() {
		// vmess
		v2rayConfig := NewV2rayConfig(withDefaultLog(), withDefaultInbound(), vmessBound(chanVmess))
		PrintV2rayOutbounds(v2rayConfig)
		err := WriteConfig(v2rayConfig)
		if err != nil {
			logger.Println(err)
			os.Exit(1)
		}
		ctx, cancelFunc := context.WithCancel(context.Background())
		time.AfterFunc(5*time.Hour, func() {
			cancelFunc()
		})
		err = ExeRealTimeOut(ctx, multiWriter, v2ray, "-config", "/Users/zhangyuheng/.config/v2ray/config.json")
		if err != nil {
			logger.Println(err)
			os.Exit(1)
		}

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
			logger.Println(err)
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
		logger.Fatal("不能访问订阅链接")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Fatal("读取订阅内容失败")
	}
	content, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		logger.Fatal("base64解码订阅内容失败")
	}
	return content
}
