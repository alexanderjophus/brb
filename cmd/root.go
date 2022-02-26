package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	twitter "github.com/g8rswimmer/go-twitter/v2"
	"github.com/gosuri/uilive"
	helix "github.com/nicklaw5/helix/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// A go template compatible message
	Message string
)

const (
	defaultMessage = `Stream will continue in {{ .Countdown }}
{{ if .TwitchFollowerCount }}Twitch followers: {{ .TwitchFollowerCount }}{{ end }}
{{ if .TwitterFollowerCount }}Twitter followers: {{ .TwitterFollowerCount }}{{ end }}`
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

type Output struct {
	Countdown            time.Duration
	TwitchFollowerCount  int
	TwitterFollowerCount int
}

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

func run(duration string) error {
	d, err := time.ParseDuration(duration)
	if err != nil {
		return fmt.Errorf("parsing duration: %w", err)
	}

	Message = "\r" + Message
	deadline := time.Now().Add(d)

	o := Output{}

	if twitchEnabled() {
		client, err := helix.NewClient(&helix.Options{
			ClientID:       viper.GetString("twitchclientid"),
			ClientSecret:   viper.GetString("twitchclientsecret"),
			AppAccessToken: viper.GetString("twitchappaccesstoken"),
		})
		if err != nil {
			return fmt.Errorf("creating helix client: %w", err)
		}

		resp, err := client.GetUsersFollows(&helix.UsersFollowsParams{
			ToID: viper.GetString("twitchuserid"),
		})
		if err != nil {
			return fmt.Errorf("getting followers: %w", err)
		}
		o.TwitchFollowerCount = resp.Data.Total
	}

	if twitterEnabled() {
		client := &twitter.Client{
			Authorizer: authorize{
				Token: viper.GetString("twitterbearertoken"),
			},
			Client: http.DefaultClient,
			Host:   "https://api.twitter.com",
		}
		tu, err := client.UserNameLookup(context.Background(), []string{viper.GetString("twitterusername")}, twitter.UserLookupOpts{
			UserFields: []twitter.UserField{
				twitter.UserFieldPublicMetrics,
			},
		})
		if err != nil {
			return fmt.Errorf("getting twitter user: %w", err)
		}

		o.TwitterFollowerCount = tu.Raw.Users[0].PublicMetrics.Followers
	}

	w := uilive.New()
	w.Start()

	for range time.Tick(1 * time.Second) {
		if deadline.Before(time.Now()) {
			break
		}
		t, err := template.New("message").Parse(Message)
		if err != nil {
			return fmt.Errorf("parsing template: %w", err)
		}

		o.Countdown = time.Since(deadline).Truncate(time.Second) * -1
		err = t.Execute(w, o)
		if err != nil {
			return fmt.Errorf("executing template: %w", err)
		}
	}
	fmt.Fprintln(w, "Stream starting imminently")

	w.Stop() // flush and stop rendering
	return nil
}

func twitchEnabled() bool {
	return viper.GetString("twitchclientid") != "" &&
		viper.GetString("twitchclientsecret") != "" &&
		viper.GetString("twitchappaccesstoken") != "" &&
		viper.GetString("twitchuserid") != ""
}

func twitterEnabled() bool {
	return viper.GetString("twitterbearertoken") != "" &&
		viper.GetString("twitterusername") != ""
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVarP(&Message, "message", "m", defaultMessage, "Message to display, must be a valid go template")

	viper.BindPFlag("twitchclientid", rootCmd.PersistentFlags().Lookup("twitchclientid"))
	viper.BindPFlag("twitchclientsecret", rootCmd.PersistentFlags().Lookup("twitchclientsecret"))
	viper.BindPFlag("twitchappaccesstoken", rootCmd.PersistentFlags().Lookup("twitchappaccesstoken"))
	viper.BindPFlag("twitchuserid", rootCmd.PersistentFlags().Lookup("twitchuserid"))
	viper.BindPFlag("twitterbearertoken", rootCmd.PersistentFlags().Lookup("twitterbearertoken"))
	viper.BindPFlag("twitterusername", rootCmd.PersistentFlags().Lookup("twitterusername"))
}

func initConfig() {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	// Search config in home directory with name ".brb" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".brb")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(fmt.Errorf("reading config: %w", err))
	}
}
