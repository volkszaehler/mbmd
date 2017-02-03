# A HTTP interface to the Eastron SDM-MODBUS smart meter series

This project provides a http interface to the Eastron SDM smart
meter series with MODBUS interface. The smart meter comes in many flavours - make sure to get the "MODBUS"
version. These meters exposes all measured values over an RS485
connection, making it very easy to integrate it into your home
automation system.

Suported Devices:

 * SDM120
 * SDM220, SDM230
 * SDM530
 * SDM630 v1 and v2, both MID and non-MID

Table of Contents:

  * [How does it look like in OpenHAB?](#how-does-it-look-like-in-openhab)
  * [Installation](#installation)
    * [Hardware installation](#hardware-installation)
    * [Installing the software from source](#installing-the-software-from-source)
    * [Testing](#testing)
    * [Installation on the Raspberry Pi](#installation-on-the-raspberry-pi)
  * [The API](#the-api)
  * [Monitoring](#monitoring)
  * [OpenHAB integration](#openhab-integration)


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

The installation consists of a hardware and a software part.
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
install the Go compiler first. Then clone this repository:

    git clone https://github.com/gonium/gosdm630.git

If you have ``make`` installed you can
use the ``Makefile`` to install the GB build tool:

    $ cd gosdm630
    $ make dep
    Installing GB build tool

Or, if you prefer to install it manually, just run 

    go get github.com/constabulary/gb/...

You can then build the software using the ``Makefile``:

    $ make
    Building for host platform
    Building binary for Raspberry Pi
    github.com/gonium/gosdm630
    github.com/gonium/gosdm630/cmd/sdm630_httpd
    Created binaries:
    sdm630_httpd
    sdm630_httpd-linux-arm

As you can see two binaries are built:

 * ``bin/sdm630_httpd`` is the software built for the host platform
 * ``bin/sdm630_httpd-linux-arm`` is the same for the Raspberry Pi.

If you prefer to build manually you can build the host software using

    gb build all

or, for the Raspberry Pi binary, use

    GOOS=linux GOARCH=arm GOARM=5 gb build all

### Testing

Now fire up the software:

    $ ./bin/sdm630_httpd --help
    NAME:
       sdm630_httpd - SDM630 power measurements via HTTP.
    
    USAGE:
       sdm630_httpd [global options] command [command options] [arguments...]
    
    COMMANDS:
         help, h  Shows a list of commands or help for one command
    
    GLOBAL OPTIONS:
       --serialadapter value, -s value  path to serial RTU device (default: "/dev/ttyUSB0")
       --url value, -u value            the URL the server should respond on (default: ":8080")
       --verbose, -v                    print verbose messages
       --device_list value, -d value    MODBUS device ID to query (default: "1")
       --help, -h                       show help
    

A typical invocation looks like this:

		./bin/sdm630_httpd -s /dev/ttyUSB0 -d 11,12,13,14,15 -v
    2017/01/25 16:34:26 Connecting to RTU via /dev/ttyUSB0
    2017/01/25 16:34:26 Starting API httpd at :8080
    RTUClientHandler: 2017/01/25 16:34:26 modbus: sending 0b 04 00 00 00 02 71 61
    RTUClientHandler: 2017/01/25 16:34:26 modbus: received 0b 04 04 43 6c 78 80 a7 bd
    RTUClientHandler: 2017/01/25 16:34:26 modbus: sending 0b 04 00 02 00 02 d0 a1
    RTUClientHandler: 2017/01/25 16:34:26 modbus: received 0b 04 04 43 6c 74 9c a3 74
    RTUClientHandler: 2017/01/25 16:34:26 modbus: sending 0b 04 00 04 00 02 30 a0
    RTUClientHandler: 2017/01/25 16:34:26 modbus: received 0b 04 04 00 00 00 00 51 84
    RTUClientHandler: 2017/01/25 16:34:26 modbus: sending 0b 04 00 06 00 02 91 60
    RTUClientHandler: 2017/01/25 16:34:26 modbus: received 0b 04 04 00 00 00 00 51 84
    RTUClientHandler: 2017/01/25 16:34:26 modbus: sending 0b 04 00 08 00 02 f0 a3
    RTUClientHandler: 2017/01/25 16:34:26 modbus: received 0b 04 04 00 00 00 00 51 84
    RTUClientHandler: 2017/01/25 16:34:26 modbus: sending 0b 04 00 0a 00 02 51 63
    RTUClientHandler: 2017/01/25 16:34:27 modbus: received 0b 04 04 00 00 00 00 51 84
    RTUClientHandler: 2017/01/25 16:34:27 modbus: sending 0b 04 00 0c 00 02 b1 62
    RTUClientHandler: 2017/01/25 16:34:27 modbus: received 0b 04 04 00 00 00 00 51 84
    RTUClientHandler: 2017/01/25 16:34:27 modbus: sending 0b 04 00 0e 00 02 10 a2
    RTUClientHandler: 2017/01/25 16:34:27 modbus: received 0b 04 04 00 00 00 00 51 84
    RTUClientHandler: 2017/01/25 16:34:27 modbus: sending 0b 04 00 10 00 02 70 a4
    RTUClientHandler: 2017/01/25 16:34:27 modbus: received 0b 04 04 00 00 00 00 51 84
    RTUClientHandler: 2017/01/25 16:34:27 modbus: sending 0b 04 00 1e 00 02 11 67
    RTUClientHandler: 2017/01/25 16:34:27 modbus: received 0b 04 04 3f 80 00 00 5c 78
    RTUClientHandler: 2017/01/25 16:34:27 modbus: sending 0b 04 00 20 00 02 70 ab
    RTUClientHandler: 2017/01/25 16:34:27 modbus: received 0b 04 04 3f 80 00 00 5c 78
    RTUClientHandler: 2017/01/25 16:34:27 modbus: sending 0b 04 00 22 00 02 d1 6b
    RTUClientHandler: 2017/01/25 16:34:27 modbus: received 0b 04 04 3f 80 00 00 5c 78
    RTUClientHandler: 2017/01/25 16:34:27 modbus: sending 0b 04 01 5a 00 02 50 8e
    RTUClientHandler: 2017/01/25 16:34:27 modbus: received 0b 04 04 3b 03 12 6f e0 2c
    RTUClientHandler: 2017/01/25 16:34:27 modbus: sending 0b 04 01 5c 00 02 b0 8f
    RTUClientHandler: 2017/01/25 16:34:27 modbus: received 0b 04 04 3b 03 12 6f e0 2c
    RTUClientHandler: 2017/01/25 16:34:27 modbus: sending 0b 04 01 5e 00 02 11 4f
    T: 2015-11-06T12:22:14+01:00 - L1: 235.84V 0.00A 0.00W 1.00cos | L2: 0.00V 0.00A 0.00W 1.00cos | L3: 0.00V 0.00A 0.00W 1.00cos

This call queries five devices with the IDs 11, 12, 13, 14 and 15. If
you use the ``-v`` commandline switch you can see modbus traffic and the
current readings on the command line.  At
[http://localhost:8080](http://localhost:8080) you should also see the
last received value printed as ASCII text:

    Modbus ID 11, last measurement taken Thursday, 12-Nov-15 14:18:10 CET:
    +-------+-------------+-------------+-----------+--------------+--------------+--------------+
    | PHASE | VOLTAGE [V] | CURRENT [A] | POWER [W] | POWER FACTOR | IMPORT [KWH] | EXPORT [KWH] |
    +-------+-------------+-------------+-----------+--------------+--------------+--------------+
    | L1    |      235.17 |        0.19 |     45.83 |         1.00 |         0.74 |         0.01 |
    | L2    |        0.00 |        0.00 |      0.00 |         1.00 |         0.00 |         0.00 |
    | L3    |        0.00 |        0.00 |      0.00 |         1.00 |         0.00 |         0.00 |
    | ALL   | n/a         |        0.19 |     45.83 | n/a          |         0.74 |         0.01 |
    +-------+-------------+-------------+-----------+--------------+--------------+--------------+

### Installation on the Raspberry Pi

You simply copy the binary from the ``bin`` subdirectory to the RPi
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

## The API

As of version 0.2 the software supports more than one device. In order
to query the API you need to provide the MODBUS ID of the device you
want to query.

The API consists of calls that return a JSON datastructure. The "GET
/last/{ID}"-call simply returns the last measurements of the device with
the Modbus ID {ID}:

    $ curl localhost:8080/last/1
    {
      "Timestamp": "2017-01-25T16:40:28.706413719+01:00",
      "Unix": 1485358828,
      "ModbusDeviceId": 11,
      "Power": {
        "L1": 0,
        "L2": 0,
        "L3": 0
      },
      "Voltage": {
        "L1": 236.27914,
        "L2": 236.26003,
        "L3": 0
      },
      "Current": {
        "L1": 0,
        "L2": 0,
        "L3": 0
      },
      "Cosphi": {
        "L1": 1,
        "L2": 1,
        "L3": 1
      },
      "Import": {
        "L1": 0.002,
        "L2": 0.002,
        "L3": 0.001
      },
      "Export": {
        "L1": 0,
        "L2": 0,
        "L3": 0
      }
    }


The "GET /minuteavg"-call returns the average measurements over the last
minute:

    $ curl localhost:8080/minuteavg/11
    {
      "Timestamp": "2017-01-25T16:41:15.835808811+01:00",
      "Unix": 1485358875,
      "ModbusDeviceId": 11,
      "Power": {
        "L1": 0,
        "L2": 0,
        "L3": 0
      },
      "Voltage": {
        "L1": 236.49846,
        "L2": 236.4804,
        "L3": 0
      },
      "Current": {
        "L1": 0,
        "L2": 0,
        "L3": 0
      },
      "Cosphi": {
        "L1": 1,
        "L2": 1,
        "L3": 1
      },
      "Import": {
        "L1": 0,
        "L2": 0,
        "L3": 0
      },
      "Export": {
        "L1": 0,
        "L2": 0,
        "L3": 0
      }
    }


If you want to receive all measurements, you can use these two calls
without the device ID:

    curl localhost:8080/last
    [{"Timestamp":"2017-01-25T16:39:56.211376135+01:00","Unix":1485358796,"ModbusDeviceId":11,"Power":{"L1":0,"L2":0,"L3":0},"Voltage":{"L1":236.50807,"L2":236.49356,"L3":0},"Current":{"L1":0,"L2":0,"L3":0},"Cosphi":{"L1":1,"L2":1,"L3":1},"Import":{"L1":0.002,"L2":0.002,"L3":0.001},"Export":{"L1":0,"L2":0,"L3":0}},{"Timestamp":"2017-01-25T16:39:56.794948625+01:00","Unix":1485358796,"ModbusDeviceId":12,"Power":{"L1":0,"L2":0,"L3":0},"Voltage":{"L1":236.40024,"L2":236.46877,"L3":0},"Current":{"L1":0,"L2":0,"L3":0},"Cosphi":{"L1":1,"L2":1,"L3":1},"Import":{"L1":0.001,"L2":0.002,"L3":0.001},"Export":{"L1":0,"L2":0,"L3":0}},{"Timestamp":"2017-01-25T16:39:57.37536849+01:00","Unix":1485358797,"ModbusDeviceId":13,"Power":{"L1":0,"L2":0,"L3":0},"Voltage":{"L1":236.50534,"L2":0,"L3":0},"Current":{"L1":0,"L2":0,"L3":0},"Cosphi":{"L1":1,"L2":0,"L3":0},"Import":{"L1":0,"L2":0,"L3":0},"Export":{"L1":0,"L2":0,"L3":0}},{"Timestamp":"2017-01-25T16:39:55.02659946+01:00","Unix":1485358795,"ModbusDeviceId":14,"Power":{"L1":0,"L2":0,"L3":0},"Voltage":{"L1":236.41461,"L2":236.50677,"L3":0},"Current":{"L1":0,"L2":0,"L3":0},"Cosphi":{"L1":1,"L2":1,"L3":1},"Import":{"L1":0,"L2":0.001,"L3":0},"Export":{"L1":0,"L2":0,"L3":0}},{"Timestamp":"2017-01-25T16:39:55.627042868+01:00","Unix":1485358795,"ModbusDeviceId":15,"Power":{"L1":0,"L2":0,"L3":0},"Voltage":{"L1":236.52869,"L2":0,"L3":0},"Current":{"L1":0,"L2":0,"L3":0},"Cosphi":{"L1":1,"L2":0,"L3":0},"Import":{"L1":0,"L2":0,"L3":0},"Export":{"L1":0,"L2":0,"L3":0}}]

and so on. I recommend the [JSON Incremental Digger
(jid)](https://github.com/simeji/jid) for exploring json datasets.

## Monitoring

In order to monitor this long running process a status report is now
available:

    $ curl http://localhost:8080/status
    {
      "Starttime": "2017-01-25T16:35:50.839829945+01:00",
      "UptimeSeconds": 65587.177092186,
      "Goroutines": 11,
      "Memory": {
        "Alloc": 1568344,
        "HeapAlloc": 1568344
      },
      "Modbus": {
        "TotalModbusRequests": 1979122,
        "ModbusRequestRatePerMinute": 1810.5264666764785,
        "TotalModbusErrors": 738,
        "ModbusErrorRatePerMinute": 0.6751319688261972
      }
    }

This is a snapshot of a process running over night, along with the error
statistics during that timeframe. The process queries 5 meters continuously,
the cabling is not a shielded, twisted wire but something that I had laying
around. With proper cabling the error rate should be lower, though.

## OpenHAB integration

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

## Firehose

The firehose enables you to observe the data read from the smart meter
in realtime: as soon as a new value is available, you will be notified.
We're using [HTTP Long Polling as described in RFC
6202](https://tools.ietf.org/html/rfc6202) for the data transfer. This
essentially means that you can connect to an HTTP endpoint. The server
will accept the connection and send you the new values as soon as they
are available. Then, you either reconnect or use the same TCP connection
for the next request. If you want to get all values, you can do the
following:

    $ while true; do curl --silent "http://localhost:8080/firehose?timeout=45&category=all" | jq; done

This requests the last values in a loop with curl and pipes the result
through jq. Of course this also closes the connection after each reply,
so this is rather costly. In production you can leave the connection
intact and reuse it.
