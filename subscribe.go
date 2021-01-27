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

	err = json.Unmarshal(decode, &v)
	if err != nil {
		log.Printf("反序列化json失败，%s", err)
		return v, err
	}
	return v, nil
}

func (v *Vmess) toV2ray() (V2ray, error) {
	port, err := strconv.Atoi(v.Port)
	if err != nil {
		return nil, err
	}
	return VmessV2ray{
		Name:    v.Ps,
		Address: v.Add,
		Port:    port,
		Users:   []VmessUserV2ray{{v.Id, v.Aid}},
	}, nil
}

type SS struct {
	Name     string
	Server   string
	Port     int
	Cipher   string
	Password string
}

func parseSS(ss string) (*SS, error) {
	ss = stripProtocol(ss)

	// auth
	auth := ss[:strings.Index(ss, "@")]
	auth = auth + "="
	server := ss[strings.Index(ss, "@")+1 : strings.Index(ss, ":")]
	decode, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		logger.Printf("base64解码失败，%s", auth)
		return nil, err
	}
	auths := strings.Split(string(decode), ":")

	//port
	port := ss[strings.Index(ss, ":")+1 : strings.Index(ss, "/?")]
	p, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}

	return &SS{
		Server:   server,
		Port:     p,
		Cipher:   auths[0],
		Password: auths[1],
	}, nil
}

func (s *SS) toV2ray() (V2ray, error) {
	// todo
	return
}

func stripProtocol(str string) string {
	return str[strings.Index(str, "//")+2:]
}
