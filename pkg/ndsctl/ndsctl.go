package ndsctl

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"time"
)

// ErrBusy means that the ndsctl thread is busy and cannot process our request.
//
// Using a lock in the NDSCTL instance is not enough to prevent this condition
// because other processes might be using ndsctl at the same time.
var ErrBusy = errors.New("ndsctl reported busy")

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
	// bin is the path to the ndsctl binary.
	bin string
	// retryAttempts is the number of times to retry a command if it fails with ErrBusy.
	retryAttempts int
	// retryDelay is the delay between retries.
	retryDelay time.Duration
}

func NewNDSCTL(bin string) *NDSCTL {
	return &NDSCTL{
		bin:           bin,
		retryAttempts: 3,
		retryDelay:    100 * time.Millisecond,
	}
}

func (n *NDSCTL) ClientInfo(id string) (*ClientInfo, error) {
	cmd := exec.Command(n.bin, "json", id)
	result := bytes.Buffer{}
	cmd.Stdout = &result
	err := n.execute(cmd)
	if err != nil {
		return nil, err
	}

	var data ClientInfo
	err = json.NewDecoder(&result).Decode(&data)
	return &data, err
}

func (n *NDSCTL) Auth(id string) error {
	cmd := exec.Command(n.bin, "auth", id, "0", "0", "0", "0", "0", "")
	return n.execute(cmd)
}

func (n *NDSCTL) Deauth(id string) error {
	cmd := exec.Command(n.bin, "deauth", id)
	return n.execute(cmd)
}

func (n *NDSCTL) execute(cmd *exec.Cmd) error {
	var err error
	for i := 0; i < n.retryAttempts; i++ {
		err = handleError(cmd.Run())
		if err == nil {
			return nil
		}

		if errors.Is(err, ErrBusy) {
			time.Sleep(n.retryDelay)
		} else {
			return err
		}
	}
	return err
}

func handleError(err error) error {
	if err == nil {
		return nil
	}
	if exitError, ok := err.(*exec.ExitError); ok {
		switch exitError.ExitCode() {
		case 4:
			return ErrBusy
		}
	}

	return err
}
