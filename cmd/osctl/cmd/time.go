// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"

	timeapi "github.com/talos-systems/talos/api/time"
	"github.com/talos-systems/talos/cmd/osctl/pkg/client"
	"github.com/talos-systems/talos/cmd/osctl/pkg/helpers"
)

// timeCmd represents the time command
var timeCmd = &cobra.Command{
	Use:   "time [--check server]",
	Short: "Gets current server time",
	Long:  ``,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return WithClient(func(ctx context.Context, c *client.Client) error {
			server, err := cmd.Flags().GetString("check")
			if err != nil {
				return fmt.Errorf("failed to parse check flag: %w", err)
			}

			var (
				resp       *timeapi.TimeResponse
				remotePeer peer.Peer
			)

			if server == "" {
				resp, err = c.Time(ctx, grpc.Peer(&remotePeer))
			} else {
				resp, err = c.TimeCheck(ctx, server, grpc.Peer(&remotePeer))
			}

			if err != nil {
				if resp == nil {
					return fmt.Errorf("error fetching time: %w", err)
				}

				helpers.Warning("%s", err)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "NODE\tNTP-SERVER\tLOCAL-TIME\tREMOTE-TIME")

			defaultNode := helpers.AddrFromPeer(&remotePeer)

			var localtime, remotetime time.Time
			for _, msg := range resp.Messages {
				node := defaultNode

				if msg.Metadata != nil {
					node = msg.Metadata.Hostname
				}

				localtime, err = ptypes.Timestamp(msg.Localtime)
				if err != nil {
					return fmt.Errorf("error parsing local time: %w", err)
				}
				remotetime, err = ptypes.Timestamp(msg.Remotetime)
				if err != nil {
					return fmt.Errorf("error parsing remote time: %w", err)
				}

				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", node, msg.Server, localtime.String(), remotetime.String())
			}

			return w.Flush()
		})
	},
}

func init() {
	timeCmd.Flags().StringP("check", "c", "pool.ntp.org", "checks server time against specified ntp server")
	rootCmd.AddCommand(timeCmd)
}
