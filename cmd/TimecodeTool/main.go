package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/marcrleonard/TimecodeTool/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

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
		Version: timecodetool.VERSION,
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
			resp := timecodetool.NewValidateTimecode(startTc, fps)
			if jsonOutput {

				if prettyPrintJsonOutput {
					prettyJSON, err := json.MarshalIndent(resp, "", "  ")
					if err != nil {
						fmt.Println("Error encoding JSON:", err)
						os.Exit(1)
					}
					fmt.Println(string(prettyJSON))
				} else if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
					panic("Error encoding json")
				}
			} else {
				PrettyPrintValidate(resp)
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
			resp := timecodetool.NewSpanTimecode(startTc, endTc, fps, excludeLastTimecode)

			if jsonOutput {

				if prettyPrintJsonOutput {
					prettyJSON, err := json.MarshalIndent(resp, "", "  ")
					if err != nil {
						fmt.Println("Error encoding JSON:", err)
						os.Exit(1)
					}
					fmt.Println(string(prettyJSON))
				} else if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
					panic("Error encoding json")
				}
			} else {
				PrettyPrintSpan(resp)
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
			resp := timecodetool.NewCalculateTimecodes(args[0], args[1:], fps, excludeLastTimecode)

			if jsonOutput {

				if prettyPrintJsonOutput {
					prettyJSON, err := json.MarshalIndent(resp, "", "  ")
					if err != nil {
						fmt.Println("Error encoding JSON:", err)
						os.Exit(1)
					}
					fmt.Println(string(prettyJSON))
				} else if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
					panic("Error encoding json")
				}
			} else {
				PrettyPrintCalc(resp)
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
	// This isn't really a secret, I just think it is probably not useful to the user, so
	// I don't want it to be confusing.
	docsCmd.Hidden = true

	rootCmd.AddCommand(validateCmd, spanCmd, calcCmd, outputSchema, docsCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Exec error: %s", err)
	}

}

func validateJson(cmd *cobra.Command, args []string) error {
	// Check the flag dependency
	if cmd.Flags().Changed("pretty-print") && !cmd.Flags().Changed("json-output") {
		return fmt.Errorf("the --pretty-print flag requires the --json-output flag to be set")
	}
	return nil
}

func printSeparator() {
	fmt.Println("─────────────────────────────")
}

const title = "🎥 TimecodeTool"

// PrettyPrint displays the timecode validation results in a user-friendly format
func PrettyPrintValidate(r *timecodetool.ValidateResponse) {
	fmt.Println(title + " Validate")
	printSeparator()
	fmt.Printf("Input Timecode:   %s\n", r.InputTimecode)
	fmt.Printf("Frame Rate (FPS): %.2f\n", r.InputFps)

	if r.Valid {
		dfIndicator := ""
		if r.IsDf {
			dfIndicator = " (Drop Frame)"
		}
		fmt.Printf("Valid Timecode:   ✅ Yes%s\n", dfIndicator)
		fmt.Printf("Next Timecode:    %s\n", r.NextTimecode)
	} else {
		fmt.Printf("Valid Timecode:   ❌ No\n")
		fmt.Printf("Error:            %s\n", r.ErrorMsg)
	}

	printSeparator()
}

// prettyPrintSpanResponse displays the span command results in a user-friendly format
func PrettyPrintSpan(r *timecodetool.SpanResponse) {
	fmt.Println(title + " Span")
	printSeparator()

	// Helper function to handle invalid timecodes
	printInvalidTimecode := func(timecode string) string {
		if timecode == "" {
			return "❌ Invalid (Empty)"
		}
		return timecode
	}

	// Print First and Last Timecodes
	fmt.Printf("First Timecode:   %s\n", printInvalidTimecode(r.InputFirstTimecode))
	fmt.Printf("Last Timecode:    %s\n", printInvalidTimecode(r.InputLastTimecode))
	fmt.Printf("Frame Rate (FPS): %.2f\n", r.InputFps)

	// Output based on the validity of the span
	if r.Valid {
		fmt.Printf("Valid Span:       ✅ Yes\n")
		fmt.Printf("Start Frame Index:    %d\n", r.StartFrameIdx)
		fmt.Printf("Last Frame Index:     %d\n", r.LastFrameIdx)
		fmt.Printf("Length (Frames):      %d\n", r.LengthFrames)
		fmt.Printf("Length (Real Time):   %s\n", r.LengthTime)
		fmt.Printf("Length (Seconds):     %.2f\n", r.LengthSeconds)
		fmt.Printf("Length (Timecode):    %s\n", r.LengthTimecode)
		fmt.Printf("Next Timecode:        %s\n", r.NextTimecode)
	} else {
		fmt.Printf("Valid Span:        ❌ No\n")
		fmt.Printf("Error:                %s\n", r.ErrorMsg)
	}

	printSeparator()
}

func PrettyPrintCalc(c *timecodetool.CalcResponse) {

	steps := c.Steps

	fmt.Println(title + " Calculate")
	printSeparator()

	// Starting timecode and frames
	fmt.Printf(" 🎬 Starting Timecode:      %s (Index %d)\n", c.InputFirstTimecode, c.StartFrameIdx)

	// Process each step
	for _, step := range steps {
		if step.Operation == "+" {
			fmt.Printf("   ➕  Add Timecode:         %s (%d frames)\n", step.Timecode, step.Frames)
		} else if step.Operation == "-" {
			fmt.Printf("   ➖  Sub Timecode:         %s (%d frames)\n", step.Timecode, step.Frames)
		}
	}

	// Resulting timecode and frames
	printSeparator()
	fmt.Printf(" 🟰  Resulting Timecode:    %s (%d total frames)\n", c.LastTimecode, c.LengthFrames)
	fmt.Printf("%d ➡️ %d frame indexes\n", c.StartFrameIdx, c.LastFrameIdx)
	printSeparator()
}
