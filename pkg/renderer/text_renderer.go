package renderer

import (
    "image"
    "image/color"
    "image/jpeg"
    "os"
    "golang.org/x/image/font"
    "golang.org/x/image/math/fixed"
    "golang.org/x/image/font/opentype"
    "io/ioutil"
    
)

type TextRenderer struct{}

func (tr *TextRenderer) RenderText(img *image.RGBA, text string, fontPath string, fontSize float64) error {
    col := color.RGBA{255, 255, 255, 255} // White color for text

    // Load the font file
    fontBytes, err := ioutil.ReadFile(fontPath)
    if err != nil {
        return err
    }

    parsedFont, err := opentype.Parse(fontBytes)
    if err != nil {
        return err
    }

    face, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{
        Size: fontSize,
        DPI:  72,
    })
    if err != nil {
        return err
    }
    defer face.Close()

    // Calculate the center position for the text
    imgWidth := img.Bounds().Dx()
    imgHeight := img.Bounds().Dy()
    textWidth := len(text) * int(fontSize/2) // Approximate width of each character
    textHeight := int(fontSize)             // Height of the font

    centerX := (imgWidth - textWidth) / 2
    centerY := (imgHeight + textHeight) / 2

    point := fixed.Point26_6{
        X: fixed.I(centerX),
        Y: fixed.I(centerY),
    }

    // Draw the text on the image
    d := &font.Drawer{
        Dst:  img,
        Src:  image.NewUniform(col),
        Face: face,
        Dot:  point,
    }
    d.DrawString(text)
    return nil
}

func SaveImage(img *image.RGBA, filename string) error {
    outFile, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer outFile.Close()

    return jpeg.Encode(outFile, img, nil)
}