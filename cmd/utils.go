package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"time"

	sh "github.com/codeskyblue/go-sh"
)

// `gpspipe -w | head -10 | grep TPV | sed -r 's/.*"time":"([^"]*)".*/\1/' | head -1`
var gpsDate = []interface{}{
	[]interface{}{"gpspipe", "-w"},
	[]interface{}{"head", "-10"},
	[]interface{}{"grep", "TPV"},
	// []interface{}{"sed", "-r", `'s/.*"time":"([^"]*)".*/\1/'`},
	[]interface{}{"head", "-1"},
}

func getCurrentTime() (*time.Time, error) {
	first := gpsDate[0].([]interface{})
	cmd := sh.Command(first[0].(string), first[1:]...)
	cmd.ShowCMD = true
	for _, v := range gpsDate[1:] {
		line := v.([]interface{})
		cmd = cmd.Command(line[0].(string), line[1:]...)
	}
	out, err := cmd.SetTimeout(time.Second * 7).Output()
	if err != nil {
		return nil, err
	}

	var tpv map[string]interface{}
	if err := json.Unmarshal(out, &tpv); err == nil {
		timestring, ok := tpv["time"].(string)
		if ok {
			t, err := time.Parse(time.RFC3339, string(timestring))
			if err != nil {
				return nil, err
			}
			if t.Year() > 2000 {
				// try system time
				return &t, nil
			}
		}
	}
	return date()
}
func date() (*time.Time, error) {
	out, err := sh.Command(`date`, `+%Y-%m-%dT%T%z`).Output()
	if err != nil {
		return nil, err
	}
	t, err := time.Parse(time.RFC3339, string(out))
	if err != nil {
		return nil, err
	}
	return &t, nil
}
func Mkdirp(dirpath string) {
	os.MkdirAll(dirpath, 0775)
}

// var bus = EventBus.New()
// findArduino looks for the file that represents the Arduino
// serial connection. Returns the fully qualified path to the
// device if we are able to find a likely candidate for an
// Arduino, otherwise an empty string if unable to find
// something that 'looks' like an Arduino device.
func findArduino() string {
	contents, _ := ioutil.ReadDir("/dev")

	// Look for what is mostly likely the Arduino device
	for _, f := range contents {
		if strings.Contains(f.Name(), "tty.usbmodem") ||
			strings.Contains(f.Name(), "ttyUSB") ||
			strings.Contains(f.Name(), "ttyACM") {
			return "/dev/" + f.Name()
		}
	}

	// Have not been able to find a USB device that 'looks'
	// like an Arduino.
	return ""
}
