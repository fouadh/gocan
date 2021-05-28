package support

import (
	"bytes"
	"github.com/spf13/cobra"
)

func RunCommand(cmd *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)

	_, err := cmd.ExecuteC()

	return buf.String(), err
}