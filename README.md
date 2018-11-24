# An HTTP interface to MODBUS smart meters

GoSDM provides an http interface to smart meters with a MODBUS interface.
Meter data, i.e. readings, is made accessible through REST API and MQTT.
Communication is possible over RS485 connections as well as TCP sockets.

A wide range of DIN-rail meters is supported (see [supported sevices](#supported-devices)).

**NOTE** Starting with version 0.7 several breaking changes were introduced. See [changelog](#changelog) for details.

# Table of Contents

* [Supported Devices](#supported-devices)
* [Installation](#installation)
  * [Hardware installation](#hardware-installation)
  * [Software installation](#software-installation)
  * [Running](#running)
  * [Installation on the Raspberry Pi](#installation-on-the-raspberry-pi)
  * [Detecting connected meters](#detecting-connected-meters)
* [API](#api)
  * [Rest API](#rest-api)
  * [Websocket API](#websocket-api)
  * [MQTT API](#mqtt-api)
* [Changelog](#changelog)

# Supported Devices

The meters have slightly different capabilities. The EASTRON SDM630 offers
a lot of features, while the smaller devices only support basic
features.  This table gives you an overview (please note: check the
manuals for yourself, I could be wrong):

| Meter | Phases | Voltage | Current | Power | Power Factor | Total Import | Total Export | Per-phase Import/Export | Line/Neutral THD |
|---|---|---|---|---|---|---|---|---|---|
| SDM120 | 1 | + | + | + | + | + | + | - | - |
| SDM220 | 1 | + | + | + | + | + | + | - | - |
| SDM220 | 1 | + | + | + | + | + | + | - | - |
| SDM530 | 3 | + | + | + | + | + | + | - | - |
| SDM630 v1 | 3 | + | + | + | + | + | + | + | + |
| SDM630 v2 | 3 | + | + | + | + | + | + | + | + |
| Janitza B23-312 | 3 | + | + | + | + | + | + | - | - |
| DZG DVH4013 | 3 | + | + | - | - | + | + | - | - |
| SBC ALE3 | 3 | + | + | + | + | + | + | - | - |
| ABB B-Series | 3 | + | + | + | + | + | + | + | + |
| SunSpec Inverters | 3 | + | + | + | + | - | + | - | - |

Please note that voltage, current, power and power factor are always
reported for each connected phase.

 * SDM120: Cheap and small (1TE), but communication parameters can only be set over MODBUS,
		 which is currently not supported by this project. You can use e.g.
		 [SDM120C](https://github.com/gianfrdp/SDM120C) to change parameters.
 * SDM220, SDM230: More comfortable (2TE), can be configured using the builtin display and
 button.
 * SDM530: Very big (7TE) - takes up a lot of space, but all connections are
 on the underside of the meter.
 * SDM630 v1 and v2, both MID and non-MID. Compact (4TE) and with lots
 of features. Can be configured for 1P2 (single phase with neutral), 3P3
 (three phase without neutral) and 3P4 (three phase with neutral)
	systems.
 * Janitza B23-312: These meters have a higher update rate than the Eastron
 devices, but they are more expensive. The -312 variant is the one with a MODBUS interface.
 * DZG DVH4013: This meter does not provide raw phase power measurements
 and only aggregated import/export measurements. The meter is only
 partially implemented and not recommended. If you want to use it: By
 default, the meter communicates using 9600 8E1 (comset 5). The meter ID
 is derived from the serial number: take the last two numbers of the
 serial number (top right of the device), e.g. 23, and add one (24).
 Assume this is a hexadecimal number and convert it to decimal (36). Use
 this as the meter ID.
 * SBC ALE3: This compact Saia Burgess Controls meter is comparable to the SDM630:
 two tariffs, both import and export depending on meter version and compact (4TE).
 It's often used with Viessmann heat pumps.
 * SunSpec: Apart from meters, SunSpec-compatible grid inverters are supported, too. This
 includes popular devices from SolarEdge (SE3000, SE9000) and SMA (planned). Grid inverters
 are typically connected using ModBus over TCP.

Some of my test devices have been provided by [B+G
E-Tech](http://bg-etech.de/) - please consider to buy your meter from
them!


# Installation

The installation consists of a hardware and a software part.
Make sure you buy/fetch the following things before starting:

* A supported Modbus/RTU smart meter.
* A USB RS485 adaptor. I use a homegrown one, please see [my
USB-ISO-RS485 project](https://github.com/gonium/usb-iso-rs485)
* Some cables to connect the adapter to the SDM630 (for testing, I use
an old speaker cable I had sitting on my workbench, for the permanent
installation, a shielded CAT5 cable seems adequate)


## Hardware installation

![SDM630 in my test setup](img/SDM630-MODBUS.jpg)

First, you should integrate the meter into your fuse box. Please ask a
professional to do this for you - I don't want you to hurt yourself!
Refer to the meter installation manual on how to do this. You need to
set the MODBUS communication parameters to ``9600 8N1``.
After this you need to connect a RS485 adaptor to the meter. This is
how I did the wiring:

![USB-SDM630 wiring](img/wiring.jpg)

You can try to use a cheap USB-RS485 adaptor, or you can [build your own
isolated adaptor](https://github.com/gonium/usb-iso-rs485). I did my
first experiments with a [Digitus USB-RS485
adaptor](http://www.digitus.info/en/products/accessories/adapter-and-converter/r-usb-serial-adapter-usb-20-da-70157/)
which comes with a handy terminal block. I mounted the [bias
network](https://en.wikipedia.org/wiki/RS-485) directly on the terminal
block:

![bias network](img/USB-RS485-Adaptor.jpg)

Since then, I tested various adaptors:

* Supercheap adaptors from China: No ground connection, one worked fine,
	another one was unstable
* Industrial adaptors like the [Meilhaus RedCOM
USB-COMi-SI](https://www.meilhaus.de/usb-comi-si.htm) or the [ADAM
4561](http://www.advantech.com/products/gf-5u7m/adam-4561/mod_92dc04b1-c0fe-4f2b-baf6-5c27e79900c6)
isolate the RS-485 bus from the USB line and work extremely reliable.
But they are really expensive.

I started to develop [my own isolated
adaptor](https://github.com/gonium/usb-iso-rs485). Please check this
link for more information.


## Software installation

### Using the precompiled binaries

You can use the [precompiled releases](https://github.com/gonium/gosdm630/releases) if you like. Just
download the right binary for your platform and unzip.

### Installing the software from source

You need a working [Golang installation](http://golang.org), the [dep
package management tool](https://github.com/golang/dep) and
[Embed](http://github.com/aprice/embed) in order to compile your binary.
Please install the Go compiler first. Then clone this repository:

    git clone https://github.com/gonium/gosdm630.git

If you have ``make`` installed you can use the ``Makefile`` to install the tools:

    $ cd gosdm630
    $ make dep
    Installing embed tool
    Installing dep tool

You can then build the software using the ``Makefile``:

    $ make
    Generating embedded assets
    Generation complete in 109.907612ms
    Building for host platform
    Created binaries:
    sdm

As you can see two sets of binaries are built:

 * ``bin/sdm630_{...}`` is the software built for the host platform
 * ``bin/sdm630_{...}-linux-arm`` is the same for the Raspberry Pi.

If you want to build for all platforms you can use

    $ make release

or, for a single platform like the Raspberry Pi binary, use

    $ GOOS=linux GOARCH=arm GOARM=5 make build


## Running

Now fire up the software:

````
$ ./bin/sdm630 --help
NAME:
   sdm - SDM modbus daemon

USAGE:
   sdm630 [global options] command [command options] [arguments...]

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --serialadapter value, -s value     path to serial RTU device (default: "/dev/ttyUSB0")
   --comset value, -c value            which communication parameter set to use. Valid sets are
                                         1:  2400 baud, 8N1
                                         2:  9600 baud, 8N1
                                         3: 19200 baud, 8N1
                                         4:  2400 baud, 8E1
                                         5:  9600 baud, 8E1
                                         6: 19200 baud, 8E1
                                            (default: 2)
   --device_list value, -d value       MODBUS device type and ID to query, separated by comma.
                                           Valid types are:
                                           "SDM" for Eastron SDM meters
                                           "JANITZA" for Janitza B-Series DIN-Rail meters
                                           "DZG" for the DZG Metering GmbH DVH4013 DIN-Rail meter
                                           Example: -d JANITZA:1,SDM:22,DZG:23 (default: "SDM:1")
   --unique_id_format value, -f value  Unique ID format.
                                           Example: -f Instrument%d
                                           The %d is replaced by the device ID (default: "Instrument%d")
   --verbose, -v                       print verbose messages
   --url value, -u value               the URL the server should respond on (default: ":8080")
   --broker value, -b value            MQTT: The broker URI. ex: tcp://10.10.1.1:1883
   --topic value, -t value             MQTT: The topic name to/from which to publish/subscribe (optional) (default: "sdm630")
   --user value                        MQTT: The User (optional)
   --password value                    MQTT: The password (optional)
   --clientid value, -i value          MQTT: The ClientID (optional) (default: "sdm630")
   --rate value, -r value              MQTT: The maximum update rate (default 0, i.e. unlimited) (after a push we will ignore more data from same device andchannel for this time) (default: 0)
   --clean, -l                         MQTT: Set Clean Session (default false)
   --qos value, -q value               MQTT: The Quality of Service 0,1,2 (default 0) (default: 0)
   --help, -h                          show help
````

A typical invocation looks like this:

    $ ./bin/sdm630 -s /dev/ttyUSB0 -d janitza:26,sdm:1
    2017/01/25 16:34:26 Connecting to RTU via /dev/ttyUSB0
    2017/01/25 16:34:26 Starting API at :8080

This call queries a Janitza B23 meter with ID 26 and an Eastron SDM
meter at ID 1. It . If you use the ``-v`` commandline switch you can see
modbus traffic and the current readings on the command line.  At
[http://localhost:8080](http://localhost:8080) you can see an embedded
web page that updates itself with the latest values:

![realtime view of incoming measurements](img/realtimeview.png)


## Installation on the Raspberry Pi

You simply copy the binary from the ``bin`` subdirectory to the RPi
and start it. I usually put the binary into ``/usr/local/bin`` and
rename it to ``sdm630``. The following sytemd unit can be used to
start the service (put this into ``/etc/systemd/system``):

    [Unit]
    Description=SDM630 via HTTP API
    After=syslog.target
    [Service]
    ExecStart=/usr/local/bin/sdm630 -s /dev/ttyAMA0
    Restart=always
    [Install]
    WantedBy=multi-user.target

You might need to adjust the ``-s`` parameter depending on where your
RS485 adapter is connected. Then, use

    # systemctl start sdm630

to test your installation. If you're satisfied use

    # systemctl enable sdm630

to start the service at boot time automatically.

*WARNING:* If you use an FTDI-based USB-RS485 adaptor you might see the
Raspberry Pi becoming unreachable after a while. This is most likely not
an issue with your RS485-USB adaptor or this software, but because of [a
bug in the Raspberry Pi kernel](https://github.com/raspberrypi/linux/issues/1187).
As mentioned there, add the following parameter to your
``/boot/cmdline.txt``:

    dwc_otg.speed=1

This switches the internal ``dwc`` USB hub of the Raspberry Pi to
USB1.1. While this reduces the available USB speed, the device now works
reliably.


## Detecting connected meters

MODBUS/RTU does not provide a mechanism to discover devices. There is no
reliable way to detect all attached devices. The ``sdm`` tool when used
with the `-detect` option attempts to read the L1 voltage from all valid
device IDs and reports which one replied correctly:

````
./bin/sdm -detect
2017/06/21 10:22:34 Starting bus scan
2017/06/21 10:22:35 Device 1: n/a
...
2017/07/27 16:16:39 Device 21: SDM type device found, L1 voltage: 234.86
2017/07/27 16:16:40 Device 22: n/a
2017/07/27 16:16:40 Device 23: n/a
2017/07/27 16:16:40 Device 24: n/a
2017/07/27 16:16:40 Device 25: n/a
2017/07/27 16:16:40 Device 26: Janitza type device found, L1 voltage: 235.10
...
2017/07/27 16:17:25 Device 247: n/a
2017/07/27 16:17:25 Found 2 active devices:
2017/07/27 16:17:25 * slave address 21: type SDM
2017/07/27 16:17:25 * slave address 26: type JANITZA
2017/07/27 16:17:25 WARNING: This lists only the devices that responded to a known L1 voltage request. Devices with different function code definitions might not be detected.
````


# API

## Rest API

GoSDM provides a convenient REST API. Supported endpoints are:

* `/last/{ID}` current data for device
* `/minuteavg/{ID}` averaged data for device
* `/status` daemon status

Both device APIs can also be called without the device id to return data for all connected devices.

The "GET /last/{ID}"-call simply returns the last measurements of the device with
the Modbus ID {ID}:

````
$ curl localhost:8080/last/11
{
  "Timestamp": "2017-03-27T15:15:09.243729874+02:00",
  "Unix": 1490620509,
  "ModbusDeviceId": 11,
  "Power": {
    "L1": 0,
    "L2": -45.28234100341797,
    "L3": 0
  },
  "Voltage": {
    "L1": 233.1257781982422,
    "L2": 233.12904357910156,
    "L3": 0
  },
  "Current": {
    "L1": 0,
    "L2": 0.19502629339694977,
    "L3": 0
  },
  "Cosphi": {
    "L1": 1,
    "L2": -0.9995147585868835,
    "L3": 1
  },
  "Import": {
    "L1": 0.16599999368190765,
    "L2": 0.10999999940395355,
    "L3": 0.0010000000474974513
  },
  "TotalImport": 0.2770000100135803,
  "Export": {
    "L1": 0,
    "L2": 0.3019999861717224,
    "L3": 0
  },
  "TotalExport": 0.3019999861717224,
  "THD": {
    "VoltageNeutral": {
      "L1": 0,
      "L2": 0,
      "L3": 0
    },
    "AvgVoltageNeutral": 0
  }
}
````

The "GET /minuteavg"-call returns the average measurements over the last
minute:

````
$ curl localhost:8080/minuteavg/11
{
  "Timestamp": "2017-03-27T15:19:06.470316939+02:00",
  "Unix": 1490620746,
  "ModbusDeviceId": 11,
  "Power": {
    "L1": 0,
    "L2": -45.333974165794174,
    "L3": 0
  },
  ...
}
````

### Monitoring

The `/status` endpoint provides the following information:

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
      },
      "ConfiguredMeters": [
        {
          "Id": 26,
          "Type": "JANITZA",
          "Status": "available"
        }
      ]
    }

This is a snapshot of a process running over night, along with the error
statistics during that timeframe. The process queries continuously,
the cabling is not a shielded, twisted wire but something that I had laying
around. With proper cabling the error rate should be lower, though.


## Websocket API

Data read from the meters can be observed by clients in realtime using the Websocket API. As soon as new readings are available, they are pushed to connected websocket clients.

The websocket API is available on `/ws`. All connected clients receive status and
meter updates for all connected meters without further subscription.


## MQTT API

Another option for receiving client updates is by using the built-in MQTT publisher.
By default, readings are published at `/sdm/<unique id>/<reading>`. Rate limiting is possible.


# Changelog

## 0.8 (unreleased)

  - Renamed `sdm630` command to `sdm` which also includes `sdm_detect` now
  - Remove legacy commands. `sdm630_logger`, `sdm630_monitor` and `sdm630_http` are no longer supported. `sdm` is now the single command provided. The legacy commands are still available in the [0.7 version](https://github.com/gonium/gosdm630/releases).
  - Parameters updated:
    - renamed `serialadapter` to `adapter` and allow TCP sockets as well
    - renamed `device_list` to `devices`
  - Add MODBUS over TCP support
  - Support SunSpec-compatible grid inverters

## 0.7

  - Support Saia Burgess Controls ALE3 meters
  - Implement modbus simulation for testing
  - Various improvements of the web UI
  - Support for go 1.11
  - Improved the README
