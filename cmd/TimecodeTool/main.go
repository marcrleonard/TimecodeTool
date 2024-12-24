package main

import (
	"encoding/json"
	"fmt"
	"os"

	"TimecodeTool"
	"TimecodeTool/handlers"
	"TimecodeTool/timecode"
	"github.com/invopop/jsonschema"

	"github.com/spf13/cobra"
)

func validateJson(cmd *cobra.Command, args []string) error {
	// Check the flag dependency
	if cmd.Flags().Changed("pretty-print") && !cmd.Flags().Changed("json-output") {
		return fmt.Errorf("the --pretty-print flag requires the --json-output flag to be set")
	}
	return nil
}

func main() {

	var (
		fps                   float64
		jsonOutput            bool
		prettyPrintJsonOutput bool
		excludeLastTimecode   bool
	)

	rootCmd := &cobra.Command{
		Use:     "TimecodeTool",
		Short:   "A timecode CLI tool",
		Version: TimecodeTool.VERSION,
		Long: fmt.Sprintf(
			"TimecodeTool --fps=29.97 [Timecode]\n" +
				"TimecodeTool --fps=29.97 [First Timecode] [Last Timecode]\n" +
				"TimecodeTool --fps=29.97 [Timecode] + [Timecode] - [Frames]\n"),
		PreRunE: validateJson,
	}

	validateCmd := &cobra.Command{
		Use:   "validate [flags] [Timecode]",
		Short: "Validate a timecode",
		Args:  cobra.ExactArgs(1),
		Long: fmt.Sprintf(
			"TimecodeTool validate --fps=29.97 [Timecode]"),
		PreRunE: validateJson,
		Run: func(cmd *cobra.Command, args []string) {
			startTc := args[0]
			resp := handlers.ValidateTimecode(startTc, fps)
			if jsonOutput {

				if prettyPrintJsonOutput {
					prettyJSON, err := json.MarshalIndent(resp, "", "  ")
					if err != nil {
						fmt.Println(os.Stderr, "Error encoding JSON:", err)
						os.Exit(1)
					}
					fmt.Println(string(prettyJSON))
				} else if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
					panic("Error encoding json")
				}
			} else {
				timecode.PrettyPrint(resp)
			}
		},
	}
	validateCmd.Flags().BoolVar(&jsonOutput, "json-output", false, "Output as JSON")
	validateCmd.Flags().BoolVar(&prettyPrintJsonOutput, "pretty-print", false, "Output indented JSON")
	validateCmd.Flags().Float64Var(&fps, "fps", 29.97, "Frame rate of timecodes")
	validateCmd.MarkFlagsOneRequired("fps")

	spanCmd := &cobra.Command{
		Use:   "span [flags] [First Timecode] [Last Timecode]",
		Short: "Get the information spanning two timecodes",
		Args:  cobra.ExactArgs(2),
		Long: fmt.Sprintf(
			"TimecodeTool span --fps=29.97 [First Timecode] [Last Timecode]"),
		PreRunE: validateJson,
		Run: func(cmd *cobra.Command, args []string) {
			startTc := args[0]
			endTc := args[1]
			resp := handlers.SpanTimecode(startTc, endTc, fps, excludeLastTimecode)

			if jsonOutput {

				if prettyPrintJsonOutput {
					prettyJSON, err := json.MarshalIndent(resp, "", "  ")
					if err != nil {
						fmt.Println(os.Stderr, "Error encoding JSON:", err)
						os.Exit(1)
					}
					fmt.Println(string(prettyJSON))
				} else if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
					panic("Error encoding json")
				}
			} else {
				timecode.PrettyPrintSpanResponse(*resp)
			}
		},
	}
	spanCmd.Flags().BoolVar(&jsonOutput, "json-output", false, "Output as JSON")
	spanCmd.Flags().BoolVar(&prettyPrintJsonOutput, "pretty-print", false, "Output indented JSON")
	spanCmd.Flags().BoolVarP(&excludeLastTimecode, "exclude-last-timecode", "e", false, `When entering a first timecode and a last timecode, the calculations will be based off the last timecode, minus one frame. This typically make it easier to read and enter timecode. For instance, with this flag set, a span of "00:00:00:00" "00:00:01:00" represents one second.`)
	spanCmd.Flags().Float64Var(&fps, "fps", 29.97, "Frame rate of timecodes")
	spanCmd.MarkFlagsOneRequired("fps")

	calcCmd := &cobra.Command{
		Use:     "calculate --fps=29.97 [First Timecode] + [Timecode] - [frame number]",
		Short:   "Timecode/Frame calculator",
		Args:    cobra.MinimumNArgs(3),
		PreRunE: validateJson,
		Long: fmt.Sprintf(
			"TimecodeTool calculate --fps=29.97 [Timecode] + [Last Timecode] - [frame number]"),
		Run: func(cmd *cobra.Command, args []string) {
			resp := handlers.CalculateTimecodes(args[0], args[1:], fps, excludeLastTimecode)

			if jsonOutput {

				if prettyPrintJsonOutput {
					prettyJSON, err := json.MarshalIndent(resp, "", "  ")
					if err != nil {
						fmt.Println(os.Stderr, "Error encoding JSON:", err)
						os.Exit(1)
					}
					fmt.Println(string(prettyJSON))
				} else if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
					panic("Error encoding json")
				}
			} else {
				timecode.PrettyPrintCalculateResponse(*resp)
			}

		},
	}
	calcCmd.Flags().BoolVar(&jsonOutput, "json-output", false, "Output as JSON")
	calcCmd.Flags().BoolVar(&prettyPrintJsonOutput, "pretty-print", false, "Output indented JSON")
	calcCmd.Flags().BoolVarP(&excludeLastTimecode, "exclude-last-timecode", "e", false, `When entering a timecode to be added or subtracted, the calculations will be based off the timecode, minus one frame. This typically make it easier to read and enter timecode. For instance, with this flag set `+"`TimecodeTool calculate \"00:00:00:00\" + \"00:00:00:01\" --fps=23.976 -e`"+"will yield `00:00:00:01`")
	calcCmd.Flags().Float64Var(&fps, "fps", 29.97, "Frame rate of timecodes")
	calcCmd.MarkFlagsOneRequired("fps")

	outputSchema := &cobra.Command{
		Use: "schema [validate|span|calculate]",
		Long: "Returns a valid json schema that describes the json output of each of the tools in the CLI. Examples:" +
			"\n  TimecodeTool schema validate" +
			"\n  TimecodeTool schema span" +
			"\n  TimecodeTool schema calculate",
		Args:      cobra.ExactArgs(1), // Expect exactly one argument
		ValidArgs: []string{"validate", "span", "calculate"},
		Run: func(cmd *cobra.Command, args []string) {
			var r *jsonschema.Schema

			switch args[0] {
			case "validate":
				r = jsonschema.Reflect(&timecode.ValidateResponse{})
			case "span":
				r = jsonschema.Reflect(&timecode.SpanResponse{})
			case "calculate":
				r = jsonschema.Reflect(&timecode.CalcResponse{})
			default:
				// Handle invalid argument, could return an error or show a message
				fmt.Println(`Invalid argument. Valid options are: "validate", "span", "calculate"`)
				return
			}

			// Marshal the schema to JSON with indentation
			schemaJSON, err := json.MarshalIndent(r, "", "  ")
			if err != nil {
				panic(err)
			}
			// Output the JSON Schema
			fmt.Println(string(schemaJSON))
		},
	}

	rootCmd.AddCommand(validateCmd, spanCmd, calcCmd, outputSchema)

	if err := rootCmd.Execute(); err != nil {
		//fmt.Printf("Error: %w", err)
	}

}
