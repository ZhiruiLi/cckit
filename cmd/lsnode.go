/*
Copyright © 2020 zhiruili zr.public@outlook.com

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
	"strings"

	"encoding/json"

	"github.com/spf13/cobra"
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

type prefabNode struct {
	Type     string `json:"__type__"`
	Name     string `json:"_name"`
	ObjFlags int    `json:"_objFlags"`
	Enabled  bool   `json:"_enabled"`
	Parent   struct {
		ID int `json:"__id__"`
	} `json:"_parent"`
	Children []struct {
		ID int `json:"__id__"`
	} `json:"_children"`
	Root struct {
		ID int `json:"__id__"`
	} `json:"root"`
	Active     bool          `json:"_active"`
	Level      int           `json:"_level"`
	Components []interface{} `json:"_components"`
	Prefab     struct {
		ID int `json:"__id__"`
	} `json:"_prefab"`
	Opacity int `json:"_opacity"`
	Color   struct {
		Type string `json:"__type__"`
		R    int    `json:"r"`
		G    int    `json:"g"`
		B    int    `json:"b"`
		A    int    `json:"a"`
	} `json:"_color"`
	ContentSize struct {
		Type   string  `json:"__type__"`
		Width  float64 `json:"width"`
		Height float64 `json:"height"`
	} `json:"_contentSize"`
	AnchorPoint struct {
		Type string  `json:"__type__"`
		X    float64 `json:"x"`
		Y    float64 `json:"y"`
	} `json:"_anchorPoint"`
	Position struct {
		Type string  `json:"__type__"`
		X    float64 `json:"x"`
		Y    float64 `json:"y"`
		Z    float64 `json:"z"`
	} `json:"_position"`
	Scale struct {
		Type string  `json:"__type__"`
		X    float64 `json:"x"`
		Y    float64 `json:"y"`
		Z    float64 `json:"z"`
	} `json:"_scale"`
	RotationX float64 `json:"_rotationX"`
	RotationY float64 `json:"_rotationY"`
	Quat      struct {
		Type string  `json:"__type__"`
		X    float64 `json:"x"`
		Y    float64 `json:"y"`
		Z    float64 `json:"z"`
		W    float64 `json:"w"`
	} `json:"_quat"`
	SkewX      float64 `json:"_skewX"`
	SkewY      float64 `json:"_skewY"`
	GroupIndex int     `json:"groupIndex"`
	ID         string  `json:"_id"`

	ccNodeChilden []*prefabNode
}

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

func findRoot(nodes []*prefabNode) *prefabNode {
	for i := range nodes {
		if nodes[i].Parent.ID == 0 && nodes[i].Type == "cc.Node" {
			return nodes[i]
		}
	}
	return nil
}

func findRootID(nodes []*prefabNode) int {
	for _, n := range nodes {
		if n.Root.ID != 0 {
			return n.Root.ID
		}
	}
	return 0
}

func listChildrenNames(myName string, children []*prefabNode) [][]string {
	if len(children) == 0 {
		return [][]string{{myName}}
	}
	var result [][]string
	for _, c := range children {
		sub := listChildrenNames(c.Name, c.ccNodeChilden)
		result = append(result, sub...)
	}
	for i := range result {
		result[i] = append(result[i], myName)
	}
	return result
}

func listChildrenNamesFromRoot(root *prefabNode) [][]string {
	reversed := listChildrenNames(root.Name, root.ccNodeChilden)
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

func makeIDMap(nodes []*prefabNode) map[int]*prefabNode {
	parentMap := make(map[int]*prefabNode)
	for i, node := range nodes {
		parentMap[i] = node
	}
	return parentMap
}

func insertChildren(nodes []*prefabNode, parentMap map[int]*prefabNode) {
	for _, node := range nodes {
		if node.Parent.ID != 0 {
			p, ok := parentMap[node.Parent.ID]
			if !ok {
				log.Printf("parent:%d no found for node %s:%s", node.Parent.ID, node.Type, node.Name)
				continue
			}
			p.ccNodeChilden = append(p.ccNodeChilden, node)
		}
	}
}

var lsnodeCmd = &cobra.Command{
	Use:   "lsnode",
	Short: "List nodes of prefab files",
	Long:  `List nodes of prefab files.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, filename := range args {
			var nodes []*prefabNode
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Fatalln(err)
			}
			if err := json.Unmarshal(data, &nodes); err != nil {
				log.Fatalln(err)
			}

			r := findRoot(nodes)
			if r == nil {
				log.Fatalln("root node no found")
			}

			parentMap := makeIDMap(nodes)
			insertChildren(nodes, parentMap)

			nestedNameList := listChildrenNamesFromRoot(r)
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
	lsnodeCmd.PersistentFlags().UintVarP(&cutHead, "cuthead", "", 0, "Remove N nodes from root")
	lsnodeCmd.PersistentFlags().UintVarP(&maxDepth, "maxdepth", "", 0, "Maximum depth level from root")
	lsnodeCmd.PersistentFlags().StringVarP(&prefix, "prefix", "", "", "Prefix string of each node lists")
	lsnodeCmd.PersistentFlags().StringVarP(&suffix, "suffix", "", "", "Suffix string of each node lists")
	lsnodeCmd.PersistentFlags().StringVarP(&sep, "sep", "", "", "Separator between each nodes")
	lsnodeCmd.PersistentFlags().StringVarP(&nPrefix, "nprefix", "", "", "Prefix string of each nodes in node lists")
	lsnodeCmd.PersistentFlags().StringVarP(&nSuffix, "nsuffix", "", "", "Suffix string of each nodes in node lists")
}
