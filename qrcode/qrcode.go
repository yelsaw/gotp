package qrcode

/*
#cgo pkg-config: zbar
#include <zbar.h>

zbar_image_t* go_zbar_scan_image(unsigned char *data, unsigned int width, unsigned int height) {
    zbar_image_scanner_t *scanner = zbar_image_scanner_create();
    zbar_image_scanner_set_config(scanner, 0, ZBAR_CFG_ENABLE, 1);

    zbar_image_t *image = zbar_image_create();
    zbar_image_set_format(image, *(unsigned int*)"Y800");
    zbar_image_set_size(image, width, height);
    zbar_image_set_data(image, data, 0, NULL);

    zbar_scan_image(scanner, image);
    zbar_image_scanner_destroy(scanner);
    return image;
}
// C will not compile, unless import is below cgo.
*/
import "C"

import (
	"image"
	"image/png"
	"os"
	// Used to create pointer in Go and pass to C
	"unsafe"
)

// Reads an image file by path and scans the image for QR codes using
// the ZBar library with compiled C code using cgo.
// In order to use this package the host machine must have libzbar-dev installed.
func ImagePath(path string) (string, error) {
	qrCode, err := os.Open(path)
	if err != nil {
		return "Unable to open image:", err
	}
	defer qrCode.Close()

	img, err := png.Decode(qrCode)
	if err != nil {
		return "Unable to decode image:", err
	}

	// Convert image to gray scale
	grayMatter, width, height := func(img image.Image) ([]byte, int, int) {
		bounds := img.Bounds()
		width, height := bounds.Dx(), bounds.Dy()
		grayScale := make([]byte, width * height)
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := img.At(x, y).RGBA()
				grayScale[y*width+x] = uint8((r + g + b) / 3 >> 8)
			}
		}
		return grayScale, width, height
	}(img)

	// Scan image with zbar using go image ref
	cData := (*C.uchar)(unsafe.Pointer(&grayMatter[0]))
	image := C.go_zbar_scan_image(cData, C.uint(width), C.uint(height))

	firstSymbol := C.zbar_image_first_symbol(image)
	for firstSymbol != nil {
		code := C.GoString(C.zbar_symbol_get_data(firstSymbol))
		return code, nil
		firstSymbol = C.zbar_symbol_next(firstSymbol)
	}

	// zbar recommends using this function--who am I to question zbar?
	// https://zbar.sourceforge.net/api/zbar_8h.html#dfa426908bd8221a59c7c333fe30d54b
	C.zbar_image_destroy(image)
	return "", nil
}
