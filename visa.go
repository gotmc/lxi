// Copyright (c) 2017-2026 The lxi developers. All rights reserved.
// Project site: https://github.com/gotmc/lxi
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package lxi

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// VisaResource represents a VISA enabled piece of test equipment.
type VisaResource struct {
	resourceString string
	interfaceType  string
	boardIndex     uint
	hostAddress    string
	port           uint
	resourceClass  string
}

// String returns the VISA resource string in canonical form.
func (v *VisaResource) String() string {
	return fmt.Sprintf("%s%d::%s::%d::%s",
		v.interfaceType, v.boardIndex, v.hostAddress, v.port, v.resourceClass)
}

var visaResourceRe = regexp.MustCompile(
	`^(?P<interfaceType>TCPIP)(?P<boardIndex>\d*)::` +
		`(?P<hostAddress>[^\s:]+)::` +
		`(?P<port>\d*)::` +
		`(?P<resourceClass>SOCKET)$`,
)

// NewVisaResource creates a new VisaResource using the given VISA resourceString.
func NewVisaResource(resourceString string) (*VisaResource, error) {
	re := visaResourceRe
	res := re.FindStringSubmatch(strings.ToUpper(resourceString))
	if res == nil {
		return nil, errors.New("visa: resource string does not match expected format")
	}
	subexpNames := re.SubexpNames()
	matchMap := map[string]string{}
	for i, n := range res {
		matchMap[subexpNames[i]] = string(n)
	}

	if matchMap["interfaceType"] != "TCPIP" {
		return nil, errors.New("visa: interface type was not TCPIP")
	}

	if matchMap["resourceClass"] != "SOCKET" {
		return nil, errors.New("visa: resource class was not SOCKET")
	}

	visa := &VisaResource{
		resourceString: resourceString,
		interfaceType:  "TCPIP",
		resourceClass:  "SOCKET",
	}

	if matchMap["boardIndex"] != "" {
		boardIndex, err := strconv.ParseUint(matchMap["boardIndex"], 0, 16)
		if err != nil {
			return nil, errors.New("visa: boardIndex error")
		}
		visa.boardIndex = uint(boardIndex)
	}

	if matchMap["hostAddress"] == "" {
		return nil, errors.New("visa: host address is required")
	}
	visa.hostAddress = matchMap["hostAddress"]

	if matchMap["port"] != "" {
		port, err := strconv.ParseUint(matchMap["port"], 10, 64)
		if err != nil {
			return nil, errors.New("visa: port error")
		}
		visa.port = uint(port)
	}
	return visa, nil
}
