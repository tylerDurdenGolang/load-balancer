package grpc

import (
	"context"
	"template-api/internal/models"
	"template-api/internal/service"
	pb "template-api/pkg/proto"
)

type Handler struct {
	pb.UnimplementedItemServiceServer
	service service.IService
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) CreateItem(ctx context.Context, req *pb.CreateItemRequest) (*pb.CreateItemResponse, error) {
	id, err := h.service.CreateItem(ctx, req.Name, req.Description, req.Price, req.Stock)
	if err != nil {
		return nil, err
	}
	return &pb.CreateItemResponse{Id: id}, nil
}

func (h *Handler) GetAllItems(ctx context.Context, req *pb.GetAllItemsRequest) (*pb.GetAllItemsResponse, error) {
	items, err := h.service.GetAllItems(ctx, req.Limit)
	if err != nil {
		return nil, err
	}

	pbItems := make([]*pb.Item, 0, len(items))
	for _, i := range items {
		pbItems = append(pbItems, &pb.Item{
			Id:          i.ID,
			Name:        i.Name,
			Description: i.Description,
			Price:       i.Price,
			Stock:       i.Stock,
		})
	}

	return &pb.GetAllItemsResponse{Items: pbItems}, nil
}

func (h *Handler) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	item, err := h.service.GetItemById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetItemResponse{
		Item: &pb.Item{
			Id:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			Stock:       item.Stock,
		},
	}, nil
}

func (h *Handler) UpdateItem(ctx context.Context, req *pb.UpdateItemRequest) (*pb.EmptyResponse, error) {
	newItem := &models.UpdateItem{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}
	if err := h.service.UpdateItem(ctx, newItem); err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}

func (h *Handler) DeleteItem(ctx context.Context, req *pb.DeleteItemRequest) (*pb.EmptyResponse, error) {
	if err := h.service.DeleteItem(ctx, req.Id); err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}
