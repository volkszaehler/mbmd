// +build darwin linux

package serial

import (
	"fmt"
	"log"
	"syscall"
	"time"
)

// port implements Port interface.
type port struct {
	fd         int
	oldTermios *syscall.Termios

	timeout time.Duration
}

// New allocates and returns a new serial port controller.
func New() Port {
	return &port{fd: -1}
}

// Open connects to the given serial port.
func (p *port) Open(c *Config) (err error) {
	termios, err := newTermios(c)
	if err != nil {
		return
	}
	// See man termios(3).
	// O_NOCTTY: no controlling terminal.
	// O_NDELAY: no data carrier detect.
	p.fd, err = syscall.Open(c.Address, syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_NDELAY|syscall.O_CLOEXEC, 0666)
	if err != nil {
		return
	}
	// Backup current termios to restore on closing.
	p.backupTermios()
	if err = p.setTermios(termios); err != nil {
		syscall.Close(p.fd)
		p.fd = -1
		p.oldTermios = nil
		return
	}
	p.timeout = c.Timeout
	return
}

func (p *port) Close() (err error) {
	if p.fd == -1 {
		return
	}
	p.restoreTermios()
	err = syscall.Close(p.fd)
	p.fd = -1
	p.oldTermios = nil
	return
}

// Read reads from serial port. Port must be opened before calling this method.
// It is blocked until all data received or timeout after p.timeout.
func (p *port) Read(b []byte) (n int, err error) {
	var rfds syscall.FdSet

	fd := p.fd
	fdSet(fd, &rfds)

	var tv *syscall.Timeval
	if p.timeout > 0 {
		timeout := syscall.NsecToTimeval(p.timeout.Nanoseconds())
		tv = &timeout
	}
	for {
		// If syscall.Select() returns EINTR (Interrupted system call), retry it
		if err = syscallSelect(fd+1, &rfds, nil, nil, tv); err == nil {
			break
		}
		if err != syscall.EINTR {
			err = fmt.Errorf("serial: could not select: %v", err)
			return
		}
	}
	if !fdIsSet(fd, &rfds) {
		// Timeout
		err = ErrTimeout
		return
	}
	n, err = syscall.Read(fd, b)
	return
}

// Write writes data to the serial port.
func (p *port) Write(b []byte) (n int, err error) {
	n, err = syscall.Write(p.fd, b)
	return
}

func (p *port) setTermios(termios *syscall.Termios) (err error) {
	if err = tcsetattr(p.fd, termios); err != nil {
		err = fmt.Errorf("serial: could not set setting: %v", err)
	}
	return
}

// backupTermios saves current termios setting.
// Make sure that device file has been opened before calling this function.
func (p *port) backupTermios() {
	oldTermios := &syscall.Termios{}
	if err := tcgetattr(p.fd, oldTermios); err != nil {
		// Warning only.
		log.Printf("serial: could not get setting: %v\n", err)
		return
	}
	// Will be reloaded when closing.
	p.oldTermios = oldTermios
}

// restoreTermios restores backed up termios setting.
// Make sure that device file has been opened before calling this function.
func (p *port) restoreTermios() {
	if p.oldTermios == nil {
		return
	}
	if err := tcsetattr(p.fd, p.oldTermios); err != nil {
		// Warning only.
		log.Printf("serial: could not restore setting: %v\n", err)
		return
	}
	p.oldTermios = nil
}

// Helpers for termios

func newTermios(c *Config) (termios *syscall.Termios, err error) {
	termios = &syscall.Termios{}
	flag := termios.Cflag
	// Baud rate
	if c.BaudRate == 0 {
		// 19200 is the required default.
		flag = syscall.B19200
	} else {
		var ok bool
		flag, ok = baudRates[c.BaudRate]
		if !ok {
			err = fmt.Errorf("serial: unsupported baud rate %v", c.BaudRate)
			return
		}
	}
	termios.Cflag |= flag
	// Input baud.
	termios.Ispeed = flag
	// Output baud.
	termios.Ospeed = flag
	// Character size.
	if c.DataBits == 0 {
		flag = syscall.CS8
	} else {
		var ok bool
		flag, ok = charSizes[c.DataBits]
		if !ok {
			err = fmt.Errorf("serial: unsupported character size %v", c.DataBits)
			return
		}
	}
	termios.Cflag |= flag
	// Stop bits
	switch c.StopBits {
	case 0, 1:
		// Default is one stop bit.
		// noop
	case 2:
		// CSTOPB: Set two stop bits.
		termios.Cflag |= syscall.CSTOPB
	default:
		err = fmt.Errorf("serial: unsupported stop bits %v", c.StopBits)
		return
	}
	switch c.Parity {
	case "N":
		// noop
	case "O":
		// PARODD: Parity is odd.
		termios.Cflag |= syscall.PARODD
		fallthrough
	case "", "E":
		// As mentioned in the modbus spec, the default parity mode must be Even parity
		// PARENB: Enable parity generation on output.
		termios.Cflag |= syscall.PARENB
		// INPCK: Enable input parity checking.
		termios.Iflag |= syscall.INPCK
	default:
		err = fmt.Errorf("serial: unsupported parity %v", c.Parity)
		return
	}
	// Control modes.
	// CREAD: Enable receiver.
	// CLOCAL: Ignore control lines.
	termios.Cflag |= syscall.CREAD | syscall.CLOCAL
	// Special characters.
	// VMIN: Minimum number of characters for noncanonical read.
	// VTIME: Time in deciseconds for noncanonical read.
	// Both are unused as NDELAY is we utilized when opening device.
	return
}

// fdGet returns index and offset of fd in fds.
func fdGet(fd int, fds *syscall.FdSet) (index, offset int) {
	index = fd / (syscall.FD_SETSIZE / len(fds.Bits)) % len(fds.Bits)
	offset = fd % (syscall.FD_SETSIZE / len(fds.Bits))
	return
}

// fdSet implements FD_SET macro.
func fdSet(fd int, fds *syscall.FdSet) {
	idx, pos := fdGet(fd, fds)
	fds.Bits[idx] = 1 << uint(pos)
}

// fdIsSet implements FD_ISSET macro.
func fdIsSet(fd int, fds *syscall.FdSet) bool {
	idx, pos := fdGet(fd, fds)
	return fds.Bits[idx]&(1<<uint(pos)) != 0
}
