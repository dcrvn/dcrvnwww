package models

import (
	"bytes"
	"encoding/gob"
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
)

type (
	gobCodec int
	Pagination struct {
		SortKey 	string		`json:"-"`
		SortVal		int			`json:"-"`
		Page 		int
		Total		int
		Records		int
	}
)

func (c gobCodec) Marshal(v interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (c gobCodec) Unmarshal(b []byte, v interface{}) error {
	r := bytes.NewReader(b)
	dec := gob.NewDecoder(r)
	return dec.Decode(v)
}

func (c gobCodec) Name() string {
	return "gob"
}

var (
	db *storm.DB
)

func Init() error {
	var err error
	Gob := new(gobCodec)
	db, err = storm.Open("dcrvnwww.db",storm.Codec(Gob), storm.Batch())
	return err
}

func (s *Pagination) Paginate() (limit int,skip int) {
	if s.Page < 1 {
		s.Page = 1
	}
	if s.Records == 0 {
		s.Records = 20
	}
	limit = s.Records
	skip = (s.Page - 1) * limit
	return
}