package cmd

import (
	"log"
	"os"
	"time"

	"git.progwebtech.com/osiloke/gotrack/gpsd"
	gps "git.progwebtech.com/osiloke/gotrack/service"
	"github.com/kardianos/service"
	dostow "github.com/osiloke/dostow-contrib/store"
	"github.com/osiloke/gostore"
	"github.com/spf13/cobra"
)

var klogger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
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
		store = dostow.NewStore(dostowUri, groupKey)
	}
	sensor := gps.NewGpsService(store, gpsdURI, deviceID)
	api := gps.NewApiService(store, sensor, ":8000")
	go api.Run()
	go sensor.Run()
	select {}
	logger.Info("exiting")
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	<-time.After(time.Second * 13)
	return nil
}

// gpsService represents the gps command
var gpsService = &cobra.Command{
	Use:   "service",
	Short: "Manage gps service",
	Long:  `Manage gps service.`,
	Run: func(cmd *cobra.Command, args []string) {
		svcConfig := &service.Config{
			Name:        "Gps",
			DisplayName: "Gps Service",
			Description: "GPS service which logs position",
			Arguments: []string{"gps", "-k", groupKey, "-g", gpsdURI, "-s",
				storageMode, "-d", deviceID, "-p", storePath},
		}

		prg := &program{}
		s, err := service.New(prg, svcConfig)
		if err != nil {
			log.Fatal(err)
		}
		if len(os.Args) > 1 {
			err = service.Control(s, args[0])
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		klogger, err = s.Logger(nil)
		if err != nil {
			log.Fatal(err)
		}
		err = s.Run()
		if err != nil {
			klogger.Error(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(gpsService)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	gpsService.Flags().StringVarP(&groupKey, "key", "k", "test", "Group access key")
	gpsService.Flags().StringVarP(&dostowUri, "dostowUri", "u", "https://worksmart.progwebtech.com/v1", "dostow Uri")
	gpsService.Flags().StringVarP(&gpsdURI, "gpsd", "g", gpsd.DefaultAddress, "Gpsd uri")
	gpsService.Flags().StringVarP(&storageMode, "storage", "s", "scribble", "switches storage mode")
	gpsService.Flags().StringVarP(&deviceID, "id", "d", "gotrack1", "id of this device")
	gpsService.Flags().StringVarP(&storePath, "path", "p", "./gotrack.scrible", "id of this device")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gpsService.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
