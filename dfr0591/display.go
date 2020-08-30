package dfr0591

import (
	"context"
	"time"

	"github.com/na4ma4/dfrobot-eink-display/bitarray"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/spi"
)

// DisplayTransport is the interface for the e-Ink Display (I think it can support i2c as well ?).
type DisplayTransport interface {
	PowerOn() (err error)
	PowerOff() (err error)
	SetRAMData(xStart, xEnd, yStart, yStart1, yEnd, yEnd1 uint8) (err error)
	SetRAMPointer(x, y, y1 uint8) (err error)
	WriteCmdAndData(cmd uint8, data []byte) (err error)
	WriteDisRAM(displayBuffer *bitarray.BitArray, sizeX, sizeY uint8) (err error)
	UpdateDis(mode Mode) (err error)
}

// DisplaySPI is the SPI provider for DisplayTransport.
type DisplaySPI struct {
	spi  spi.Conn
	cs   gpio.PinIO
	cd   gpio.PinIO
	busy gpio.PinIO
}

// NewDisplaySPI returns a new DisplaySPI that provides the DisplayTransport interface.
func NewDisplaySPI(s spi.Conn) DisplayTransport {
	return &DisplaySPI{
		spi:  s,
		cs:   gpioreg.ByName("GPIO27"),
		cd:   gpioreg.ByName("GPIO17"),
		busy: gpioreg.ByName("GPIO4"),
	}
}

// WriteCmdAndData writes a command and optional data to the device.
func (d *DisplaySPI) WriteCmdAndData(cmd uint8, data []byte) (err error) {
	d.waitBusy(context.Background())

	_ = d.cs.Out(gpio.Low)
	_ = d.cd.Out(gpio.Low)

	if err = d.spi.Tx([]byte{cmd}, nil); err != nil {
		return
	}

	_ = d.cd.Out(gpio.High)

	if data != nil {
		if err = d.spi.Tx(data, nil); err != nil {
			return
		}
	}

	_ = d.cs.Out(gpio.High)

	return
}

// WriteDisRAM writes to the display RAM on the device.
// def _writeDisRam(self, sizeX, sizeY):
// 	if sizeX % 8 != 0:
// 		sizeX = sizeX + (8 - sizeX % 8)
// 	sizeX = sizeX // 8
// 	self.writeCmdAndData(0x24, self._displayBuffer[0: sizeX * sizeY])
//nolint:gomnd // Lots of magic numbers, maybe I'll extract them later.
func (d *DisplaySPI) WriteDisRAM(displayBuffer *bitarray.BitArray, sizeX, sizeY uint8) (err error) {
	// log.Printf("WriteDisRAM(buf, %d, %d)[before]", sizeX, sizeY)
	if sizeX%8 != 0 {
		sizeX += (8 - (sizeX % 8))
	}

	sizeX /= 8

	// log.Printf("WriteDisRAM(buf, %d, %d)[after] buf[0:%d]", sizeX, sizeY, uint(uint(sizeX)*uint(sizeY)))
	return d.WriteCmdAndData(0x24, displayBuffer.ByteSlice()[0:int(sizeX)*int(sizeY)])
}

// UpdateDis updates the display of the device.
// def _updateDis(self, mode):
// 	if mode == self.FULL:
// 		self.writeCmdAndData(0x22, [0xc7])
// 	elif mode == self.PART:
// 		self.writeCmdAndData(0x22, [0x04])
// 	else:
// 		return
// 	self.writeCmdAndData(0x20, [])
// 	self.writeCmdAndData(0xff, [])
//nolint:gomnd // Lots of magic numbers, maybe I'll extract them later.
func (d *DisplaySPI) UpdateDis(mode Mode) (err error) {
	switch mode {
	case ModeFull:
		if err = d.WriteCmdAndData(0x22, []byte{0xc7}); err != nil {
			return err
		}
	case ModePart:
		if err = d.WriteCmdAndData(0x22, []byte{0x04}); err != nil {
			return err
		}
	}

	if err = d.WriteCmdAndData(0x20, nil); err != nil {
		return err
	}

	return d.WriteCmdAndData(0xff, nil)
}

// PowerOn wakes the device from sleep.
// def _powerOn(self):
//    self.writeCmdAndData(0x22, [0xc0])
//    self.writeCmdAndData(0x20, [])
//nolint:gomnd // Lots of magic numbers, maybe I'll extract them later.
func (d *DisplaySPI) PowerOn() (err error) {
	//    self.writeCmdAndData(0x22, [0xc0])
	if err = d.WriteCmdAndData(0x22, []byte{0xc0}); err != nil {
		return
	}

	//    self.writeCmdAndData(0x20, [])
	return d.WriteCmdAndData(0x20, nil)
}

// PowerOff puts the device to sleep.
// def _powerOff(self):
//    self.writeCmdAndData(0x12, [])
// 	  self.writeCmdAndData(0x82, [0x00])
// 	  self.writeCmdAndData(0x01, [0x02, 0x00, 0x00, 0x00, 0x00])
// 	  self.writeCmdAndData(0x02, [])
//nolint:gomnd // Lots of magic numbers, maybe I'll extract them later.
func (d *DisplaySPI) PowerOff() (err error) {
	if err = d.WriteCmdAndData(0x12, nil); err != nil {
		return
	}

	if err = d.WriteCmdAndData(0x82, []byte{0x00}); err != nil {
		return
	}

	if err = d.WriteCmdAndData(0x01, []byte{0x02, 0x00, 0x00, 0x00, 0x00}); err != nil {
		return
	}

	return d.WriteCmdAndData(0x02, nil)
}

// SetRAMData does something.
//nolint:gomnd // Lots of magic numbers, maybe I'll extract them later.
func (d *DisplaySPI) SetRAMData(xStart, xEnd, yStart, yStart1, yEnd, yEnd1 uint8) (err error) {
	if err = d.WriteCmdAndData(0x44, []byte{xStart, xEnd}); err != nil {
		return
	}

	return d.WriteCmdAndData(0x45, []byte{yStart, yStart1, yEnd, yEnd1})
}

// SetRAMPointer also does something.
//nolint:gomnd // Lots of magic numbers, maybe I'll extract them later.
func (d *DisplaySPI) SetRAMPointer(x, y, y1 uint8) (err error) {
	if err = d.WriteCmdAndData(0x4e, []byte{x}); err != nil {
		return
	}

	return d.WriteCmdAndData(0x4f, []byte{y, y1})
}

func (d *DisplaySPI) readBusy() gpio.Level {
	return d.busy.Read()
}

func (d *DisplaySPI) waitBusy(ctx context.Context) {
	for {
		select {
		case <-time.After(time.Millisecond):
			if !d.readBusy() {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
