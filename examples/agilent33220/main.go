// Copyright (c) 2017-2020 The lxi developers. All rights reserved.
// Project site: https://github.com/gotmc/lxi
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"fmt"
	"log"

	"github.com/gotmc/lxi"
)

func main() {

	fg, err := lxi.NewDevice("TCPIP0::10.12.112.7::5025::SOCKET")
	if err != nil {
		log.Fatalf("NewDevice error: %s", err)
	}

	// Configure function generator
	fg.WriteString("*CLS\n")
	fg.WriteString("burst:state off\n")
	fg.Write([]byte("apply:sinusoid 100, 0.1, 0.0\n")) // Write using byte slice
	fg.WriteString("burst:internal:period 0.224\n")    // WriteString
	fg.WriteString("burst:ncycles 11\n")
	fg.WriteString("BURS:STAT ON\n")

	// Query using the query method
	queries := []string{"volt", "freq", "volt:offs", "volt:unit"}
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
		log.Printf("Querying %s", ws)
		s, err := fg.Query(ws)
		log.Printf("Completed %s query", ws)
		if err != nil {
			log.Printf("Error reading: %v", err)
		} else {
			log.Printf("Query %s? = %s", q, s)
		}
	}
}
