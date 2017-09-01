package brands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/Financial-Times/go-fthealth/v1a"
	"github.com/Financial-Times/service-status-go/gtg"
	status "github.com/Financial-Times/service-status-go/httphandlers"
	log "github.com/sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const (
	testUUID              = "bba39990-c78d-3629-ae83-808c333c6dbc"
	testUUID2             = "be2e7e2b-0fa2-3969-a69b-74c46e754032"
	getBrandResponse      = "{\"uuid\":\"bba39990-c78d-3629-ae83-808c333c6dbc\",\"alternativeIdentifiers\":{}}\n{\"uuid\":\"be2e7e2b-0fa2-3969-a69b-74c46e754032\",\"alternativeIdentifiers\":{}}\n"
	getBrandUUIDsResponse = `{"ID":"bba39990-c78d-3629-ae83-808c333c6dbc"}
{"ID":"be2e7e2b-0fa2-3969-a69b-74c46e754032"}
`
	getBrandByUUIDResponse = "{\"uuid\":\"bba39990-c78d-3629-ae83-808c333c6dbc\",\"prefLabel\":\"Financial Times\",\"type\":\"Brand\",\"alternativeIdentifiers\":{\"TME\":[\"RlQK-QnJhbmRzCg==\"],\"uuids\":[\"bba39990-c78d-3629-ae83-808c333c6dbc\"]}}\n"
)

