package user_handlers

import (
	"api-for-people/internal/user/model"
	"api-for-people/internal/user/service"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Handler struct {
	service *service.Service
	ctx     context.Context
	mu      *sync.Mutex
}

type Countries struct {
	CountryId   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type PersonCountry struct {
	Country []Countries `json:"country"`
}

func (h *Handler) FetchData(name string) map[string]any {
	set := make(map[string]any)
	client := http.Client{
		Timeout: time.Second * 10,
	}
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		slog.Debug("Fetching nationality")
		defer slog.Debug("Fetching nationality finished")
		defer wg.Done()
		var country PersonCountry
		resp, err := client.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", name))
		if err != nil {
			return
		}
		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				slog.Error("Error closing body")
			}
		}(resp.Body)
		err = json.NewDecoder(resp.Body).Decode(&country)
		if err != nil {
			return
		}
		var nation string
		probability := 0.0
		for _, item := range country.Country {
			if item.Probability > probability {
				probability = item.Probability
				nation = item.CountryId
			}
		}
		h.mu.Lock()
		set["nation"] = nation
		h.mu.Unlock()
	}()
	go func() {
		defer wg.Done()
		slog.Debug("Fetching age")
		defer slog.Debug("Fetching age finished")
		resp, err := client.Get(fmt.Sprintf("https://api.agify.io/?name=%s", name))
		if err != nil {
		}
		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				slog.Error("Error closing response body")
			}
		}(resp.Body)
		var result struct {
			Age int `json:"age"`
		}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return
		}
		h.mu.Lock()
		set["age"] = result.Age
		h.mu.Unlock()
	}()
	go func() {
		slog.Debug("Fetching gender")
		defer slog.Debug("Fetching gender finished")
		defer wg.Done()
		resp, err := client.Get(fmt.Sprintf("https://api.genderize.io/?name=%s", name))
		if err != nil {
		}
		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				slog.Error("Error closing response body")
			}
		}(resp.Body)
		var result struct {
			Gender string `json:"gender"`
		}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return
		}
		h.mu.Lock()
		set["gender"] = result.Gender
		h.mu.Unlock()
	}()
	wg.Wait()
	slog.Debug("Fetched data from API")
	return set
}

func NewHandler(ctx context.Context, service *service.Service) *Handler {
	slog.Info("Creating new handler")
	return &Handler{service: service, ctx: ctx, mu: &sync.Mutex{}}
}

// @Summary People
// @Tag create
// @Router / [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var person model.PersonFromRequest
	if err := render.DecodeJSON(r.Body, &person); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("Error decoding body", err)
	}
	set := h.FetchData(person.Name)
	personToDB := &model.Person{}
	personToDB.Name = person.Name
	personToDB.Surname = person.Surname
	if person.Patronymic != nil {
		personToDB.Patronymic = person.Patronymic
	} else {
		personToDB.Patronymic = nil
	}
	personToDB.Age = set["age"].(int)
	personToDB.Gender = set["gender"].(string)
	personToDB.Nationality = set["nation"].(string)
	slog.Debug("Creating new person")
	err := h.service.Create(h.ctx, personToDB)
	if err != nil {
		slog.Error("Error creating person", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// @Summary People
// @Tag getPerson
// @Router /persons/{:id} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Getting person")
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	fmt.Println(id)
	person, err := h.service.Get(h.ctx, id)
	fmt.Println(person)
	if err != nil {
		slog.Error("Error getting person", err)
	}
	render.JSON(w, r, person)
}

// @Summary People
// @Tag getAllPerson
// @Router /persons [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Getting all persons")
	params := model.ParseQueryParams(r)
	persons, err := h.service.GetAll(h.ctx, params)
	if err != nil {
		slog.Error("Error getting all persons", err)
	}
	render.JSON(w, r, persons)
}

// @Summary People
// @Tag opdatePerson
// @Router /persons/{:id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Updating person")
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	var person model.Person
	if err := render.DecodeJSON(r.Body, &person); err != nil {
		slog.Error("Error decoding body", err)
	}
	err := h.service.Update(h.ctx, person, id)
	if err != nil {
		slog.Error("Error updating person", err)
	}
}

// @Summary People
// @Tag deletePerson
// @Router /persons/{:id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Deleting person")
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	err := h.service.Delete(h.ctx, id)
	if err != nil {
		slog.Error("Error deleting person", err)
	}
}
