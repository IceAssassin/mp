package main

import (
	"fmt"
	"github.com/ugorji/go/codec"
	"encoding/json"
	"github.com/google/flatbuffers/go"
	"mp/fb/proto/fbpkg"
)

type Score struct {
	A int
	B int
}

type Payload struct {
	Uid string
	Name string
	Gender int
	S Score
}

func Marshal(data interface{}, id uint32) ([]byte, error) {
	bts := make([]byte, 0, 64)
	h := new(codec.MsgpackHandle)
	enc := codec.NewEncoderBytes(&bts, h)
	if err := enc.Encode(data); err != nil {
		fmt.Println("encode error", err)
		return nil, err
	}

	fb := flatbuffers.NewBuilder(1024)
	//fb.Reset()

	boft := fb.CreateByteVector(bts)

	fbpkg.FBpkgStart(fb)


	fbpkg.FBpkgAddPayload(fb, boft)

	fbpkg.FBpkgAddId(fb, id)

	finish := fbpkg.FBpkgEnd(fb)
	fb.Finish(finish)

	return fb.Bytes[fb.Head():], nil
}

func UnMarshal(data []byte) error {
	fb := fbpkg.GetRootAsFBpkg(data, 0)

	fmt.Println("id = ", fb.Id())

	bts := fb.PayloadBytes()

	ddata :=  new(Payload)
	h := new(codec.MsgpackHandle)
	dec := codec.NewDecoderBytes(bts, h)
	if err := dec.Decode(ddata); err != nil {
		fmt.Println("decode error ", err)
		return err
	}


	fmt.Println(ddata.Name, ddata.S.A)
	return nil
}

func main () {
	data := Payload{
		Uid: "123",
		Name:"test",
		Gender:1,
		S: Score {
			A:10,
			B:100,
		},
	}


	jbts, err := json.Marshal(data)
	if err != nil {
		fmt.Println("json marsha error ", err)
		return
	}

	fmt.Println("jbts = ", string(jbts), "len = ", len(jbts))


	bts, err := Marshal(data, 101)
	if err != nil {
		return
	}

	UnMarshal(bts)


	return
}
