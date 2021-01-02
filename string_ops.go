package blocktemplates

func trimEndWhitespace(b []byte) ([]byte, bool) {
	clean := true
	ri := len(b) - 1
	for ; ri >= 0; ri-- {
		if b[ri] == ' ' || b[ri] == '\t' {
			continue
		} else if b[ri] == '\n' {
			break
		} else {
			clean = false
			break
		}
	}
	return b[:ri], clean
}

func trimStartWhitespace(b []byte) ([]byte, bool) {
	clean := true
	ri := 0
	for ; ri <= len(b) - 1; ri++ {
		if b[ri] == ' ' || b[ri] == '\t' {
			continue
		} else if b[ri] == '\n' {
			break
		} else {
			clean = false
			break
		}
	}
	return b[ri:], clean
}
