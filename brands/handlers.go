package brands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Financial-Times/go-fthealth/v1a"
	"github.com/Financial-Times/service-status-go/gtg"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

// BrandHandler - struct for the handlers
type BrandHandler struct {
	service BrandService
}

// NewBrandHandler - Create a new BrandHandler
func NewBrandHandler(service BrandService) BrandHandler {
	return BrandHandler{service}
}

// GetBrands - Return a JSON encoded list of all brands
func (h *BrandHandler) GetBrands(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Add("Content-Type", "application/json")
	if !h.service.isInitialised() {
		writeStatusServiceUnavailable(writer)
		return
	}

	if c, _ := h.service.getCount(); c == 0 {
		writeJSONMessageWithStatus(writer, "Brands not found", http.StatusNotFound)
		return
	}

	pv, err := h.service.getBrands()

	if err != nil {
		writeJSONMessageWithStatus(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	defer pv.Close()
	writer.WriteHeader(http.StatusOK)
	io.Copy(writer, &pv)
}

// GetBrandUUIDs - Get a list of JSON objects (not a JSON list) giving each id.
func (h *BrandHandler) GetBrandUUIDs(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Add("Content-Type", "application/json")
	if !h.service.isInitialised() {
		writeStatusServiceUnavailable(writer)
		return
	}

	if c, _ := h.service.getCount(); c == 0 {
		writeJSONMessageWithStatus(writer, "Brands not found", http.StatusNotFound)
		return
	}

	pv, err := h.service.getBrandUUIDs()

	if err != nil {
		writeJSONMessageWithStatus(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	defer pv.Close()
	writer.WriteHeader(http.StatusOK)
	io.Copy(writer, &pv)
}

// GetCount - Get a count of the number of available brands
func (h *BrandHandler) GetCount(writer http.ResponseWriter, req *http.Request) {
	if !h.service.isInitialised() {
		writer.Header().Add("Content-Type", "application/json")
		writeStatusServiceUnavailable(writer)
		return
	}
	count, err := h.service.getCount()
	if err != nil {
		writer.Header().Add("Content-Type", "application/json")
		writeJSONMessageWithStatus(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Write([]byte(strconv.Itoa(count)))
}

// HealthCheck - Return FT standard healthcheck
func (h *BrandHandler) HealthCheck() v1a.Check {

	return v1a.Check{
		BusinessImpact:   "Unable to respond to requests",
		Name:             "Check service has finished initilising.",
		PanicGuide:       "https://sites.google.com/a/ft.com/ft-technology-service-transition/home/run-book-library/v1-brands-transformer",
		Severity:         1,
		TechnicalSummary: "Cannot serve any content as data not loaded.",
		Checker: func() (string, error) {
			if h.service.isInitialised() {
				return "Service is up and running", nil
			}
			return "Error as service initilising", errors.New("Service is initilising.")
		},
	}
}

// G2GCheck - Return FT standard good-to-go check
func (h *BrandHandler) G2GCheck() gtg.Status {
	count, err := h.service.getCount()
	if h.service.isInitialised() && err == nil && count > 0 {
		return gtg.Status{GoodToGo: true}
	}
	return gtg.Status{GoodToGo: false}
}

// GetBrandByUUID - Return the JSON for a single brand
func (h *BrandHandler) GetBrandByUUID(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Add("Content-Type", "application/json")
	if !h.service.isInitialised() {
		writeStatusServiceUnavailable(writer)
		return
	}

	vars := mux.Vars(req)
	uuid := vars["uuid"]

	obj, found, err := h.service.getBrandByUUID(uuid)
	if err != nil {
		writeJSONMessageWithStatus(writer, err.Error(), http.StatusInternalServerError)
	}
	writeJSONResponse(obj, found, writer)
}

// Reload - Reload the cache with fresh information
func (h *BrandHandler) Reload(writer http.ResponseWriter, req *http.Request) {
	if !h.service.isInitialised() || !h.service.isDataLoaded() {
		writeStatusServiceUnavailable(writer)
		return
	}

	go func() {
		if err := h.service.reloadDB(); err != nil {
			log.Errorf("ERROR opening db: %v", err.Error())
		}
	}()
	writeJSONMessageWithStatus(writer, "Reloading brands", http.StatusAccepted)
}

func writeJSONResponse(obj interface{}, found bool, writer http.ResponseWriter) {
	if !found {
		writeJSONMessageWithStatus(writer, "Brand not found", http.StatusNotFound)
		return
	}

	enc := json.NewEncoder(writer)
	if err := enc.Encode(obj); err != nil {
		log.Errorf("Error on json encoding=%v", err)
		writeJSONMessageWithStatus(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func writeJSONMessageWithStatus(w http.ResponseWriter, msg string, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprintln(w, fmt.Sprintf("{\"message\": \"%s\"}", msg))
}

func writeStatusServiceUnavailable(w http.ResponseWriter) {
	writeJSONMessageWithStatus(w, "Service Unavailable", http.StatusServiceUnavailable)
}
