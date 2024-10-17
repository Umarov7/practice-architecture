package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"practice/internal/controller/http/responder"
	"practice/internal/repository/mongodb/computer"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateComputer godoc
// @Summary Computer creation
// @Description Adds a new computer instance
// @Tags Computer
// @Router /computer [post]
// @Accept			json
// @Produce			json
// @Param computer body ComputerReq true "Computer object"
// @Success 201 {object} computer.Computer
// @Failure 400 {object} responder.Response
// @Failure 500 {object} responder.Response
func (h *Handler) CreateComputer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var (
		req      ComputerReq
		response = &responder.Response{}
	)

	defer responder.Send(w, response)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(fmt.Sprintf("wrong body format: %v", err))
		responder.WrongBodyFormat(response, err)
		return
	}

	res, err := h.serviceComputer.Create(ctx, &computer.Computer{
		IP:           req.IP,
		Manufacturer: req.Manufacturer,
		CPU:          req.CPU,
		RAM:          req.RAM,
		HDD:          req.HDD,
		GPU:          req.GPU,
		OS:           req.OS,
		IsDeleted:    false,
	})
	if err != nil {
		h.logger.Error(fmt.Sprintf("internal server error: %v", err))
		responder.InternalServerError(response, err)
		return
	}

	response.Code = http.StatusCreated
	response.Payload = res
	response.ContentType = "application/json"
}

// @Summary Computer reading
// @Description Returns a computer instance
// @Tags Computer
// @Router /computer/{id} [get]
// @Param id path string true "Computer ID"
// @Success 200 {object} computer.Computer
// @Failure 400 {object} responder.Response
// @Failure 404 {object} responder.Response
// @Failure 500 {object} responder.Response
func (h *Handler) GetComputer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var response responder.Response
	defer responder.Send(w, &response)

	id := chi.URLParam(r, "id")

	res, err := h.serviceComputer.Read(ctx, id)
	if err != nil {
		h.logger.Error(fmt.Sprintf("internal server error: %v", err))
		responder.InternalServerError(&responder.Response{}, err)
		return
	}

	response.Code = http.StatusOK
	response.Payload = res
	response.ContentType = "application/json"
}

// @Summary Computer update
// @Description Updates a computer instance
// @Tags Computer
// @Router /computer/{id} [put]
// @Accept			json
// @Produce			json
// @Param id path string true "Computer ID"
// @Param computer body ComputerReq true "Computer object"
// @Success 200 {object} computer.Computer
// @Failure 400 {object} responder.Response
// @Failure 500 {object} responder.Response
func (h *Handler) UpdateComputer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var (
		req      ComputerReq
		response = &responder.Response{}
	)

	defer responder.Send(w, response)

	idStr := chi.URLParam(r, "id")

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		h.logger.Error(fmt.Sprintf("wrong body format: %v", err))
		responder.WrongBodyFormat(response, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(fmt.Sprintf("wrong body format: %v", err))
		responder.WrongBodyFormat(response, err)
		return
	}

	h.logger.Info("body: ", req)
	res, err := h.serviceComputer.Update(ctx, &computer.Computer{
		ID:           &id,
		IP:           req.IP,
		Manufacturer: req.Manufacturer,
		CPU:          req.CPU,
		RAM:          req.RAM,
		HDD:          req.HDD,
		GPU:          req.GPU,
		OS:           req.OS,
	})
	if err != nil {
		h.logger.Error(fmt.Sprintf("internal server error: %v", err))
		responder.InternalServerError(response, err)
		return
	}

	response.Code = http.StatusOK
	response.Payload = res
	response.ContentType = "application/json"
}

// DeleteComputer godoc
// @Summary Computer deletion
// @Description Deletes a computer instance
// @Tags Computer
// @Router /computer/{id} [delete]
// @Param id path string true "Computer ID"
// @Success 200 {object} string
// @Failure 400 {object} responder.Response
// @Failure 500 {object} responder.Response
func (h *Handler) DeleteComputer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var response responder.Response
	defer responder.Send(w, &response)

	id := chi.URLParam(r, "id")

	id, err := h.serviceComputer.Delete(ctx, id)
	if err != nil {
		h.logger.Error(fmt.Sprintf("internal server error: %v", err))
		responder.InternalServerError(&responder.Response{}, err)
		return
	}

	response.Code = http.StatusOK
	response.Payload = id
	response.ContentType = "application/json"
}

// ListComputers godoc
// @Summary Computer list
// @Description Returns a list of computer instances
// @Tags Computer
// @Router /computer [get]
// @Success 200 {object} []computer.Computer
// @Failure 400 {object} responder.Response
// @Failure 500 {object} responder.Response
func (h *Handler) ListComputers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var response responder.Response
	defer responder.Send(w, &response)

	res, err := h.serviceComputer.GetAll(ctx)
	if err != nil {
		h.logger.Error(fmt.Sprintf("internal server error: %v", err))
		responder.InternalServerError(&responder.Response{}, err)
		return
	}

	response.Code = http.StatusOK
	response.Payload = res
	response.ContentType = "application/json"
}
