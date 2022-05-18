package encoding

// StringLsbFirst decodes bytes as string with words in little endian encoding
func StringLsbFirst(s []byte) string {
	b := make([]byte, len(s))
	_ = copy(b, s)

	for i := 0; i < len(b); i += 2 {
		b[i], b[i+1] = b[i+1], b[i]
	}

	return string(b)
}
