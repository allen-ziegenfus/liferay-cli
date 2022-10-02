/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package ext

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"liferay.com/lcectl/constants"
	"liferay.com/lcectl/docker"
	"liferay.com/lcectl/flags"
	"liferay.com/lcectl/prereq"
	"liferay.com/lcectl/spinner"
)

var dir string

// downCmd represents the down command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops all client-extension workloads",
	Long:  `Stops localdev server and DXP after shutting down all client-extension workloads.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		prereq.Prereq(flags.Verbose)

		config := container.Config{
			Image: "localdev-server",
			Cmd:   []string{"/repo/scripts/ext/down.sh"},
		}
		host := container.HostConfig{
			Binds: []string{
				fmt.Sprintf("%s:%s", viper.GetString(constants.Const.RepoDir), "/repo"),
				docker.GetDockerSocket() + ":/var/run/docker.sock",
				fmt.Sprintf("%s:/workspace/client-extensions", dir),
				"localdevGradleCache:/root/.gradle",
				"localdevLiferayCache:/root/.liferay",
			},
			NetworkMode: container.NetworkMode(viper.GetString(constants.Const.DockerNetwork)),
		}

		spinner.Spin(
			spinner.SpinOptions{
				Doing: "Stopping", Done: "Stopped", On: "'localdev' extension environment", Enable: flags.Verbose,
			},
			func(fior func(io.ReadCloser, bool)) int {
				return docker.InvokeCommandInLocaldev("localdev-down", config, host, true, flags.Verbose, fior)
			})
	},
}

func init() {
	extCmd.AddCommand(stopCmd)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working dir")
	}
	stopCmd.Flags().StringVarP(&dir, "dir", "d", wd, "Set the base dir for down command")
}
