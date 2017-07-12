// Copyright (c) 2015-2017 The lxi developers. All rights reserved.
// Project site: https://github.com/gotmc/lxi
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"fmt"
	"io"
	"log"

	"github.com/gotmc/lxi"
)

func main() {

	fg, err := lxi.NewDevice("TCPIP0::10.12.100.150::5025::SOCKET")
	if err != nil {
		log.Fatalf("NewDevice error: %s", err)
	}

	// Configure function generator
	fg.WriteString("*CLS\n")
	fg.WriteString("burst:state off\n")
	fg.Write([]byte("apply:sinusoid 2340, 0.1, 0.0\n")) // Write using byte slice
	io.WriteString(fg, "burst:internal:period 0.112\n") // WriteString using io's Writer interface
	fg.WriteString("burst:internal:period 0.112\n")     // WriteString
	fg.WriteString("burst:ncycles 131\n")
	fg.WriteString("burst:state on\n")

	// Query using a write and then a read.
	queries := []string{"volt", "freq", "volt:offs", "volt:unit"}
	for _, q := range queries {
		ws := fmt.Sprintf("%s?\n", q)
		fg.WriteString(ws)
		var p [512]byte
		bytesRead, err := fg.Read(p[:])
		if err != nil {
			log.Printf("Error reading: %v", err)
		} else {
			log.Printf("Read %d bytes for %s? = %s", bytesRead, q, p)
		}
	}

	// Query using the query method
	queryRange(fg, queries)

	// Close the function generator and check for errors.
	err = fg.Close()
	if err != nil {
		log.Printf("error closing fg: %s", err)
	}
}

func queryRange(fg *lxi.Device, r []string) {
	for _, q := range r {
		ws := fmt.Sprintf("%s?", q)
		s, err := fg.Query(ws)
		if err != nil {
			log.Printf("Error reading: %v", err)
		} else {
			log.Printf("Query %s? = %s", q, s)
		}
	}
}
