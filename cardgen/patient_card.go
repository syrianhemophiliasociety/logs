package cardgen

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"shs/app/models"

	"github.com/01walid/goarabic"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const (
	cardWidth  = 1062
	cardHeight = 590
)

var (
	//go:embed IBMPlexSansArabic-Regular.ttf
	fontBytes []byte
	//go:embed IBMPlexSansArabic-Bold.ttf
	boldFontBytes []byte
)

type Buffer struct {
	*bytes.Buffer
}

func NewBuffer(buf []byte) *Buffer {
	return &Buffer{
		Buffer: bytes.NewBuffer(buf),
	}
}

func (mwc *Buffer) Close() error {
	// Noop
	return nil
}

type PatientCardGenerator struct {
	writer  io.WriteCloser
	patient models.Patient
	ttf     *opentype.Font
	boldttf *opentype.Font

	baseImage   draw.Image
	lastDrawnAt image.Point
}

func New(writer io.WriteCloser, patient models.Patient) (*PatientCardGenerator, error) {
	patient.FirstName = fixArabicText(patient.FirstName)
	patient.LastName = fixArabicText(patient.LastName)
	patient.Nationality = fixArabicText(patient.Nationality)

	p := &PatientCardGenerator{
		writer:  writer,
		patient: patient,
	}

	ttf, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}
	p.ttf = ttf

	boldttf, err := opentype.Parse(boldFontBytes)
	if err != nil {
		return nil, err
	}
	p.boldttf = boldttf

	p.baseImage = image.NewRGBA(image.Rect(0, 0, cardWidth, cardHeight))
	for x := range cardWidth {
		for y := range cardHeight {
			p.baseImage.Set(x, y, image.White)
		}
	}

	p.lastDrawnAt = image.Point{X: 0, Y: 0}

	return p, nil
}

type qrVersionData struct {
	maxLength    int
	totalModules int
}

func calculateQRModuleSize(urlLength int, maxDimension int) uint8 {
	qrMap := []qrVersionData{
		{20, 21},
		{40, 25},
		{70, 29},
		{100, 33},
		{140, 37},
		{180, 41},
		{220, 45},
	}

	totalModules := 0

	for _, data := range qrMap {
		if urlLength <= data.maxLength {
			totalModules = data.totalModules
			break
		}
	}

	if totalModules == 0 {
		return 1
	}

	return uint8(maxDimension / totalModules)
}

func (p *PatientCardGenerator) generateQrCode(text string, rtl bool) error {
	qrc, err := qrcode.New(text)
	if err != nil {
		return err
	}

	qrBuf := NewBuffer(nil)
	qrWriter := standard.NewWithWriter(qrBuf,
		standard.WithBorderWidth(0),
		standard.WithQRWidth(calculateQRModuleSize(len(text), 400)),
		standard.WithBuiltinImageEncoder(standard.PNG_FORMAT),
		standard.WithBgTransparent(),
	)
	if err = qrc.Save(qrWriter); err != nil {
		return err
	}

	qrImage, err := png.Decode(qrBuf)
	if err != nil {
		return err
	}

	qrImageX := 45
	if rtl {
		qrImageX = p.baseImage.Bounds().Max.X - 45 - qrImage.Bounds().Max.X
	}
	qrImageY := (p.baseImage.Bounds().Max.Y - qrImage.Bounds().Max.Y) / 2
	if qrImageY < 0 {
		qrImageY *= -1
	}
	draw.Draw(p.baseImage, p.baseImage.Bounds(), qrImage, image.Point{X: -1 * qrImageX, Y: -1 * qrImageY}, draw.Over)

	p.lastDrawnAt = image.Point{
		X: qrImageX + qrImage.Bounds().Max.X,
		Y: qrImageY + qrImage.Bounds().Max.Y,
	}
	if rtl {
		p.lastDrawnAt.X = qrImageX
	}

	return nil
}

func (p *PatientCardGenerator) drawText(text string, bold bool, drawAt image.Point) error {
	ttf := p.ttf
	if bold {
		ttf = p.boldttf
	}

	face, err := opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    35,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return err
	}

	drawer := &font.Drawer{
		Dst:  p.baseImage,
		Src:  image.NewUniform(color.Black),
		Face: face,
		Dot: fixed.Point26_6{
			X: fixed.I(drawAt.X),
			Y: fixed.I(drawAt.Y),
		},
	}

	drawer.DrawString(text)

	newX := drawer.MeasureString(text)
	p.lastDrawnAt.X += newX.Round()

	return nil
}

func (p *PatientCardGenerator) Generate(rtl bool) error {
	if err := p.generateQrCode("https://logs.syrianhemophiliasociety.com/login?username="+p.patient.PublicId, rtl); err != nil {
		return err
	}

	rtlMultiplier := 1
	if rtl {
		rtlMultiplier *= -1
	}

	p.lastDrawnAt.X += 35 * rtlMultiplier
	p.lastDrawnAt.Y = 250

	oldX := p.lastDrawnAt.X
	if err := p.drawText("Patient ID: ", true, p.lastDrawnAt); err != nil {
		return err
	}
	if err := p.drawText(p.patient.PublicId, false, p.lastDrawnAt); err != nil {
		return err
	}

	p.lastDrawnAt.Y += 50
	p.lastDrawnAt.X = oldX

	if err := p.drawText("First name: ", true, p.lastDrawnAt); err != nil {
		return err
	}
	if err := p.drawText(p.patient.FirstName, false, p.lastDrawnAt); err != nil {
		return err
	}

	p.lastDrawnAt.Y += 50
	p.lastDrawnAt.X = oldX

	if err := p.drawText("Last name: ", true, p.lastDrawnAt); err != nil {
		return err
	}
	if err := p.drawText(p.patient.LastName, false, p.lastDrawnAt); err != nil {
		return err
	}

	p.lastDrawnAt.Y += 50
	p.lastDrawnAt.X = oldX

	if err := p.drawText("Nationality: ", true, p.lastDrawnAt); err != nil {
		return err
	}
	if err := p.drawText(p.patient.Nationality, false, p.lastDrawnAt); err != nil {
		return err
	}

	return nil
}

func (p *PatientCardGenerator) Finalize() error {
	if err := png.Encode(p.writer, p.baseImage); err != nil {
		return err
	}

	p.writer.Close()
	return nil
}

func isArabic(text string) bool {
	for _, chr := range text {
		if chr >= 0x600 && chr <= 0x6FF {
			return true
		}
	}

	return false
}

func fixArabicText(text string) string {
	if isArabic(text) {
		return goarabic.Reverse(goarabic.ToGlyph(text))
	}

	return text
}
