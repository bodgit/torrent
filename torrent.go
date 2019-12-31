package torrent

import (
	"bytes"

	"github.com/zeebo/bencode"
)

type Torrent struct {
	Announce string `bencode:"announce"`
	Info     struct {
		Files []struct {
			Length int64    `bencode:"length"`
			Path   []string `bencode:"path"`
		} `bencode:"files"`
		Length      int64  `bencode:"length"`
		Name        string `bencode:"name"`
		PieceLength int64  `bencode:"piece length"`
		Pieces      []byte `bencode:"pieces"`
	} `bencode:"info"`
}

func New() (*Torrent, error) {
	return &Torrent{}, nil
}

func (t *Torrent) MarshalBinary() ([]byte, error) {
	// TODO
	return nil, nil
}

func (t *Torrent) UnmarshalBinary(b []byte) error {
	r := bytes.NewReader(b)

	dec := bencode.NewDecoder(r)
	if err := dec.Decode(t); err != nil {
		return err
	}

	return nil
}
