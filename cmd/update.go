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
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"

	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"github.com/spf13/cobra"
)

var (
	updateBin, deviceKey, deviceUsername, deviceUri string
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the gotrack binary",
	Long:  `Update the gotrack binary.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Use SSH key authentication from the auth package
		authorizedKeysBytes, err := ioutil.ReadFile("~/.ssh/known_hosts")
		if err != nil {
			// log.Fatalf("Failed to load authorized_keys, err: %v", err)
		} else {
			authorizedKeysMap := map[string]bool{}
			for len(authorizedKeysBytes) > 0 {
				pubKey, _, _, rest, err := ssh.ParseAuthorizedKey(authorizedKeysBytes)
				if err != nil {
					log.Fatal(err)
				}

				authorizedKeysMap[string(pubKey.Marshal())] = true
				authorizedKeysBytes = rest
			}
		}
		clientConfig, _ := auth.PrivateKey(deviceUsername, deviceKey)
		clientConfig.HostKeyCallback = func(hostname string, remote net.Addr, pubKey ssh.PublicKey) error {
			// if authorizedKeysMap[string(pubKey.Marshal())] {
			// 	return &ssh.Permissions{
			// 		// Record the public key used for authentication.
			// 		Extensions: map[string]string{
			// 			"pubkey-fp": ssh.FingerprintSHA256(pubKey),
			// 		},
			// 	}
			// }
			return nil
		}
		clientConfig.Auth = append(clientConfig.Auth, ssh.Password("password"))
		// For other authentication methods see ssh.ClientConfig and ssh.AuthMethod

		// Create a new SCP client
		client := scp.NewClient(deviceUri, &clientConfig)

		// Connect to the remote server
		err = client.Connect()
		if err != nil {
			fmt.Println("Couldn't establish a connection to the remote server ", err)
			return
		}

		// Open a file
		f, _ := os.Open(updateBin)

		// Close session after the file has been copied
		defer client.Session.Close()

		// Close the file after it has been copied
		defer f.Close()

		// Finaly, copy the file over
		// Usage: CopyFile(fileReader, remotePath, permission)

		client.CopyFile(f, "~/gotrack_linux_arm", "0655")
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	updateCmd.Flags().StringVarP(&deviceUsername, "username", "u", "pi", "username")
	updateCmd.Flags().StringVarP(&deviceKey, "key", "k", "~/.ssh/id_rsa", "device ssh key")
	updateCmd.Flags().StringVarP(&deviceUri, "uri", "r", "192.168.1.8:22", "device ssh uri")
	updateCmd.Flags().StringVarP(&updateBin, "bin", "b", "./gotrack_linux_arm", "binary file to use to update")
}
