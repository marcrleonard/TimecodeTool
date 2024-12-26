package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/marcrleonard/TimecodeTool"
	"github.com/marcrleonard/TimecodeTool/timecodetool"
	"github.com/spf13/cobra/doc"

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

	var rootCmd = &cobra.Command{
		Use:     "TimecodeTool [validate|span|calculate|schema] [args] [flags]",
		Short:   "A timecode CLI tool.",
		Version: TimecodeTool.VERSION,
		Long: "A timecode CLI tool.\n\n`TimecodeTool validate [args] [flags]` For timecode validation\n\n" +
			"`TimecodeTool span [args] [flags]` For timecode span length information\n\n" +
			"`TimecodeTool calculator [args] [flags]` for timecode calculations",
		PreRunE: validateJson,
	}

	validateCmd := &cobra.Command{
		Use:     "validate [flags] [Timecode]",
		Short:   "Returns a timecodes validity and information regarding the frame.",
		Args:    cobra.ExactArgs(1),
		Long:    "Returns a timecodes validity and information regarding the frame.",
		PreRunE: validateJson,
		Run: func(cmd *cobra.Command, args []string) {
			startTc := args[0]
			resp := timecodetool.ValidateTimecode(startTc, fps)
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
				resp.PrettyPrint()
			}
		},
	}
	validateCmd.Flags().BoolVar(&jsonOutput, "json-output", false, "Output as JSON")
	validateCmd.Flags().BoolVar(&prettyPrintJsonOutput, "pretty-print", false, "Output indented JSON")
	validateCmd.Flags().Float64Var(&fps, "fps", 29.97, "Frame rate of timecodes")
	validateCmd.MarkFlagsOneRequired("fps")

	spanCmd := &cobra.Command{
		Use:     "span [flags] [First Timecode] [Last Timecode]",
		Short:   "Get duration information spanning two timecodes.",
		Args:    cobra.ExactArgs(2),
		Long:    "Get duration information spanning two timecodes. Returns durations in frames, seconds, time, and more.",
		PreRunE: validateJson,
		Run: func(cmd *cobra.Command, args []string) {
			startTc := args[0]
			endTc := args[1]
			resp := timecodetool.SpanTimecode(startTc, endTc, fps, excludeLastTimecode)

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
				resp.PrettyPrint()
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
		Short:   "Timecode/Frame calculator. Enter either timecode strings or frame numbers. ",
		Args:    cobra.MinimumNArgs(3),
		PreRunE: validateJson,
		Long:    "Timecode/Frame calculator. Enter either timecode strings or frame numbers. It will add all these all together and generate span. When entering a timecode the amount of frames added/subtracted is relative to 00:00:00:00, with the timecode entered being inclusive. Use the `-e` flag to make it exclusive.",
		Run: func(cmd *cobra.Command, args []string) {
			resp := timecodetool.CalculateTimecodes(args[0], args[1:], fps, excludeLastTimecode)

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
				resp.PrettyPrint()
			}

		},
	}
	calcCmd.Flags().BoolVar(&jsonOutput, "json-output", false, "Output as JSON")
	calcCmd.Flags().BoolVar(&prettyPrintJsonOutput, "pretty-print", false, "Output indented JSON")
	calcCmd.Flags().BoolVarP(&excludeLastTimecode, "exclude-last-timecode", "e", false, `When entering a timecode to be added or subtracted, the calculations will be based off the timecode, minus one frame. This typically make it easier to read and enter timecode."`)
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
				r = jsonschema.Reflect(&timecodetool.ValidateResponse{})
			case "span":
				r = jsonschema.Reflect(&timecodetool.SpanResponse{})
			case "calculate":
				r = jsonschema.Reflect(&timecodetool.CalcResponse{})
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

	docsCmd := &cobra.Command{
		Use:   "gendocs",
		Short: "Generate CLI documentation",
		Run: func(cmd *cobra.Command, args []string) {
			// Generate Markdown files

			//outputDir := "./dist/md"
			outputDir := args[0]
			if _, err := os.Stat(outputDir); os.IsNotExist(err) {
				if err = os.Mkdir(outputDir, os.ModePerm); err != nil {
					return
				}
			}

			err := doc.GenMarkdownTree(rootCmd, outputDir)
			if err != nil {
				fmt.Println("Error generating markdown docs:", err)
				return
			}
			fmt.Printf("Markdown documentation generated in %s\n", outputDir)

		},
	}

	rootCmd.AddCommand(validateCmd, spanCmd, calcCmd, outputSchema, docsCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Exec error: %w", err)
	}

}
