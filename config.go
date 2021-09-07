package gsn

import (
	"time"
)

var GlobalConfig *config

func init() {
	GlobalConfig = &config{
		HeartbeatMode:     false,
		HeartbeatInterval: 5 * time.Second,
	}
}

type config struct {
	HeartbeatMode     bool
	HeartbeatInterval time.Duration
}

func (g *config) SetHeartBit(heartBitOn bool, interval time.Duration) {
	g.HeartbeatMode = heartBitOn
	g.HeartbeatInterval = interval
}
