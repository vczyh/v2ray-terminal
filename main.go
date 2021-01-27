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
	"path"
	"strings"
	"sync"
)

var (
	url             string
	v2ray           string
	v2rayConfigFile string
	logPath         string
	socksPort       int
)

var logger *log.Logger

func main() {
	flag.StringVar(&url, "url", "", "订阅链接")
	flag.StringVar(&v2ray, "v2ray", "", "v2ray可执行文件路径")
	flag.StringVar(&v2rayConfigFile, "v2rayConfig", "", "v2ray配置文件路径")
	flag.StringVar(&logPath, "logPath", "", "日志路径")
	flag.IntVar(&socksPort, "socksPort", 1080, "socks端口")
	flag.Parse()

	if url == "" {
		fmt.Print("订阅链接不能为空，请输入订阅链接：")
		_, _ = fmt.Scanln(&url)
	}

	if v2ray == "" {
		fmt.Print("v2ray路径不能为空，请输入v2ray可执行文件路径：")
		_, _ = fmt.Scanln(&v2ray)
	}

	if v2rayConfigFile == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("获取Home路径失败")
			os.Exit(1)
		}
		v2rayConfigFile = path.Join(homeDir, ".config/v2ray", "config.json")
	}

	logWriter, err := LogWriter(logPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	multiWriter := io.MultiWriter(os.Stdout, logWriter)

	logger = log.New(multiWriter, "[v2rayT] ", log.LUTC)

	content := subscribeContent(url)
	scanner := bufio.NewScanner(bytes.NewReader(content))
	chanVMESS := make(chan Vmess)
	chanSS := make(chan *SS)
	var wg sync.WaitGroup
	for scanner.Scan() {
		wg.Add(1)
		line := scanner.Text()
		go func(l string) {
			defer wg.Done()
			protocol := line[:strings.Index(l, ":")]
			switch protocol {
			case "ss":
				ss, err := parseSS(l)
				if err == nil {
					chanSS <- ss
				}
			case "ssr":
				//parseSSR(l)
			case "vmess":
				v, err := parseVMESS(l)
				if err == nil {
					chanVMESS <- v
				}
			default:
				logger.Printf("解析不了的协议：%s", protocol)
			}

		}(line)
	}

	// 将订阅转为 v2ray 格式
	done := make(chan bool)
	go func() {
		v2rayConfig := NewV2rayConfig(withDefaultLog(), withDefaultSocksInbound(socksPort), VMESSBound(chanVMESS))
		PrintV2rayOutbounds(v2rayConfig)
		err := WriteConfig(v2rayConfigFile, v2rayConfig)
		if err != nil {
			logger.Println(err)
			os.Exit(1)
		}
		// TODO cancel()
		ctx, _ := context.WithCancel(context.Background())
		err = ExeRealTimeOut(ctx, multiWriter, v2ray, "-config", v2rayConfigFile)
		if err != nil {
			logger.Println(err)
			os.Exit(1)
		}

		done <- true
	}()

	wg.Wait()
	close(chanVMESS)

	select {
	case <-done:
	}

}

func VMESSBound(ch <-chan Vmess) OutboundV2ray {
	var v2rays []V2ray
	for v := range ch {
		v2ray, err := v.toV2ray()
		if err != nil {
			logger.Println(err)
		}
		v2rays = v2ray.Join(v2rays)
	}
	outbound := NewVMessOutBound(v2rays)
	return outbound
}

func subscribeContent(url string) []byte {
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
