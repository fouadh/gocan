package shell

import "os/exec"

func ExecuteCommand(command string, args []string, directory string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = directory
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return out, nil
}
