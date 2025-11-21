package cmd

import "github.com/spf13/cobra"

var rootCMD = &cobra.Command{
	Use:   "FlutterGen",
	Short: "A Advanced CLI to build and release APK Faster",
}

func Execute() {
	rootCMD.Execute()
}
