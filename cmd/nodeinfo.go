// Copyright © 2020 Weald Technology Trading
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"fmt"
	"os"

	eth2client "github.com/attestantio/go-eth2-client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wealdtech/ethdo/util"
)

// defaultBeaconNode is used if no other connection is supplied.
var defaultBeaconNode = "http://mainnet-consensus.attestant.io/"

var nodeInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Obtain information about a node",
	Long: `Obtain information about a node.  For example:

    ethdo node info

In quiet mode this will return 0 if the node information can be obtained, otherwise 1.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		eth2Client, err := util.ConnectToBeaconNode(ctx, viper.GetString("connection"), viper.GetDuration("timeout"), viper.GetBool("allow-insecure-connections"))
		if err != nil {
			if viper.GetString("connection") != "" {
				// The user provided a connection, so don't second-guess them by using a different node.
				fmt.Fprintln(os.Stderr, err.Error())
				return
			}

			// The user did not provide a connection, so attempt to use the default node.
			if viper.GetBool("debug") {
				fmt.Fprintf(os.Stderr, "No node connection, attempting to use %s\n", defaultBeaconNode)
			}
			eth2Client, err = util.ConnectToBeaconNode(ctx, defaultBeaconNode, viper.GetDuration("timeout"), viper.GetBool("allow-insecure-connections"))
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				return
			}
			if !viper.GetBool("quiet") {
				fmt.Fprintf(os.Stderr, "No connection supplied; using mainnet public access endpoint\n")
			}
		}

		if quiet {
			os.Exit(_exitSuccess)
		}

		if verbose {
			version, err := eth2Client.(eth2client.NodeVersionProvider).NodeVersion(ctx)
			errCheck(err, "Failed to obtain node version")
			fmt.Printf("Version: %s\n", version)
		}

		syncState, err := eth2Client.(eth2client.NodeSyncingProvider).NodeSyncing(ctx)
		errCheck(err, "failed to obtain node sync state")
		fmt.Printf("Syncing: %t\n", syncState.SyncDistance != 0)

		os.Exit(_exitSuccess)
	},
}

func init() {
	nodeCmd.AddCommand(nodeInfoCmd)
	nodeFlags(nodeInfoCmd)
}
