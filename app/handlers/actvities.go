package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ungame/command-time-track/app/httpext"
	"github.com/ungame/command-time-track/app/service"
	"github.com/ungame/command-time-track/app/types"
	"io/ioutil"
	"net/http"
	"strconv"
)

type activitiesHandler struct {
	activitiesService service.ActivitiesService
}

func NewActivitiesHandler(activitiesService service.ActivitiesService) Handler {
	return &activitiesHandler{activitiesService: activitiesService}
}

func (h *activitiesHandler) Register(router *mux.Router) {
	router.Path("/activities").HandlerFunc(h.PostStartActivity).Methods(http.MethodPost)
	router.Path("/activities/{id}/stop").HandlerFunc(h.PutStopActivity).Methods(http.MethodPut)
	router.Path("/activities/{id}/category").HandlerFunc(h.PutActivityCategory).Methods(http.MethodPut)
	router.Path("/activities/{id}/description").HandlerFunc(h.PutActivityDescription).Methods(http.MethodPut)
	router.Path("/activities/{id}").HandlerFunc(h.GetActivity).Methods(http.MethodGet)
	router.Path("/activities/_/search").HandlerFunc(h.SearchActivity).Methods(http.MethodGet)
	router.Path("/activities").HandlerFunc(h.GetActivities).Methods(http.MethodGet)
	router.Path("/activities/{id}").HandlerFunc(h.DeleteActivity).Methods(http.MethodDelete)
}

func (h *activitiesHandler) PostStartActivity(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httpext.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}
	input := new(types.StartActivityInput)
	err = json.Unmarshal(body, input)
	if err != nil {
		httpext.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}
	output, err := h.activitiesService.StartActivity(r.Context(), input)
	if err != nil {
		httpext.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s/%d", r.RequestURI, output.ID))
	httpext.WriteJson(w, http.StatusCreated, output)
}

func (h *activitiesHandler) PutStopActivity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		httpext.WriteError(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httpext.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}
	input := new(types.UpdateActivityInput)
	err = json.Unmarshal(body, input)
	if err != nil {
		httpext.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}
	input.ID = id
	output, err := h.activitiesService.StopActivity(r.Context(), input)
	if err != nil {
		httpext.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}
	httpext.WriteJson(w, http.StatusOK, output)
}

func (h *activitiesHandler) PutActivityCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		httpext.WriteError(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httpext.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}
	input := new(types.UpdateActivityInput)
	err = json.Unmarshal(body, input)
	if err != nil {
		httpext.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}
	input.ID = id
	output, err := h.activitiesService.UpdateActivityCategory(r.Context(), input)
	if err != nil {
		httpext.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}
	httpext.WriteJson(w, http.StatusOK, output)
}

func (h *activitiesHandler) PutActivityDescription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		httpext.WriteError(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httpext.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}
	input := new(types.UpdateActivityInput)
	err = json.Unmarshal(body, input)
	if err != nil {
		httpext.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}
	input.ID = id
	output, err := h.activitiesService.UpdateActivityDescription(r.Context(), input)
	if err != nil {
		httpext.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}
	httpext.WriteJson(w, http.StatusOK, output)
}

func (h *activitiesHandler) GetActivity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		httpext.WriteError(w, http.StatusBadRequest, err)
		return
	}
	activity, err := h.activitiesService.GetActivityByID(r.Context(), &types.GetActivityInput{ID: id})
	if err != nil {
		httpext.WriteError(w, http.StatusBadRequest, err)
		return
	}
	httpext.WriteJson(w, http.StatusOK, activity)
}

func (h *activitiesHandler) SearchActivity(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("term")
	activities, err := h.activitiesService.SearchActivities(r.Context(), term)
	if err != nil {
		httpext.WriteError(w, http.StatusBadRequest, err)
		return
	}
	httpext.WriteJson(w, http.StatusOK, activities)
}

func (h *activitiesHandler) GetActivities(w http.ResponseWriter, r *http.Request) {
	activities, err := h.activitiesService.ListActivities(r.Context())
	if err != nil {
		httpext.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	httpext.WriteJson(w, http.StatusOK, activities)
}

func (h *activitiesHandler) DeleteActivity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		httpext.WriteError(w, http.StatusBadRequest, err)
		return
	}
	id, err = h.activitiesService.DeleteActivityByID(r.Context(), &types.DeleteActivityInput{ID: id})
	if err != nil {
		httpext.WriteError(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprint(id))
	w.WriteHeader(http.StatusNoContent)
}
