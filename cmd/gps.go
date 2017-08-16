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
	"git.progwebtech.com/osiloke/gotrack/gpsd"
	"git.progwebtech.com/osiloke/gotrack/service"
	"github.com/mgutz/logxi/v1"
	dostow "github.com/osiloke/dostow-contrib/store"
	"github.com/osiloke/gostore"
	"github.com/spf13/cobra"
)

var (
	gpsdURI     string
	groupKey    string
	dostowUri   string
	storageMode string
	deviceID    string
	storePath   string
	logger      = log.New("[gps]")
)

// gpsCmd represents the gps command
var gpsCmd = &cobra.Command{
	Use:   "gps",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err   error
			store gostore.ObjectStore
		)
		if storageMode == "bolt" {
			store, err = gostore.NewBoltObjectStore(storePath + ".bolt")
			if err != nil {
				logger.Info("unable to create bolt store", "err", err)
			}
		} else if storageMode == "scribble" {
			store = gostore.NewScribbleStore(storePath + ".scribble")
		} else {
			store = dostow.NewStore("https://worksmart.progwebtech.com/v1", groupKey)
		}
		sensor := service.NewGpsService(store, gpsdURI, deviceID)
		api := service.NewApiService(store, sensor, ":8000")
		go api.Run()
		go sensor.Run()
		select {}
		logger.Info("exiting")
	},
}

func init() {
	RootCmd.AddCommand(gpsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	gpsCmd.Flags().StringVarP(&groupKey, "key", "k", "test", "Group access key")
	gpsCmd.Flags().StringVarP(&gpsdURI, "gpsd", "g", gpsd.DefaultAddress, "Gpsd uri")
	gpsCmd.Flags().StringVarP(&storageMode, "storage", "s", "scribble", "switches storage mode")
	gpsCmd.Flags().StringVarP(&deviceID, "id", "d", "gotrack1", "id of this device")
	gpsCmd.Flags().StringVarP(&storePath, "path", "p", "./gotrack.scrible", "id of this device")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gpsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
