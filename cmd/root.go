package cmd

import (
	"flag"
	"fmt"
	"os"

	tl "github.com/JoelOtter/termloop"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog"

	"github.com/sudermanjr/text-game/pkg/game"
)

var (
	version = "development"
	commit  = "n/a"
	height  int
	width   int
	fps     float64
)

func init() {
	// Flags
	rootCmd.PersistentFlags().IntVar(&height, "height", 80, "The height of the arena")
	rootCmd.PersistentFlags().IntVar(&width, "width", 200, "The width of the arena")
	rootCmd.PersistentFlags().Float64Var(&fps, "framerate", 30, "The framerate of the game for termloop")

	//Commands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(startCmd)

	environmentVariables := map[string]string{
		"GAME_HEIGHT": "height",
		"GAME_WIDTH":  "width",
	}

	for env, flag := range environmentVariables {
		flag := rootCmd.PersistentFlags().Lookup(flag)
		flag.Usage = fmt.Sprintf("%v [%v]", flag.Usage, env)
		if value := os.Getenv(env); value != "" {
			err := flag.Value.Set(value)
			if err != nil {
				klog.Errorf("Error setting flag %v to %s from environment variable %s", flag, value, env)
			}
		}
	}

	klog.InitFlags(nil)
	flag.Parse()
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
}

var rootCmd = &cobra.Command{
	Use:   "text-game",
	Short: "text-game",
	Long:  `A text game`,
	Run: func(cmd *cobra.Command, args []string) {
		klog.Error("You must specify a sub-command.")
		err := cmd.Help()
		if err != nil {
			klog.Error(err)
		}
		os.Exit(1)
	},
}

// Execute the stuff
func Execute(VERSION string, COMMIT string) {
	version = VERSION
	commit = COMMIT
	if err := rootCmd.Execute(); err != nil {
		klog.Error(err)
		os.Exit(1)
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the current version of the tool.",
	Long:  `Prints the current version.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:" + version + " Commit:" + commit)
	},
}

var startCmd = &cobra.Command{
	Use:     "start",
	Short:   "Starts the game",
	Long:    "Starts the game",
	Aliases: []string{"run"},
	Run: func(cmd *cobra.Command, args []string) {
		instance := tl.NewGame()
		instance.Screen().SetFps(fps)
		level := game.BuildLevel(instance, width, height)

		instance.Screen().SetLevel(level)
		instance.Start()
	},
}
