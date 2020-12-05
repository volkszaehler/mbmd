package encoding

// BigEndianUint32Swapped converts bytes to uint32 wrapped as uint64 with swapped word order.
// To use the result as int32 value make sure to convert to uint32 first before converting to int32.
func BigEndianUint32Swapped(b []byte) uint32 {
	_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
	return uint32(b[3])<<16 | uint32(b[2])<<24 | uint32(b[1]) | uint32(b[0])<<8
}
