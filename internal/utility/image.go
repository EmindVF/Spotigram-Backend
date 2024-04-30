package utility

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
)

// Converts a png or a jpg image to webp format.
// new_width and new_height must be lower than 4096.
func ConvertAndResizeImageToWebP(input []byte, new_width uint, new_height uint) ([]byte, error) {

	if new_width > 4096 || new_height > 4096 {
		return nil,
			fmt.Errorf("invalid image size, dimensions must be lower than 4096")
	}

	_, format, err := image.DecodeConfig(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}

	var im image.Image
	if format == "jpeg" {
		im, err = jpeg.Decode(bytes.NewReader(input))
		if err != nil {
			return nil, err
		}
	} else if format == "png" {
		im, err = png.Decode(bytes.NewReader(input))
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("invalid image format")
	}

	resizedImg := resize.Resize(new_width, new_height, im, resize.Lanczos3)

	var output bytes.Buffer
	err = webp.Encode(&output, resizedImg, &webp.Options{Lossless: false, Quality: 50})
	if err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}
