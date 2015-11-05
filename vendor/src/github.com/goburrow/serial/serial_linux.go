package serial

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"time"
	"unsafe"
)

var baudRates = map[int]uint32{
	50:      syscall.B50,
	75:      syscall.B75,
	110:     syscall.B110,
	134:     syscall.B134,
	150:     syscall.B150,
	200:     syscall.B200,
	300:     syscall.B300,
	600:     syscall.B600,
	1200:    syscall.B1200,
	1800:    syscall.B1800,
	2400:    syscall.B2400,
	4800:    syscall.B4800,
	9600:    syscall.B9600,
	19200:   syscall.B19200,
	38400:   syscall.B38400,
	57600:   syscall.B57600,
	115200:  syscall.B115200,
	230400:  syscall.B230400,
	460800:  syscall.B460800,
	500000:  syscall.B500000,
	576000:  syscall.B576000,
	921600:  syscall.B921600,
	1000000: syscall.B1000000,
	1152000: syscall.B1152000,
	1500000: syscall.B1500000,
	2000000: syscall.B2000000,
	2500000: syscall.B2500000,
	3000000: syscall.B3000000,
	3500000: syscall.B3500000,
	4000000: syscall.B4000000,
}

var charSizes = map[int]uint32{
	5: syscall.CS5,
	6: syscall.CS6,
	7: syscall.CS7,
	8: syscall.CS8,
}

// port implements Port interface.
type port struct {
	// Should use fd directly by using syscall.Open() ?
	file       *os.File
	oldTermios *syscall.Termios

	timeout time.Duration
}

// New allocates and returns a new serial port controller.
func New() Port {
	return &port{}
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
	p.file, err = os.OpenFile(c.Address, syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_NDELAY, os.FileMode(0666))
	if err != nil {
		return
	}
	// Backup current termios to restore on closing.
	p.backupTermios()
	if err = p.setTermios(termios); err != nil {
		p.file.Close()
		p.file = nil
		p.oldTermios = nil
		return
	}
	p.timeout = c.Timeout
	return
}

func (p *port) Close() (err error) {
	if p.file == nil {
		return
	}
	p.restoreTermios()
	err = p.file.Close()
	p.file = nil
	p.oldTermios = nil
	return
}

// Read reads from serial port. Port must be opened before calling this method.
// It is blocked until all data received or timeout after p.timeout.
func (p *port) Read(b []byte) (n int, err error) {
	var rfds syscall.FdSet

	fd := int(p.file.Fd())
	fdSet(fd, &rfds)

	var tv *syscall.Timeval
	if p.timeout > 0 {
		timeout := syscall.NsecToTimeval(p.timeout.Nanoseconds())
		tv = &timeout
	}
	for {
		// If syscall.Select() returns EINTR (Interrupted system call), retry it
		if _, err = syscall.Select(fd+1, &rfds, nil, nil, tv); err == nil {
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
	n, err = p.file.Read(b)
	return
}

// Write writes data to the serial port.
func (p *port) Write(b []byte) (n int, err error) {
	n, err = p.file.Write(b)
	return
}

func (p *port) setTermios(termios *syscall.Termios) (err error) {
	if err = tcsetattr(int(p.file.Fd()), termios); err != nil {
		err = fmt.Errorf("serial: could not set setting: %v", err)
	}
	return
}

// backupTermios saves current termios setting.
// Make sure that device file has been opened before calling this function.
func (p *port) backupTermios() {
	oldTermios := &syscall.Termios{}
	if err := tcgetattr(int(p.file.Fd()), oldTermios); err != nil {
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
	if err := tcsetattr(int(p.file.Fd()), p.oldTermios); err != nil {
		// Warning only.
		log.Printf("serial: could not restore setting: %v\n", err)
		return
	}
	p.oldTermios = nil
}

// Helpers for termios

func newTermios(c *Config) (termios *syscall.Termios, err error) {
	termios = &syscall.Termios{}
	var flag uint32
	// Baud rate
	if c.BaudRate == 0 {
		// 19200 is the required default.
		flag = syscall.B19200
	} else {
		flag = baudRates[c.BaudRate]
		if flag == 0 {
			err = fmt.Errorf("serial: unsupported baud rate %v", c.BaudRate)
			return
		}
	}
	// Input baud.
	termios.Ispeed = flag
	// Output baud.
	termios.Ospeed = flag
	// Character size.
	if c.DataBits == 0 {
		flag = syscall.CS8
	} else {
		flag = charSizes[c.DataBits]
		if flag == 0 {
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

// tcsetattr sets terminal file descriptor parameters.
// See man tcsetattr(3).
func tcsetattr(fd int, termios *syscall.Termios) (err error) {
	r, _, errno := syscall.Syscall(uintptr(syscall.SYS_IOCTL),
		uintptr(fd), uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(termios)))
	if errno != 0 {
		err = errno
		return
	}
	if r != 0 {
		err = fmt.Errorf("tcsetattr failed %v", r)
	}
	return
}

// tcgetattr gets terminal file descriptor parameters.
// See man tcgetattr(3).
func tcgetattr(fd int, termios *syscall.Termios) (err error) {
	r, _, errno := syscall.Syscall(uintptr(syscall.SYS_IOCTL),
		uintptr(fd), uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(termios)))
	if errno != 0 {
		err = errno
		return
	}
	if r != 0 {
		err = fmt.Errorf("tcgetattr failed %v", r)
		return
	}
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
