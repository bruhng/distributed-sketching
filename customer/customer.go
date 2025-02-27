package customer

import (
	"bufio"
	"context"
	"fmt"
	"image/color"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	pb "github.com/bruhng/distributed-sketching/proto"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
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

		case "PlotKll":
			if len(words) < 2 {
				fmt.Println("PlotKll requires an int")
				continue
			}

			numBins, err := strconv.Atoi(words[1])
			if err != nil {
				fmt.Println("PlotKll requires an int")
				continue
			}
			res, err := c.PlotKll(ctx, &pb.OrderedValue{Value: &pb.OrderedValue_IntVal{IntVal: int32(numBins)}})
			if err != nil {
				fmt.Println("Could not fetch: ", err)
				continue
			}
			pmf := res.Pmf
			pHist := plot.New()
			pHist.Title.Text = "KLL Sketch Histogram"
			pHist.X.Label.Text = "Value"
			pHist.Y.Label.Text = "Probability Mass"

			bars := make(plotter.Values, numBins)
			labels := make([]string, len(pmf))
			step := math.Round(float64(res.Step))
			for i, v := range pmf {
				bars[i] = float64(v)
				labels[i] = strconv.Itoa(int(step) * i)
			}

			hist, err := plotter.NewBarChart(bars, vg.Points(float64(step)))
			if err != nil {
				panic(err)
			}

			hist.Width = vg.Points(float64(step) * 6.5)
			hist.LineStyle.Width = vg.Points(2)
			hist.LineStyle.Color = color.RGBA{R: 0, B: 0, G: 0, A: 255}
			hist.Color = color.RGBA{R: 135, G: 206, B: 250, A: 255}

			pHist.Add(hist)
			pHist.NominalX(labels...)
			if err := pHist.Save(10*vg.Inch, 5*vg.Inch, "histogram.png"); err != nil {
				panic(err)
			}

			fmt.Println("Histogram saved as histogram.png")

		case "help":
			fmt.Println("ReverseQueryKll [float]")
			fmt.Println("Returns value at quantile [float]\n")
			fmt.Println("QueryKll x")
			fmt.Println("Returns quantlie of value [int]\n")
			fmt.Println("PlotKll [int]")
			fmt.Println("Returns a histogram of the sketch\n")
			fmt.Println("help")
			fmt.Println("Prints Help")

		default:
			continue
		}

	}

}
