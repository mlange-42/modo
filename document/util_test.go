package document

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppendNew(t *testing.T) {
	sl1 := make([]int, 0, 32)
	sl1 = append(sl1, 1, 2)

	sl2 := appendNew(sl1, 3, 4)

	assert.Equal(t, []int{1, 2}, sl1)
	assert.Equal(t, []int{1, 2, 3, 4}, sl2)
}
