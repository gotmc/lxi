// Copyright (c) 2017-2026 The lxi developers. All rights reserved.
// Project site: https://github.com/gotmc/lxi
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package lxi

import (
	"testing"
)

func TestParsingVisaResourceString(t *testing.T) {
	testCases := []struct {
		name           string
		resourceString string
		interfaceType  string
		boardIndex     uint
		hostAddress    string
		port           uint
		resourceClass  string
		isError        bool
		errorSubstring string
	}{
		{
			name:           "valid with board index 0",
			resourceString: "TCPIP0::10.12.100.15::5025::SOCKET",
			interfaceType:  "TCPIP",
			boardIndex:     0,
			hostAddress:    "10.12.100.15",
			port:           5025,
			resourceClass:  "SOCKET",
		},
		{
			name:           "valid without board index",
			resourceString: "TCPIP::192.168.1.100::5025::SOCKET",
			interfaceType:  "TCPIP",
			boardIndex:     0,
			hostAddress:    "192.168.1.100",
			port:           5025,
			resourceClass:  "SOCKET",
		},
		{
			name:           "valid with board index 1",
			resourceString: "TCPIP1::10.0.0.1::7777::SOCKET",
			interfaceType:  "TCPIP",
			boardIndex:     1,
			hostAddress:    "10.0.0.1",
			port:           7777,
			resourceClass:  "SOCKET",
		},
		{
			name:           "valid lowercase input",
			resourceString: "tcpip0::10.12.100.15::5025::socket",
			interfaceType:  "TCPIP",
			boardIndex:     0,
			hostAddress:    "10.12.100.15",
			port:           5025,
			resourceClass:  "SOCKET",
		},
		{
			name:           "valid hostname instead of IP",
			resourceString: "TCPIP0::MY-INSTRUMENT::5025::SOCKET",
			interfaceType:  "TCPIP",
			boardIndex:     0,
			hostAddress:    "MY-INSTRUMENT",
			port:           5025,
			resourceClass:  "SOCKET",
		},
		{
			name:           "invalid interface type",
			resourceString: "GPIB0::10.12.100.15::5025::SOCKET",
			isError:        true,
			errorSubstring: "interface type",
		},
		{
			name:           "invalid resource class",
			resourceString: "TCPIP0::10.12.100.15::5025::INSTR",
			isError:        true,
			errorSubstring: "interface type",
		},
		{
			name:           "empty string",
			resourceString: "",
			isError:        true,
			errorSubstring: "interface type",
		},
		{
			name:           "missing port and resource class",
			resourceString: "TCPIP0::10.12.100.15",
			isError:        true,
			errorSubstring: "interface type",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resource, err := NewVisaResource(tc.resourceString)
			if tc.isError {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tc.errorSubstring)
				}
				if tc.errorSubstring != "" {
					if got := err.Error(); !contains(got, tc.errorSubstring) {
						t.Errorf("error = %q, want substring %q", got, tc.errorSubstring)
					}
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resource.interfaceType != tc.interfaceType {
				t.Errorf("interfaceType = %q, want %q", resource.interfaceType, tc.interfaceType)
			}
			if resource.boardIndex != tc.boardIndex {
				t.Errorf("boardIndex = %d, want %d", resource.boardIndex, tc.boardIndex)
			}
			if resource.hostAddress != tc.hostAddress {
				t.Errorf("hostAddress = %q, want %q", resource.hostAddress, tc.hostAddress)
			}
			if resource.port != tc.port {
				t.Errorf("port = %d, want %d", resource.port, tc.port)
			}
			if resource.resourceClass != tc.resourceClass {
				t.Errorf("resourceClass = %q, want %q", resource.resourceClass, tc.resourceClass)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
