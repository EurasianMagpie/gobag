package icmp

import (
	"errors"
	"fmt"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

func RunPinger(host string) (string, error) {
	pinger, err := probing.NewPinger(host)
	if err != nil {
		return "", err
	}

	pinger.Timeout = time.Second * 10

	// https://github.com/go-ping/ping/issues/17
	// Windows doesn't allow ICMP messages over UDP. Setting SetPrivileged(true) on Windows
	// will workaround this problem by using ICMP over IP.
	pinger.SetPrivileged(true)

	pinger.Count = 4
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		return "", err
	}
	stats := pinger.Statistics()
	r := fmt.Sprintf("Pinger Stats [Sent:%d  Recv:%d  Loss:%.0f%%]", stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)

	if stats.PacketsRecv > 0 {
		return r, nil
	} else {
		return "", errors.New(r)
	}
	return r, nil
}
