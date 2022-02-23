package cmd

import (
	"fmt"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/spf13/cobra"
)

// A go template compatible message
var Message string

const (
	defaultMessage = "Stream will start again in {{ . }}"
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

	Message = "\r" + Message
	deadline := time.Now().Add(d)

	for range time.Tick(1 * time.Second) {
		if deadline.Before(time.Now()) {
			break
		}
		t, err := template.New("message").Parse(Message)
		if err != nil {
			return fmt.Errorf("parsing template: %w", err)
		}

		err = t.Execute(os.Stdout, time.Since(deadline).Truncate(time.Second)*-1)
		if err != nil {
			return fmt.Errorf("executing template: %w", err)
		}
	}
	fmt.Println()
	fmt.Println("Stream starting imminently")
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().StringVarP(&Message, "message", "m", defaultMessage, "Message to display, must be a valid go template")
}
