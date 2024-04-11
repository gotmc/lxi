// Copyright (c) 2017-2024 The lxi developers. All rights reserved.
// Project site: https://github.com/gotmc/lxi
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package lxi

import (
	"bufio"
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
	tcpAddress := fmt.Sprintf("%s:%d", v.hostAddress, v.port)
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

// WriteString writes a string using the underlying network connection. An end
// mark character, such as a newline, is not added to the string being written.
func (d *Device) WriteString(s string) (n int, err error) {
	return d.Write([]byte(s))
}

// Command sends the SCPI/ASCII command to the underlying network connection.
// The Device's EndMark character (newline by default) is automatically added
// to the end of the string.
func (d *Device) Command(format string, a ...any) error {
	cmd := format
	if a != nil {
		cmd = fmt.Sprintf(format, a...)
	}
	_, err := d.WriteString(strings.TrimSpace(cmd) + string(d.EndMark))
	return err
}

// Query writes the given string to the underlying network connection and
// returns a string. A newline character is automatically added to the query
// command sent to the instrument.
func (d *Device) Query(cmd string) (string, error) {
	err := d.Command(cmd)
	if err != nil {
		return "", err
	}
	return d.rd.ReadString(d.EndMark)
}
