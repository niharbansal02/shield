package main

type grpcClient struct {

}

//func main() {
//	var conn *grpc.ClientConn
//	conn, err := grpc.Dial(":5556", grpc.WithInsecure())
//	if err != nil {
//		fmt.Println("did not connect: %s", err)
//	}
//	defer conn.Close()
//	c := helloworld.NewGreeterClient(conn)
//	response, err := c.SayHello(context.Background(), &helloworld.HelloRequest{Name: "1"})
//	if err != nil {
//		fmt.Println("Error when calling SayHello: %s", err)
//	}
//	fmt.Println("response from server: %s", response.Message)
//}
