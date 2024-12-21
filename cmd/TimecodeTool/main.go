package main

import (
	"fmt"
	"os"
	"strconv"

	"TimecodeTool"
	"TimecodeTool/timecode"

	"github.com/spf13/cobra"
)

func main() {

	var fps float64
	var jsonOutput bool

	rootCmd := &cobra.Command{
		Use:     "TimecodeTool",
		Short:   "A timecode CLI tool",
		Version: TimecodeTool.VERSION,
		Long: fmt.Sprintf(
			"TimecodeTool --fps=29.97 [Timecode]\n" +
				"TimecodeTool --fps=29.97 [First Timecode] [Last Timecode]\n" +
				"TimecodeTool --fps=29.97 [Timecode] + [Timecode] - [Frames]\n"),
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) == 0 {
				cmd.Help()
				os.Exit(1)
			} else if len(args) == 1 {

			} else if len(args) == 2 {

			} else {

			}
		},
	}

	validateCmd := &cobra.Command{
		Use:   "validate [flags] [Timecode]",
		Short: "Validate a timecode",
		Args:  cobra.ExactArgs(1),
		Long: fmt.Sprintf(
			"TimecodeTool validate --fps=29.97 [Timecode]"),
		Run: func(cmd *cobra.Command, args []string) {
			startTc := args[0]
			firstTc, err := timecode.NewTimecodeFromString(startTc, fps)
			if err != nil {
				fmt.Printf("Error: %w", err)
			} else {
				fmt.Printf("Timecode: %s@%f (%s)\n", startTc, fps, firstTc.GetValid())

			}
		},
	}
	validateCmd.PersistentFlags().Float64Var(&fps, "fps", 29.97, "Frame rate of timecodes")
	validateCmd.MarkFlagsOneRequired("fps")

	spanCmd := &cobra.Command{
		Use:   "span [flags] [First Timecode] [Last Timecode]",
		Short: "Get the information spanning two timecodes",
		Args:  cobra.ExactArgs(2),
		Long: fmt.Sprintf(
			"TimecodeTool span --fps=29.97 [First Timecode] [Last Timecode]"),
		Run: func(cmd *cobra.Command, args []string) {
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
		},
	}
	spanCmd.PersistentFlags().Float64Var(&fps, "fps", 29.97, "Frame rate of timecodes")
	spanCmd.MarkFlagsOneRequired("fps")

	calcCmd := &cobra.Command{
		Use:   "calculate",
		Short: "Timecode/Frame calculator",
		Args:  cobra.MinimumNArgs(3),
		Long: fmt.Sprintf(
			"TimecodeTool calculate --fps=29.97 [Timecode] + [Last Timecode] - [frame number]"),
		Run: func(cmd *cobra.Command, args []string) {

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
		},
	}
	calcCmd.PersistentFlags().Float64Var(&fps, "fps", 29.97, "Frame rate of timecodes")
	calcCmd.MarkFlagsOneRequired("fps")

	rootCmd.AddCommand(validateCmd, spanCmd, calcCmd)
	rootCmd.Flags().BoolVar(&jsonOutput, "json-output", false, "Output as JSON")

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %w", err)
	}

}
