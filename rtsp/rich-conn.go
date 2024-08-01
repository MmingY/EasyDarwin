package rtsp

import (
	"net"
	"time"
)

type RichConn struct {
	net.Conn
	timeout time.Duration
}

func (conn *RichConn) Read(b []byte) (n int, err error) {
	if conn.timeout > 0 {
		location, _ := time.LoadLocation("Asia/Shanghai")
		conn.Conn.SetReadDeadline(time.Now().In(location).Add(conn.timeout))
	} else {
		var t time.Time
		conn.Conn.SetReadDeadline(t)
	}
	return conn.Conn.Read(b)
}

func (conn *RichConn) Write(b []byte) (n int, err error) {
	if conn.timeout > 0 {
		location, _ := time.LoadLocation("Asia/Shanghai")
		conn.Conn.SetWriteDeadline(time.Now().In(location).Add(conn.timeout))
	} else {
		var t time.Time
		conn.Conn.SetWriteDeadline(t)
	}
	return conn.Conn.Write(b)
}