func TestHandlers(t *testing.T) {
	var wg sync.WaitGroup
	tests := []struct {
		name         string
		req          *http.Request
		dummyService BrandService
		statusCode   int
		contentType  string // Contents of the Content-Type header
		body         string
	}{
		{"Success - get brand by uuid",
			newRequest("GET", fmt.Sprintf("/transformers/brands/%s", testUUID)),
			&dummyService{
				found:       true,
				initialised: true,
				dataLoaded:  true,
				brands:      []brand{{UUID: testUUID, PrefLabel: "Financial Times", AlternativeIdentifiers: alternativeIdentifiers{UUIDs: []string{testUUID}, TME: []string{"RlQK-QnJhbmRzCg=="}}, Type: "Brand"}}},
			http.StatusOK,
			"application/json",
			getBrandByUUIDResponse},
		{"405 - get brand by uuid",
			newRequest("POST", fmt.Sprintf("/transformers/brands/%s", testUUID)),
			&dummyService{
				found:       true,
				initialised: true,
				dataLoaded:  true,
				brands:      []brand{{UUID: testUUID, PrefLabel: "Financial Times", AlternativeIdentifiers: alternativeIdentifiers{UUIDs: []string{testUUID}, TME: []string{"RlQK-QnJhbmRzCg=="}}, Type: "Brand"}}},
			http.StatusMethodNotAllowed,
			"application/json",
			""},
		{"Not found - get brand by uuid",
			newRequest("GET", fmt.Sprintf("/transformers/brands/%s", testUUID)),
			&dummyService{
				found:       false,
				initialised: true,
				dataLoaded:  true,
				brands:      []brand{{}}},
			http.StatusNotFound,
			"application/json",
			"{\"message\": \"Brand not found\"}\n"},
		{"Service unavailable - get brand by uuid",
			newRequest("GET", fmt.Sprintf("/transformers/brands/%s", testUUID)),
			&dummyService{
				found:       false,
				initialised: false,
				brands:      []brand{}},
			http.StatusServiceUnavailable,
			"application/json",
			"{\"message\": \"Service Unavailable\"}\n"},
		{"Success - get brands count",
			newRequest("GET", "/transformers/brands/__count"),
			&dummyService{
				found:       true,
				count:       1,
				initialised: true,
				dataLoaded:  true,
				brands:      []brand{{UUID: testUUID}}},
			http.StatusOK,
			"application/json",
			"1"},
		{"405 - get brands count",
			newRequest("POST", "/transformers/brands/__count"),
			&dummyService{
				found:       true,
				count:       1,
				initialised: true,
				dataLoaded:  true,
				brands:      []brand{{UUID: testUUID}}},
			http.StatusMethodNotAllowed,
			"application/json",
			""},
		{"Failure - get brands count",
			newRequest("GET", "/transformers/brands/__count"),
			&dummyService{
				err:         errors.New("Something broke"),
				found:       true,
				count:       1,
				initialised: true,
				dataLoaded:  true,
				brands:      []brand{{UUID: testUUID}}},
			http.StatusInternalServerError,
			"application/json",
			"{\"message\": \"Something broke\"}\n"},
		{"Failure - get brands count not init",
			newRequest("GET", "/transformers/brands/__count"),
			&dummyService{
				err:         errors.New("Something broke"),
				found:       true,
				count:       1,
				initialised: false,
				brands:      []brand{{UUID: testUUID}}},
			http.StatusServiceUnavailable,
			"application/json", "{\"message\": \"Service Unavailable\"}\n"},
		{"get brands - success",
			newRequest("GET", "/transformers/brands"),
			&dummyService{
				found:       true,
				initialised: true,
				dataLoaded:  true,
				count:       2,
				brands:      []brand{{UUID: testUUID}, {UUID: testUUID2}}},
			http.StatusOK,
			"application/json",
			getBrandResponse},
		{"get brands - 405",
			newRequest("POST", "/transformers/brands"),
			&dummyService{
				found:       true,
				initialised: true,
				dataLoaded:  true,
				count:       2,
				brands:      []brand{{UUID: testUUID}, {UUID: testUUID2}}},
			http.StatusMethodNotAllowed,
			"application/json",
			""},
		{"get brands - Not found",
			newRequest("GET", "/transformers/brands"),
			&dummyService{
				initialised: true,
				dataLoaded:  true,
				count:       0,
				brands:      []brand{}},
			http.StatusNotFound,
			"application/json",
			"{\"message\": \"Brands not found\"}\n"},
		{"get brands - Service unavailable",
			newRequest("GET", "/transformers/brands"),
			&dummyService{
				found:       false,
				initialised: false,
				brands:      []brand{}},
			http.StatusServiceUnavailable,
			"application/json",
			"{\"message\": \"Service Unavailable\"}\n"},
		{"get brands IDS - Success",
			newRequest("GET", "/transformers/brands/__ids"),
			&dummyService{
				found:       true,
				initialised: true,
				dataLoaded:  true,
				count:       1,
				brands:      []brand{{UUID: testUUID}, {UUID: testUUID2}}},
			http.StatusOK,
			"application/json",
			getBrandUUIDsResponse},
		{"get brands IDS - Not found",
			newRequest("GET", "/transformers/brands/__ids"),
			&dummyService{
				initialised: true,
				dataLoaded:  true,
				count:       0,
				brands:      []brand{}},
			http.StatusNotFound,
			"application/json",
			"{\"message\": \"Brands not found\"}\n"},
		{"get brands IDS - Service unavailable",
			newRequest("GET", "/transformers/brands/__ids"),
			&dummyService{
				found:       false,
				initialised: false,
				brands:      []brand{}},
			http.StatusServiceUnavailable,
			"application/json",
			"{\"message\": \"Service Unavailable\"}\n"},
		{"GTG unavailable - get GTG",
			newRequest("GET", status.GTGPath),
			&dummyService{
				found:       false,
				initialised: false,
				dataLoaded:  false,
				brands:      []brand{}},
			http.StatusServiceUnavailable,
			"application/json",
			""},
		{"GTG unavailable - get GTG but no brands",
			newRequest("GET", status.GTGPath),
			&dummyService{
				found:       false,
				initialised: true,
				dataLoaded:  true,
			},
			http.StatusOK,
			"application/json",
			"OK"},
		{"GTG OK - get GTG",
			newRequest("GET", status.GTGPath),
			&dummyService{
				found:       true,
				initialised: true,
				dataLoaded:  true,
				count:       2},
			http.StatusOK,
			"application/json",
			"OK"},
		{"Health bad - get Health check",
			newRequest("GET", "/__health"),
			&dummyService{
				found:       false,
				initialised: false,
				dataLoaded:  true},
			http.StatusOK,
			"application/json",
			"regex=Service is initilising"},
		{"Health good - get Health check",
			newRequest("GET", "/__health"),
			&dummyService{
				found:       false,
				initialised: true,
				dataLoaded:  true},
			http.StatusOK,
			"application/json",
			"regex=Service is up and running"},
		{"Reload accepted - request reload",
			newRequest("POST", "/transformers/brands/__reload"),
			&dummyService{
				wg:          &wg,
				initialised: true,
				dataLoaded:  true},
			http.StatusAccepted,
			"application/json",
			"{\"message\": \"Reloading brands\"}\n"},
		{"Reload accepted even though error loading data in background.",
			newRequest("POST", "/transformers/brands/__reload"),
			&dummyService{
				wg:          &wg,
				err:         errors.New("Boom goes the backend..."),
				initialised: true,
				dataLoaded:  true},
			http.StatusAccepted,
			"application/json",
			"{\"message\": \"Reloading brands\"}\n"},
	}
	for _, test := range tests {
		wg.Add(1)
		rec := httptest.NewRecorder()
		router(test.dummyService).ServeHTTP(rec, test.req)
		assert.Equal(t, test.statusCode, rec.Code, fmt.Sprintf("%s: Wrong response code, was %d, should be %d", test.name, rec.Code, test.statusCode))

		b, err := ioutil.ReadAll(rec.Body)
		assert.NoError(t, err)
		body := string(b)
		if strings.HasPrefix(test.body, "regex=") {
			regex := strings.TrimPrefix(test.body, "regex=")
			matched, err := regexp.MatchString(regex, body)
			assert.NoError(t, err)
			assert.True(t, matched, fmt.Sprintf("Could not match regex:\n %s \nin body:\n %s", regex, body))
		} else {
			assert.Equal(t, test.body, body, fmt.Sprintf("%s: Wrong body", test.name))
		}
	}
}

