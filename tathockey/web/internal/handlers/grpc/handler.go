package grpc

import (
	"context"
	"google.golang.org/grpc"
	"time"

	grpc_service "tat_hockey_pack/api/grpc-client"
)

func ProcessWithTriton(videoName string) ([]float32, []int32, error) {
	conn, err := grpc.Dial("localhost:8001", grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	defer conn.Close()

	client := grpc_service.NewGRPCInferenceServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Формирование запроса в Triton
	request := &grpc_service.ModelInferRequest{
		ModelName: "resnet18",
		Inputs: []*grpc_service.ModelInferRequest_InferInputTensor{
			{
				Name:     "video_name",
				Datatype: "BYTES",
				Shape:    []int64{1},
				Contents: &grpc_service.InferTensorContents{
					BytesContents: [][]byte{[]byte(videoName)},
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
	secondsTensor := response.Outputs[0].Contents.Fp32Contents
	intsTensor := response.Outputs[1].Contents.IntContents

	return secondsTensor, intsTensor, nil
}
