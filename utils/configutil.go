package utils

import
(
	"strings"
	"../snapshot"
)

//This module for construction data from config

const (
	SNAPSHOTS = "snapshots"
	HISTORY_LIMIT = "history_limit"
)

//ActionsNamesToFuncs related with option
// 'Every' from configuration
//This module gets names of actions and return functions
func ActionsNamesToFuncs(actions []string)[] func() {
	result := []func(){}
	for _, action := range actions{
		if lower(action) == SNAPSHOTS{
			f := func(){
				snap:= snapshot.NewSnapshotObject()
				if snap != nil {

				}
			}

			result = append(result, f)
		}

		if lower(action) == HISTORY_LIMIT {
			
		}
	}

	return result
}

func lower(str string) string{
	return strings.ToLower(str)
}