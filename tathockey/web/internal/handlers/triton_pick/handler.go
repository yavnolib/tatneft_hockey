package grpc

import (
	"context"
	"google.golang.org/grpc"
	"time"

	grpcServ "tat_hockey_pack/api/grpc-client"
)

func ProcessWithTriton(videoName string) ([]uint64, []int32, error) {
	conn, err := grpc.Dial("localhost:8001", grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	defer conn.Close()

	client := grpcServ.NewGRPCInferenceServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Формирование запроса в Triton
	request := &grpcServ.ModelInferRequest{
		ModelName: "resnet18",
		Parameters: map[string]*grpcServ.InferParameter{
			"video_name": {
				ParameterChoice: &grpcServ.InferParameter_StringParam{
					StringParam: videoName,
				},
			},
		},
	}

	// Отправка запроса
	response, err := client.ModelInfer(ctx, request)
	if err != nil {
		return nil, nil, err
	}

	// Обработка ответа Triton
	secondsTensor := response.Outputs[0].Contents.Uint64Contents
	intsTensor := response.Outputs[1].Contents.IntContents

	return secondsTensor, intsTensor, nil
}
