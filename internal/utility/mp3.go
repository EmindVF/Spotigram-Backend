package utility

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	dhowden "github.com/dhowden/tag"
	tcolgate "github.com/tcolgate/mp3"
)

func IsValidMP3(file []byte) bool {
	return strings.HasPrefix(string(file), "ID3")
}

func GetMP3Length(file []byte) (int, error) {
	decoder := tcolgate.NewDecoder(bytes.NewReader(file))
	if decoder == nil {
		return 0, fmt.Errorf("cannot decode mp3 file")
	}
	var frame tcolgate.Frame
	var time float64 = 0
	skipped := 0
	for {
		if err := decoder.Decode(&frame, &skipped); err != nil {
			if err == io.EOF {
				break
			}
			return 0, err
		}
		time = time + frame.Duration().Seconds()
	}
	return int(time), nil
}

func GetMP3AlbumCover(file []byte) ([]byte, string, error) {
	tag, err := dhowden.ReadFrom(bytes.NewReader(file))
	if err != nil {
		return nil, "", err
	}

	pic := tag.Picture()
	if pic == nil {
		return nil, "", nil
	}
	return pic.Data, pic.Type, nil
}
