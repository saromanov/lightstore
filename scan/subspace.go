package scan

//Subspace needs for range over key-value pairs, but not over all keys, but only for
//part of this keys.

type Subspace struct {
	title string
	keys []string
}

func CreateSubspace(title string) *Subspace{
	ss := new(Subspace)
	ss.title = title
	return ss
}