package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// FileMeta holds metadata about an uploaded file.
type FileMeta struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Type string `json:"type"`
}

// GetFileMeta retrieves file metadata from the sender's /meta endpoint.
func GetFileMeta(senderURL string) (*FileMeta, error) {
	resp, err := http.Get(senderURL + "/meta")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("file not available")
	}

	var meta FileMeta
	err = json.NewDecoder(resp.Body).Decode(&meta)
	return &meta, err
}
