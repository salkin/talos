// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/talos-systems/talos/cmd/installer/pkg"
	"github.com/talos-systems/talos/cmd/installer/pkg/ova"
	"github.com/talos-systems/talos/cmd/installer/pkg/qemuimg"
	"github.com/talos-systems/talos/internal/pkg/runtime"
	"github.com/talos-systems/talos/internal/pkg/runtime/platform"
	"github.com/talos-systems/talos/pkg/cmd"
	"github.com/talos-systems/talos/pkg/config/types/v1alpha1"
	"github.com/talos-systems/talos/pkg/constants"
)

var outputArg string

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runImageCmd(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	imageCmd.Flags().StringVar(&outputArg, "output", "/out", "The output path")
	rootCmd.AddCommand(imageCmd)
}

//nolint: gocyclo
func runImageCmd() (err error) {
	p, err := platform.NewPlatform(options.Platform)
	if err != nil {
		return err
	}

	log.Printf("creating image for %s", p.Name())

	log.Print("creating RAW disk")

	img, err := pkg.CreateRawDisk()
	if err != nil {
		return err
	}

	log.Print("attaching loopback device ")

	if options.Disk, err = pkg.Loattach(img); err != nil {
		return err
	}

	// nolint: errcheck
	defer func() {
		log.Println("detaching loopback device")

		if err = pkg.Lodetach(options.Disk); err != nil {
			log.Println(err)
		}
	}()

	config := &v1alpha1.Config{
		ClusterConfig: &v1alpha1.ClusterConfig{
			ControlPlane: &v1alpha1.ControlPlaneConfig{},
		},
		MachineConfig: &v1alpha1.MachineConfig{
			MachineInstall: &v1alpha1.InstallConfig{
				InstallForce:           true,
				InstallBootloader:      options.Bootloader,
				InstallDisk:            options.Disk,
				InstallExtraKernelArgs: options.ExtraKernelArgs,
			},
		},
	}

	if options.ConfigSource == "" {
		switch p.Name() {
		case "aws", "azure", "digital-ocean", "gcp":
			options.ConfigSource = "none"
		case "vmware":
			options.ConfigSource = constants.ConfigGuestInfo
		default:
		}
	}

	if err = pkg.Install(p, config, runtime.None, options); err != nil {
		return err
	}

	if err := finalize(p, img); err != nil {
		return err
	}

	return nil
}

//nolint: gocyclo,interfacer
func finalize(platform runtime.Platform, img string) (err error) {
	dir := filepath.Dir(img)

	file := filepath.Base(img)
	name := strings.TrimSuffix(file, filepath.Ext(file))

	switch platform.Name() {
	case "aws":
		if err = tar("aws.tar.gz", file, dir); err != nil {
			return err
		}
	case "azure":
		file = name + ".vhd"

		if err = qemuimg.Convert("raw", "vpc", "subformat=fixed,force_size", img, filepath.Join(dir, file)); err != nil {
			return err
		}

		if err = tar("azure.tar.gz", file, dir); err != nil {
			return err
		}
	case "digital-ocean":
		if err = tar("digital-ocean.tar.gz", file, dir); err != nil {
			return err
		}
	case "gcp":
		if err = tar("gcp.tar.gz", file, dir); err != nil {
			return err
		}
	case "vmware":
		if err = ova.CreateOVAFromRAW(name, img, outputArg); err != nil {
			return err
		}
	}

	return nil
}

func tar(filename, src, dir string) error {
	if _, err := cmd.Run("tar", "-czvf", filepath.Join(outputArg, filename), src, "-C", dir); err != nil {
		return err
	}

	return nil
}
