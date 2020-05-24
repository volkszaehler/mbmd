## mbmd read

Read register (EXPERIMENTAL)

### Synopsis

Read reads a single register (input, holding, coil, discrete input)
and will return it according to defined format. Read will ignore the config file
and requires adapter configuration using command line.

```
mbmd read [flags] register length
```

### Options

```
  -d, --device string     MODBUS device ID to query. Only single device allowed. (default "1")
  -e, --encoding string   Data encoding: bit|int|uint|int32s|uint32s|hex|float|floats|string (default "int")
  -t, --type string       Register type to read: holding|input|coil|discrete (default "holding")
```

### Options inherited from parent commands

```
  -a, --adapter string   Default MODBUS adapter. This option can be used if all devices are attached to a single adapter.
                         Can be either an RTU device (/dev/ttyUSB0) or TCP socket (localhost:502).
                         The default adapter can be overridden per device
  -b, --baudrate int     Serial interface baud rate (default 9600)
      --comset string    Communication parameters for default adapter, either 8N1 or 8E1.
                         Only applicable if the default adapter is an RTU device (default "8N1")
  -c, --config string    Config file (default is $HOME/mbmd.yaml)
  -h, --help             Help for mbmd
      --raw              Log raw device data
      --rtu              Use RTU over TCP for default adapter.
                         Typically used with RS485 to Ethernet adapters that don't perform protocol conversion (e.g. USR-TCP232).
                         Only applicable if the default adapter is a TCP connection
  -v, --verbose          Verbose mode
```

### SEE ALSO

* [mbmd](mbmd.md)	 - ModBus Measurement Daemon

