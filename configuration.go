package lightstore
import
(
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"fmt"
)


type Config struct {
	//Address for server
	address string
	//port for server
	port uint
	//Directory for storing log data
	logdir string
	//Directory for storing cache data
	cachedir string
	//Diretcory for storing db data
	dbdir string
	//Cluster name
	cluster string
	//Sync for commits
	commitsync bool
	//Size for cache
	cachesize int
}

func LoadConfigData(path string){
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	conf := Config{}
	yamlerr := yaml.Unmarshal([]byte(data), &conf)
	if yamlerr != nil {
		panic(yamlerr)
	}
}