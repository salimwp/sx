package command

import (
	"context"
	"os"
	"os/signal"
	"strings"

	"github.com/spf13/cobra"
	"github.com/v-byte-cpu/sx/pkg/scan/tcp"
)

func newTCPFINCmd() *tcpFINCmd {
	c := &tcpFINCmd{}

	cmd := &cobra.Command{
		Use:     "fin [flags] subnet",
		Example: strings.Join([]string{"tcp fin -p 22 192.168.0.1/24", "tcp fin -p 22-4567 10.0.0.1"}, "\n"),
		Short:   "Perform TCP FIN scan",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
			defer cancel()

			if err = c.opts.parseRawOptions(); err != nil {
				return
			}

			scanName := tcp.FINScanType
			if err = c.opts.parseOptions(scanName, args); err != nil {
				return
			}

			m := c.opts.newTCPScanMethod(ctx,
				withTCPScanName(scanName),
				withTCPPacketFillerOptions(tcp.WithFIN()),
				withTCPPacketFilterFunc(tcp.TrueFilter),
				withTCPPacketFlags(tcp.AllFlags),
			)

			return startPacketScanEngine(ctx, newPacketScanConfig(
				withPacketScanMethod(m),
				withPacketBPFFilter(tcp.BPFFilter),
				withRateCount(c.opts.rateCount),
				withRateWindow(c.opts.rateWindow),
				withPacketVPNmode(c.opts.vpnMode),
				withPacketEngineConfig(newEngineConfig(
					withLogger(c.opts.logger),
					withScanRange(c.opts.scanRange),
					withExitDelay(c.opts.exitDelay),
				)),
			))
		},
	}

	c.opts.initCliFlags(cmd)

	c.cmd = cmd
	return c
}

type tcpFINCmd struct {
	cmd  *cobra.Command
	opts tcpCmdOpts
}
