package lightstore


//Recording data
type Record struct {
	key, value string
	key_size, value_size uint
	key_address, value_address uint
}