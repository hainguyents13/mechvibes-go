package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	hook "github.com/robotn/gohook"
)

func main() {
	add()
}

func add() {
	f, err := os.Open("./audio/down.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Press any key ---")
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/100))
	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()

	done := make(chan bool)
	hook.Register(hook.KeyDown, []string{}, func(e hook.Event) {
		fmt.Println(e.Keychar, e.Keycode)
		shot := buffer.Streamer(0, buffer.Len())
		speaker.Play(beep.Seq(shot, beep.Callback(func() {
			done <- true
		})))
	})

	<-done
	s := hook.Start()
	<-hook.Process(s)
}
