package handlers

import (
	m "cards/internal/api/v1/middlewares"
	"cards/internal/models"
	filters "cards/internal/pkg/filters"
	writeJSONresponse "cards/internal/pkg/writeJSONresponse"
	"cards/internal/service"
	"context"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
)

const healthPath = "/api/v1/ping"
const cardsPath = "/api/v1/cards"

type HandlerAdapter struct {
	core.RequestAccessor
	handler http.Handler
	serviceCards service.Cards
}

func NewCardsHandler(handler http.Handler, s service.Cards) HandlerAdapter {
	return HandlerAdapter{
		handler: handler,
		serviceCards: s,
	}
}

func (h *HandlerAdapter) Handlers() {
	// Health check
	healthCheck := http.HandlerFunc(h.handleHealthCheck)
	http.Handle(healthPath, m.EnforceJSONHandler(healthCheck))
	
	// Cards
	hc := http.HandlerFunc(h.handleCards)
	http.Handle(cardsPath, m.EnforceJSONHandler(hc))

    // Serve Swagger UI files
    http.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir("swagger-ui"))))
    // Handler for serving Swagger JSON
    http.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, ".docs/openapi.json")
    })
}

func (h *HandlerAdapter) handleHealthCheck(response http.ResponseWriter, request *http.Request) {
	writeJSONresponse.WriteJSONresponse(response, http.StatusOK, http.StatusText(http.StatusOK), "pong", models.Pagination{}, "")
}

func (h *HandlerAdapter) handleCards(response http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	switch request.Method {
	case http.MethodGet:
		h.handleGetCards(ctx, response, request)

	case http.MethodPost:
		h.handleCreateCard(ctx, response, request)

	case http.MethodPatch:
		h.handleUpdateCard(ctx, response, request) // PATCH only modifies resource contents, PUT updates a resource by replacing its content entirely.
	
	case http.MethodDelete:
		h.handleDeleteCards(ctx, response, request)

	default:
		writeJSONresponse.WriteJSONresponse(response, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed), nil, models.Pagination{}, "")
	}
}

func (h *HandlerAdapter) handleGetCards(ctx context.Context, response http.ResponseWriter, request *http.Request) {
	ids := strings.Split(request.URL.Query().Get("id"), ",")

	p := strings.Split(request.URL.Query().Get("page_number"), ",")[0]
	ps := strings.Split(request.URL.Query().Get("page_size"), ",")[0]
	pagination := models.Pagination{
		PageNumber:  p,
		PageSize: ps,
	}

	rawFilters := request.URL.Query().Get("filter") // with a value as -1 for gorms Limit method, we'll get a request without limit as default
	if rawFilters == "" && len(ids) == 0 {
		writeJSONresponse.WriteJSONresponse(response, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), []models.Card{}, models.Pagination{}, "Please provide at least one filter or id")
		return
	}
	filtersMap := map[string]string{}
    if rawFilters != "" {
		var err error
        filtersMap, err = filters.ValidateAndReturnFilterMap(rawFilters)
        if err != nil {
			writeJSONresponse.WriteJSONresponse(response, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), []models.Card{}, models.Pagination{}, err.Error())
			return
		}
    }

	cards, err := h.serviceCards.GetCards(ctx, ids, pagination, filtersMap)
	if err != nil {
		writeJSONresponse.WriteJSONresponse(response, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), []models.Card{}, models.Pagination{}, err.Error())
		return	
	}

	writeJSONresponse.WriteJSONresponse(response, http.StatusOK, http.StatusText(http.StatusOK), cards, pagination, "")
}

func (h *HandlerAdapter) handleCreateCard(ctx context.Context, response http.ResponseWriter, request *http.Request) {
	card, err := h.serviceCards.CreateCard(ctx, request)
	if err != nil {
		writeJSONresponse.WriteJSONresponse(response, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), []models.Card{}, models.Pagination{}, err.Error())
		return
	}

	writeJSONresponse.WriteJSONresponse(response, http.StatusCreated, http.StatusText(http.StatusCreated), card, models.Pagination{}, "")
}

func (h *HandlerAdapter) handleUpdateCard(ctx context.Context, response http.ResponseWriter, request *http.Request) {
	cards, err := h.serviceCards.UpdateCard(ctx, request)
	if err != nil {
		writeJSONresponse.WriteJSONresponse(response, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), []models.Card{}, models.Pagination{}, err.Error())
		return	
	}

	writeJSONresponse.WriteJSONresponse(response, http.StatusOK, http.StatusText(http.StatusOK), cards, models.Pagination{}, "")
}

func (h *HandlerAdapter) handleDeleteCards(ctx context.Context, response http.ResponseWriter, request *http.Request) {
	id := strings.Split(request.URL.Query().Get("id"), ",")[0]
	err := h.serviceCards.DeleteCard(ctx, id)
	if err != nil {
		writeJSONresponse.WriteJSONresponse(response, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), []models.Card{}, models.Pagination{}, err.Error())
		return	
	}

	writeJSONresponse.WriteJSONresponse(response, http.StatusOK, http.StatusText(http.StatusOK), models.Card{}, models.Pagination{}, "")
}

func (h *HandlerAdapter) ProxyWithContext(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	req, err := h.EventToRequestWithContext(ctx, event)
	return h.proxyInternal(req, err)
}

func (h *HandlerAdapter) proxyInternal(req *http.Request, err error) (events.APIGatewayProxyResponse, error) {
	if err != nil {
		return core.GatewayTimeout(), core.NewLoggedError("Could not convert proxy event to request: %v", err)
	}

	w := core.NewProxyResponseWriter()
	h.handler.ServeHTTP(http.ResponseWriter(w), req)

	resp, err := w.GetProxyResponse()
	if err != nil {
		return core.GatewayTimeout(), core.NewLoggedError("Error while generating proxy response: %v", err)
	}

	return resp, nil
}
