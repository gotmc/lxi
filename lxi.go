// Copyright (c) 2017-2026 The lxi developers. All rights reserved.
// Project site: https://github.com/gotmc/lxi
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

// Package lxi provides a LAN eXtensions for Instrumentation (LXI) interface
// for controlling test equipment over Ethernet ports using SCPI commands. It
// implements the VISA LXI resource string format and serves as an instrument
// driver for the ivi and visa packages.
package lxi

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
)

// Device models an LXI device, which is currently just a TCPIP socket
// interface. An LXI Device also implements the ivi.Driver interface.
type Device struct {
	EndMark byte
	conn    net.Conn
	rd      *bufio.Reader
}

// NewDevice opens a TCPIP Device using the given VISA address resource string.
func NewDevice(address string) (*Device, error) {
	v, err := NewVisaResource(address)
	if err != nil {
		return nil, err
	}
	tcpAddress := net.JoinHostPort(v.hostAddress, fmt.Sprintf("%d", v.port))
	c, err := net.Dial("tcp", tcpAddress)
	if err != nil {
		return nil, err
	}
	return &Device{EndMark: '\n', conn: c, rd: bufio.NewReader(c)}, nil
}

// Write writes the given data to the network connection.
func (d *Device) Write(p []byte) (n int, err error) {
	return d.conn.Write(p)
}

// Read reads from the network connection into the given byte slice.
func (d *Device) Read(p []byte) (n int, err error) {
	return d.rd.Read(p)
}

// Close closes the underlying network connection.
func (d *Device) Close() error {
	d.rd = nil
	c := d.conn
	d.conn = nil
	return c.Close()
}

// WriteString writes a string to the underlying network connection. An endmark
// character, such as a newline, is not automatically added to the end of the
// string.
func (d *Device) WriteString(s string) (n int, err error) {
	return d.Write([]byte(s))
}

// Command sends the SCPI/ASCII command to the underlying network connection.
// An endmark character (newline by default) is automatically added to the end
// of the string. The context deadline, if set, is applied to the underlying
// network connection.
func (d *Device) Command(ctx context.Context, format string, a ...any) error {
	deadline, ok := ctx.Deadline()
	if ok {
		if err := d.conn.SetWriteDeadline(deadline); err != nil {
			return err
		}
	}
	cmd := format
	if a != nil {
		cmd = fmt.Sprintf(format, a...)
	}
	_, err := d.WriteString(strings.TrimSpace(cmd) + string(d.EndMark))
	return err
}

// Query writes the given SCPI/ASCII command to the underlying network
// connection and returns the response string. The device's endmark character
// (newline by default) is automatically added to the query command. The string
// returned is not stripped of any whitespace. The context deadline, if set, is
// applied to the underlying network connection for both the write and read
// operations.
func (d *Device) Query(ctx context.Context, cmd string) (string, error) {
	err := d.Command(ctx, cmd)
	if err != nil {
		return "", err
	}
	deadline, ok := ctx.Deadline()
	if ok {
		if err := d.conn.SetReadDeadline(deadline); err != nil {
			return "", err
		}
	}
	return d.rd.ReadString(d.EndMark)
}
