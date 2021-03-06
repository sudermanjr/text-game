package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog"

	"github.com/sudermanjr/text-game/pkg/game"
)

var (
	version    = "development"
	commit     = "n/a"
	height     int
	width      int
	fps        float64
	mapType    string
	showFPS    bool
	customSeed int64
	numRooms   int
)

func init() {
	//Commands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(startCmd)

	startCmd.PersistentFlags().IntVar(&height, "height", 70, "The height of the arena.")
	startCmd.PersistentFlags().IntVar(&width, "width", 200, "The width of the arena.")
	startCmd.PersistentFlags().Float64Var(&fps, "framerate", 30, "The framerate of the game for termloop")
	startCmd.PersistentFlags().StringVar(&mapType, "map-type", "rooms", "The type of map. Must be one of (drunkwalk|rooms)")
	startCmd.PersistentFlags().BoolVar(&showFPS, "show-fps", false, "Enables the FPS text")
	startCmd.PersistentFlags().Int64VarP(&customSeed, "seed", "s", 0, "If non-zero, use this seed for the level generation. 0 will cause a random seed to be used.")
	startCmd.PersistentFlags().IntVarP(&numRooms, "room-count", "n", 0, "If non-zero, the number of rooms in a room-type map.")

	klog.InitFlags(nil)
	flag.Set("logtostderr", "false")
	flag.Set("log_file", "game.log")

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
		player := game.NewPlayer('@', height, width, "Andy", mapType, fps, customSeed, numRooms)
		player.Start()

		klog.Flush()
		fmt.Println("Thanks for playing!")
	},
}
