package main

import (
	"encoding/json"
	"os"
)

const currentVersion = "0.3.1"

type VersionResponse struct {
	CniVersion        string   `json:"cniVersion"`
	SupportedVersions []string `json:"supportedVersions"`
}

func HandleVersion() error {
	return json.NewEncoder(os.Stdout).Encode(VersionResponse{
		CniVersion:        currentVersion,
		SupportedVersions: []string{currentVersion},
	})
}
