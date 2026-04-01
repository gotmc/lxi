// Copyright (c) 2017-2026 The lxi developers. All rights reserved.
// Project site: https://github.com/gotmc/lxi
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package lxi

import (
	"bufio"
	"context"
	"io"
	"net"
	"testing"
	"time"
)

// newTestDevice creates a Device backed by a net.Pipe for testing.
// Both ends are closed automatically when the test finishes.
func newTestDevice(t *testing.T) (*Device, net.Conn) {
	t.Helper()
	client, server := net.Pipe()
	dev := &Device{
		EndMark: '\n',
		conn:    client,
		rd:      bufio.NewReader(client),
	}
	t.Cleanup(func() {
		_ = server.Close()
		if dev.conn != nil {
			_ = dev.Close()
		}
	})
	return dev, server
}

func TestNewDeviceCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := NewDevice(ctx, "TCPIP0::192.0.2.1::5025::SOCKET")
	if err == nil {
		t.Fatal("expected error from NewDevice with canceled context")
	}
}

func TestNewDeviceInvalidAddress(t *testing.T) {
	ctx := context.Background()
	_, err := NewDevice(ctx, "INVALID::ADDRESS")
	if err == nil {
		t.Fatal("expected error from NewDevice with invalid VISA address")
	}
}

func TestWriteString(t *testing.T) {
	dev, remote := newTestDevice(t)

	msg := "hello"
	go func() {
		buf := make([]byte, 64)
		n, _ := remote.Read(buf)
		if got := string(buf[:n]); got != msg {
			t.Errorf("WriteString sent %q, want %q", got, msg)
		}
	}()

	n, err := dev.WriteString(msg)
	if err != nil {
		t.Fatalf("WriteString error: %v", err)
	}
	if n != len(msg) {
		t.Errorf("WriteString returned n=%d, want %d", n, len(msg))
	}
}

func TestCommand(t *testing.T) {
	dev, remote := newTestDevice(t)

	done := make(chan string, 1)
	go func() {
		buf := make([]byte, 128)
		n, _ := remote.Read(buf)
		done <- string(buf[:n])
	}()

	ctx := context.Background()
	err := dev.Command(ctx, "*RST")
	if err != nil {
		t.Fatalf("Command error: %v", err)
	}

	got := <-done
	want := "*RST\n"
	if got != want {
		t.Errorf("Command sent %q, want %q", got, want)
	}
}

func TestCommandWithFormatArgs(t *testing.T) {
	dev, remote := newTestDevice(t)

	done := make(chan string, 1)
	go func() {
		buf := make([]byte, 128)
		n, _ := remote.Read(buf)
		done <- string(buf[:n])
	}()

	ctx := context.Background()
	err := dev.Command(ctx, "VOLT %d", 5)
	if err != nil {
		t.Fatalf("Command error: %v", err)
	}

	got := <-done
	want := "VOLT 5\n"
	if got != want {
		t.Errorf("Command sent %q, want %q", got, want)
	}
}

func TestQuery(t *testing.T) {
	dev, remote := newTestDevice(t)

	go func() {
		buf := make([]byte, 128)
		_, _ = remote.Read(buf)
		_, _ = remote.Write([]byte("1.234\n"))
	}()

	ctx := context.Background()
	resp, err := dev.Query(ctx, "*IDN?")
	if err != nil {
		t.Fatalf("Query error: %v", err)
	}
	want := "1.234"
	if resp != want {
		t.Errorf("Query response = %q, want %q", resp, want)
	}
}

func TestQueryCanceledContext(t *testing.T) {
	dev, _ := newTestDevice(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := dev.Query(ctx, "*IDN?")
	if err == nil {
		t.Fatal("expected error from Query with canceled context")
	}
}

func TestCommandCanceledContext(t *testing.T) {
	dev, _ := newTestDevice(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := dev.Command(ctx, "*RST")
	if err == nil {
		t.Fatal("expected error from Command with canceled context")
	}
}

func TestCommandWithDeadline(t *testing.T) {
	dev, remote := newTestDevice(t)

	done := make(chan string, 1)
	go func() {
		buf := make([]byte, 128)
		n, _ := remote.Read(buf)
		done <- string(buf[:n])
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := dev.Command(ctx, "*CLS")
	if err != nil {
		t.Fatalf("Command error: %v", err)
	}

	got := <-done
	want := "*CLS\n"
	if got != want {
		t.Errorf("Command sent %q, want %q", got, want)
	}
}

func TestClose(t *testing.T) {
	dev, _ := newTestDevice(t)

	err := dev.Close()
	if err != nil {
		t.Fatalf("Close error: %v", err)
	}
	if dev.conn != nil {
		t.Error("expected conn to be nil after Close")
	}
	if dev.rd != nil {
		t.Error("expected rd to be nil after Close")
	}
}

func TestDoubleClose(t *testing.T) {
	dev, _ := newTestDevice(t)

	if err := dev.Close(); err != nil {
		t.Fatalf("first Close error: %v", err)
	}
	if err := dev.Close(); err != nil {
		t.Fatalf("second Close error: %v", err)
	}
}

func TestWriteAndRead(t *testing.T) {
	dev, remote := newTestDevice(t)

	msg := []byte("test data")
	go func() {
		_, _ = remote.Write(msg)
	}()

	buf := make([]byte, 64)
	n, err := dev.Read(buf)
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}
	if string(buf[:n]) != string(msg) {
		t.Errorf("Read got %q, want %q", string(buf[:n]), string(msg))
	}
}

func TestDeviceImplementsInterfaces(t *testing.T) {
	dev, _ := newTestDevice(t)

	var _ io.Reader = dev
	var _ io.Writer = dev
	var _ io.Closer = dev
	var _ io.StringWriter = dev
}

func TestApplyContextBackground(t *testing.T) {
	dev, _ := newTestDevice(t)

	called := false
	setter := func(_ time.Time) error {
		called = true
		return nil
	}

	cleanup, err := dev.applyContext(context.Background(), setter)
	if err != nil {
		t.Fatalf("applyContext error: %v", err)
	}
	cleanup()
	if called {
		t.Error("expected setDeadline to not be called for context.Background()")
	}
}

func TestApplyContextWithDeadline(t *testing.T) {
	dev, _ := newTestDevice(t)

	var setTime time.Time
	var resetTime time.Time
	calls := 0
	setter := func(tm time.Time) error {
		calls++
		if calls == 1 {
			setTime = tm
		} else {
			resetTime = tm
		}
		return nil
	}

	deadline := time.Now().Add(10 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	cleanup, err := dev.applyContext(ctx, setter)
	if err != nil {
		t.Fatalf("applyContext error: %v", err)
	}
	if setTime.IsZero() {
		t.Error("expected setDeadline to be called with non-zero time")
	}

	cleanup()
	if !resetTime.IsZero() {
		t.Error("expected cleanup to reset deadline to zero time")
	}
	if resetTime != (time.Time{}) {
		t.Error("expected cleanup to reset deadline to zero time")
	}
}

func TestApplyContextCancelableCleanup(t *testing.T) {
	dev, _ := newTestDevice(t)

	resetCalled := false
	setter := func(tm time.Time) error {
		if tm.IsZero() {
			resetCalled = true
		}
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cleanup, err := dev.applyContext(ctx, setter)
	if err != nil {
		t.Fatalf("applyContext error: %v", err)
	}

	cleanup()
	if !resetCalled {
		t.Error("expected cleanup to reset the deadline")
	}
}
