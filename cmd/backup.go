// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jhoonb/archivex"
	"github.com/osiloke/dostow-contrib/api"
	"github.com/spf13/cobra"
)

var backupPath string

//Backup defines a geo json backup
type Backup struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Device string `json:"device"`
	Size   int64  `json:"size"`
}

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backups up gps data",
	Long: `Backups up gps data. 
	Compresses and backs up json data. 
	It can also delete existing data`,
	Run: func(cmd *cobra.Command, args []string) {
		Mkdirp(backupPath)
		cl := api.NewClient(dostowUri, groupKey)

		if storageMode == "scribble" {
			storePath = storePath + ".scribble"
			arch := new(archivex.TarFile)
			filepath.Walk(storePath, func(path string, info os.FileInfo, err error) error {
				if !strings.Contains(info.Name(), filepath.Base(storePath)) {
					if info.IsDir() {
						archiveName := info.Name() + ".tar.gz"
						safeName := strings.Replace(archiveName, "-", "_", -1)
						fullPath := filepath.Join(backupPath, safeName)
						logger.Info("archiving " + info.Name() + " at " + path + " to " + fullPath)
						if _, err := os.Stat(fullPath); os.IsNotExist(err) {
							arch.Create(fullPath)
							defer arch.Close()
							arch.AddAll(path, false)
							//decide to delete old database files
						} else {
							logger.Warn("existing backup at " + fullPath)
						}
						//check if we can upload and delete backup
						// check if database files should be deleted
						backup := &Backup{
							Name:   info.Name(),
							Size:   info.Size(),
							Device: deviceID,
						}
						raw, err := cl.Store.Create("backup", backup)
						if err != nil {
							logger.Warn(err.Error())
							if e, ok := err.(*api.APIError); ok {
								if e.Message == "already exists" {
									//get id and upload
								}
							}
						} else {
							json.Unmarshal(*raw, backup)
							f, err := os.Open(fullPath)
							if err != nil {
								return err
							}
							defer f.Close()
							if _, raw, err := cl.File.Create("backup", backup.ID, "file", info.Name(), f); err != nil {
								logger.Warn(info.Name()+" was not uploaded", "err", err, "rsp", raw)
							} else {
								t := time.Now()
								diff := t.Sub(info.ModTime())
								if diff.Hours() > 0 {
									days := diff.Hours() / 24
									if days > 1 {
										logger.Warn("deleting "+path, "info", info)
									}
								}
							}
						}

					}
				}
				return nil
			})

		}
	},
}

func init() {
	RootCmd.AddCommand(backupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// backupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// backupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	backupCmd.Flags().StringVarP(&dostowUri, "dostowUri", "u", "https://cloud.dostow.com/v1/", "dostow Uri")
	backupCmd.Flags().StringVarP(&deviceID, "id", "i", "gotrack1", "id of this device")
	backupCmd.Flags().StringVarP(&groupKey, "groupKey", "g", "groupKey", "dostow group key")
	backupCmd.Flags().StringVarP(&storageMode, "storage", "s", "scribble", "switches storage mode")
	backupCmd.Flags().StringVarP(&storePath, "store-path", "t", "./gotrack", "id of this device")
	backupCmd.Flags().StringVarP(&backupPath, "path", "p", "./gotrack_backup", "id of this device")
}
