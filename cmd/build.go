package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var apkName string

func init() {
	buildCMD.Flags().StringVar(&apkName, "NAME", "output", "APK output filename (without .apk)")
	rootCMD.AddCommand(buildCMD)
}

var buildCMD = &cobra.Command{
	Use:   "build apk",
	Short: "build flutter apk and rename it",
	RunE: func(cmd *cobra.Command, args []string) error {
		flutterCMD := exec.Command("flutter", "build", "apk")
		var out bytes.Buffer
		flutterCMD.Stdout = &out
		flutterCMD.Stderr = os.Stderr
		fmt.Println("🚀 Running flutter build apk...")
		if err := flutterCMD.Run(); err != nil {
			return fmt.Errorf("build failed: %v ", err)
		}

		var builtAPKPath string
		scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))

		for scanner.Scan() {
			line := scanner.Text()
			// Flutter prints: "Built build/app/outputs/flutter-apk/app-release.apk."
			if strings.HasPrefix(line, "Built ") {
				builtAPKPath = strings.TrimPrefix(line, "Built ")
				builtAPKPath = strings.TrimSuffix(builtAPKPath, ".")
				break
			}
		}

		if builtAPKPath == "" {
			return fmt.Errorf("❌ Could not detect built APK path from Flutter output")
		}

		finalAPK := apkName
		if !strings.HasSuffix(finalAPK, ".apk") {
			finalAPK += ".apk"
		}

		fmt.Printf("📦 Renaming %s → %s\n", builtAPKPath, finalAPK)
		if err := os.Rename(builtAPKPath, finalAPK); err != nil {
			return fmt.Errorf("rename failed: %v", err)
		}

		abs, _ := filepath.Abs(finalAPK)
		fmt.Printf("🎉 APK created: %s\n", abs)

		return nil

	},
}
