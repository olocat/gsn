package gsn

import (
	"time"
)

var GlobalConfig *config

func init() {
	GlobalConfig = &config{
		HeartbeatMode:     false,
		HeartbeatInterval: 5 * time.Second,

		MaxConnLimit: 65535,
	}
}

type config struct {
	HeartbeatMode     bool
	HeartbeatInterval time.Duration

	MaxConnLimit int
}

func (g *config) SetHeartBit(heartBitOn bool, interval time.Duration) {
	g.HeartbeatMode = heartBitOn
	g.HeartbeatInterval = interval
}
