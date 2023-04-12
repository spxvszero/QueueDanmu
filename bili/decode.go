package bili

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"io"
)

type IDecode interface {
	Decode([]byte) ([][]byte, error)
}

type DecodeV0 struct {
	IDecode
}

// Decode v0没有压缩
func (d *DecodeV0) Decode(input []byte) (result [][]byte, err error) {
	return [][]byte{input}, nil
}

type DecodeV2 struct {
	IDecode
}

// Decode v2zip压缩
func (d *DecodeV2) Decode(input []byte) (result [][]byte, err error) {
	result = [][]byte{}
	b := bytes.NewReader(input)
	r, err := zlib.NewReader(b)
	if err != nil {
		err = errors.Wrapf(err, "[DecodeV2 | Decode] unzip err")
		return
	}
	var bodyBuf bytes.Buffer
	_, err = io.Copy(&bodyBuf, r)
	if err != nil {
		err = errors.Wrapf(err, "[DecodeV2 | Decode] io.Copy err")
		return
	}
	bodyLen := int32(bodyBuf.Len())
	var offset int32
	// ｜int32(一个cmd的长度+4位)｜cmd内容｜int32(一个cmd的长度+4位)｜cmd内容｜int32(一个cmd的长度+4位)｜cmd内容｜....
	for offset < bodyLen {
		cmdSize := int32(binary.BigEndian.Uint32(bodyBuf.Bytes()[offset : offset+CmdSize]))
		// 协议长度大于body长度
		if offset+cmdSize > bodyLen {
			err = fmt.Errorf("[DecodeV2 | Decode] offset:%d + cmdSize:%d > bodyLen:%d", offset, cmdSize, bodyLen)
			return
		}
		cmd := bodyBuf.Bytes()[offset+CmdSize : offset+cmdSize]
		result = append(result, cmd)
		offset += cmdSize
	}
	return
}

type DecodeManager struct {
	m map[int64]IDecode
}

func NewDecodeManager() (manager *DecodeManager) {
	return &DecodeManager{
		m: map[int64]IDecode{
			ProtoVersion0: &DecodeV0{},
			ProtoVersion2: &DecodeV2{},
		},
	}
}

func (m *DecodeManager) Decode(version int64, input []byte) (result [][]byte, err error) {
	d, exist := m.m[version]
	if !exist {
		err = errors.Wrapf(err, "[DecodeManager | Decode] version not found")
		return
	}
	result, err = d.Decode(input)
	if err != nil {
		err = errors.Wrapf(err, "[DecodeManager | Decode] Decode err")
		return
	}
	return
}
