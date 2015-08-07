package distributed

import
(
	"fmt"
)

type Distributed struct {
	nodes  []*Node
}

func New() *Distributed {
	distr := new(Distributed)
	distr.nodes = []*Node{}
	return distr
}

//Join provides connect new node to the cluster
func (distr*Distributed) Join(node *Node) {
	fmt.Println(fmt.Sprintf("Join %s to the cluster", node.addr))
}


//Start new leader election
func (distr* Distributed) elect (){

}
