// Copyright © 2020 zhiruili zr.public@outlook.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package cmd 实现了 cckit 的命令行参数
package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zhiruili/cckit/prefab"
)

var (
	cutHead  uint = 0
	maxDepth uint = 0
)

var (
	prefix  = ""
	suffix  = ""
	sep     = ""
	nPrefix = ""
	nSuffix = ""
)

func cutNameList(nameList []string, cutHead, maxDepth uint) []string {
	beg := cutHead
	end := uint(len(nameList))
	if maxDepth > 0 && maxDepth <= end {
		end = maxDepth
	}
	if beg < end {
		return nameList[beg:end]
	}
	return nil
}

func listChildrenNames(myName string, children []*prefab.Node) [][]string {
	if len(children) == 0 {
		return [][]string{{myName}}
	}
	var result [][]string
	for _, c := range children {
		sub := listChildrenNames(c.Name, c.CCChildren)
		result = append(result, sub...)
	}
	for i := range result {
		result[i] = append(result[i], myName)
	}
	return result
}

func listChildrenNamesFromRoot(root *prefab.Node) [][]string {
	reversed := listChildrenNames(root.Name, root.CCChildren)
	for idx := range reversed {
		for i, j := 0, len(reversed[idx])-1; i < j; i, j = i+1, j-1 {
			reversed[idx][i], reversed[idx][j] = reversed[idx][j], reversed[idx][i]
		}
	}
	return reversed
}

type fmtOptions struct {
	prefix  string
	suffix  string
	sep     string
	nPrefix string
	nSuffix string
}

func fmtNameList(nameList []string, opts *fmtOptions) string {
	var newNameList []string
	for _, name := range nameList {
		newNameList = append(newNameList, opts.nPrefix+name+opts.nSuffix)
	}
	return opts.prefix + strings.Join(newNameList, opts.sep) + opts.suffix
}

var lsnodeCmd = &cobra.Command{
	Use:   "lsnode",
	Short: "List nodes of prefab files",
	Long:  `List nodes of prefab files.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, filename := range args {
			pf, err := prefab.ParseFile(filename)
			if err != nil {
				log.Fatalf("%s:%s", filename, err.Error())
			}
			if pf.Root == nil {
				log.Fatalf("%s:root node not found", filename)
			}
			nestedNameList := listChildrenNamesFromRoot(pf.Root)
			fmtOpt := fmtOptions{
				prefix:  prefix,
				suffix:  suffix,
				sep:     sep,
				nPrefix: nPrefix,
				nSuffix: nSuffix,
			}
			hasPrinted := make(map[string]bool)
			for _, nameList := range nestedNameList {
				newNameList := cutNameList(nameList, cutHead, maxDepth)
				if len(newNameList) > 0 {
					str := fmtNameList(newNameList, &fmtOpt)
					if printed := hasPrinted[str]; !printed {
						fmt.Println(str)
						hasPrinted[str] = true
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(lsnodeCmd)
	lsnodeCmd.PersistentFlags().UintVarP(&cutHead, "cuthead", "h", 0, "Remove N nodes from root")
	lsnodeCmd.PersistentFlags().UintVarP(&maxDepth, "maxdepth", "d", 0, "Maximum depth level from root")
	lsnodeCmd.PersistentFlags().StringVarP(&prefix, "prefix", "P", "", "Prefix string of each node lists")
	lsnodeCmd.PersistentFlags().StringVarP(&suffix, "suffix", "S", "", "Suffix string of each node lists")
	lsnodeCmd.PersistentFlags().StringVarP(&sep, "sep", "e", "", "Separator between each nodes")
	lsnodeCmd.PersistentFlags().StringVarP(&nPrefix, "nprefix", "p", "", "Prefix string of each nodes in node lists")
	lsnodeCmd.PersistentFlags().StringVarP(&nSuffix, "nsuffix", "s", "", "Suffix string of each nodes in node lists")
}
