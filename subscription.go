package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"strconv"
	"strings"
)

type Subscriber interface {
	toV2ray() (V2ray, error)
}

type Vmess struct {
	V    string `json:"v"`
	Ps   string `json:"ps"`
	Add  string `json:"add"`
	Port string `json:"port"`
	Id   string `json:"id"`
	Aid  int    `json:"aid"`
	Net  string `json:"net"`
	Type string `json:"type"`
	Host string `json:"host"`
	Path string `json:"path"`
	Tls  string `json:"tls"`
}

func parseVMESS(str string) (Vmess, error) {
	cont := stripProtocol(str)
	var v Vmess

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
	//fmt.Println(v)
	return v, nil
}

func (v *Vmess) toV2ray() (V2ray, error) {
	port, err := strconv.Atoi(v.Port)
	if err != nil {
		return nil, err
	}
	return VmessV2ray{
		Address: v.Add,
		Port:    port,
		Users:   []VmessUserV2ray{{v.Id, v.Aid}},
	}, nil
}

// tool
func stripProtocol(str string) string {
	return str[strings.Index(str, "//")+2:]
}
