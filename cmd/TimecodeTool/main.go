package main

import (
	"fmt"
	"os"
	"strconv"

	"TimecodeTool/timecode"
	"github.com/spf13/cobra"
)

func main() {

	var fps float64
	var jsonOutput bool

	version := "v0.1.0"

	rootCmd := &cobra.Command{
		Use:     "TimecodeTool",
		Short:   "A timecode CLI tool",
		Version: version,
		Long: fmt.Sprintf(
			"TimecodeTool --fps=29.97 [Timecode]\n" +
				"TimecodeTool --fps=29.97 [First Timecode] [Last Timecode]\n" +
				"TimecodeTool --fps=29.97 [Timecode] + [Timecode] - [Frames]\n"),
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) == 0 {
				cmd.Help()
				os.Exit(1)
			} else if len(args) == 1 {
				startTc := args[0]
				firstTc, _ := timecode.NewTimecodeFromString(startTc, fps)
				fmt.Printf("Timecode: %s@%f (%s)\n", startTc, fps, firstTc.GetValid())
			} else if len(args) == 2 {

				startTc := args[0]
				endTc := args[1]

				firstTc, _ := timecode.NewTimecodeFromString(startTc, fps)
				lastTimecode, _ := timecode.NewTimecodeFromString(endTc, fps)

				fmt.Printf("First Timecode: %s (%s)\n", firstTc.GetTimecode(), firstTc.GetValid())
				fmt.Printf("Last Timecode: %s (%s)\n", lastTimecode.GetTimecode(), lastTimecode.GetValid())

				dfness := "NDF"
				if lastTimecode.DropFrame {
					dfness = "DF"
				}
				println("Last Timecode Frame Index (0 based):", lastTimecode.GetFrameIdx())
				fmt.Printf(
					"Framerate: %s%s\n",
					strconv.FormatFloat(fps, 'f', -1, 64),
					dfness,
				)

			} else {

				firstTc, _ := timecode.NewTimecodeFromString(args[0], fps)
				curIdx := 1
				for {
					if curIdx >= len(args)-1 {
						break
					}

					opperator := args[curIdx]

					nextTime := args[curIdx+1]

					frames, err := timecode.ParseStringToFrames(nextTime, fps)
					if err != nil {
						return
					}

					switch opperator {
					case "-":
						firstTc.AddFrames(int(frames) * -1)
					case "+":
						firstTc.AddFrames(int(frames))
					}

					curIdx += 2
				}

				fmt.Println(firstTc.GetTimecode())
			}
		},
	}

	rootCmd.PersistentFlags().Float64Var(&fps, "fps", 29.97, "Frame rate of timecodes")
	rootCmd.Flags().BoolVar(&jsonOutput, "json-output", false, "Output as JSON")

	rootCmd.MarkFlagsOneRequired("fps")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}

}
