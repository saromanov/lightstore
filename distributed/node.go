package distributed

import
(
	"time"
)

type Node struct {
	addr string
	knownnodes []*Node
	server     *Server
}