package main

import "strconv"

type V2ray interface {
	Join([]V2ray) []V2ray
	Print() []string
}

type Option interface {
	apply(config *V2rayConfig)
}

type V2rayConfig struct {
	Log       LogV2ray        `json:"log"`
	Inbounds  []InboundV2ray  `json:"inbounds"`
	Outbounds []OutboundV2ray `json:"outbounds"`
}

func NewV2rayConfig(opts ...Option) V2rayConfig {
	var v2rayConfig V2rayConfig
	for _, opt := range opts {
		opt.apply(&v2rayConfig)
	}
	return v2rayConfig
}

// log
type LogV2ray struct {
	Loglevel string `json:"loglevel"`
}

func (l LogV2ray) apply(config *V2rayConfig) {
	config.Log = l
}

func withDefaultLog() LogV2ray {
	return LogV2ray{
		Loglevel: "DEBUG",
	}
}

// inbound
type InboundV2ray struct {
	Port     int                    `json:"port"`
	Protocol string                 `json:"protocol"`
	Sniffing InboundSniffingV2ray   `json:"sniffing"`
	Settings map[string]interface{} `json:"settings"`
}
type InboundSniffingV2ray struct {
	Enabled      bool     `json:"enabled"`
	DestOverride []string `json:"destOverride"`
}

func (i InboundV2ray) apply(config *V2rayConfig) {
	config.Inbounds = append(config.Inbounds, i)
}

func withDefaultInbound() InboundV2ray {
	return InboundV2ray{
		Port:     1080,
		Protocol: "socks",
		Sniffing: InboundSniffingV2ray{
			Enabled:      true,
			DestOverride: []string{"http", "tls"},
		},
		Settings: map[string]interface{}{
			"auth": "noauth",
		},
	}
}

// outbound
type OutboundV2ray struct {
	Protocol string                 `json:"protocol"`
	Settings map[string]interface{} `json:"settings"`
}

func (o OutboundV2ray) apply(config *V2rayConfig) {
	config.Outbounds = append(config.Outbounds, o)
}

func NewV2rayOutBound(protocol string, v2rays []V2ray) OutboundV2ray {
	return OutboundV2ray{
		Protocol: protocol,
		Settings: map[string]interface{}{
			"vnext": v2rays,
		},
	}
}

type VmessV2ray struct {
	Name    string           `json:"-"`
	Address string           `json:"address"`
	Port    int              `json:"port"`
	Users   []VmessUserV2ray `json:"users"`
}
type VmessUserV2ray struct {
	Id      string `json:"id"`
	AlterId int    `json:"alterId"`
}

func (vv VmessV2ray) Join(v2rays []V2ray) []V2ray {
	return append(v2rays, vv)
}

func (vv VmessV2ray) Print() []string {
	return []string{
		vv.Name,
		vv.Address,
		strconv.Itoa(vv.Port),
		vv.Users[0].Id,
		strconv.Itoa(vv.Users[0].AlterId),
	}
}
