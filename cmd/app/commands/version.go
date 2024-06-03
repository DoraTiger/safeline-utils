package commands

import (
	"encoding/json"
	"fmt"

	"github.com/DoraTiger/safeline-utils/version"
	"github.com/spf13/cobra"
)

var (
	verbose bool
)

func init() {
	registerFlagsVersionCmd(VersionCmd)
}

func registerFlagsVersionCmd(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "show more info")
}

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "show the version info",
	Run: func(cmd *cobra.Command, args []string) {

		// If verbose is true, more information is displayed in json format
		if verbose {
			values, _ := json.MarshalIndent(struct {
				SafelineUtils string `json:"safeline_utils"`
				Build         string `json:"build"`
				Repo          string `json:"repo"`
			}{
				SafelineUtils: version.Version,
				Build:         version.Build,
				Repo:          version.Repo,
			}, "", "  ")
			fmt.Println(string(values))
		} else {
			fmt.Println(version.Version)
		}
	},
}
