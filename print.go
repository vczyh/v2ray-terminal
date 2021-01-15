package main

import (
	"github.com/pterm/pterm"
)

func PrintV2rayOutbounds(config V2rayConfig) {
	outbounds := config.Outbounds
	var table = pterm.TableData{
		{"address", "port", "id", "alterId"},
	}

	for _, outbound := range outbounds {
		switch outbound.Protocol {
		case "vmess":
			for _, vmess := range outbound.Settings["vnext"].([]V2ray) {
				table = append(table, vmess.Print())
			}
		}
	}
	pterm.DefaultTable.WithHasHeader().WithData(table).Render()
}
