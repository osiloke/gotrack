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
	"git.progwebtech.com/osiloke/gotrack/component"
	"git.progwebtech.com/osiloke/gotrack/gpsd"
	"github.com/mgutz/logxi/v1"
	dostow "github.com/osiloke/dostow-contrib/store"
	"github.com/osiloke/gostore"
	"github.com/spf13/cobra"
	"os"
)

var (
	gpsdURI     string
	groupKey    string
	storageMode string
	deviceID    string
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
		var s gostore.ObjectStore
		if storageMode == "local" {
			s = gostore.NewBoltDatabaseStore("./gotrack.db")
		} else {
			s = dostow.NewStore("https://worksmart.progwebtech.com/v1", groupKey)
		}
		gps, err := gpsd.Dial(gpsdURI)
		if err != nil {
			os.Exit(-1)
		}

		in := make(chan interface{})
		gpsComponent := component.NewComponent("gps", func(val interface{}, out chan<- interface{}) error {
			storedReport := gpsd.TPVToGeoJSON(val.(*gpsd.TPVReport))
			storedReport["deviceID"] = deviceID
			k, err := s.Save("geo", storedReport)
			if err != nil {
				log.Warn("report not saved", "report", storedReport, "err", err.Error())
				return err
			}
			log.Debug("gps report", "key", k, "Report", storedReport)
			return nil
		})
		gpsNet := component.NewLinearGraph(in, gpsComponent)

		gps.AddFilter("TPV", func(r interface{}) {
			in <- r
		})
		done := gps.Watch()

		<-done
		<-gpsNet.Wait()
	},
}

func init() {
	RootCmd.AddCommand(gpsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	gpsCmd.Flags().StringVarP(&groupKey, "key", "k", "test", "Group access key")
	gpsCmd.Flags().StringVarP(&gpsdURI, "gpsd", "g", gpsd.DefaultAddress, "Gpsd uri")
	gpsCmd.Flags().StringVarP(&storageMode, "storage", "s", "local", "switches storage mode")
	gpsCmd.Flags().StringVarP(&deviceID, "id", "d", "gotrack1", "id of this device")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gpsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
