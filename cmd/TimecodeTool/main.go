package main

import (
	"TimecodeTool/internal"
	"fmt"
	"os"
	"strconv"
)

func main() {

	if len(os.Args) == 1 {
		println("TimecodeTool v0.1.0\nAuthor: Marc Leonard")
		os.Exit(0)
	} else if len(os.Args) == 2 {
		println("Please provide a timecode and an a framerate.")
	} else if len(os.Args) == 3 {

		tc := os.Args[1]
		_fr := os.Args[2]

		fr, err := strconv.ParseFloat(_fr, 64)
		if err != nil {
			println("Please provide a valid framerate.")
			os.Exit(1)
		}

		tcObj2 := internal.TimecodeFromString(tc, fr)

		println("Frame Index (0 based):", tcObj2.GetFrameIdx())
		fmt.Printf("Framerate: %f\n", fr)
		println(tc, tcObj2.Validate().Error())

		os.Exit(0)
	}

	os.Exit(0)

	// Todo Begin adding arg parsing!!!

	tcObj2 := internal.TimecodeFromString("00:06:59;25", 29.97)
	for x := 0; x < 10; x++ {
		tcObj2.AddFrames(1)
		tcObj2.Print()
	}

	fmt.Println("------")

	s := internal.NewTimecodeSpan("00:00:00;00", "00:10:00;17", 29.97)

	println("00:10:00;07", s.GetTotalFrames(), s.GetSpanRealtime())

	ws := internal.NewTimecodeSpan("00:00:00:00", "00:09:59:29", 29.97)

	println("00:09:59:29", ws.GetTotalFrames(), ws.GetSpanRealtime())

	fmt.Println("---Invalid DF Timecode---")

	// This is not a valid timecode.
	etc := internal.TimecodeFromString("00:07:00;00", 29.97)
	e := etc.Validate()
	if e != nil {
		println(e.Error())
	}

	correct_df := internal.TimecodeFromString("00:07:00;02", 29.97)
	e = correct_df.Validate()
	if e != nil {
		println(e.Error())
	}

	fmt.Println("------")

	intc := "00:01:00:24"
	etcc := internal.TimecodeFromString(intc, 24)
	etcc.Print()
	ee := etcc.Validate()
	if ee != nil {
		println(ee.Error())
	} else {
		println(intc, "is valid")
	}

}
