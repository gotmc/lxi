// Copyright (c) 2015-2017 The lxi developers. All rights reserved.
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
		boardIndex     uint16
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
		if resource.manufacturerID != testCase.manufacturerID {
			t.Errorf(
				"manufacturerID == %d, want %d for resource %s",
				resource.manufacturerID,
				testCase.manufacturerID,
				testCase.resourceString,
			)
		}
		if resource.modelCode != testCase.modelCode {
			t.Errorf(
				"modelCode == %d, want %d for resource %s",
				resource.modelCode,
				testCase.modelCode,
				testCase.resourceString,
			)
		}
		if resource.serialNumber != testCase.serialNumber {
			t.Errorf(
				"serialNumber == %s, want %s for resource %s",
				resource.serialNumber,
				testCase.serialNumber,
				testCase.resourceString,
			)
		}
		if resource.interfaceIndex != testCase.interfaceIndex {
			t.Errorf(
				"interfaceIndex == %d, want %d for resource %s",
				resource.interfaceIndex,
				testCase.interfaceIndex,
				testCase.resourceString,
			)
		}
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
