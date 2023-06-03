package ndsctl

import (
	"encoding/json"
	"os/exec"
)

type ClientInfo struct {
	ID          int    `json:"id"`
	IpAddress   string `json:"ip"`
	MacAddress  string `json:"mac"`
	Added       int    `json:"added"`
	Active      int    `json:"active"`
	Duration    int    `json:"duration"`
	Token       string `json:"token"`
	State       string `json:"state"`
	DownloadKiB string `json:"download_this_session"`
}

type NDSCTL struct {
	//
}

func (ndsctl *NDSCTL) ClientInfo(id string) (*ClientInfo, error) {
	cmd := exec.Command("ndsctl", "json", id)
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var data ClientInfo
	err = json.Unmarshal(stdout, &data)
	return &data, err
}

func (ndsctl *NDSCTL) Auth(id string) error {
	cmd := exec.Command("ndsctl", "auth", id, "0", "0", "0", "0", "0", "")
	return cmd.Run()
}

func (ndsctl *NDSCTL) Deauth(id string) error {
	cmd := exec.Command("ndsctl", "deauth", id)
	return cmd.Run()
}
