package guide

import (
	"fmt"
	"github.com/Viverov/guideliner/internal/cradle"
	"github.com/Viverov/guideliner/internal/domains/guide/service"
	"github.com/Viverov/guideliner/internal/domains/util"
	"github.com/Viverov/guideliner/internal/domains/util/uservice"
	"github.com/Viverov/guideliner/internal/server/controller/utils"
	"github.com/Viverov/guideliner/internal/server/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Controller struct {
	httpResponder utils.HttpResponder
}

func NewGuideController(httpResponder utils.HttpResponder) *Controller {
	return &Controller{httpResponder: httpResponder}
}

func (c *Controller) Init(router *gin.Engine, cradle *cradle.Cradle, prefix string) {
	router.GET(prefix + "/guides", createFindHandler(cradle, c.httpResponder))
	router.GET(prefix + "/guides/:id", createFindByIdHandler(cradle, c.httpResponder))
	router.POST(prefix + "/guides", middleware.CreateAuthMiddleware(cradle, c.httpResponder), createNewGuideHandler(cradle, c.httpResponder))
	router.PATCH(prefix + "/guides/:id", middleware.CreateAuthMiddleware(cradle, c.httpResponder), createUpdateHandler(cradle, c.httpResponder))
}

type guideResponse struct {
	ID          uint   `json:"id"`
	Description string `json:"description"`
	Nodes       string `json:"nodes"`
}

type findQuery struct {
	Limit       int    `json:"limit"`
	Offset      int    `json:"offset"`
	Order       string `json:"order"`
	ParsedOrder []util.Order
	Search      string `json:"search"`
}

type findResult struct {
	Guides []guideResponse `json:"guides"`
	Total  int64           `json:"total"`
}

func createFindHandler(cradle *cradle.Cradle, responder utils.HttpResponder) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		query := &findQuery{}
		if err := ctx.ShouldBindQuery(query); err != nil {
			responder.Response(ctx, http.StatusBadRequest, "Validation error", err.Error(), "")
			return
		}
		parsedOrder, err := utils.ParseOrderQuery(query.Order, map[string]string{"id": "id", "created_at": "createdAt"})
		if err != nil {
			responder.Response(ctx, http.StatusBadRequest, "Validation error", err.Error(), "")
		}
		query.ParsedOrder = parsedOrder

		gService := cradle.GetGuideService()

		count, err := gService.Count(service.CountConditions{
			Search: query.Search,
		})
		if err != nil {
			responder.InternalError(ctx, err.Error())
		}

		dtos, err := gService.Find(service.FindConditions{
			DefaultFindConditions: util.DefaultFindConditions{
				Limit:  query.Limit,
				Offset: query.Offset,
				Order:  query.ParsedOrder,
			},
			Search: query.Search,
		})
		if err != nil {
			responder.InternalError(ctx, err.Error())
		}

		var formattedGuides []guideResponse
		for _, dto := range dtos {
			var g guideResponse

			g.ID = dto.ID()
			g.Description = dto.Description()
			g.Nodes = dto.NodesJson()

			formattedGuides = append(formattedGuides, g)
		}

		ctx.JSON(http.StatusOK, findResult{Guides: formattedGuides, Total: count})
	}
}

func createFindByIdHandler(cradle *cradle.Cradle, responder utils.HttpResponder) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		sID := ctx.Param("id")
		uID, err := strconv.ParseUint(sID, 10, 32)
		if err != nil {
			responder.Response(ctx, http.StatusBadRequest, "Invalid ID", fmt.Sprintf("Expected uint, got %s", sID), err.Error())
			return
		}
		id := uint(uID)

		dto, err := cradle.GetGuideService().FindById(id)
		if err != nil {
			switch err.(type) {
			case *uservice.NotFoundError:
				responder.Response(ctx, http.StatusNotFound, "Not found", fmt.Sprintf("Guide with id %d not found", id), err.Error())
				return
			default:
				responder.InternalError(ctx, err.Error())
				return
			}
		}

		ctx.JSON(http.StatusOK, guideResponse{
			ID:          dto.ID(),
			Description: dto.Description(),
			Nodes:       dto.NodesJson(),
		})
	}
}

type createBody struct {
	Description string `json:"description" binding:"required"`
	Nodes       string `json:"nodes" binding:"required"`
}

type createResponse struct {
	ID          uint   `json:"id"`
	Description string `json:"description"`
	Nodes       string `json:"nodes"`
}

func createNewGuideHandler(cradle *cradle.Cradle, responder utils.HttpResponder) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		body := &createBody{}
		if err := ctx.ShouldBindJSON(body); err != nil {
			responder.Response(ctx, http.StatusBadRequest, "Validation error", err.Error(), "")
			return
		}

		dto, err := cradle.GetGuideService().Create(body.Description, body.Nodes)
		if err != nil {
			switch err.(type) {
			case *service.InvalidNodesJsonError:
				responder.Response(ctx, http.StatusBadRequest, "Invalid nodes format", err.Error(), "")
				return
			default:
				responder.InternalError(ctx, err.Error())
			}
		}

		ctx.JSON(http.StatusCreated, createResponse{
			ID:          dto.ID(),
			Description: dto.Description(),
			Nodes:       dto.NodesJson(),
		})
	}
}

type updateBody struct {
	Description string `json:"description" binding:"required"`
	Nodes       string `json:"nodes" binding:"required"`
}

type updateResponse struct {
	ID          uint   `json:"id"`
	Description string `json:"description"`
	Nodes       string `json:"nodes"`
}

func createUpdateHandler(cradle *cradle.Cradle, responder utils.HttpResponder) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		sID := ctx.Param("id")
		uID, err := strconv.ParseUint(sID, 10, 32)
		if err != nil {
			responder.Response(ctx, http.StatusBadRequest, "Invalid ID", fmt.Sprintf("Expected uint, got %s", sID), err.Error())
			return
		}
		id := uint(uID)

		body := &updateBody{}
		if err := ctx.ShouldBindJSON(body); err != nil {
			responder.Response(ctx, http.StatusBadRequest, "Validation error", err.Error(), "")
			return
		}

		dto, err := cradle.GetGuideService().Update(id, service.UpdateParams{
			Description: body.Description,
			NodesJson:   body.Nodes,
		})
		if err != nil {
			switch err.(type) {
			case *uservice.NotFoundError:
				responder.Response(ctx, http.StatusNotFound, "Not found", fmt.Sprintf("Guide with id %d not found", id), err.Error())
				return
			case *service.InvalidNodesJsonError:
				responder.Response(ctx, http.StatusBadRequest, "Invalid nodes format", err.Error(), "")
				return
			default:
				responder.InternalError(ctx, err.Error())
				return
			}
		}

		ctx.JSON(http.StatusCreated, updateResponse{
			ID:          dto.ID(),
			Description: dto.Description(),
			Nodes:       dto.NodesJson(),
		})
	}
}
