// +build ignore

package serial

// #include <windows.h>
import "C"

const (
	c_MAXDWORD    = C.MAXDWORD
	c_ONESTOPBIT  = C.ONESTOPBIT
	c_TWOSTOPBITS = C.TWOSTOPBITS
	c_EVENPARITY  = C.EVENPARITY
	c_ODDPARITY   = C.ODDPARITY
	c_NOPARITY    = C.NOPARITY
)

type c_COMMTIMEOUTS C.COMMTIMEOUTS

type c_DCB C.DCB

func toDWORD(val int) C.DWORD {
	return C.DWORD(val)
}

func toBYTE(val int) C.BYTE {
	return C.BYTE(val)
}