func TestReloadIsCalled(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	rec := httptest.NewRecorder()
	s := &dummyService{
		wg:          &wg,
		found:       true,
		initialised: true,
		dataLoaded:  true,
		count:       2,
		brands:      []brand{}}
	log.Infof("s.loadDBCalled: %v", s.loadDBCalled)
	router(s).ServeHTTP(rec, newRequest("POST", "/transformers/brands/__reload"))
	wg.Wait()
	assert.True(t, s.loadDBCalled)
}

func newRequest(method, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	return req
}

type dummyService struct {
	found        bool
	brands       []brand
	initialised  bool
	dataLoaded   bool
	count        int
	err          error
	loadDBCalled bool
	wg           *sync.WaitGroup
}

func (s *dummyService) loadCuratedBrands(bBrands []berthaBrand) error {
	return nil
}

func (s *dummyService) getBrands() (io.PipeReader, error) {
	pv, pw := io.Pipe()
	go func() {
		encoder := json.NewEncoder(pw)
		for _, sub := range s.brands {
			encoder.Encode(sub)
		}
		pw.Close()
	}()
	return *pv, nil
}

func (s *dummyService) getBrandUUIDs() (io.PipeReader, error) {
	pv, pw := io.Pipe()
	go func() {
		encoder := json.NewEncoder(pw)
		for _, sub := range s.brands {
			encoder.Encode(brandUUID{UUID: sub.UUID})
		}
		pw.Close()
	}()
	return *pv, nil
}

func (s *dummyService) getBrandLinks() (io.PipeReader, error) {
	pv, pw := io.Pipe()
	go func() {
		var links []brandLink
		for _, sub := range s.brands {
			links = append(links, brandLink{APIURL: "http://localhost:8080/transformers/brands/" + sub.UUID})
		}
		b, _ := json.Marshal(links)
		log.Infof("Writing bytes... %v", string(b))
		pw.Write(b)
		pw.Close()
	}()
	return *pv, nil
}

func (s *dummyService) getCount() (int, error) {
	return s.count, s.err
}

func (s *dummyService) getBrandByUUID(uuid string) (brand, bool, error) {
	return s.brands[0], s.found, nil
}

func (s *dummyService) isInitialised() bool {
	return s.initialised
}

func (s *dummyService) isDataLoaded() bool {
	return s.dataLoaded
}

func (s *dummyService) Shutdown() error {
	return s.err
}

func (s *dummyService) reloadDB() error {
	defer s.wg.Done()
	s.loadDBCalled = true
	return s.err
}

func router(s BrandService) *mux.Router {
	handler := NewBrandHandler(s)
	servicesRouter := mux.NewRouter()

	getBrandsSubrouter := servicesRouter.Path("/transformers/brands").Subrouter()
	getBrandsSubrouter.Methods("GET").HandlerFunc(handler.GetBrands)
	getBrandsSubrouter.NewRoute().HandlerFunc(handler.OnlyGetAllowed)

	brandCountSubrouter := servicesRouter.Path("/transformers/brands/__count").Subrouter()
	brandCountSubrouter.Methods("GET").HandlerFunc(handler.GetCount)
	brandCountSubrouter.NewRoute().HandlerFunc(handler.OnlyGetAllowed)

	brandIDsSubrouter := servicesRouter.Path("/transformers/brands/__ids").Subrouter()
	brandIDsSubrouter.Methods("GET").HandlerFunc(handler.GetBrandUUIDs)
	brandIDsSubrouter.NewRoute().HandlerFunc(handler.OnlyGetAllowed)

	brandByUUIDSubrouter := servicesRouter.Path("/transformers/brands/{uuid:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}").Subrouter()
	brandByUUIDSubrouter.Methods("GET").HandlerFunc(handler.GetBrandByUUID)
	brandByUUIDSubrouter.NewRoute().HandlerFunc(handler.OnlyGetAllowed)

	reloadSubrouter := servicesRouter.Path("/transformers/brands/__reload").Subrouter()
	reloadSubrouter.Methods("POST").HandlerFunc(handler.Reload)
	reloadSubrouter.NewRoute().HandlerFunc(handler.OnlyPostAllowed)

	servicesRouter.HandleFunc("/__health", v1a.Handler("V1 Brands Transformer Healthchecks", "Checks for the health of the service", handler.HealthCheck()))
	g2gHandler := status.NewGoodToGoHandler(gtg.StatusChecker(handler.G2GCheck))
	servicesRouter.HandleFunc(status.GTGPath, g2gHandler)
	return servicesRouter
}
