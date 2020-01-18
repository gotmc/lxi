// Copyright (c) 2017-2020 The lxi developers. All rights reserved.
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

	fg, err := lxi.NewDevice("TCPIP0::10.12.112.7::5025::SOCKET")
	if err != nil {
		log.Fatalf("NewDevice error: %s", err)
	}

	// Configure function generator
	numCycles := 131
	period := 0.112
	fg.WriteString("*CLS\n")                              // Write using lxi.WriteString
	io.WriteString(fg, "burst:state off\n")               // Write using io.WriteString
	fg.Write([]byte("apply:sinusoid 2340, 0.1, 0.0\n"))   // Write using lxi.Write
	fmt.Fprintf(fg, "burst:internal:period %f\n", period) // Write using fmt.Fprint
	fg.Command("burst:ncycles %d", numCycles)             // Write using lxi.Command
	fg.Command("burst:state on")                          // Command appends a newline.

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
