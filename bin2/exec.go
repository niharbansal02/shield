package main

import "fmt"

func main() {
	if err := startTestGRPCServer(grpcBackendPort, &greetServer{}); err != nil {
		fmt.Println(err)
	}
}

