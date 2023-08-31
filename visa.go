// Copyright (c) 2017-2023 The lxi developers. All rights reserved.
// Project site: https://github.com/gotmc/lxi
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package lxi

import (
	"errors"
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

// NewVisaResource creates a new VisaResource using the given VISA resourceString.
func NewVisaResource(resourceString string) (visa *VisaResource, err error) {
	visa = &VisaResource{
		resourceString: resourceString,
	}
	regString := `^(?P<interfaceType>TCPIP)(?P<boardIndex>\d*)::` +
		`(?P<hostAddress>[^\s:]+)::` +
		`(?P<port>\d*)::` +
		`(?P<resourceClass>SOCKET)$`

	re := regexp.MustCompile(regString)
	res := re.FindStringSubmatch(strings.ToUpper(resourceString))
	subexpNames := re.SubexpNames()
	matchMap := map[string]string{}
	for i, n := range res {
		matchMap[subexpNames[i]] = string(n)
	}

	if matchMap["interfaceType"] != "TCPIP" {
		return visa, errors.New("visa: interface type was not TCPIP")
	}
	visa.interfaceType = "TCPIP"

	if matchMap["resourceClass"] != "SOCKET" {
		return visa, errors.New("visa: resource class was not SOCKET")
	}
	visa.resourceClass = "SOCKET"

	if matchMap["boardIndex"] != "" {
		boardIndex, err := strconv.ParseUint(matchMap["boardIndex"], 0, 16)
		if err != nil {
			return visa, errors.New("visa: boardIndex error")
		}
		visa.boardIndex = uint(boardIndex)
	}

	if matchMap["hostAddress"] != "" {
		visa.hostAddress = matchMap["hostAddress"]
	}

	if matchMap["port"] != "" {
		port, err := strconv.ParseUint(matchMap["port"], 10, 64)
		if err != nil {
			return visa, errors.New("visa: port error")
		}
		visa.port = uint(port)
	}
	return visa, nil
}
