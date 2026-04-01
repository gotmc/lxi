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

// Sentinel errors returned by NewVisaResource.
var (
	ErrResourceFormat = errors.New("visa: resource string does not match expected format")
	ErrInterfaceType  = errors.New("visa: interface type was not TCPIP")
	ErrResourceClass  = errors.New("visa: resource class was not SOCKET")
	ErrBoardIndex     = errors.New("visa: board index error")
	ErrHostAddress    = errors.New("visa: host address is required")
	ErrPort           = errors.New("visa: port error")
)

// VisaResource represents a VISA enabled piece of test equipment.
type VisaResource struct {
	resourceString string
	interfaceType  string
	boardIndex     int
	hostAddress    string
	port           int
	resourceClass  string
}

// InterfaceType returns the VISA interface type (e.g., "TCPIP").
func (v *VisaResource) InterfaceType() string { return v.interfaceType }

// BoardIndex returns the VISA board index.
func (v *VisaResource) BoardIndex() int { return v.boardIndex }

// HostAddress returns the hostname or IP address of the instrument.
func (v *VisaResource) HostAddress() string { return v.hostAddress }

// Port returns the TCP port number of the instrument.
func (v *VisaResource) Port() int { return v.port }

// ResourceClass returns the VISA resource class (e.g., "SOCKET").
func (v *VisaResource) ResourceClass() string { return v.resourceClass }

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
		return nil, ErrResourceFormat
	}
	subexpNames := re.SubexpNames()
	matchMap := map[string]string{}
	for i, n := range res {
		matchMap[subexpNames[i]] = string(n)
	}

	if matchMap["interfaceType"] != "TCPIP" {
		return nil, ErrInterfaceType
	}

	if matchMap["resourceClass"] != "SOCKET" {
		return nil, ErrResourceClass
	}

	visa := &VisaResource{
		resourceString: resourceString,
		interfaceType:  "TCPIP",
		resourceClass:  "SOCKET",
	}

	if matchMap["boardIndex"] != "" {
		boardIndex, err := strconv.Atoi(matchMap["boardIndex"])
		if err != nil {
			return nil, fmt.Errorf("visa: parsing board index: %w", err)
		}
		visa.boardIndex = boardIndex
	}

	if matchMap["hostAddress"] == "" {
		return nil, ErrHostAddress
	}
	visa.hostAddress = matchMap["hostAddress"]

	if matchMap["port"] != "" {
		port, err := strconv.Atoi(matchMap["port"])
		if err != nil {
			return nil, fmt.Errorf("visa: parsing port: %w", err)
		}
		visa.port = port
	}
	return visa, nil
}
