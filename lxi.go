// Project site: https://github.com/gotmc/visa
// Copyright (c) 2017 The visa developers. All rights reserved.
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package lxi

import (
	"errors"
	"fmt"
	"net"
)

// Device models an LXI device, which is currently just a TCPIP socket
// interface. An LXI Device also implements the ivi.Driver interface.
type Device struct {
	conn net.Conn
}

// NewDevice opens a TCPIP Device using the given VISA address resource string.
func NewDevice(address string) (*Device, error) {
	var d *Device
	v, err := NewVisaResource(address)
	if err != nil {
		return d, err
	}
	tcpAddress := fmt.Sprintf("%s:%d", v.hostAddress, v.port)
	conn, err := net.Dial("tcp", tcpAddress)
	if err != nil {
		return d, err
	}
	d.conn = conn
	return d, nil
}

// Write writes the given data to the network connection.
func (d *Device) Write(p []byte) (n int, err error) {
	return d.conn.Write(p)
}

// Read reads from the network connection into the given byte slice.
func (d *Device) Read(p []byte) (n int, err error) {
	return d.conn.Read(p)
}

// Close closes the underlying network connection.
func (d *Device) Close() error {
	return d.conn.Close()
}

// WriteString writes a string using the underlying network connection.
func (d *Device) WriteString(s string) (n int, err error) {
	// FIXME(mdr): Is WriteString required?
	return d.Write([]byte(s))
}

// Query queries the device using the Read and Write methods.
func (d *Device) Query(s string) (value string, err error) {
	return "", errors.New("query method not implemented for lxi device")
}
