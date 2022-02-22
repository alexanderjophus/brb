package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "brb",
	Short: "A tool for letting viewers know when your stream is starting again",
	Run: func(cmd *cobra.Command, args []string) {
		if err := run(args[0]); err != nil {
			log.Fatal(err)
		}
	},
	Args: cobra.ExactArgs(1),
}

func run(duration string) error {
	d, err := time.ParseDuration(duration)
	if err != nil {
		return fmt.Errorf("parsing duration: %w", err)
	}

	deadline := time.Now().Add(d)

	for range time.Tick(1 * time.Second) {
		if deadline.Before(time.Now()) {
			break
		}
		fmt.Printf("\rStream will start again in: %s", time.Now().Sub(deadline).Truncate(time.Second))
	}
	fmt.Println()
	fmt.Println("Hello")
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
