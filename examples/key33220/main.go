// Copyright (c) 2017-2026 The lxi developers. All rights reserved.
// Project site: https://github.com/gotmc/lxi
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/gotmc/lxi"
)

func main() {

	// Get IP address from CLI flag.
	var ip string
	flag.StringVar(
		&ip,
		"ip",
		"192.168.1.100",
		"IP address of Keysight 33220A",
	)
	flag.Parse()

	// Create a new LXI device
	address := fmt.Sprintf("TCPIP0::%s::5025::SOCKET", ip)
	log.Printf("Using VISA address: %s", address)
	fg, err := lxi.NewDevice(address)
	if err != nil {
		log.Fatalf("NewDevice error: %s", err)
	}

	// Configure function generator
	ctx := context.Background()
	numCycles := 131
	period := 0.112

	if _, err = fg.WriteString("*CLS\n"); err != nil { // Write using lxi.WriteString
		log.Fatal(err)
	}
	if _, err = io.WriteString(fg, "burst:state off\n"); err != nil { // Write using io.WriteString
		log.Fatal(err)
	}
	// Write using lxi.Write
	if _, err = fg.Write([]byte("apply:sinusoid 2340, 0.1, 0.0\n")); err != nil {
		log.Fatal(err)
	}
	// Write using fmt.Fprint
	if _, err = fmt.Fprintf(fg, "burst:internal:period %f\n", period); err != nil {
		log.Fatal(err)
	}
	if err = fg.Command(ctx, "burst:ncycles %d", numCycles); err != nil { // Write using lxi.Command
		log.Fatal(err)
	}
	if err = fg.Command(ctx, "burst:state on"); err != nil { // Command appends a newline.
		log.Fatal(err)
	}

	// Query using the query method
	queries := []string{"volt", "freq", "volt:offs", "volt:unit"}
	queryRange(ctx, fg, queries)

	// Close the function generator and check for errors.
	err = fg.Close()
	if err != nil {
		log.Printf("error closing fg: %s", err)
	}
}

func queryRange(ctx context.Context, fg *lxi.Device, r []string) {
	for _, q := range r {
		ws := fmt.Sprintf("%s?", q)
		log.Printf("Querying %s", ws)
		s, err := fg.Query(ctx, ws)
		log.Printf("Completed %s query", ws)
		if err != nil {
			log.Printf("Error reading: %v", err)
		} else {
			log.Printf("Query %s? = %s", q, s)
		}
	}
}
