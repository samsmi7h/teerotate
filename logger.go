package teerotate

// tmpLogWriter is short lived & immutable
// it corresponds to a log file
// and is rotated out
type tmpLogWriter struct {
	w                    WriteCloseSizer
	rotateConditionCheck rotateConditionCheck
}

func newTmpLogger(w WriteCloseSizer, rotateConditionCheck rotateConditionCheck) *tmpLogWriter {
	return &tmpLogWriter{
		w:                    w,
		rotateConditionCheck: rotateConditionCheck,
	}
}

func (l tmpLogWriter) isDone() bool {
	return l.rotateConditionCheck(l.w)
}

func (l tmpLogWriter) Write(b []byte) (done bool) {
	l.w.Write(b)
	return l.isDone()
}

func (l tmpLogWriter) Close() {
	l.w.Close()
}
