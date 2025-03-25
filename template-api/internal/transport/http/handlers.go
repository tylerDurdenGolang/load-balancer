package http

import (
	"io"
	"net/http"

	"template-api/internal/models"
	pb "template-api/pkg/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

func (h *Handler) createItem(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			return
		}

		st := status.New(codes.Unknown, err.Error())
		h.writeProtoResponse(c, http.StatusInternalServerError, st.Proto())
		return
	}
	defer c.Request.Body.Close()

	req := pb.CreateItemRequest{}
	if err = protojson.Unmarshal(body, &req); err != nil {
		st := status.New(codes.InvalidArgument, err.Error())
		h.writeProtoResponse(c, http.StatusBadRequest, st.Proto())
		return
	}

	id, err := h.service.CreateItem(
		c.Request.Context(),
		req.Name,
		req.Description,
		req.Price,
		req.Stock,
	)

	if err != nil {
		st := status.New(codes.Unknown, err.Error())
		h.writeProtoResponse(c, http.StatusInternalServerError, st.Proto())
		return
	}

	h.writeProtoResponse(
		c,
		http.StatusOK,
		&pb.CreateItemResponse{
			Id: id,
		},
	)
}

func (h *Handler) getAllItems(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			return
		}

		st := status.New(codes.Unknown, err.Error())
		h.writeProtoResponse(c, http.StatusInternalServerError, st.Proto())
		return
	}
	defer c.Request.Body.Close()

	req := pb.GetAllItemsRequest{}
	if err = protojson.Unmarshal(body, &req); err != nil {
		st := status.New(codes.InvalidArgument, err.Error())
		h.writeProtoResponse(c, http.StatusBadRequest, st.Proto())
		return
	}

	items, err := h.service.GetAllItems(c.Request.Context(), req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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

	h.writeProtoResponse(
		c,
		http.StatusOK,
		&pb.GetAllItemsResponse{
			Items: pbItems,
		},
	)
}

func (h *Handler) getItemById(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			return
		}

		st := status.New(codes.Unknown, err.Error())
		h.writeProtoResponse(c, http.StatusInternalServerError, st.Proto())
		return
	}
	defer c.Request.Body.Close()

	req := pb.GetItemRequest{}
	if err = protojson.Unmarshal(body, &req); err != nil {
		st := status.New(codes.InvalidArgument, err.Error())
		h.writeProtoResponse(c, http.StatusBadRequest, st.Proto())
		return
	}

	item, err := h.service.GetItemById(c.Request.Context(), req.Id)
	if err != nil {
		st := status.New(codes.Unknown, err.Error())
		h.writeProtoResponse(c, http.StatusInternalServerError, st.Proto())
		return
	}
	h.writeProtoResponse(
		c,
		http.StatusOK,
		&pb.GetItemResponse{
			Item: &pb.Item{
				Id:          item.ID,
				Name:        item.Name,
				Description: item.Description,
				Price:       item.Price,
				Stock:       item.Stock,
			},
		},
	)
}

func (h *Handler) updateItem(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			return
		}

		st := status.New(codes.Unknown, err.Error())
		h.writeProtoResponse(c, http.StatusInternalServerError, st.Proto())
		return
	}
	defer c.Request.Body.Close()

	req := pb.UpdateItemRequest{}
	if err = protojson.Unmarshal(body, &req); err != nil {
		st := status.New(codes.InvalidArgument, err.Error())
		h.writeProtoResponse(c, http.StatusBadRequest, st.Proto())
		return
	}
	newItem := &models.UpdateItem{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}
	err = h.service.UpdateItem(c, newItem)
	if err != nil {
		st := status.New(codes.Unknown, err.Error())
		h.writeProtoResponse(c, http.StatusInternalServerError, st.Proto())
		return
	}
	h.writeProtoResponse(
		c,
		http.StatusOK,
		&pb.EmptyResponse{},
	)
}

func (h *Handler) deleteItem(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			return
		}

		st := status.New(codes.Unknown, err.Error())
		h.writeProtoResponse(c, http.StatusInternalServerError, st.Proto())
		return
	}
	defer c.Request.Body.Close()

	req := pb.DeleteItemRequest{}
	if err = protojson.Unmarshal(body, &req); err != nil {
		st := status.New(codes.InvalidArgument, err.Error())
		h.writeProtoResponse(c, http.StatusBadRequest, st.Proto())
		return
	}
	err = h.service.DeleteItem(c.Request.Context(), req.Id)
	if err != nil {
		st := status.New(codes.Unknown, err.Error())
		h.writeProtoResponse(c, http.StatusInternalServerError, st.Proto())
		return
	}
	h.writeProtoResponse(
		c,
		http.StatusOK,
		&pb.EmptyResponse{},
	)
}
