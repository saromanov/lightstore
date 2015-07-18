package lightstore

import
(
	"strings"
)

//This module for construction data from config

const (
	SNAPSHOTS = "snapshots"
)

//ActionsNamesToFuncs related with option
// 'Every' from configuration
//This module gets names of actions and return functions
func ActionsNamesToFuncs(actions []string)[] func() {
	result := []func(){}
	for _, action := range actions{
		if lower(action) == SNAPSHOTS{
			f := func(){
				snapshot:= NewSnapshotObject()
				if snapshot != nil {

				}
			}

			result = append(result, f)
		}
	}

	return result
}

func lower(str string) string{
	return strings.ToLower(str)
}