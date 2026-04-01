// Copyright (c) 2017-2026 The lxi developers. All rights reserved.
// Project site: https://github.com/gotmc/lxi
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package lxi_test

import (
	"fmt"
	"log"

	"github.com/gotmc/lxi"
)

func ExampleNewVisaResource() {
	resource, err := lxi.NewVisaResource("TCPIP0::192.168.1.100::5025::SOCKET")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resource)
	// Output:
	// TCPIP0::192.168.1.100::5025::SOCKET
}
