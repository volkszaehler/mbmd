## mbmd completion

Generate the autocompletion script for the specified shell

### Synopsis

Generate the autocompletion script for mbmd for the specified shell.
See each sub-command's help for details on how to use the generated script.


### Options inherited from parent commands

```
  -a, --adapter string   Default MODBUS adapter. This option can be used if all devices are attached to a single adapter.
                         Can be either an RTU device (/dev/ttyUSB0) or TCP socket (localhost:502).
                         The default adapter can be overridden per device
  -b, --baudrate int     Serial interface baud rate (default 9600)
      --comset string    Communication parameters for default adapter, either 8N1 or 8E1.
                         Only applicable if the default adapter is an RTU device (default "8N1")
  -c, --config string    Config file (default is $HOME/mbmd.yaml, ./mbmd.yaml, /etc/mbmd.yaml)
  -h, --help             Help for mbmd
      --raw              Log raw device data
      --rtu              Use RTU over TCP for default adapter.
                         Typically used with RS485 to Ethernet adapters that don't perform protocol conversion (e.g. USR-TCP232).
                         Only applicable if the default adapter is a TCP connection
  -v, --verbose          Verbose mode
```

### SEE ALSO

* [mbmd](mbmd.md)	 - ModBus Measurement Daemon
* [mbmd completion bash](mbmd_completion_bash.md)	 - Generate the autocompletion script for bash
* [mbmd completion fish](mbmd_completion_fish.md)	 - Generate the autocompletion script for fish
* [mbmd completion powershell](mbmd_completion_powershell.md)	 - Generate the autocompletion script for powershell
* [mbmd completion zsh](mbmd_completion_zsh.md)	 - Generate the autocompletion script for zsh

