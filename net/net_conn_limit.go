package net

import (
	"github.com/zhangheng0027/ratelimit-plus"
	"io"
	"net"
	"time"
)

// DefaultReadLimit 50 MB/s
const DefaultReadLimit = 50 * 1024 * 1024

// DefaultWriteLimit 50 MB/s
const DefaultWriteLimit = 50 * 1024 * 1024

type ConnLimit struct {
	conn       net.Conn
	r          io.Reader
	w          io.Writer
	readLimit  ratelimit.BucketI
	writeLimit ratelimit.BucketI
}

func ToLimitNetConn(conn net.Conn, rBucket ratelimit.BucketI, wBucket ratelimit.BucketI) *ConnLimit {
	n := &ConnLimit{
		conn: conn,
	}
	if rBucket == nil {
		rBucket = ratelimit.NewBucketWithRate(DefaultReadLimit, DefaultReadLimit*10)
	}
	n.readLimit = rBucket
	n.r = ratelimit.Reader(conn, rBucket)

	if wBucket == nil {
		wBucket = ratelimit.NewBucketWithRate(DefaultWriteLimit, DefaultWriteLimit*10)
	}
	n.writeLimit = wBucket
	n.w = ratelimit.Writer(conn, wBucket)
	return n
}

func (n2 ConnLimit) addUpStream(readLimit *ratelimit.Bucket, writeLimit *ratelimit.Bucket) {
	if nil != readLimit {
		n2.readLimit.AddUpstream(readLimit)
	}
	if nil != writeLimit {
		n2.writeLimit.AddUpstream(writeLimit)
	}
}

func (n2 ConnLimit) Read(b []byte) (n int, err error) {
	return n2.r.Read(b)
}

func (n2 ConnLimit) Write(b []byte) (n int, err error) {
	return n2.w.Write(b)
}

func (n2 ConnLimit) Close() error {
	return n2.conn.Close()
}

func (n2 ConnLimit) LocalAddr() net.Addr {
	return n2.conn.LocalAddr()
}

func (n2 ConnLimit) RemoteAddr() net.Addr {
	return n2.conn.RemoteAddr()
}

func (n2 ConnLimit) SetDeadline(t time.Time) error {
	return n2.conn.SetDeadline(t)
}

func (n2 ConnLimit) SetReadDeadline(t time.Time) error {
	return n2.conn.SetReadDeadline(t)
}

func (n2 ConnLimit) SetWriteDeadline(t time.Time) error {
	return n2.conn.SetWriteDeadline(t)
}
