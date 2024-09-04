package teerotate

import "time"

type ByteSize int64

const (
	Byte     ByteSize = 1
	Kilobyte ByteSize = Byte * 1_000
	Megabyte ByteSize = Kilobyte * 1_000
	Gigabyte ByteSize = Megabyte * 1_000
)

type makeRotateConditionCheck func() rotateConditionCheck
type rotateConditionCheck func(f Sizer) (done bool)

func rotateConditionFactory(opts Opts) makeRotateConditionCheck {
	if opts.MinimumLifespan == 0 {
		panic("MinimumLifespan can't be 0")
	}

	return func() rotateConditionCheck {
		start := time.Now()
		minEnd := start.Add(opts.MinimumLifespan)
		maxEnd := start.Add(opts.MaximumLifespan)
		return func(f Sizer) bool {
			// must breach minLifespan AND minFileSize
			// OR
			// maxLifespan

			now := time.Now()

			if now.After(maxEnd) {
				return true
			}

			return time.Now().After(minEnd) && f.SizeInBytes() >= opts.MinimumByteSize
		}
	}
}
