package customer

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	pb "github.com/bruhng/distributed-sketching/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Init(port string, adr string) {
	conn, err := grpc.NewClient(adr+":"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		panic("Could not connect to server")
	}
	defer conn.Close()
	c := pb.NewSketcherClient(conn)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Write help for help")
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			panic("could not read string")
		}
		input = strings.TrimSpace(input)
		words := strings.Split(input, " ")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		switch words[0] {
		case "QueryKll":
			if len(words) < 2 {
				fmt.Println("QueryKll requires an int")
				continue
			}
			x, err := strconv.Atoi(words[1])
			if err != nil {
				fmt.Println("QueryKll requires an int")
				continue
			}
			res, err := c.QueryKll(ctx, &pb.OrderedValue{Value: &pb.OrderedValue_IntVal{IntVal: int32(x)}})
			if err != nil {
				fmt.Println("Could not fetch: ", err)
			}
			fmt.Println(res)

		case "ReverseQueryKll":
			if len(words) < 2 {
				fmt.Println("ReverseQueryKll requires an float")
				continue
			}
			x, err := strconv.ParseFloat(words[1], 32)
			if err != nil {
				fmt.Println("ReverseQueryKll requires an float")
				continue
			}
			res, err := c.ReverseQueryKll(ctx, &pb.ReverseQuery{Phi: float32(x)})

			if err != nil {
				fmt.Println("Could not fetch: ", err)
			}
			fmt.Println(res)

		case "help":
			fmt.Println("ReverseQueryKll [float]")
			fmt.Println("Returns value at quantile [float]\n")
			fmt.Println("QueryKll x")
			fmt.Println("Returns quantlie of value [int]\n")
			fmt.Println("help")
			fmt.Println("Prints Help")

		default:
			continue
		}

	}

}
