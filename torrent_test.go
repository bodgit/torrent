package torrent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalBinary(t *testing.T) {
	_, err := New()
	assert.Nil(t, err)
}

func TestUnmarshalBinary(t *testing.T) {
	torrent, err := New()
	assert.Nil(t, err)

	err = torrent.UnmarshalBinary([]byte("d8:announce19:http://example.com/4:infod5:filesld6:lengthi123e4:pathl4:testeee4:name4:test12:piece lengthi123e6:pieces20:01234567890123456789ee"))
	assert.Nil(t, err)
	assert.Equal(t, "http://example.com/", torrent.Announce)
	assert.Equal(t, 1, len(torrent.Info.Files))
	assert.Equal(t, int64(123), torrent.Info.Files[0].Length)
	assert.Equal(t, []string{"test"}, torrent.Info.Files[0].Path)
	assert.Equal(t, int64(0), torrent.Info.Length)
	assert.Equal(t, "test", torrent.Info.Name)
	assert.Equal(t, int64(123), torrent.Info.PieceLength)
	assert.Equal(t, []byte("01234567890123456789"), torrent.Info.Pieces)
}
