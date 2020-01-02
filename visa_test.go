// Copyright (c) 2017-2020 The lxi developers. All rights reserved.
// Project site: https://github.com/gotmc/lxi
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package lxi

import (
	"errors"
	"testing"
)

func TestParsingVisaResourceString(t *testing.T) {
	testCases := []struct {
		resourceString string
		interfaceType  string
		boardIndex     uint
		hostAddress    string
		port           uint
		resourceClass  string
		isError        bool
		errorString    error
	}{
		{
			"TCPIP0::10.12.100.15::5025::SOCKET",
			"TCPIP", 0, "10.12.100.15", 5025, "SOCKET",
			false, errors.New(""),
		},
	}
	for _, testCase := range testCases {
		resource, err := NewVisaResource(testCase.resourceString)
		if resource.interfaceType != testCase.interfaceType {
			t.Errorf(
				"interfaceType == %s, want %s for resource %s",
				resource.interfaceType,
				testCase.interfaceType,
				testCase.resourceString,
			)
		}
		if resource.boardIndex != testCase.boardIndex {
			t.Errorf(
				"boardIndex == %d, want %d for resource %s",
				resource.boardIndex,
				testCase.boardIndex,
				testCase.resourceString,
			)
		}
		if resource.hostAddress != testCase.hostAddress {
			t.Errorf(
				"hostAddress == %s, want %s for resource %s",
				resource.hostAddress,
				testCase.hostAddress,
				testCase.resourceString,
			)
		}
		// if resource.port != testCase.port {
		// t.Errorf(
		// "port == %d, want %d for resource %s",
		// resource.port,
		// testCase.port,
		// testCase.resourceString,
		// )
		// }
		if resource.resourceClass != testCase.resourceClass {
			t.Errorf(
				"resourceClass == %s, want %s for resource %s",
				resource.resourceClass,
				testCase.resourceClass,
				testCase.resourceString,
			)
		}
		if err != nil && testCase.isError {
			if err.Error() != testCase.errorString.Error() {
				t.Errorf(
					"err == %s, want %s for resource %s",
					err,
					testCase.errorString,
					testCase.resourceString,
				)
			}
		}
		if err != nil && !testCase.isError {
			t.Errorf("Unhandled error: %q for resource %s", err, testCase.resourceString)
		}
	}
}
