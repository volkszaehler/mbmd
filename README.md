# A HTTP interface to the SDM630-MODBUS smart meter

This project provides a http interface to the Eastron SDM630 smart
meter. The device comes in many flavour - make sure to get the "MODBUS"
version. The SDM630-MODBUS exposes all measured values over an RS485
connection, making it very easy to integrate it into your home
automation system.

## How does it look like in OpenHAB?

I use [OpenHAB](http://openhab.org) to record various measurements at
home. In the classic ui, this is how one of the graphs looks like:

![OpenHAB interface screenshot](img/openhab.png)

Everything is in German, but the "Verlauf Strombezug" graph shows my
power consumption for three phases. I have a SDM630 installed in my
distribution cabinet. A serial connection links it to a Raspberry Pi
(RPi).
This is where this piece of software runs and exposes the measurements
via a RESTful API. OpenHAB connects to it and stores the values, just as
it does with other sensors in my home.

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
build tool](http://getgb.io/) in order to compile your binary. Please
install the Go compiler first. Afterwards you can install GB like this:

    go get github.com/constabulary/gb/...

Clone this repository:

    git clone https://github.com/gonium/gosdm630.git

and build it:

TODO: Mention Makefile


    cd gosdm630
    gb build all

Now, there should be a binary in the ````bin```` subfolder.

### Testing

Now fire up the software:

    ./bin/sdm630_httpd -s /dev/ttyUSB1 -u localhost:8080 -v
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

If you use the ``-v`` commandline switch you can see modbus traffic and
the current readings on the command line.  If you visit
[http://localhost:8080](http://localhost:8080) you should also see the
last received value printed as ASCII text:

    Last measurement taken Thursday, 12-Nov-15 14:18:10 CET:
    +-------+-------------+-------------+-----------+--------------+--------------+--------------+
    | PHASE | VOLTAGE [V] | CURRENT [A] | POWER [W] | POWER FACTOR | IMPORT [KWH] | EXPORT [KWH] |
    +-------+-------------+-------------+-----------+--------------+--------------+--------------+
    | L1    |      235.17 |        0.19 |     45.83 |         1.00 |         0.74 |         0.01 |
    | L2    |        0.00 |        0.00 |      0.00 |         1.00 |         0.00 |         0.00 |
    | L3    |        0.00 |        0.00 |      0.00 |         1.00 |         0.00 |         0.00 |
    | ALL   | n/a         |        0.19 |     45.83 | n/a          |         0.74 |         0.01 |
    +-------+-------------+-------------+-----------+--------------+--------------+--------------+

### Crosscompiling e.g. for Raspberry Pi

Go has very good crosscompilation support. Typically, I develop under
Mac OS and crosscompile a binary for my RPi. It is easy:

    # clear whatever old binaries I have
    rm -rf pkg bin
    # start crosscompilation
    GOOS=linux GOARCH=arm GOARM=5 gb build all

You can then copy the binary from the ``bin`` subdirectory to the RPi
and start it. I usually put the binary into ``/usr/local/bin`` and
rename it to ``sdm630_httpd``. The following sytemd unit can be used to
start the service (put this into ``/etc/systemd/system``):

    [Unit]
    Description=SDM630 via HTTP API
    After=syslog.target
    [Service]
    ExecStart=/usr/local/bin/sdm630_httpd -s /dev/ttyAMA0
    Restart=always
    [Install]
    WantedBy=multi-user.target

You might need to adjust the ``-s`` parameter depending on where your
RS485 adapter is connected. Then, use

    # systemctl start sdm630

to test your installation. If you're satisfied use

    # systemctl enable sdm630

to start the service at boot time automatically.

## The API: OpenHAB integration

TODO: Mention the use of several IDs

The API consists of two calls that return a JSON array. The "GET
/last/{ID}"-call simply returns the last measurements of the device with
the Modbus ID {ID}:

    $ curl localhost:8080/last/1
		{"Timestamp":"2017-01-12T20:14:44.375777188Z","Unix":1484252084,"ModbusDeviceId":1,"Power":{"L1":309.1006,"L2":81.13953,"L3":0},"Voltage":{"L1":230.93193,"L2":231.07565,"L3":232.11107},"Current":{"L1":1.4531646,"L2":0.36446536,"L3":0},"Cosphi":{"L1":0.9210861,"L2":0.9629059,"L3":1},"Import":{"L1":2282.521,"L2":693.844,"L3":238.658},"Export":{"L1":0,"L2":0.005,"L3":0.024}}

The "GET /minuteavg"-call returns the average measurements over the last
minute:

    $ curl localhost:8080/minuteavg/1
    {"Timestamp":"2017-01-12T20:15:17.17833005Z","Unix":1484252117,"ModbusDeviceId":1,"Power":{"L1":300.57672,"L2":81.01144,"L3":0},"Voltage":{"L1":231.22124,"L2":231.07536,"L3":231.94556},"Current":{"L1":1.3945557,"L2":0.3644369,"L3":0},"Cosphi":{"L1":0.9347229,"L2":0.961964,"L3":1},"Import":{"L1":0,"L2":0,"L3":0},"Export":{"L1":0,"L2":0,"L3":0}}

It is very easy to translate this into OpenHAB items. I run the SDM630
software on a Raspberry Pi with the IP ``192.168.1.44``. My items look
like this:

    Group Power_Chart
    Number Power_L1 "Strombezug L1 [%.1f W]" <power> (Power, Power_Chart) { http="<[http://192.168.1.44:8080/last/1:60000:JS(SDM630GetL1Power.js)]" }

I'm using the http plugin to call the ``/last/1`` endpoint every 60
seconds. Then, I feed the result into a JSON transform stored in
``SDM630GetL1Power.js``. The contents of
``transform/SDM630GetL1Power.js`` looks like this:

    JSON.parse(input).Power.L1;

Just repeat these lines for each measurement you want to track. Finally,
my sitemap contains the following lines:

    Chart item=Power_Chart period=D refresh=1800

This draws a chart of all items in the ``Power_Chart`` group.

