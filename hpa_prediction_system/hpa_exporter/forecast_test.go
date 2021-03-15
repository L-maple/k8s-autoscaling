package main

import (
	"context"
	"fmt"
	pb "github.com/k8s-autoscaling/hpa_prediction_system/time_series_forecast"
	"google.golang.org/grpc"
	"log"
	"testing"
)

var (
	serverAddress string = "localhost:50000"
)

func getForecastServiceClient() (pb.ForecastServiceClient, *grpc.ClientConn) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := pb.NewForecastServiceClient(conn)

	return client, conn
}


func TestForecastClient(t *testing.T) {
	client, conn := getForecastServiceClient()
	defer conn.Close()

	response, err := client.GetForeCastValue(context.TODO(), &pb.ForecastRequest{
		Data: "hello",
		Minutes: 1,
	})
	if err != nil {
		log.Fatal("error: ", err)
	}

	fmt.Println("value: ", response.Value)
}
