// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sys/unix"

	"github.com/talos-systems/go-procfs/procfs"

	"github.com/talos-systems/talos/internal/pkg/kmsg"
	"github.com/talos-systems/talos/internal/pkg/mount/manager"
	"github.com/talos-systems/talos/internal/pkg/mount/manager/pseudo"
	"github.com/talos-systems/talos/internal/pkg/mount/manager/squashfs"
	"github.com/talos-systems/talos/internal/pkg/mount/switchroot"
	"github.com/talos-systems/talos/pkg/constants"
)

// nolint: gocyclo
func run() (err error) {
	// Mount the pseudo devices.
	mountpoints, err := pseudo.MountPoints()
	if err != nil {
		return err
	}

	pseudo := manager.NewManager(mountpoints)
	if err = pseudo.MountAll(); err != nil {
		return err
	}

	// Setup logging to /dev/kmsg.
	err = kmsg.Setup("[talos] [initramfs]", false)
	if err != nil {
		return err
	}

	// Mount the rootfs.
	log.Println("mounting the rootfs")

	mountpoints, err = squashfs.MountPoints(constants.NewRoot)
	if err != nil {
		return err
	}

	squashfs := manager.NewManager(mountpoints)
	if err = squashfs.MountAll(); err != nil {
		return err
	}

	// Switch into the new rootfs.
	log.Println("entering the rootfs")

	if err = switchroot.Switch(constants.NewRoot, pseudo); err != nil {
		return err
	}

	return nil
}

func recovery() {
	// If panic is set in the kernel flags, we'll hang instead of rebooting.
	// But we still allow users to hit CTRL+ALT+DEL to try and restart when they're ready.
	// Listening for these signals also keep us from deadlocking the goroutine.
	if r := recover(); r != nil {
		log.Printf("recovered from: %+v\n", r)

		p := procfs.ProcCmdline().Get(constants.KernelParamPanic).First()
		if p != nil && *p == "0" {
			log.Printf("panic=0 kernel flag found. sleeping forever")

			exitSignal := make(chan os.Signal, 1)
			signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
			<-exitSignal
		}

		for i := 10; i >= 0; i-- {
			log.Printf("rebooting in %d seconds\n", i)
			time.Sleep(1 * time.Second)
		}
	}

	// nolint: errcheck
	unix.Reboot(unix.LINUX_REBOOT_CMD_RESTART)
}

func main() {
	defer recovery()

	if err := run(); err != nil {
		panic(fmt.Errorf("early boot failed: %w", err))
	}

	// We should never reach this point if things are working as intended.
	panic(errors.New("unknown error"))
}
