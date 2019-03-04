package cmd
 
import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var (
	// Version 版本号
	Version string
	Date   	string
	Commit 	string
)


var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if Version == "" {
			Version = "Unknow"
		}
		if Commit == "" {
			Commit = "Unknow"
		}
		fmt.Printf("version: %s\n", Version)
		if Date != "" {
			fmt.Printf("date: %s\n", Date)
		}
		if Commit != "" {
			fmt.Printf("commit: %s\n", Commit)
		}
		os.Exit(1)
	},
}

