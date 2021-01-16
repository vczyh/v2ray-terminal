package main

import (
	"github.com/pterm/pterm"
	"strconv"
)

func PrintV2rayOutbounds(config V2rayConfig) {
	outbounds := config.Outbounds
	var table = pterm.TableData{
		{"Index", "Name", "Address", "Port", "ID", "Alter ID"},
	}

	index := 1
	for _, outbound := range outbounds {
		switch outbound.Protocol {
		case "vmess":
			for _, vmess := range outbound.Settings["vnext"].([]V2ray) {
				line := append([]string{strconv.Itoa(index)}, vmess.Print()...)
				table = append(table, line)
				index++
			}
		}

	}
	pterm.DefaultTable.WithHasHeader().WithData(table).Render()
}
