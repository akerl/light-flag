package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/user"
	"strings"
	"time"

	"github.com/zshift/luxafor"
)

type Color struct {
	R uint8
	G uint8
	B uint8
}

var colors = map[string]Color{
	"red":     {255, 0, 0},
	"green":   {0, 255, 0},
	"blue":    {0, 0, 255},
	"cyan":    {0, 255, 255},
	"yellow":  {255, 255, 0},
	"magenta": {255, 0, 255},
	"black":   {0, 0, 0},
	"grey":    {15, 15, 15},
	"white":   {255, 255, 255},
}

func stateFile() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir + "/.light-flag"
}

func lastColors() (string, string) {
	file, err := ioutil.ReadFile(stateFile())
	if err != nil {
		return "off", "off"
	}
	s := strings.Split(string(file), " ")
	return s[0], s[1]
}

func saveLastColor(color, mini string) {
	_ = ioutil.WriteFile(stateFile(), []byte(fmt.Sprintf("%s %s", color, mini)), 0644)
}

func main() {
	luxs := luxafor.Enumerate()
	if len(luxs) == 0 {
		fmt.Println("No attached devices. Exiting")
		return
	}
	lux := luxs[1]
	lastSolid, lastMini := lastColors()
	var solid = flag.String("solid", lastSolid, "solid color to set")
	var blink = flag.String("blink", "", "blink color to set")
	var mini = flag.String("mini", lastMini, "mini color to set")
	var duration = flag.Int("duration", 1, "blink duration")
	var count = flag.Int("count", 1, "blink count")
	var side = flag.String("side", "", "blink front or back")
	flag.Parse()

	var blinkSide luxafor.LED

	switch *side {
	case "front":
		blinkSide = luxafor.FrontAll
	case "back":
		blinkSide = luxafor.BackAll
	default:
		blinkSide = luxafor.All

	}

	if *blink != "" {
		c := colors[*blink]
		for i := 0; i< *count; i++ {
			err := lux.Fade(blinkSide, c.R, c.G, c.B, 128)
			time.Sleep(time.Duration(*duration) * time.Second / time.Duration(*count))
			if err != nil {
				fmt.Println(err.Error())
			}
			c := colors[*solid]
			if err := lux.Solid(c.R, c.G, c.B); err != nil {
				fmt.Println(err.Error())
			}
			if *mini != "" {
				c := colors[*mini]
				if err := lux.Set(luxafor.BackTop, c.R, c.G, c.B); err != nil {
					fmt.Println(err.Error())
				}
			}
			time.Sleep(time.Duration(*duration) * time.Second / time.Duration(*count))
		}
	} else if *solid != "" {
		c := colors[*solid]
		if err := lux.Solid(c.R, c.G, c.B); err != nil {
			fmt.Println(err.Error())
		}
	}
	if *mini != "" {
		c := colors[*mini]
		if err := lux.Set(luxafor.BackTop, c.R, c.G, c.B); err != nil {
			fmt.Println(err.Error())
		}
	}
	saveLastColor(*solid, *mini)
}
