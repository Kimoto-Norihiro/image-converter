package batchservice

import (
	"context"
	"errors"

	batchservicepb "github.com/Kimoto-Norihiro/image-converter/pkg/grpc"
)

type BatchService struct {
	batchservicepb.UnimplementedBatchServiceServer
}

func NewBatchService() *BatchService {
	return &BatchService{}
}

func (s *BatchService) Convert(ctx context.Context, req *batchservicepb.ConvertRequest) (*batchservicepb.ConvertResponse, error) {
	return nil, errors.New("not implemented")
}