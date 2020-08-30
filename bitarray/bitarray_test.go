package bitarray_test

import (
	"github.com/na4ma4/dfrobot-eink-display/bitarray"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BitArray", func() {
	Context("Creating array", func() {
		It("size not multiple of 8", func() {
			_, err := bitarray.NewBitArray(7)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(bitarray.ErrInvalidSize))
		})

		It("negative size", func() {
			_, err := bitarray.NewBitArray(-10)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(bitarray.ErrInvalidSize))
		})

		It("negative size but multiple of 8", func() {
			_, err := bitarray.NewBitArray(-64)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(bitarray.ErrInvalidSize))
		})
	})

	Context("Setting Bits", func() {
		It("set all bits to on", func() {
			b, err := bitarray.NewBitArray(64)
			Expect(err).NotTo(HaveOccurred())
			for i := 0; i < 64; i++ {
				err = b.SetBit(i, true)
				Expect(err).NotTo(HaveOccurred())
			}

			o := b.ByteSlice()
			Expect(o).To(Equal([]byte{
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			}))
		})

		It("set all bits to on then off", func() {
			b, err := bitarray.NewBitArray(64)
			Expect(err).NotTo(HaveOccurred())
			for i := 0; i < 64; i++ {
				err = b.SetBit(i, true)
				Expect(err).NotTo(HaveOccurred())
			}

			o := b.ByteSlice()
			Expect(o).To(Equal([]byte{
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			}))

			for i := 0; i < 64; i++ {
				err = b.SetBit(i, false)
				Expect(err).NotTo(HaveOccurred())
			}

			o = b.ByteSlice()
			Expect(o).To(Equal([]byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			}))
		})

		It("set every second bit to on", func() {
			b, err := bitarray.NewBitArray(64)
			Expect(err).NotTo(HaveOccurred())
			for i := 0; i < 64; i += 2 {
				err = b.SetBit(i, true)
				Expect(err).NotTo(HaveOccurred())
			}

			o := b.ByteSlice()
			Expect(o).To(Equal([]byte{
				0b01010101, 0b01010101, 0b01010101, 0b01010101,
				0b01010101, 0b01010101, 0b01010101, 0b01010101,
			}))
		})
	})
})
