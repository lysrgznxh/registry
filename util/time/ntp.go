package time

import (
	"agent-network-protocol/registry/core"
	"encoding/binary"
	"net"
	"time"
)

const ntpEpochOffset = 2208988800

type packet struct {
	Settings       uint8
	Stratum        uint8
	Poll           int8
	Precision      int8
	RootDelay      uint32
	RootDispersion uint32
	ReferenceID    uint32
	RefTimeSec     uint32
	RefTimeFrac    uint32
	OrigTimeSec    uint32
	OrigTimeFrac   uint32
	RxTimeSec      uint32
	RxTimeFrac     uint32
	TxTimeSec      uint32
	TxTimeFrac     uint32
}

func getRemoteTime() (time.Time, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmServNtp)
	var host = "time.windows.com:123"

	conn, err := net.Dial("udp", host)
	if err != nil {
		log.WithError(err).Error("err dialing ntp server")
		return time.Now(), err
	}
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(15 * time.Second)); err != nil {
		log.WithError(err).Error("err setting conn's deadline")
		return time.Now(), err
	}
	req := &packet{Settings: 0x1B}
	if err := binary.Write(conn, binary.BigEndian, req); err != nil {
		log.WithError(err).Error("err sending request")
		return time.Now(), err
	}
	rsp := &packet{}
	if err := binary.Read(conn, binary.BigEndian, rsp); err != nil {
		log.WithError(err).Error("err reading response")
		return time.Now(), err
	}
	secs := float64(rsp.TxTimeSec) - ntpEpochOffset
	nanos := (int64(rsp.TxTimeFrac) * 1e9) >> 32
	showtime := time.Unix(int64(secs), nanos)
	return showtime, nil
}

func GetTimeDiffWithNtp() (int64, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmServNtp)
	ntime, err := getRemoteTime()
	if err != nil {
		log.WithError(err).Error("err checking ntp server's time")
		return 0, err
	}
	timeMilDiff := time.Now().Sub(ntime).Milliseconds()
	if timeMilDiff > 30000 { // 时间差大于半分钟时提醒节点同步
		log.Error("time offset too large")
	}
	return timeMilDiff, nil
}
