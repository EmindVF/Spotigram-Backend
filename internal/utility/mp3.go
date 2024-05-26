package utility

import (
	"bytes"
	"strconv"
	"strings"

	dhowden "github.com/dhowden/tag"
)

func IsValidMP3(file []byte) bool {
	//return strings.HasPrefix(string(file), "ID3")
	return len(file) > 2 && file[0] == 'I' && file[1] == 'D' && file[2] == '3'
}

func GetSongLengthFromM3U8(file []byte) (int, error) {
	lines := strings.Split(string(file), "\n")
	var length float64 = 0
	for i, line := range lines {
		if strings.HasPrefix(line, "#EXTINF:") {
			partLength, err := strconv.ParseFloat(strings.Split(lines[i][8:], ",")[0], 64)
			if err != nil {
				return 0, err
			}
			length += partLength
		}
	}
	return int(length), nil
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
