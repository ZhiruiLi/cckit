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

// Package cmd 实现了 cckit 的命令行参数
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/zhiruili/cckit/meta"
	"github.com/zhiruili/cckit/prefab"
	"github.com/zhiruili/cckit/slices"
)

var scopes []string

type location struct {
	filename string
	node     *prefab.Node
	prefab   *prefab.Prefab
}

type target struct {
	filename string
	uuids    []string
	meta     *meta.Meta
}

func allUUIDs(m *meta.Meta, uuids []string) []string {
	if m == nil {
		return uuids
	}
	if m.UUID != "" {
		if _, ok := slices.FindString(uuids, m.UUID); !ok {
			uuids = append(uuids, m.UUID)
		}
	}
	for _, sub := range m.SubMetas {
		uuids = allUUIDs(sub, uuids)
	}
	return uuids
}

type info struct {
	tar *target
	loc *location
}

func findInPrefab(tars []*target, path string, pf *prefab.Prefab) ([]*info, error) {
	var is []*info
	for _, node := range pf.Nodes {
		for _, tar := range tars {
			for _, uuid := range tar.uuids {
				if uuid == node.SpriteFrame.UUID {
					is = append(is, &info{tar, &location{path, node, pf}})
					break
				}
			}
		}
	}
	return is, nil
}

func findInDir(tars []*target, path string) ([]*info, error) {
	var is []*info
	fs, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		return nil, err
	}
	for _, f := range fs {
		tmp, err := findInFile(tars, f)
		if err != nil {
			return nil, err
		}
		is = append(is, tmp...)
	}
	return is, nil
}

func findInFile(tars []*target, path string) ([]*info, error) {
	st, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if st.IsDir() {
		return findInDir(tars, path)
	}
	if filepath.Ext(path) != ".prefab" {
		return nil, nil
	}
	pf, err := prefab.ParseFile(path)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", path, err)
	}
	is, err := findInPrefab(tars, path, pf)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", path, err)
	}
	return is, nil
}

func findAll(tars []*target, filenames []string) ([]*info, error) {
	var is []*info
	for _, f1 := range filenames {
		subs, err := filepath.Glob(f1)
		if err != nil {
			return nil, err
		}
		for _, f2 := range subs {
			tmp, err := findInFile(tars, f2)
			if err != nil {
				return nil, err
			}
			is = append(is, tmp...)
		}
	}
	return is, nil
}

func fmtLocation(l *location) string {
	node := l.node
	name := ""
	for node != nil {
		name = node.Name
		if name != "" {
			break
		}
		if node.CCNode != nil {
			node = node.CCNode
		} else if node.CCParent != nil {
			node = node.CCParent
		} else {
			break
		}
	}
	return fmt.Sprintf("%s: %s", filepath.Base(l.filename), name)
}

func fmtInfo(i *info, onlyOneTarget bool) string {
	ls := fmtLocation(i.loc)
	if onlyOneTarget {
		return ls
	}
	return i.tar.filename + ": " + ls
}

func loadMeta(filename string) (*target, error) {
	m, err := meta.ParseFile(filename)
	if err != nil {
		return nil, err
	}
	return &target{
		filename: filename,
		uuids:    allUUIDs(m, nil),
		meta:     m,
	}, nil
}

func loadTarget(filename string) (*target, error) {
	if filepath.Ext(filename) != ".meta" {
		filename = filename + ".meta"
	}
	return loadMeta(filename)
}

func loadTargets(filenames []string) ([]*target, error) {
	var tars []*target
	for _, f := range filenames {
		t, err := loadTarget(f)
		if err != nil {
			return nil, fmt.Errorf("%s:%w", f, err)
		}
		tars = append(tars, t)
	}
	return tars, nil
}

// findrefCmd represents the findref command
var findrefCmd = &cobra.Command{
	Use:   "findref",
	Short: "Find references of the given resource in the given prefabs",
	Long:  `Find references of the given resource in the given prefabs.`,
	Run: func(cmd *cobra.Command, args []string) {
		tars, err := loadTargets(args)
		if err != nil {
			log.Fatal(err)
		}
		if len(tars) <= 0 {
			return
		}
		is, err := findAll(tars, scopes)
		if err != nil {
			log.Fatalln(err)
		}
		for _, i := range is {
			s := fmtInfo(i, len(tars) <= 1)
			fmt.Println(s)
		}
	},
}

func init() {
	rootCmd.AddCommand(findrefCmd)
	findrefCmd.PersistentFlags().StringArrayVarP(&scopes, "scope", "s", []string{"."}, "Search scopes")
}
