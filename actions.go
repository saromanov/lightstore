package lightstore
import
(
	"strings"
)
//Actions provides state of DB

//ReadyToSet provides checking of kv before set
type ReadyToSet struct {
	kvitems KVITEM
	syskeys KVITEM
	ready bool
}

func NewReadyToSet(rawitems KVITEM) *ReadyToSet {
	rts := new(ReadyToSet)
	rts.syskeys = kvItems(rawitems, func(str string) bool {
		return checkSystemKeys(str)
	})
	rts.kvitems = kvItems(rawitems, func(str string)bool {
		return !checkSystemKeys(str)
	})
	return rts
}

func kvItems(rawitems KVITEM, pred func(string) bool) KVITEM {
	result := KVITEM{}
	for key, value := range rawitems {
		if pred(key) {
			result[key] = value
		}
	}
	return result
}

//System keys starts with _ (_index, for example)
func checkSystemKeys(key string) bool {
	return strings.HasPrefix(key, "_")
}
