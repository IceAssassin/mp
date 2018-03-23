package main

import (
	"github.com/ugorji/go/codec"
	"fmt"
	"encoding/json"
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

	bts := make([]byte, 0, 64)
	h := new(codec.MsgpackHandle)
	enc := codec.NewEncoderBytes(&bts, h)
	if err := enc.Encode(data); err != nil {
		fmt.Println("encode error", err)
		return
	}

	jbts, err := json.Marshal(data)
	if err != nil {
		fmt.Println("json marsha error ", err)
		return
	}

	fmt.Println("jbts = ", string(jbts), "len = ", len(jbts))
	fmt.Println("bts = ", string(bts), "len = ", len(bts))

	ddata :=  new(Payload)
	dec := codec.NewDecoderBytes(bts, h)
	if err := dec.Decode(ddata); err != nil {
		fmt.Println("decode error ", err)
		return
	}

	fmt.Println(ddata.Name, ddata.S.A)

	return
}
