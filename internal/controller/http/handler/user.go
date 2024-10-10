package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"practice/internal/controller/http/responder"
	"practice/internal/repository/postgres/user"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// CreateUser godoc
// @Summary User creataion
// @Description Adds a new user instance
// @Tags User
// @Router /user [post]
// @Accept			json
// @Produce			json
// @Param userData body UserReq true "User object"
// @Success 201 {object} user.User
// @Failure 400 {object} responder.Response
// @Failure 500 {object} responder.Response
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var (
		req      UserReq
		response = &responder.Response{}
	)

	defer responder.Send(w, response)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(fmt.Sprintf("wrong body format: %v", err))
		responder.WrongBodyFormat(response, err)
		return
	}

	res, err := h.serviceUser.Create(ctx, &user.User{
		ID:        uuid.NewString(),
		Name:      req.Name,
		Age:       req.Age,
		Email:     req.Email,
		IsDeleted: false,
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

// GetUser godoc
// @Summary User reading
// @Description Returns a user instance
// @Tags User
// @Router /user/{id} [get]
// @Param id path string true "User ID"
// @Success 200 {object} user.User
// @Failure 400 {object} responder.Response
// @Failure 404 {object} responder.Response
// @Failure 500 {object} responder.Response
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var response responder.Response
	defer responder.Send(w, &response)

	id := chi.URLParam(r, "id")

	res, err := h.serviceUser.Read(ctx, id)
	if err != nil {
		h.logger.Error(fmt.Sprintf("internal server error: %v", err))
		responder.InternalServerError(&responder.Response{}, err)
		return
	}

	response.Code = http.StatusOK
	response.Payload = res
	response.ContentType = "application/json"
}

// UpdateUser godoc
// @Summary User update
// @Description Updates a user instance
// @Tags User
// @Router /user/{id} [put]
// @Accept			json
// @Produce			json
// @Param id path string true "User ID"
// @Param userData body UserReq true "User object"
// @Success 200 {object} user.User
// @Failure 400 {object} responder.Response
// @Failure 500 {object} responder.Response
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var (
		req      UserReq
		response = &responder.Response{}
	)

	defer responder.Send(w, response)

	id := chi.URLParam(r, "id")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(fmt.Sprintf("wrong body format: %v", err))
		responder.WrongBodyFormat(response, err)
		return
	}

	res, err := h.serviceUser.Update(ctx, &user.User{
		ID:    id,
		Name:  req.Name,
		Age:   req.Age,
		Email: req.Email,
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

// DeleteUser godoc
// @Summary User deletion
// @Description Deletes a user instance
// @Tags User
// @Router /user/{id} [delete]
// @Param id path string true "User ID"
// @Success 200 {object} user.User
// @Failure 400 {object} responder.Response
// @Failure 500 {object} responder.Response
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var response responder.Response
	defer responder.Send(w, &response)

	id := chi.URLParam(r, "id")

	id, err := h.serviceUser.Delete(ctx, id)
	if err != nil {
		h.logger.Error(fmt.Sprintf("internal server error: %v", err))
		responder.InternalServerError(&responder.Response{}, err)
		return
	}

	response.Code = http.StatusOK
	response.Payload = id
	response.ContentType = "application/json"
}
