/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at


    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"log"

	"encoding/json"

	"github.com/spf13/cobra"
)

var cutHead uint32

// lsnodeCmd represents the lsnode command
var lsnodeCmd = &cobra.Command{
	Use:   "lsnode",
	Short: "List nodes of prefab files",
	Long:  `List nodes of prefab files.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("lsnode called")
		for _, filename := range args {
			var nodes []map[string]interface{}
			log.Println("got: ", filename)
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Fatalln(err)
			}
			if err := json.Unmarshal(data, &nodes); err != nil {
				log.Fatalln(err)
			}
			for _, node := range nodes {
				for k, v := range node {
					log.Println(k, v)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(lsnodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsnodeCmd.PersistentFlags().String("foo", "", "A help for foo")
	lsnodeCmd.PersistentFlags().Uint32VarP(&cutHead, "cuthead", "", 1, "Remove N level from root")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsnodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
