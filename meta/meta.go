package meta

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type Meta struct {
	Ver              string `json:"ver"`
	UUID             string `json:"uuid"`
	Type             string `json:"type"`
	WrapMode         string `json:"wrapMode"`
	FilterMode       string `json:"filterMode"`
	PremultiplyAlpha bool   `json:"premultiplyAlpha"`
	RawTextureUUID   string `json:"rawTextureUuid"`
	TrimType         string `json:"trimType"`
	TrimThreshold    int    `json:"trimThreshold"`
	Rotated          bool   `json:"rotated"`
	OffsetX          int    `json:"offsetX"`
	OffsetY          int    `json:"offsetY"`
	TrimX            int    `json:"trimX"`
	TrimY            int    `json:"trimY"`
	Width            int    `json:"width"`
	Height           int    `json:"height"`
	RawWidth         int    `json:"rawWidth"`
	RawHeight        int    `json:"rawHeight"`
	BorderTop        int    `json:"borderTop"`
	BorderBottom     int    `json:"borderBottom"`
	BorderLeft       int    `json:"borderLeft"`
	BorderRight      int    `json:"borderRight"`

	SubMetas map[string]*Meta `json:"subMetas"`
}

// ParseData 解析字节数组形式的 meta 数据
func ParseData(dat []byte) (*Meta, error) {
	var m Meta
	if err := json.Unmarshal(dat, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// Parse 从 io.Reader 中解析 meta 数据
func Parse(r io.Reader) (*Meta, error) {
	dat, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return ParseData(dat)
}

// ParseFile 从给定文件中解析 meta 数据
func ParseFile(filename string) (*Meta, error) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParseData(dat)
}
