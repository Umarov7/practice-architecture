package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"practice/internal/controller/http/responder"
	"practice/internal/repository/mongodb/computer"
	"practice/internal/repository/postgres/user"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateUserRabbit godoc
// @Summary User creation through RabbitMQ
// @Description Adds a new user instance via RabbitMQ
// @Tags RabbitMQ
// @Router /user/rabbit [post]
// @Accept			json
// @Produce			json
// @Param userData body UserReq true "User object"
// @Success 201 {object} responder.Response
// @Failure 400 {object} responder.Response
// @Failure 500 {object} responder.Response
func (h *Handler) CreateUserRabbit(w http.ResponseWriter, r *http.Request) {
	var (
		req      UserReq
		response = &responder.Response{}
	)
	defer responder.Send(w, response)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(fmt.Sprintf("wrong body format: %v", err))
		responder.WrongBodyFormat(&responder.Response{}, err)
		return
	}

	msg, err := json.Marshal(user.User{
		ID:        uuid.NewString(),
		Name:      req.Name,
		Age:       req.Age,
		Email:     req.Email,
		IsDeleted: false,
	})
	if err != nil {
		h.logger.Error(fmt.Sprintf("wrong body format: %v", err))
		responder.WrongBodyFormat(&responder.Response{}, err)
		return
	}

	err = h.rabbitProducer.Publish(h.queueUserCreated, msg)
	if err != nil {
		h.logger.Error(fmt.Sprintf("internal server error: %v", err))
		responder.InternalServerError(&responder.Response{}, err)
		return
	}

	response.Code = http.StatusCreated
	response.Payload = "User created"
	response.ContentType = "application/json"
}

// UpdateUserRabbit godoc
// @Summary User update through RabbitMQ
// @Description Updates a user instance via RabbitMQ
// @Tags RabbitMQ
// @Router /user/rabbit/{id} [put]
// @Param id path string true "User ID"
// @Param userData body UserReq true "User object"
// @Success 200 {object} responder.Response
// @Failure 400 {object} responder.Response
// @Failure 500 {object} responder.Response
func (h *Handler) UpdateUserRabbit(w http.ResponseWriter, r *http.Request) {
	var (
		req      UserReq
		response = &responder.Response{}
	)
	defer responder.Send(w, response)

	id := chi.URLParam(r, "id")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(fmt.Sprintf("wrong body format: %v", err))
		responder.WrongBodyFormat(&responder.Response{}, err)
		return
	}

	msg, err := json.Marshal(user.User{
		ID:    id,
		Name:  req.Name,
		Age:   req.Age,
		Email: req.Email,
	})
	if err != nil {
		h.logger.Error(fmt.Sprintf("wrong body format: %v", err))
		responder.WrongBodyFormat(&responder.Response{}, err)
		return
	}

	err = h.rabbitProducer.Publish(h.queueUserUpdated, msg)
	if err != nil {
		h.logger.Error(fmt.Sprintf("internal server error: %v", err))
		responder.InternalServerError(&responder.Response{}, err)
		return
	}

	response.Code = http.StatusOK
	response.Payload = "User updated"
	response.ContentType = "application/json"
}

// DeleteUserRabbit godoc
// @Summary User deletion through RabbitMQ
// @Description Deletes a user instance via RabbitMQ
// @Tags RabbitMQ
// @Router /user/rabbit/{id} [delete]
// @Param id path string true "User ID"
// @Success 200 {object} responder.Response
// @Failure 400 {object} responder.Response
// @Failure 500 {object} responder.Response
func (h *Handler) DeleteUserRabbit(w http.ResponseWriter, r *http.Request) {
	var response responder.Response
	defer responder.Send(w, &response)

	id := chi.URLParam(r, "id")

	err := h.rabbitProducer.Publish(h.queueUserDeleted, []byte(id))
	if err != nil {
		h.logger.Error(fmt.Sprintf("internal server error: %v", err))
		responder.InternalServerError(&responder.Response{}, err)
		return
	}

	response.Code = http.StatusOK
	response.Payload = "User deleted"
	response.ContentType = "application/json"
}

// CreateComputerRabbit godoc
// @Summary Computer creation through RabbitMQ
// @Description Creates a computer instance via RabbitMQ
// @Tags RabbitMQ
// @Router /computer/rabbit [post]
// @Accept			json
// @Produce			json
// @Param computer body ComputerReq true "Computer object"
// @Success 201 {object} responder.Response
// @Failure 400 {object} responder.Response
// @Failure 500 {object} responder.Response
func (h *Handler) CreateComputerRabbit(w http.ResponseWriter, r *http.Request) {
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

	id := primitive.NewObjectID()

	msg, err := json.Marshal(computer.Computer{
		ID:           &id,
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
		h.logger.Error(fmt.Sprintf("wrong body format: %v", err))
		responder.WrongBodyFormat(response, err)
		return
	}

	err = h.rabbitProducer.Publish(h.queueComputerCreated, msg)
	if err != nil {
		h.logger.Error(fmt.Sprintf("internal server error: %v", err))
		responder.InternalServerError(response, err)
		return
	}

	response.Code = http.StatusCreated
	response.Payload = "Computer created"
	response.ContentType = "application/json"
}

// UpdateComputerRabbit godoc
// @Summary Computer update through RabbitMQ
// @Description Updates a computer instance via RabbitMQ
// @Tags RabbitMQ
// @Router /computer/rabbit/{id} [put]
// @Param id path string true "Computer ID"
// @Param computer body ComputerReq true "Computer object"
// @Success 200 {object} responder.Response
// @Failure 400 {object} responder.Response
// @Failure 500 {object} responder.Response
func (h *Handler) UpdateComputerRabbit(w http.ResponseWriter, r *http.Request) {
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

	msg, err := json.Marshal(computer.Computer{
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
		h.logger.Error(fmt.Sprintf("wrong body format: %v", err))
		responder.WrongBodyFormat(response, err)
		return
	}

	err = h.rabbitProducer.Publish(h.queueComputerUpdated, msg)
	if err != nil {
		h.logger.Error(fmt.Sprintf("internal server error: %v", err))
		responder.InternalServerError(response, err)
		return
	}

	response.Code = http.StatusOK
	response.Payload = "Computer updated"
	response.ContentType = "application/json"
}

// DeleteComputerRabbit godoc
// @Summary Computer deletion through RabbitMQ
// @Description Deletes a computer instance via RabbitMQ
// @Tags RabbitMQ
// @Router /computer/rabbit/{id} [delete]
// @Param id path string true "Computer ID"
// @Success 200 {object} responder.Response
// @Failure 400 {object} responder.Response
// @Failure 500 {object} responder.Response
func (h *Handler) DeleteComputerRabbit(w http.ResponseWriter, r *http.Request) {
	var response responder.Response
	defer responder.Send(w, &response)

	id := chi.URLParam(r, "id")

	err := h.rabbitProducer.Publish(h.queueComputerDeleted, []byte(id))
	if err != nil {
		h.logger.Error(fmt.Sprintf("internal server error: %v", err))
		responder.InternalServerError(&responder.Response{}, err)
		return
	}

	response.Code = http.StatusOK
	response.Payload = "Computer deleted"
	response.ContentType = "application/json"
}
