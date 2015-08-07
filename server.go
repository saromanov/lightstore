package lightstore

import
(
	"net"
	"bufio"
	"fmt"
)

type Server struct {
	addr     string
	listener net.Listener
	store    *Store
}


//NewServer provides initialization of single node(server)
func NewServer(addr, typestore string) *Server {
	conn, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	serv := new (Server)
	serv.addr = addr
	serv.listener = conn
	serv.store = InitStore(Settings{Innerdata: typestore})
	return serv
}

func (serv*Server) RunServer() {
	defer serv.listener.Close()
	go func() {
		for {
			conn, err := serv.listener.Accept()
			if err != nil {
				//
			}

			go handleRequest(conn)
		}
	}()
}

func handleRequest(conn net.Conn) {
	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		panic(err)
	}

	fmt.Println(status)
}