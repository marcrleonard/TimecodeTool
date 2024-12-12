package main

import (
	"fmt"
	"os"
	"strconv"

	"TimecodeTool/timecode"
)

func main() {

	tt := timecode.TimecodeFromFrames(0, 29.97, true)
	println(tt.GetTimecode())
	println(tt.GetFrameIdx())
	// println(tt.GetFrameIdx())
	// println(tt.GetValid())
	// println("--------")

	//var firstTimecode internal.Timecode;
	// var lastTimecode internal.Timecode;

	if len(os.Args) == 1 {
		println("TimecodeTool v0.1.0\nAuthor: Marc Leonard")
		println("	TimecodeTool [End Timecode] [FPS]")
		println("	TimecodeTool [Start Timecode] [End Timecode] [FPS]")
		os.Exit(0)
	} else if len(os.Args) == 2 {
		println("Please provide a timecode and an a framerate.")
	} else if len(os.Args) > 2 {

		tc := os.Args[2]

		_fr := os.Args[1]

		fr, err := strconv.ParseFloat(_fr, 64)
		if err != nil {
			println("Please provide a valid framerate.")
			os.Exit(1)
		}

		lastTimecode := timecode.NewTimecodeFromString(tc, fr)

		firstTimecode := timecode.TimecodeFromFrames(1, lastTimecode.FrameRate, lastTimecode.DropFrame)
		if len(os.Args) == 4 {
			firstTimecode = timecode.NewTimecodeFromString(os.Args[3], fr)
		}

		fmt.Printf("First Timecode: %s (%s)\n", firstTimecode.GetTimecode(), firstTimecode.GetValid())
		fmt.Printf("Last Timecode: %s (%s)\n", lastTimecode.GetTimecode(), lastTimecode.GetValid())

		println("Last Timecode Frame Index (0 based):", lastTimecode.GetFrameIdx())
		fmt.Printf("Framerate: %s\n", lastTimecode.GetFramerateString())
		fmt.Println("Dropframe Timecode: ", lastTimecode.DropFrame)

	}

	os.Exit(0)

	// Todo Begin adding arg parsing!!!

	tcObj2 := timecode.NewTimecodeFromString("00:06:59;25", 29.97)
	for x := 0; x < 10; x++ {
		tcObj2.AddFrames(1)
		tcObj2.Print()
	}

	fmt.Println("------")

	s := timecode.NewTimecodeSpan("00:00:00;00", "00:10:00;17", 29.97)

	println("00:10:00;07", s.GetTotalFrames(), s.GetSpanRealtime())

	ws := timecode.NewTimecodeSpan("00:00:00:00", "00:09:59:29", 29.97)

	println("00:09:59:29", ws.GetTotalFrames(), ws.GetSpanRealtime())

	fmt.Println("---Invalid DF Timecode---")

	// This is not a valid timecode.
	etc := timecode.NewTimecodeFromString("00:07:00;00", 29.97)
	e := etc.Validate()
	if e != nil {
		println(e.Error())
	}

	correct_df := timecode.NewTimecodeFromString("00:07:00;02", 29.97)
	e = correct_df.Validate()
	if e != nil {
		println(e.Error())
	}

	fmt.Println("------")

	intc := "00:01:00:24"
	etcc := timecode.NewTimecodeFromString(intc, 24)
	etcc.Print()
	ee := etcc.Validate()
	if ee != nil {
		println(ee.Error())
	} else {
		println(intc, "is valid")
	}

}
