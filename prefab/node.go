// Copyright © 2020 zhiruili zr.public@outlook.com
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

// Package prefab 实现了 prefab 的相关程序内表示和 prefab 数据解析
package prefab

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

// Prefab 代表一个 Cocos Creator prefab 文件
type Prefab struct {
	// Root 是一个 prefab 的根节点
	Root *Node
	// Nodes 表示一个 prefab 中所有的节点
	Nodes []*Node
}

// Node 表示 prefab 中的一个节点
type Node struct {
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
	Active     bool `json:"_active"`
	Level      int  `json:"_level"`
	Components []struct {
		ID int `json:"__id__"`
	} `json:"_components"`
	Prefab struct {
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

	LocalContext interface{} `json:"-"`
	CCParent     *Node       `json:"-"`
	CCChildren   []*Node     `json:"-"`
	CCComponents []*Node     `json:"-"`
}

// ParseData 解析字节数组形式的 prefab 数据
func ParseData(dat []byte) (*Prefab, error) {
	var nodes []*Node
	if err := json.Unmarshal(dat, &nodes); err != nil {
		return nil, err
	}
	if err := insertCCData(nodes); err != nil {
		return nil, err
	}
	var r *Node
	if len(nodes) > 1 {
		r = nodes[1]
	}
	return &Prefab{r, nodes}, nil
}

// Parse 从 io.Reader 中解析 prefab 数据
func Parse(r io.Reader) (*Prefab, error) {
	dat, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return ParseData(dat)
}

// ParseFile 从给定文件中解析 prefab 数据
func ParseFile(filename string) (*Prefab, error) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParseData(dat)
}

func inRange(idx int, nodes []*Node) bool {
	return idx >= 0 && idx < len(nodes)
}

func insertCCData(nodes []*Node) error {
	for _, node := range nodes {
		if node.Parent.ID != 0 {
			if !inRange(node.Parent.ID, nodes) {
				return fmt.Errorf("parent:%d no found for node %s:%s", node.Parent.ID, node.Type, node.Name)
			} else {
				node.CCParent = nodes[node.Parent.ID]
			}
		}
		for _, ch := range node.Children {
			if !inRange(ch.ID, nodes) {
				return fmt.Errorf("child:%d no found for node %s:%s", ch.ID, node.Type, node.Name)
			}
			node.CCChildren = append(node.CCChildren, nodes[ch.ID])
		}
		for _, co := range node.Components {
			if !inRange(co.ID, nodes) {
				return fmt.Errorf("component:%d no found for node %s:%s", co.ID, node.Type, node.Name)
			}
			node.CCComponents = append(node.CCComponents, nodes[co.ID])
		}
	}
	return nil
}
