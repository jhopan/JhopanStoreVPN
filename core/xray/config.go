package xray

import (
	"encoding/json"
	"fmt"
	"strconv"

	"jhovpn/core/vless"
)

// GenerateConfig builds a Xray JSON config from a VLESS Config.
// dns1/dns2: DNS servers to use. If empty, defaults to Google DNS.
// allowInsecure: skip TLS certificate verification.
func GenerateConfig(vc vless.Config, dns1, dns2 string, allowInsecure bool) ([]byte, error) {
	if dns1 == "" {
		dns1 = "8.8.8.8"
	}
	if dns2 == "" {
		dns2 = "8.8.4.4"
	}

	domain, portStr, err := vless.SplitAddress(vc.Address)
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid port: %s", portStr)
	}

	config := map[string]interface{}{
		"log": map[string]interface{}{
			"loglevel": "none",
		},
		"dns": map[string]interface{}{
			"servers": []string{dns1, dns2},
		},
		"inbounds": []map[string]interface{}{
			{
				"tag":      "socks-in",
				"port":     10808,
				"listen":   "127.0.0.1",
				"protocol": "socks",
				"settings": map[string]interface{}{
					"udp": true,
				},
				// sniffing disabled: all traffic routes through proxy, no rule-based override needed
			},
		},
		"outbounds": []map[string]interface{}{
			{
				"tag":      "proxy",
				"protocol": "vless",
				"settings": map[string]interface{}{
					"vnext": []map[string]interface{}{
						{
							"address": domain,
							"port":    port,
							"users": []map[string]interface{}{
								{
									"id":         vc.UUID,
									"encryption": "none",
									"level":      0,
								},
							},
						},
					},
				},
				"streamSettings": map[string]interface{}{
					"network":  "ws",
					"security": "tls",
					"tlsSettings": map[string]interface{}{
						"serverName":    vc.SNI,
						"allowInsecure": allowInsecure,
					},
					"wsSettings": map[string]interface{}{
						"path": vc.Path,
						"host": vc.Host,
					},
					// sockopt: tune the Xray→server TCP connection
					"sockopt": map[string]interface{}{
						"tcpFastOpen":          true, // -1 RTT on reconnect
						"tcpNoDelay":           true, // disable Nagle → lower latency
						"tcpKeepAliveIdle":     60,   // start probes after 60s idle
						"tcpKeepAliveInterval": 30,   // probe every 30s
					},
				},
			},
			{
				"tag":      "direct",
				"protocol": "freedom",
				"settings": map[string]interface{}{},
			},
		},
		"routing": map[string]interface{}{
			"domainStrategy": "AsIs",
			"rules": []map[string]interface{}{
				{
					"type": "field",
					"ip": []string{
						"10.0.0.0/8",
						"172.16.0.0/12",
						"192.168.0.0/16",
						"127.0.0.0/8",
					},
					"outboundTag": "direct",
				},
			},
		},
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}
	return data, nil
}
