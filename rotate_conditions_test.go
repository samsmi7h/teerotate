package teerotate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockSizer struct {
	size ByteSize
}

func (m mockSizer) SizeInBytes() ByteSize {
	return m.size
}

func TestRotateConditionFactory(t *testing.T) {
	t.Run("nothing reached -> no rotate", func(t *testing.T) {
		factory := rotateConditionFactory(Opts{
			MinimumLifespan: time.Hour,
			MaximumLifespan: time.Hour * 2,
			MinimumByteSize: Kilobyte,
		})

		check := factory()
		result := check(mockSizer{Byte})
		assert.False(t, result)
	})

	t.Run("min time has past BUT not min size -> no rotate", func(t *testing.T) {
		factory := rotateConditionFactory(Opts{
			MinimumLifespan: -time.Hour,
			MaximumLifespan: time.Hour * 2,
			MinimumByteSize: Kilobyte,
		})

		check := factory()
		result := check(mockSizer{Byte})
		assert.False(t, result)
	})

	t.Run("min size reched BUT not min time -> no rotate", func(t *testing.T) {
		factory := rotateConditionFactory(Opts{
			MinimumLifespan: time.Hour,
			MaximumLifespan: time.Hour * 2,
			MinimumByteSize: Kilobyte,
		})

		check := factory()
		result := check(mockSizer{Megabyte})
		assert.False(t, result)
	})

	t.Run("min time AND min size reached -> yes rotate", func(t *testing.T) {
		factory := rotateConditionFactory(Opts{
			MinimumLifespan: -time.Hour,
			MaximumLifespan: time.Hour * 2,
			MinimumByteSize: Kilobyte,
		})

		check := factory()
		result := check(mockSizer{Megabyte})
		assert.True(t, result)
	})

	t.Run("max time reached BUT not min size -> yes rotate", func(t *testing.T) {
		factory := rotateConditionFactory(Opts{
			MinimumLifespan: -time.Hour * 2,
			MaximumLifespan: -time.Hour,
			MinimumByteSize: Kilobyte,
		})

		check := factory()
		result := check(mockSizer{Byte})
		assert.True(t, result)
	})
}
