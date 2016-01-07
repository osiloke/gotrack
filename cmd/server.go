// Copyright © 2015 Osiloke Emoekpere <osi@progwebtech.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/everdev/mack"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/spf13/cobra"

	"git.progwebtech.com/osiloke/gotrack/process"
)

func setupArduino() *gobot.Robot {
	//Arduino bot
	pCtrl := process.NewPowerControl(true)

	device := findArduino()
	arduino := firmata.NewFirmataAdaptor("arduino", device)
	statusLed := gpio.NewLedDriver(arduino, "led", "13")
	powerButton := gpio.NewButtonDriver(arduino, "powerButton", "2")
	voltageMeasure := gpio.NewAnalogSensorDriver(arduino, "vin", "0")

	log.Println(time.Now().String())
	work := func() {
		gobot.On(powerButton.Event("push"), func(data interface{}) {
			log.Println("power button pushed")
		})
		gobot.On(powerButton.Event("release"), func(data interface{}) {
			log.Println("power button released")
			pCtrl.In("release")

		})
		gobot.Every(1*time.Second, func() {
			statusLed.Toggle()
		})
		gobot.On(voltageMeasure.Event("data"), func(data interface{}) {
			// _data := data.(int)
			// vin := uint8(
			//   gobot.ToScale(gobot.FromScale(float64(_data), 0, 1023), 0, 5),
			// )
			// log.Println(fmt.Sprintf("analog voltage in: %d", _data))
			// log.Println(fmt.Sprintf("vin: %d", vin) )
		})
	}

	return gobot.NewRobot("arduinoExpansion",
		[]gobot.Connection{arduino},
		[]gobot.Device{voltageMeasure, powerButton, statusLed},
		work,
	)
}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("server called")

		gbot := gobot.NewGobot()
		gbot.AddRobot(setupArduino())

		mack.Notify("Bot is running")
		mack.Say("Car system is on!", "Victoria")
		gbot.Start()
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
