package input

// inputTimeSubtickBits is the number of bits for a counter in a tick.
// An input time consists of a tick and a counter in a tick.
// This means that an input time will be invalid when 2^20 = 1048576 inputs are handled in a tick,
// but this should unlikely happen.
const inputTimeSubtickBits = 20

type InputTime int64

func NewInputTimeFromTick(tick int64, subtick int64) InputTime {
	return InputTime((tick << inputTimeSubtickBits) | int64(subtick&((1<<inputTimeSubtickBits)-1)))
}

func (i InputTime) Tick() int64 {
	return int64(i >> inputTimeSubtickBits)
}

func (i InputTime) Subtick() int64 {
	return int64(i & ((1 << inputTimeSubtickBits) - 1))
}
