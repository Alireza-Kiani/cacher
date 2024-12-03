package controller

import (
	"cache/internal/dto"
	"cache/internal/model"
	"cache/internal/service"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Controller struct {
	service *service.CacheService
	port    int64
}

func New(s *service.CacheService, port int64) *Controller {
	c := &Controller{
		service: s,
		port:    port,
	}

	return c
}

func (c *Controller) Init() {
	mux := http.NewServeMux()
	mux.HandleFunc("/set", c.set)
	mux.HandleFunc("/get", c.get)
	mux.HandleFunc("/del", c.delete)

	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", c.port), mux)
	if err != nil {
		panic(err)
	}
}

func (c *Controller) get(w http.ResponseWriter, r *http.Request) {
	n := time.Now()
	key := ""

	defer func() {
		log.Println("HTTP request get value", "response_time", time.Since(n), "key", key)
	}()

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	q, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, "Bad Request, Query Cannot Be Parsed", http.StatusBadRequest)
		return
	}

	key = q.Get("key")
	resultChan := c.service.ReceiveCommand(model.NewGetCommand(model.Key(key)))

	result := <-resultChan
	if !result.Ok {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if result.Error != nil {
		http.Error(w, fmt.Sprintf("Cannot Retreive Data, Reason: %s", result.Error.Error()), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(result.Value)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot Encode Data, Reason: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")

	_, err = w.Write(resp)
	if err != nil {
		log.Println("error on write http response", "key", key, "error", err)
	}
}

func (c *Controller) set(w http.ResponseWriter, r *http.Request) {
	n := time.Now()
	body := []byte{}

	defer func() {
		log.Println("HTTP request set value", "response_time", time.Since(n), "body", string(body))
	}()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()
	body, _ = io.ReadAll(r.Body)

	dto := dto.SetRequest{}
	err := json.Unmarshal(body, &dto)
	if err != nil {
		log.Println("error on reading body", "error", err)
		http.Error(w, "Bad Request, Body Cannot Be Parsed, Please Provide Valid JSON", http.StatusBadRequest)
		return
	}

	err = dto.Validate()
	if err != nil {
		log.Println("error on reading body", "error", err)
		http.Error(w, fmt.Sprintf("Bad Request, Reason: %s", err.Error()), http.StatusBadRequest)
		return
	}

	resultChan := c.service.ReceiveCommand(model.NewSetCommand(model.Key(dto.Key), dto.Value, dto.TTL))

	result := <-resultChan
	if result.Error != nil {
		log.Println("error on setting", "error", result.Error.Error(), "dto", string(body))
		http.Error(w, fmt.Sprintf("Internal Server Error, Reason: %s", result.Error.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *Controller) delete(w http.ResponseWriter, r *http.Request) {
	n := time.Now()
	key := ""

	defer func() {
		log.Println("HTTP request delete value", "response_time", time.Since(n), "key", key)
	}()

	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	q, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, "Bad Request, Query Cannot Be Parsed", http.StatusBadRequest)
		return
	}

	key = q.Get("key")
	resultChan := c.service.ReceiveCommand(model.NewDeleteCommand(model.Key(key)))

	<-resultChan

	w.WriteHeader(200)
}
