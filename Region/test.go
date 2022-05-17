package Region

import (
	"fmt"
	"net"
)

func main() {
	_, err := net.Dial("tcp", "127.0.0.1:3037")
	if err != nil {
		fmt.Println("net.Dial failed")
	}
}
