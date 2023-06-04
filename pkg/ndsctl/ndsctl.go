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
	bin string
}

func NewNDSCTL(bin string) *NDSCTL {
	return &NDSCTL{
		bin: bin,
	}
}

func (n *NDSCTL) ClientInfo(id string) (*ClientInfo, error) {
	cmd := exec.Command(n.bin, "json", id)
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var data ClientInfo
	err = json.Unmarshal(stdout, &data)
	return &data, err
}

func (n *NDSCTL) Auth(id string) error {
	cmd := exec.Command(n.bin, "auth", id, "0", "0", "0", "0", "0", "")
	return cmd.Run()
}

func (n *NDSCTL) Deauth(id string) error {
	cmd := exec.Command(n.bin, "deauth", id)
	return cmd.Run()
}
