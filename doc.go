// Copyright (c) 2017-2026 The lxi developers. All rights reserved.
// Project site: https://github.com/gotmc/lxi
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

// Package lxi provides a LAN eXtensions for Instrumentation (LXI) interface
// for controlling test equipment over Ethernet ports using SCPI commands. It
// implements the VISA LXI resource string format and serves as an instrument
// driver for the ivi and visa packages.
//
// This package is part of the gotmc ecosystem. The visa package
// (github.com/gotmc/visa) defines a common interface for instrument
// communication across different transports (GPIB, USB, TCP/IP, serial). The
// asrl package provides the serial transport implementation. The ivi package
// (github.com/gotmc/ivi) builds on top of visa to provide standardized,
// instrument-class-specific APIs following the IVI Foundation specifications.
//
// Devices are addressed using VISA resource strings of the form:
//
//	TCPIP<boardIndex>::<hostAddress>::<port>::SOCKET
//
// For example:
//
//	TCPIP0::192.168.1.101::5025::SOCKET
package lxi
