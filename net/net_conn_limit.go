package net

import (
	"github.com/zhangheng0027/ratelimit-plus"
	"io"
	"net"
	"time"
)

// DefaultReadLimit 20 MB/s
const DefaultReadLimit = 20 * 1024 * 1024

// DefaultWriteLimit 20 MB/s
const DefaultWriteLimit = 20 * 1024 * 1024

type netConnLimit struct {
	conn net.Conn
	r    io.Reader
	w    io.Writer
}

func ToLimitNetConn(conn net.Conn, rBucket ratelimit.BucketI, wBucket ratelimit.BucketI) net.Conn {
	n := &netConnLimit{
		conn: conn,
	}
	if rBucket == nil {
		n.r = ratelimit.Reader(conn, ratelimit.NewBucketWithRate(DefaultReadLimit, DefaultReadLimit*10))
	} else {
		n.r = ratelimit.Reader(conn, rBucket)
	}

	if wBucket == nil {
		n.w = ratelimit.Writer(conn, ratelimit.NewBucketWithRate(DefaultWriteLimit, DefaultWriteLimit*10))
	} else {
		n.w = ratelimit.Writer(conn, wBucket)
	}
	return n
}

func (n2 netConnLimit) Read(b []byte) (n int, err error) {
	if n2.r != nil {
		return n2.r.Read(b)
	}
	return n2.conn.Read(b)
}

func (n2 netConnLimit) Write(b []byte) (n int, err error) {
	if n2.w != nil {
		return n2.w.Write(b)
	}
	return n2.conn.Write(b)
}

func (n2 netConnLimit) Close() error {
	return n2.conn.Close()
}

func (n2 netConnLimit) LocalAddr() net.Addr {
	return n2.conn.LocalAddr()
}

func (n2 netConnLimit) RemoteAddr() net.Addr {
	return n2.conn.RemoteAddr()
}

func (n2 netConnLimit) SetDeadline(t time.Time) error {
	return n2.conn.SetDeadline(t)
}

func (n2 netConnLimit) SetReadDeadline(t time.Time) error {
	return n2.conn.SetReadDeadline(t)
}

func (n2 netConnLimit) SetWriteDeadline(t time.Time) error {
	return n2.conn.SetWriteDeadline(t)
}
