## mbmd run

Read and publish measurements from all configured devices

### Synopsis

Read and publish measurements from all configured devices

```
mbmd run [flags]
```

### Options

```
      --api string                   REST API url. Use 127.0.0.1:8080 to limit to localhost. (default "0.0.0.0:8080")
  -d, --devices strings              MODBUS device type and ID to query, multiple devices separated by comma or by repeating the flag.
                                       Example: -d SDM:1,SDM:2 -d DZG:1.
                                     Valid types are:
                                       RTU
                                         ABB      ABB A/B-Series meters
                                         DZG      DZG Metering GmbH DVH4013 meters
                                         IEM3000  Schneider Electric iEM3000 series
                                         INEPRO   Inepro Metering Pro 380
                                         JANITZA  Janitza B-Series meters
                                         MPM      Bernecker Engineering MPM3PM meters
                                         ORNO1P   ORNO WE-514 & WE-515
                                         ORNO3P   ORNO WE-516 & WE-517
                                         SBC      Saia Burgess Controls ALE3 meters
                                         SDM      Eastron SDM630
                                         SDM220   Eastron SDM220
                                         SDM230   Eastron SDM230
                                         SDM72    Eastron SDM72
                                       TCP
                                         SUNS     Sunspec-compatible MODBUS TCP device (SMA, SolarEdge, KOSTAL, etc)
                                     To use an adapter different from default, append RTU device or TCP address separated by @.
                                     If the adapter is a TCP connection (identified by :port), the device type (SUNS) is ignored and
                                     any type is considered valid.
                                       Example: -d SDM:1@/dev/USB11 -d SMA:126@localhost:502
      --influx-database string       InfluxDB database
      --influx-measurement string    InfluxDB measurement (default "data")
      --influx-organization string   InfluxDB organization
      --influx-password string       InfluxDB password (optional)
      --influx-token string          InfluxDB token (optional)
  -i, --influx-url string            InfluxDB URL. ex: http://10.10.1.1:8086
      --influx-user string           InfluxDB user (optional)
  -m, --mqtt-broker string           MQTT broker URI. ex: tcp://10.10.1.1:1883
      --mqtt-clientid string         MQTT client id (default "mbmd")
      --mqtt-homie string            MQTT Homie IoT discovery base topic (homieiot.github.io). Set empty to disable. (default "homie")
      --mqtt-password string         MQTT password (optional)
      --mqtt-qos int                 MQTT quality of service 0,1,2 (default 0)
      --mqtt-topic string            MQTT root topic. Set empty to disable publishing. (default "mbmd")
      --mqtt-user string             MQTT user (optional)
  -r, --rate duration                Rate limit. Devices will not be queried more often than rate limit. (default 1s)
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

