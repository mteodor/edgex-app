//
// Copyright (c) 2018
// Mainflux
//
// SPDX-License-Identifier: Apache-2.0
//

package exapp

import (
	"encoding/json"
	"net/http"
)

const version string = "0.0.1"

// VersionInfo - application version
type VersionInfo struct {
	Service string `json:"service"`
	Version string `json:"version"`
}

// Version exposes an HTTP handler for retrieving service version.
func Version(service string) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		res := VersionInfo{service, version}

		data, _ := json.Marshal(res)

		rw.Write(data)
	})
}
