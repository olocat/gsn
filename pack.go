package gsn

import (
	"encoding/json"
	"errors"
	"math"
)

type ReceivePack struct {
	UnPacker
	Ctx *Context
}

type SendPack struct {
	Packer
	Ctx *Context
}

func NewReceivePack(ctx *Context, dataStream []byte) *ReceivePack {
	return &ReceivePack{Ctx: ctx, UnPacker: *NewUnPacker(dataStream)}
}

func NewSendPack(ctx *Context) *SendPack {
	return &SendPack{Ctx: ctx, Packer: *NewPacker()}
}

func (r *ReceivePack) UnPackString() (string, error) {
	strLength, err := r.UnPackUint32()
	if err != nil {
		return "", err
	}

	return r.UnPacker.UnPackString(int(strLength))
}

func (r *ReceivePack) UnPackJsonData() (map[string]interface{}, error) {
	jsonLength, err := r.UnPackUint32()
	if err != nil {
		return nil, err
	}

	return r.UnPacker.UnPackJsonData(int(jsonLength))
}

func (s *SendPack) Send() error {
	dataStream := s.GetData()
	if len(dataStream) > math.MaxUint32 {
		return errors.New("pack data too long")
	}
	_, err := s.Ctx.Write(dataStream)
	return err
}

func (s *SendPack) PackString(str string) {
	strBytes := []byte(str)
	s.PackUint32(uint32(len(strBytes)))
	s.Packer.PackString(str)
}

func (s *SendPack) PackJsonData(jsonData map[string]interface{}) error {
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	s.PackInt32(int32(len(jsonBytes)))
	return s.Packer.PackJsonData(jsonData)
}
