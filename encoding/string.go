package encoding

// StringSwapped does not modify the underlying byte slice
func StringSwapped(s []byte) string {
	b := make([]byte, len(s))
	_ = copy(b, s)

	for i := 0; i < len(b); i += 2 {
		c := b[i]
		b[i] = b[i+1]
		b[i+1] = c
	}

	return string(b)
}
