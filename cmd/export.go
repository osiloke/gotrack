// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
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
	"fmt"
	"github.com/go-resty/resty"
	"github.com/mgutz/logxi/v1"
	"github.com/spf13/cobra"
	// "time"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export geo data to an external database",
	Long: `Export geo data to an external database.
	Only dostow is supported for now`,
}

type dostowData struct {
	Data       []interface{} `json:"data"`
	TotalCount int           `json:"total_count"`
}

var (
	pageCount       int
	pages           int
	gotrackUrl      string
	dostowStores    []string
	dostowUrl       string
	dostowAccessKey string
	dry             bool
	dostowCmd       = &cobra.Command{
		Use:   "dostow",
		Short: "Export geo data to a dostow account",
		Long:  `Export geo data to a dostow account`,
		Run: func(cmd *cobra.Command, args []string) {

			var (
				last string
			)
			var count int = 0
			for pages > 0 {
				resp, err := resty.R().
					SetQueryParams(map[string]string{
					"size":    "100",
					"_before": "20",
					"before":  last,
				}).
					SetHeader("Accept", "application/json").
					Get(gotrackUrl + "/data/geo")
				if err != nil {
					log.Warn("error retrieving data", "err", err)
					return
				}
				var current dostowData
				err = json.Unmarshal([]byte(resp.String()), &current)
				if err != nil {
					fmt.Println(err)
					break
				}
				_last := current.Data[len(current.Data)-1].(map[string]interface{})["id"].(string)
				if _last == last {
					break
				}
				last = _last
				cur := len(current.Data)
				count = count + cur
				if dry {
					// explore response object
					fmt.Printf("\nError: %v", err)
					fmt.Printf("\nResponse Status Code: %v", resp.StatusCode())
					fmt.Printf("\nResponse Status: %v", resp.Status())
					fmt.Printf("\nResponse Time: %v", resp.Time())
					fmt.Printf("\nResponse Recevied At: %v", resp.ReceivedAt())
					// fmt.Printf("\nResponse Body: %v", resp) // or resp.String() or string(resp.Body())

					fmt.Printf("\nResponse rows count: %v", cur)
					fmt.Printf("\nResponse rows retrieved: %v/%v", count, current.TotalCount)

				} else {
					for _, v := range current.Data {
						data, _ := json.Marshal(v)
						fmt.Println(string(data))
					}
				}
				pages--
			}
		},
	}
)

func init() {
	exportCmd.AddCommand(dostowCmd)
	RootCmd.AddCommand(exportCmd)
	dostowCmd.Flags().IntVarP(&pageCount, "page-count", "c", 20, "no of entries per page")
	dostowCmd.Flags().IntVarP(&pages, "pages", "p", 10, "no of pages of data")
	dostowCmd.Flags().StringVarP(&gotrackUrl, "device-url", "e", "http://localhost:8000", "device url")
	dostowCmd.Flags().StringVarP(&dostowUrl, "url", "u", "http://localhost:3001/v1", "url to dostow api endpoint")
	dostowCmd.Flags().StringVarP(&dostowAccessKey, "access-key", "a", "", "group access key for dostow account")
	dostowCmd.Flags().StringSliceVarP(&dostowStores, "stores", "s", []string{"trips", "vehicle", "crums"}, "group access key for dostow account")
	dostowCmd.Flags().BoolVarP(&dry, "dry", "d", true, "dry run dont save")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
