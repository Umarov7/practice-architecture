package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"practice/internal/controller/http/responder"
	"practice/internal/repository/mongodb/computer"
	"practice/internal/repository/postgres/user"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateUserKafka godoc
// @Summary User creation through kafka
// @Description Adds a new user instance via kafka
// @Tags Kafka
// @Router /user/kafka [post]
// @Accept			json
// @Produce			json
// @Param userData body UserReq true "User object"
// @Success 201 {object} responder.Response
// @Failure 400 {object} responder.Response
// @Failure 500 {object} responder.Response
func (h *Handler) CreateUserKafka(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

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

	err = h.kafkaProducer.Produce(ctx, h.topicUserCreated, msg)
	if err != nil {
		h.logger.Error(fmt.Sprintf("internal server error: %v", err))
		responder.InternalServerError(&responder.Response{}, err)
		return
	}

	response.Code = http.StatusCreated
	response.Payload = "User created"
	response.ContentType = "application/json"
}

// UpdateUserKafka godoc
// @Summary User update through kafka
// @Description Updates a user instance via kafka
// @Tags Kafka
// @Router /user/kafka/{id} [put]
// @Param id path string true "User ID"
// @Param userData body UserReq true "User object"
// @Success 200 {object} responder.Response
// @Failure 400 {object} responder.Response
// @Failure 500 {object} responder.Response
func (h *Handler) UpdateUserKafka(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

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

	err = h.kafkaProducer.Produce(ctx, h.topicUserUpdated, msg)
	if err != nil {
		h.logger.Error(fmt.Sprintf("internal server error: %v", err))
		responder.InternalServerError(&responder.Response{}, err)
		return
	}

	response.Code = http.StatusOK
	response.Payload = "User updated"
	response.ContentType = "application/json"
}

// DeleteUserKafka godoc
// @Summary User deletion through kafka
// @Description Deletes a user instance via kafka
// @Tags Kafka
// @Router /user/kafka/{id} [delete]
// @Param id path string true "User ID"
// @Success 200 {object} responder.Response
// @Failure 400 {object} responder.Response
// @Failure 500 {object} responder.Response
func (h *Handler) DeleteUserKafka(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var response responder.Response
	defer responder.Send(w, &response)

	id := chi.URLParam(r, "id")

	err := h.kafkaProducer.Produce(ctx, h.topicUserDeleted, []byte(id))
	if err != nil {
		h.logger.Error(fmt.Sprintf("internal server error: %v", err))
		responder.InternalServerError(&responder.Response{}, err)
		return
	}

	response.Code = http.StatusOK
	response.Payload = "User deleted"
	response.ContentType = "application/json"
}

// CreateComputerKafka godoc
// @Summary Computer creation through kafka
// @Description Creates a computer instance via kafka
// @Tags Kafka
// @Router /computer/kafka [post]
// @Accept			json
// @Produce			json
// @Param computer body ComputerReq true "Computer object"
// @Success 201 {object} responder.Response
// @Failure 400 {object} responder.Response
// @Failure 500 {object} responder.Response
func (h *Handler) CreateComputerKafka(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
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

	err = h.kafkaProducer.Produce(ctx, h.topicComputerCreated, msg)
	if err != nil {
		h.logger.Error(fmt.Sprintf("internal server error: %v", err))
		responder.InternalServerError(response, err)
		return
	}

	response.Code = http.StatusCreated
	response.Payload = "Computer created"
	response.ContentType = "application/json"
}

// UpdateComputerKafka godoc
// @Summary Computer update through kafka
// @Description Updates a computer instance via kafka
// @Tags Kafka
// @Router /computer/kafka/{id} [put]
// @Param id path string true "Computer ID"
// @Param computer body ComputerReq true "Computer object"
// @Success 200 {object} responder.Response
// @Failure 400 {object} responder.Response
// @Failure 500 {object} responder.Response
func (h *Handler) UpdateComputerKafka(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
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

	err = h.kafkaProducer.Produce(ctx, h.topicComputerUpdated, msg)
	if err != nil {
		h.logger.Error(fmt.Sprintf("internal server error: %v", err))
		responder.InternalServerError(response, err)
		return
	}

	response.Code = http.StatusOK
	response.Payload = "Computer updated"
	response.ContentType = "application/json"
}

// DeleteComputerKafka godoc
// @Summary Computer deletion through kafka
// @Description Deletes a computer instance via kafka
// @Tags Kafka
// @Router /computer/kafka/{id} [delete]
// @Param id path string true "Computer ID"
// @Success 200 {object} responder.Response
// @Failure 400 {object} responder.Response
// @Failure 500 {object} responder.Response
func (h *Handler) DeleteComputerKafka(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var response responder.Response
	defer responder.Send(w, &response)

	id := chi.URLParam(r, "id")

	err := h.kafkaProducer.Produce(ctx, h.topicComputerDeleted, []byte(id))
	if err != nil {
		h.logger.Error(fmt.Sprintf("internal server error: %v", err))
		responder.InternalServerError(&responder.Response{}, err)
		return
	}

	response.Code = http.StatusOK
	response.Payload = "Computer deleted"
	response.ContentType = "application/json"
}
