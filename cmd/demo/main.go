package main

import (
	"log"
	"math/rand"

	"github.com/na4ma4/dfrobot-eink-display/dfr0591"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/host"
)

// nolint: gochecknoglobals // cobra uses globals in main
var rootCmd = &cobra.Command{
	Use:  "demo",
	Run:  mainCommand,
	Args: cobra.NoArgs,
}

// nolint:gochecknoinits // init is used in main for cobra
func init() {
	cobra.OnInitialize(configInit)

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Debug output")
	_ = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	_ = viper.BindEnv("debug", "DEBUG")
}

func main() {
	_ = rootCmd.Execute()
}

const (
	xMax       = 250
	yMax       = 128
	randMax    = 10
	randMiddle = 5
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func mainCommand(cmd *cobra.Command, args []string) {
	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// gpiocs := gpioreg.ByName("27")
	// gpiocd := gpioreg.ByName("17")
	// gpiobusy := gpioreg.ByName("4")

	// pins := gpioreg.All()
	// for _, pin := range pins {
	// 	log.Printf("%s(%d): %s", pin.Name(), pin.Number(), pin.String())
	// }

	// Use spireg SPI port registry to find the first available SPI bus.
	p, err := spireg.Open("")
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close()

	// Convert the spi.Port into a spi.Conn so it can be used for communication.
	c, err := p.Connect(physic.MegaHertz, spi.HalfDuplex+spi.NoCS, 8)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Duplex: %s", c.Duplex().String())

	dt := dfr0591.NewDisplaySPI(c)
	ep := dfr0591.NewEpaper(dt, xMax, yMax)

	err = ep.Flush(dfr0591.ModeFull)
	checkErr(err)

	for y := 0; y < yMax; y++ {
		for x := 0; x < xMax; x++ {
			if rand.Intn(randMax) >= randMiddle { //nolint:gosec // not concerned with security in this rand function
				err = ep.Pixel(x, y, true) // white
				checkErr(err)
			} else {
				err = ep.Pixel(x, y, false) // black
				checkErr(err)
			}
		}
	}

	err = ep.Flush(dfr0591.ModeFull)
	checkErr(err)
}
