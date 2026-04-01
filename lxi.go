// Copyright (c) 2017-2026 The lxi developers. All rights reserved.
// Project site: https://github.com/gotmc/lxi
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package lxi

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

// Device models an LXI device, which is currently just a TCPIP socket
// interface. An LXI Device also implements the ivi.Driver interface.
type Device struct {
	EndMark byte
	conn    net.Conn
	rd      *bufio.Reader
}

// NewDevice opens a TCPIP Device using the given VISA address resource string.
// The context controls the timeout/deadline for the TCP connection attempt.
func NewDevice(ctx context.Context, address string) (*Device, error) {
	v, err := NewVisaResource(address)
	if err != nil {
		return nil, err
	}
	tcpAddress := net.JoinHostPort(v.hostAddress, fmt.Sprintf("%d", v.port))
	var dialer net.Dialer
	c, err := dialer.DialContext(ctx, "tcp", tcpAddress)
	if err != nil {
		return nil, err
	}
	return &Device{EndMark: '\n', conn: c, rd: bufio.NewReader(c)}, nil
}

// Close closes the underlying network connection.
func (d *Device) Close() error {
	if d.conn == nil {
		return nil
	}
	d.rd = nil
	c := d.conn
	d.conn = nil
	return c.Close()
}

// Read reads from the network connection into the given byte slice.
func (d *Device) Read(p []byte) (n int, err error) {
	return d.rd.Read(p)
}

// Write writes the given data to the network connection.
func (d *Device) Write(p []byte) (n int, err error) {
	return d.conn.Write(p)
}

// WriteString writes a string to the underlying network connection. An endmark
// character, such as a newline, is not automatically added to the end of the
// string.
func (d *Device) WriteString(s string) (n int, err error) {
	return d.Write([]byte(s))
}

// ReadContext reads from the network connection in the given byte slice in a
// context aware manner.
func (d *Device) ReadContext(ctx context.Context, p []byte) (n int, err error) {
	cleanup, err := d.applyContext(ctx, d.conn.SetReadDeadline)
	if err != nil {
		return 0, err
	}
	defer cleanup()
	return d.rd.Read(p)
}

// WriteContext writes the given data to the network connection in a context
// aware manner.
func (d *Device) WriteContext(ctx context.Context, p []byte) (n int, err error) {
	cleanup, err := d.applyContext(ctx, d.conn.SetWriteDeadline)
	if err != nil {
		return 0, err
	}
	defer cleanup()
	return d.conn.Write(p)
}

// Command sends a SCPI/ASCII command to the underlying network connection. The
// command can be optionally formatted according to a format specifier. An
// endmark character, such as newline, is automatically added to the end of the
// string. The context deadline, if set, is applied to the underlying network
// connection.
func (d *Device) Command(ctx context.Context, cmd string, a ...any) error {
	cleanup, err := d.applyContext(ctx, d.conn.SetWriteDeadline)
	if err != nil {
		return err
	}
	defer cleanup()
	if len(a) > 0 {
		cmd = fmt.Sprintf(cmd, a...)
	}
	_, err = d.WriteString(strings.TrimSpace(cmd) + string(d.EndMark))
	return err
}

// Query writes the given SCPI/ASCII command to the underlying network
// connection and returns the response string. The device's endmark character
// (newline by default) is automatically added to the query command. The
// trailing endmark character is stripped from the response. The context
// deadline, if set, is applied to the underlying network connection for both
// the write and read operations.
func (d *Device) Query(ctx context.Context, cmd string) (string, error) {
	err := d.Command(ctx, cmd)
	if err != nil {
		return "", err
	}
	cleanup, err := d.applyContext(ctx, d.conn.SetReadDeadline)
	if err != nil {
		return "", err
	}
	defer cleanup()
	s, err := d.rd.ReadString(d.EndMark)
	if err != nil {
		return s, err
	}
	return strings.TrimSuffix(s, string(d.EndMark)), nil
}

// applyContext sets a deadline on the connection using the provided setter and
// returns a cleanup function that must be called when the I/O operation
// completes. If the context has a deadline, it is used directly. If the context
// has no deadline but is already done, an error is returned. Otherwise, a
// goroutine watches for context cancellation and forces an immediate deadline
// to unblock any pending I/O. The cleanup function stops the goroutine and
// resets the deadline.
func (d *Device) applyContext(
	ctx context.Context,
	setDeadline func(time.Time) error,
) (func(), error) {
	noop := func() {}
	if deadline, ok := ctx.Deadline(); ok {
		return func() { _ = setDeadline(time.Time{}) }, setDeadline(deadline)
	}
	if ctx.Done() == nil {
		return noop, nil
	}
	select {
	case <-ctx.Done():
		return noop, ctx.Err()
	default:
	}
	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			_ = setDeadline(time.Now())
		case <-done:
		}
	}()
	cleanup := func() {
		close(done)
		_ = setDeadline(time.Time{})
	}
	return cleanup, nil
}
