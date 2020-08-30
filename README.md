# dfrobot-eink-display

DFRobot e-Ink Display library for golang

Work in progress, basic functionality is working (set pixel, flush to display).

## Development

Something isn't deactivating the device properly, so to get the python script working again on the RPi Zero W:

```shell
sudo rmmod spi_bcm2835 && sudo modprobe spi_bcm2835
```

## Documentation

- [PDF Documentation](docs/20180622175013limdcz.pdf)
- [DFRobot/DFRobot_RPi_Display](https://github.com/DFRobot/DFRobot_RPi_Display)
