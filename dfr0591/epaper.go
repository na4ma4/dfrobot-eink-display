package dfr0591

import (
	"errors"

	"github.com/na4ma4/dfrobot-eink-display/bitarray"
)

// Epaper is the epaper device and abstracted commands.
type Epaper struct {
	d             DisplayTransport
	w, h          int
	displayBuffer *bitarray.BitArray
}

// Mode is the update mode (Full or Part).
type Mode int

const (
	// ModeFull is the full update mode.
	ModeFull Mode = iota
	// ModePart is the partial update mode.
	ModePart
)

const (
	// Xmax is the max size of the device (X).
	Xmax = 250
	// Ymax is the max size of the device (Y).
	Ymax = 128
)

// NewEpaper returns a new Epaper struct.
func NewEpaper(d DisplayTransport, width, height int) *Epaper {
	ba, err := bitarray.NewBitArray(width * height)
	if err != nil {
		panic(err)
	}

	return &Epaper{
		d:             d,
		w:             width,
		h:             height,
		displayBuffer: ba,
	}
}

//nolint:gomnd // Lots of magic numbers, maybe I'll extract them later.
func (e *Epaper) _init() (err error) {
	if err = e.d.WriteCmdAndData(0x01, []byte{(Xmax - 1) % 256, (Xmax - 1) / 256, 0x00}); err != nil {
		return
	}

	if err = e.d.WriteCmdAndData(0x0c, []byte{0xd7, 0xd6, 0x9d}); err != nil {
		return
	}

	if err = e.d.WriteCmdAndData(0x2c, []byte{0xa8}); err != nil {
		return
	}

	if err = e.d.WriteCmdAndData(0x3a, []byte{0x1a}); err != nil {
		return
	}

	if err = e.d.WriteCmdAndData(0x3b, []byte{0x08}); err != nil {
		return
	}

	if err = e.d.WriteCmdAndData(0x11, []byte{0x01}); err != nil {
		return
	}

	if err = e.d.SetRAMData(0x00,
		uint8((Ymax-1)/8),
		uint8((Xmax-1)%256),
		uint8((Xmax-1)/256),
		0x00,
		0x00,
	); err != nil {
		return
	}

	return e.d.SetRAMPointer(0x00, uint8((Xmax-1)%256), uint8((Xmax-1)/256))
}

//nolint:gomnd // Lots of magic numbers, maybe I'll extract them later.
func (e *Epaper) _initLut(mode Mode) error {
	switch mode {
	case ModeFull:
		return e.d.WriteCmdAndData(0x32, []byte{
			0x22, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA,
			0x11, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x1E, 0x1E, 0x1E, 0x1E, 0x1E,
			0x1E, 0x1E, 0x1E, 0x01, 0x00, 0x00, 0x00,
			0x00,
		})
	case ModePart:
		return e.d.WriteCmdAndData(0x32, []byte{
			0x18, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x0F, 0x01, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00,
		})
	}

	return nil
}

// def _disPart(self, xStart, xEnd, yStart, yEnd):
// 	self._setRamData(xStart // 8, xEnd // 8, yEnd % 256, yEnd // 256, yStart % 256, yStart // 256)
// 	self._setRamPointer(xStart // 8, yEnd % 256, yEnd // 256)
// 	self._writeDisRam(xEnd - xStart, yEnd - yStart + 1)
// 	self._updateDis(self.PART)
//nolint:gomnd // Lots of magic numbers, maybe I'll extract them later.
func (e *Epaper) _disPart(xStart, xEnd, yStart, yEnd uint8) (err error) {
	err = e.d.SetRAMData(
		xStart/8,
		xEnd/8,
		uint8(int(yEnd)%256),
		uint8(int(yEnd)/256),
		uint8(int(yStart)%256),
		uint8(int(yStart)/256),
	)
	if err != nil {
		return err
	}

	err = e.d.SetRAMPointer(xStart/8, uint8(int(yEnd)%256), uint8(int(yEnd)/256))
	if err != nil {
		return err
	}

	err = e.d.WriteDisRAM(e.displayBuffer, xEnd-xStart, yEnd-yStart+1)
	if err != nil {
		return err
	}

	err = e.d.UpdateDis(ModePart)
	if err != nil {
		return err
	}

	return
}

// Flush the buffer to the display.
//nolint:gomnd // Lots of magic numbers, maybe I'll extract them later.
func (e *Epaper) Flush(mode Mode) (err error) {
	//   self._init()
	if err = e._init(); err != nil {
		return
	}
	//   self._initLut(mode)
	if err = e._initLut(mode); err != nil {
		return
	}
	//   self._powerOn()
	if err = e.d.PowerOn(); err != nil {
		return
	}

	switch mode {
	case ModePart:
		return e._disPart(0, Ymax-1, 0, Xmax-1)
	case ModeFull:
		if err := e.d.SetRAMPointer(0x00, (Xmax-1)%256, (Xmax-1)/256); err != nil {
			return err
		}

		if err := e.d.WriteDisRAM(e.displayBuffer, Ymax, Xmax); err != nil {
			return err
		}

		return e.d.UpdateDis(mode)
	}

	return err
}

// Pixel sets individual pixels on or off.
//   def pixel(self, x, y, color):
//     if x < 0 or x >= self._width:
//       return
//     if y < 0 or y >= self._height:
//       return
//     x = int(x)
//     y = int(y)
//     m = int(x * 16 + (y + 1) / 8)
//     sy = int((y + 1) % 8)
//     if color == self.WHITE:
//       if sy != 0:
//         self._displayBuffer[m] = self._displayBuffer[m] | int(pow(2, 8 - sy))
//       else:
//         self._displayBuffer[m - 1] = self._displayBuffer[m - 1] | 1
//     elif color == self.BLACK:
//       if sy != 0:
//         self._displayBuffer[m] = self._displayBuffer[m] & (0xff - int(pow(2, 8 - sy)))
//       else:
//         self._displayBuffer[m - 1] = self._displayBuffer[m - 1] & 0xfe
func (e *Epaper) Pixel(x, y int, on bool) (err error) {
	if x < 0 || x > e.w {
		return ErrOutOfBounds
	}

	if y < 0 || y > e.h {
		return ErrOutOfBounds
	}

	return e.displayBuffer.SetBit((y*Xmax)+x, on)
}

// ErrOutOfBounds is returned when a Pixel x/y is outside the window.
var ErrOutOfBounds = errors.New("position out of bounds")
