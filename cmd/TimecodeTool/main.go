package main

import (
	"fmt"
	"strconv"

	"TimecodeTool"
	"TimecodeTool/timecode"

	"github.com/spf13/cobra"
)

func main() {

	var fps float64
	var jsonOutput bool
	var prettyPrintJsonOutput bool

	rootCmd := &cobra.Command{
		Use:     "TimecodeTool",
		Short:   "A timecode CLI tool",
		Version: TimecodeTool.VERSION,
		Long: fmt.Sprintf(
			"TimecodeTool --fps=29.97 [Timecode]\n" +
				"TimecodeTool --fps=29.97 [First Timecode] [Last Timecode]\n" +
				"TimecodeTool --fps=29.97 [Timecode] + [Timecode] - [Frames]\n"),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Check the flag dependency
			if cmd.Flags().Changed("pretty-print") && !cmd.Flags().Changed("json-output") {
				return fmt.Errorf("the --pretty-print flag requires the --json-output flag to be set")
			}
			return nil
		},
	}

	rootCmd.PersistentFlags().Float64Var(&fps, "fps", 29.97, "Frame rate of timecodes")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json-output", false, "Output as JSON")
	rootCmd.PersistentFlags().BoolVar(&prettyPrintJsonOutput, "pretty-print", false, "Output indented JSON")
	rootCmd.MarkFlagsOneRequired("fps")

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

	var excludeLastTimecode bool
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
			if excludeLastTimecode {
				lastTimecode.AddFrames(-1)
			}

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
	spanCmd.Flags().BoolVarP(&excludeLastTimecode, "exclude-last-timecode", "e", false, `When entering a first timecode and a last timecode, the calculations will be based off the last timecode, minus one frame. This typically make it easier to read and enter timecode. For instance, with this flag set, a span of "00:00:00:00" "00:00:01:00" represents one second.`)

	calcCmd := &cobra.Command{
		Use:   "calculate --fps=29.97 [First Timecode] + [Timecode] - [frame number]",
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

				frames, err := timecode.ParseStringToFrames(nextTime, fps, excludeLastTimecode)
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
	calcCmd.Flags().BoolVarP(&excludeLastTimecode, "exclude-last-timecode", "e", false, `When entering a timecode to be added or subtracted, the calculations will be based off the timecode, minus one frame. This typically make it easier to read and enter timecode. For instance, with this flag set `+"`TimecodeTool calculate \"00:00:00:00\" + \"00:00:00:01\" --fps=23.976 -e`"+"will yield `00:00:00:01`")

	rootCmd.AddCommand(validateCmd, spanCmd, calcCmd)

	if err := rootCmd.Execute(); err != nil {
		//fmt.Printf("Error: %w", err)
	}

}
