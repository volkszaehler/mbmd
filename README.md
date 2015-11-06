# A HTTP interface to the SDM630-MODBUS smart meter

This project provides a http interface to the Eastron SDM630 smart
meter. The device comes in many flavour - make sure to get the "MODBUS"
version. The SDM630-MODBUS exposes all measured values over an RS485
connection, making it very easy to integrate it into your home
automation system.

## How does it look like in OpenHAB?

TODO

## Installation

The installation consists necessarily of a hardware and a software part.
Make sure you buy/fetch the following things before starting:

* A SDM630 smart meter, the MODBUS version.
* A USB RS485 adaptor. Mine is from Digitus.
* Some cables to connect the adapter to the SDM630 (for testing, I use
an old speaker cable I had sitting on my workbench, for the permanent
installation, a shielded CAT5 cable seems adequate)
* Two 120 Ohm and two 680 Ohm resistors (1/4W metal).

### Hardware installation

![SDM630 in my test setup](img/SDM630-MODBUS.jpg)

First, you should integrate the SDM630 into your fuse box. Please ask a
professional to do this for you - I don't want you to hurt yourself!
Refer to the SDM630 installation manual on how to do this. You need to
set the MODBUS communication parameters to ``9600 8N1``. I obtained my
SDM630 from the German distributor [B+G
E-Tech](http://bg-etech.de/os/product_info.php/cPath/25_28/products_id/50).

After this you need to connect the USB adaptor to the SDM630. This is
how I did the wiring:

![USB-SDM630 wiring](img/wiring.jpg)

I got this [Digitus USB-RS485
adaptor](http://www.digitus.info/en/products/accessories/adapter-and-converter/r-usb-serial-adapter-usb-20-da-70157/)
which comes with a handy terminal block. I mounted the [bias
network](https://en.wikipedia.org/wiki/RS-485) directly on the terminal
block:

![bias network](img/USB-RS485-Adaptor.jpg)

### Installing the software from source

You need a working [Golang installation](http://golang.org) and the [GB
build tool](http://getgb.io/) in order to compile your binary. You can
install GB like this:

    go get github.com/constabulary/gb/...

Afterwards you clone this repository:

    git clone https://github.com/gonium/gosdm630.git

and build it:

    gb build all

Now, there should be a binary in the ````bin```` subfolder.

### Testing

Now fire up the software:

    ./bin/sdm630_httpd -d /dev/ttyUSB1 -u localhost:8080 -v
    RTUClientHandler: 2015/11/06 12:22:14 modbus: sending 01 04 00 00 00 02 71 cb
    RTUClientHandler: 2015/11/06 12:22:14 modbus: received 01 04 04 43 6b d7 3d 01 fd
    RTUClientHandler: 2015/11/06 12:22:14 modbus: sending 01 04 00 02 00 02 d0 0b
    RTUClientHandler: 2015/11/06 12:22:14 modbus: received 01 04 04 00 00 00 00 fb 84
    RTUClientHandler: 2015/11/06 12:22:14 modbus: sending 01 04 00 04 00 02 30 0a
    RTUClientHandler: 2015/11/06 12:22:14 modbus: received 01 04 04 00 00 00 00 fb 84
    RTUClientHandler: 2015/11/06 12:22:14 modbus: sending 01 04 00 06 00 02 91 ca
    RTUClientHandler: 2015/11/06 12:22:14 modbus: received 01 04 04 00 00 00 00 fb 84
    RTUClientHandler: 2015/11/06 12:22:14 modbus: sending 01 04 00 08 00 02 f0 09
    RTUClientHandler: 2015/11/06 12:22:14 modbus: received 01 04 04 00 00 00 00 fb 84
    RTUClientHandler: 2015/11/06 12:22:14 modbus: sending 01 04 00 0a 00 02 51 c9
    RTUClientHandler: 2015/11/06 12:22:14 modbus: received 01 04 04 00 00 00 00 fb 84
    RTUClientHandler: 2015/11/06 12:22:14 modbus: sending 01 04 00 0c 00 02 b1 c8
    RTUClientHandler: 2015/11/06 12:22:14 modbus: received 01 04 04 00 00 00 00 fb 84
    RTUClientHandler: 2015/11/06 12:22:14 modbus: sending 01 04 00 0e 00 02 10 08
    RTUClientHandler: 2015/11/06 12:22:14 modbus: received 01 04 04 00 00 00 00 fb 84
    RTUClientHandler: 2015/11/06 12:22:14 modbus: sending 01 04 00 10 00 02 70 0e
    RTUClientHandler: 2015/11/06 12:22:14 modbus: received 01 04 04 00 00 00 00 fb 84
    RTUClientHandler: 2015/11/06 12:22:14 modbus: sending 01 04 00 1e 00 02 11 cd
    RTUClientHandler: 2015/11/06 12:22:14 modbus: received 01 04 04 3f 80 00 00 f6 78
    RTUClientHandler: 2015/11/06 12:22:14 modbus: sending 01 04 00 20 00 02 70 01
    RTUClientHandler: 2015/11/06 12:22:14 modbus: received 01 04 04 3f 80 00 00 f6 78
    RTUClientHandler: 2015/11/06 12:22:14 modbus: sending 01 04 00 22 00 02 d1 c1
    RTUClientHandler: 2015/11/06 12:22:14 modbus: received 01 04 04 3f 80 00 00 f6 78
    T: 2015-11-06T12:22:14+01:00 - L1: 235.84V 0.00A 0.00W 1.00cos | L2: 0.00V 0.00A 0.00W 1.00cos | L3: 0.00V 0.00A 0.00W 1.00cos

You can see modbus traffic and the current readings on the command line.
If you visit [http://localhost:8080](http://localhost:8080) you should
also see the last received value printed as ASCII text. Please note: The
exact representation might change in future releases.
