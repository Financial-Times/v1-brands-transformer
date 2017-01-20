package brands

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"

	"bytes"
	"github.com/Financial-Times/tme-reader/tmereader"
	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"sort"
)

type testSuiteForBrands struct {
	name  string
	uuid  string
	found bool
	err   error
}

type mockClient struct {
	resp []berthaBrand
	err  error
}

func (c *mockClient) Do(req *http.Request) (*http.Response, error) {
	b, e := json.Marshal(c.resp)
	if c.err == nil {
		c.err = e
	}
	cb := ioutil.NopCloser(bytes.NewReader(b))
	return &http.Response{Body: cb}, c.err
}

const (
	fredUuid = "132a00d6-966c-3afb-b5c6-35da4f0dd70e"
	bobUuid  = "89400620-0727-3b07-b39e-3e614c115706"
)

var defaultTypes = []string{"Thing", "Concept", "Brand"}

var bobTMEBrand = brand{
	UUID:                   bobUuid,
	PrefLabel:              "Bob",
	ParentUUID:             financialTimesBrandUuid,
	AlternativeIdentifiers: alternativeIdentifiers{TME: []string{"Ym9i-QnJhbmRz"}, UUIDs: []string{bobUuid}},
	Aliases:                []string{"Bob"},
	Type:                   "Brand",
}

var fredTMEBrand = brand{
	UUID:                   fredUuid,
	PrefLabel:              "Fred",
	ParentUUID:             financialTimesBrandUuid,
	AlternativeIdentifiers: alternativeIdentifiers{TME: []string{"ZnJlZA==-QnJhbmRz"}, UUIDs: []string{fredUuid}},
	Aliases:                []string{"Fred"},
	Type:                   "Brand",
}

var testBerthaBrand = berthaBrand{
	Active:              true,
	PrefLabel:           "Financial Times",
	Strapline:           "Make the right connections",
	DescriptionXML:      "<p>The Financial Times (FT) is one of the world’s leading business news and information organisations.</p>",
	ImageURL:            "http://aboutus.ft.com/files/2010/11/ft-logo.gif",
	TmeParentIdentifier: "TmeParentIdentifier",
	TmeIdentifier:       "1234567890",
}

var expectedBrand = brand{
	UUID:           "e807f1fc-f82d-332f-9bb0-18ca6738a19f",
	ParentUUID:     "17b1538f-eda4-3402-9304-98853fb58c4d",
	PrefLabel:      "Financial Times",
	Type:           "Brand",
	Strapline:      "Make the right connections",
	Description:    "The Financial Times (FT) is one of the world’s leading business news and information organisations.",
	DescriptionXML: "<p>The Financial Times (FT) is one of the world’s leading business news and information organisations.</p>",
	ImageURL:       "http://aboutus.ft.com/files/2010/11/ft-logo.gif",
	AlternativeIdentifiers: alternativeIdentifiers{
		UUIDs: []string{"e807f1fc-f82d-332f-9bb0-18ca6738a19f"},
		TME:   []string{"1234567890"},
	},
}

var testBerthaBrandFixedUUID = berthaBrand{
	Active:         true,
	PrefLabel:      "FT Data",
	Strapline:      "Strapline",
	DescriptionXML: "<p>DescriptionXML</p>",
	ImageURL:       "http://some.ft.com/image/url",
	TmeIdentifier:  "MGY2ZTQ3MTYtYjJiNS00ODVhLTlkYTktNzZlNzc3YTcxOWYy-QnJhbmRz",
}

var expectedBrandFixedUUID = brand{
	UUID:           "b8513403-7892-4901-bb97-1765fc0ba190",
	ParentUUID:     financialTimesBrandUuid,
	PrefLabel:      "FT Data",
	Type:           "Brand",
	Strapline:      "Strapline",
	Description:    "DescriptionXML",
	DescriptionXML: "<p>DescriptionXML</p>",
	ImageURL:       "http://some.ft.com/image/url",
	Aliases:        []string{"FT Data"},
	AlternativeIdentifiers: alternativeIdentifiers{
		UUIDs: []string{"c4316c4a-da19-3a29-bf48-75761174756f", "b8513403-7892-4901-bb97-1765fc0ba190"},
		TME:   []string{"MGY2ZTQ3MTYtYjJiNS00ODVhLTlkYTktNzZlNzc3YTcxOWYy-QnJhbmRz"},
	},
}

var testBerthaBrandForFT = berthaBrand{
	Active:         true,
	PrefLabel:      "Financial Times",
	Strapline:      "Make the right connections",
	DescriptionXML: "<p>The Financial Times (FT) is one of the world’s leading business news and information organisations, recognised internationally for its authority, integrity and accuracy.</p>",
	ImageURL:       "http://aboutus.ft.com/files/2010/11/ft-logo.gif",
	TmeIdentifier:  financialTimesBrandUuid,
}

var expectedBrandFT = brand{
	UUID:           financialTimesBrandUuid,
	PrefLabel:      "Financial Times",
	Type:           "Brand",
	Strapline:      "Make the right connections",
	Description:    "The Financial Times (FT) is one of the world’s leading business news and information organisations, recognised internationally for its authority, integrity and accuracy.",
	DescriptionXML: "<p>The Financial Times (FT) is one of the world’s leading business news and information organisations, recognised internationally for its authority, integrity and accuracy.</p>",
	ImageURL:       "http://aboutus.ft.com/files/2010/11/ft-logo.gif",
	AlternativeIdentifiers: alternativeIdentifiers{
		UUIDs: []string{financialTimesBrandUuid},
	},
}

func TestInit(t *testing.T) {
	repo := blockingRepo{}
	repo.Add(1)
	tmpfile := getTempFile(t)
	defer os.Remove(tmpfile.Name())
	service := createTestBrandService(&repo, tmpfile.Name())
	defer func() {
		repo.Done()
		service.Shutdown()
	}()
	assert.False(t, service.isDataLoaded())
	assert.True(t, service.isInitialised())
}

func TestGetBrands(t *testing.T) {
	tmpfile := getTempFile(t)
	defer os.Remove(tmpfile.Name())
	repo := dummyRepo{terms: []term{{CanonicalName: "Bob", RawID: "bob"}, {CanonicalName: "Fred", RawID: "fred"}}}
	service := createTestBrandService(&repo, tmpfile.Name())
	defer service.Shutdown()
	waitTillInit(t, service)
	waitTillDataLoaded(t, service)
	pv, err := service.getBrands()

	var wg sync.WaitGroup
	var res []brand
	wg.Add(1)
	go func(reader io.Reader, w *sync.WaitGroup) {
		var err error
		scan := bufio.NewScanner(reader)
		for scan.Scan() {
			var p brand
			assert.NoError(t, err)
			err = json.Unmarshal(scan.Bytes(), &p)
			assert.NoError(t, err)
			res = append(res, p)
		}
		wg.Done()
	}(&pv, &wg)
	wg.Wait()

	assert.NoError(t, err)
	assert.Len(t, res, 2)

	compareBrands(res[1], bobTMEBrand, t)
	compareBrands(res[0], fredTMEBrand, t)
}

func TestGetBrandUUIDs(t *testing.T) {
	tmpfile := getTempFile(t)
	defer os.Remove(tmpfile.Name())
	repo := dummyRepo{terms: []term{{CanonicalName: "Bob", RawID: "bob"}, {CanonicalName: "Fred", RawID: "fred"}}}
	service := createTestBrandService(&repo, tmpfile.Name())
	defer service.Shutdown()
	waitTillInit(t, service)
	waitTillDataLoaded(t, service)
	pv, err := service.getBrandUUIDs()

	var wg sync.WaitGroup
	var res []brandUUID
	wg.Add(1)
	go func(reader io.Reader, w *sync.WaitGroup) {
		var err error
		scan := bufio.NewScanner(reader)
		for scan.Scan() {
			var p brandUUID
			assert.NoError(t, err)
			err = json.Unmarshal(scan.Bytes(), &p)
			assert.NoError(t, err)
			res = append(res, p)
		}
		wg.Done()
	}(&pv, &wg)
	wg.Wait()

	assert.NoError(t, err)
	assert.Len(t, res, 2)

	assert.Equal(t, "132a00d6-966c-3afb-b5c6-35da4f0dd70e", res[0].UUID)
	assert.Equal(t, "89400620-0727-3b07-b39e-3e614c115706", res[1].UUID)
}

func TestGetBrandLinks(t *testing.T) {
	tmpfile := getTempFile(t)
	defer os.Remove(tmpfile.Name())
	repo := dummyRepo{terms: []term{{CanonicalName: "Bob", RawID: "bob"}, {CanonicalName: "Fred", RawID: "fred"}}}
	service := createTestBrandService(&repo, tmpfile.Name())
	defer service.Shutdown()
	waitTillInit(t, service)
	waitTillDataLoaded(t, service)
	pv, err := service.getBrandLinks()

	var wg sync.WaitGroup
	var res []brandLink
	wg.Add(1)
	go func(reader io.Reader, w *sync.WaitGroup) {
		var err error
		jsonBlob, err := ioutil.ReadAll(reader)
		assert.NoError(t, err)
		log.Infof("Got bytes: %v", string(jsonBlob[:]))
		err = json.Unmarshal(jsonBlob, &res)
		assert.NoError(t, err)
		wg.Done()
	}(&pv, &wg)
	wg.Wait()

	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, "/base/url/132a00d6-966c-3afb-b5c6-35da4f0dd70e", res[0].APIURL)
	assert.Equal(t, "/base/url/89400620-0727-3b07-b39e-3e614c115706", res[1].APIURL)
}

func TestGetCount(t *testing.T) {
	tmpfile := getTempFile(t)
	defer os.Remove(tmpfile.Name())
	repo := dummyRepo{terms: []term{{CanonicalName: "Bob", RawID: "bob"}, {CanonicalName: "Fred", RawID: "fred"}}}
	service := createTestBrandService(&repo, tmpfile.Name())
	defer service.Shutdown()
	waitTillInit(t, service)
	waitTillDataLoaded(t, service)
	assertCount(t, service, 2)
}

func TestReload(t *testing.T) {
	tmpfile := getTempFile(t)
	defer os.Remove(tmpfile.Name())
	repo := dummyRepo{terms: []term{{CanonicalName: "Bob", RawID: "bob"}, {CanonicalName: "Fred", RawID: "fred"}}}
	service := createTestBrandService(&repo, tmpfile.Name())
	defer service.Shutdown()
	waitTillInit(t, service)
	waitTillDataLoaded(t, service)
	assertCount(t, service, 2)
	repo.terms = append(repo.terms, term{CanonicalName: "Third", RawID: "third"})
	repo.count = 0
	assert.NoError(t, service.reloadDB())
	waitTillInit(t, service)
	waitTillDataLoaded(t, service)
	assertCount(t, service, 3)
}

func TestGetBrandByUUID(t *testing.T) {
	tmpfile := getTempFile(t)
	defer os.Remove(tmpfile.Name())
	repo := dummyRepo{terms: []term{{CanonicalName: "Bob", RawID: "bob"}, {CanonicalName: "Fred", RawID: "fred"}}}
	service := createTestBrandService(&repo, tmpfile.Name())
	defer service.Shutdown()
	waitTillInit(t, service)
	waitTillDataLoaded(t, service)

	tests := []testSuiteForBrands{
		{"Success", "132a00d6-966c-3afb-b5c6-35da4f0dd70e", true, nil},
		{"Success", "xxxxxxxx-bb56-363d-80c1-f2d957ef58cf", false, nil}}
	for _, test := range tests {
		brand, found, err := service.getBrandByUUID(test.uuid)
		if test.err != nil {
			assert.Equal(t, test.err, err)
		} else if test.found {
			assert.True(t, found)
			assert.NotNil(t, brand)
		} else {
			assert.False(t, found)
		}
	}
}

func TestFailingOpeningDB(t *testing.T) {
	dir, err := ioutil.TempDir("", "service_test")
	assert.NoError(t, err)
	service := createTestBrandService(&dummyRepo{}, dir)
	defer service.Shutdown()
	for i := 1; i <= 1000; i++ {
		if !service.isInitialised() {
			log.Info("isInitialised was false")
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	assert.False(t, service.isInitialised(), "isInitialised should be false")
}

func TestBerthaToBrand(t *testing.T) {
	actualBrand, err := berthaToBrand(testBerthaBrand, "e807f1fc-f82d-332f-9bb0-18ca6738a19f")
	assert.Equal(t, expectedBrand, actualBrand)
	assert.Nil(t, err)
}

func TestAddingInformationFromBerthaWorks(t *testing.T) {
	emptyBrand := brand{
		UUID:      "e807f1fc-f82d-332f-9bb0-18ca6738a19f",
		PrefLabel: "Fred Black",
		AlternativeIdentifiers: alternativeIdentifiers{
			UUIDs: []string{"e807f1fc-f82d-332f-9bb0-18ca6738a19f"},
			TME:   []string{"1234567890"},
		},
	}

	actualBrand, err := addBerthaInformation(emptyBrand, testBerthaBrand)
	assert.Equal(t, expectedBrand, actualBrand)
	assert.Nil(t, err)
}

func TestLoadingCuratedBrands(t *testing.T) {
	tmpfile := getTempFile(t)
	defer os.Remove(tmpfile.Name())

	brandService := NewBrandService(&dummyRepo{}, "/base/url", "Brands", 1, tmpfile.Name(), "/bertha/url", &mockClient{})
	input := []berthaBrand{testBerthaBrand}

	waitTillInit(t, brandService)
	waitTillDataLoaded(t, brandService)

	brandService.loadCuratedBrands(input)
	actualOutput, found, err := brandService.getBrandByUUID("e807f1fc-f82d-332f-9bb0-18ca6738a19f")
	assert.Equal(t, true, found)
	assert.EqualValues(t, expectedBrand, actualOutput)
	assert.Nil(t, err)
}

func TestNoTmeIdentifierWhenLoadingCuratedBrandsThenBrandIsIgnored(t *testing.T) {
	tmpfile := getTempFile(t)
	defer os.Remove(tmpfile.Name())

	brandService := NewBrandService(&dummyRepo{}, "/base/url", "Brands", 1, tmpfile.Name(), "/bertha/url", &mockClient{})

	testBerthaBrandWithTme := berthaBrand{
		Active:              true,
		PrefLabel:           "Funky Chicken",
		TmeParentIdentifier: "some TmeParentIdentifier",
	}

	input := []berthaBrand{testBerthaBrandWithTme}

	waitTillInit(t, brandService)
	waitTillDataLoaded(t, brandService)

	brandService.loadCuratedBrands(input)

	actualOutput, err := brandService.getCount()
	assert.NoError(t, err)
	assert.EqualValues(t, 0, actualOutput)
}

func TestSetsParentBrandToBeFinacialTimesIfNoOtherParentSetFromBertha(t *testing.T) {
	tmpfile := getTempFile(t)
	defer os.Remove(tmpfile.Name())

	brandService := NewBrandService(&dummyRepo{}, "/base/url", "Brands", 1, tmpfile.Name(), "/bertha/url", &mockClient{})
	input := []berthaBrand{testBerthaBrand}
	waitTillInit(t, brandService)
	waitTillDataLoaded(t, brandService)

	brandService.loadCuratedBrands(input)
	actualOutput, found, err := brandService.getBrandByUUID("e807f1fc-f82d-332f-9bb0-18ca6738a19f")
	assert.Equal(t, true, found)
	assert.EqualValues(t, expectedBrand, actualOutput)
	assert.Nil(t, err)
}

//  Replaces data from TME with Bertha updates
func TestBerthaLoadedCorrectly(t *testing.T) {
	tmpfile := getTempFile(t)
	defer os.Remove(tmpfile.Name())
	repo := dummyRepo{terms: []term{{CanonicalName: "awesome brand", RawID: "some tme identifier"}, {CanonicalName: "FT Data", RawID: "0f6e4716-b2b5-485a-9da9-76e777a719f2"}}}
	client := mockClient{resp: []berthaBrand{testBerthaBrand, testBerthaBrandFixedUUID}}
	brandService := NewBrandService(&repo, "/base/url", "Brands", 1, tmpfile.Name(), "/bertha/url", &client)

	waitTillInit(t, brandService)
	waitTillDataLoaded(t, brandService)

	reader, err := brandService.getBrands()
	assert.NoError(t, err)
	printBrands(&reader)
	assertCount(t, brandService, 3)
	actualBrand, found, err := brandService.getBrandByUUID(expectedBrandFixedUUID.UUID)
	assert.True(t, found)
	assert.NoError(t, err)
	assert.EqualValues(t, expectedBrandFixedUUID, actualBrand)
}

func TestBerthaLoadedFTCorrectly(t *testing.T) {
	tmpfile := getTempFile(t)
	defer os.Remove(tmpfile.Name())
	repo := dummyRepo{terms: []term{}}
	client := mockClient{resp: []berthaBrand{testBerthaBrandForFT}}
	brandService := NewBrandService(&repo, "/base/url", "Brands", 1, tmpfile.Name(), "/bertha/url", &client)

	waitTillInit(t, brandService)
	waitTillDataLoaded(t, brandService)

	reader, err := brandService.getBrands()
	assert.NoError(t, err)
	printBrands(&reader)
	assertCount(t, brandService, 1)
	actualBrand, found, err := brandService.getBrandByUUID(financialTimesBrandUuid)
	assert.True(t, found)
	assert.NoError(t, err)
	assert.EqualValues(t, expectedBrandFT, actualBrand)
}

func assertCount(t *testing.T, s BrandService, expected int) {
	count, err := s.getCount()
	assert.NoError(t, err)
	assert.Equal(t, expected, count)
}

func createTestBrandService(repo tmereader.Repository, cacheFileName string) BrandService {
	return NewBrandService(repo, "/base/url", "Brands", 1, cacheFileName, "http://bertha/url", &mockClient{})
}

func getTempFile(t *testing.T) *os.File {
	tmpfile, err := ioutil.TempFile("", "example")
	assert.NoError(t, err)
	assert.NoError(t, tmpfile.Close())
	log.Debug("File:%s", tmpfile.Name())
	return tmpfile
}

func waitTillInit(t *testing.T, s BrandService) {
	for i := 1; i <= 1000; i++ {
		if s.isInitialised() {
			log.Info("isInitialised was true")
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	assert.True(t, s.isInitialised())
}

func waitTillDataLoaded(t *testing.T, s BrandService) {
	for i := 1; i <= 1000; i++ {
		if s.isDataLoaded() {
			log.Info("isDataLoaded was true")
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	assert.True(t, s.isDataLoaded())
}

type dummyRepo struct {
	sync.Mutex
	terms []term
	err   error
	count int
}

func (d *dummyRepo) GetTmeTermsFromIndex(startRecord int) ([]interface{}, error) {
	defer func() {
		d.count++
	}()
	if len(d.terms) == d.count {
		return nil, d.err
	}
	return []interface{}{d.terms[d.count]}, d.err
}

// Never used
func (d *dummyRepo) GetTmeTermById(uuid string) (interface{}, error) {
	return nil, nil
}

type blockingRepo struct {
	sync.WaitGroup
	err  error
	done bool
}

func (d *blockingRepo) GetTmeTermsFromIndex(startRecord int) ([]interface{}, error) {
	d.Wait()
	if d.done {
		return nil, d.err
	}
	d.done = true
	return []interface{}{term{CanonicalName: "Bob", RawID: "bob"}}, d.err
}

// Never used
func (d *blockingRepo) GetTmeTermById(uuid string) (interface{}, error) {
	return nil, nil
}

func printBrands(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		log.Infof("%s \n", scanner.Text())
	}

}
func compareBrands(actual, expected brand, t *testing.T) {
	sort.Strings(expected.AlternativeIdentifiers.TME)
	sort.Strings(expected.AlternativeIdentifiers.UUIDs)
	sort.Strings(actual.AlternativeIdentifiers.TME)
	sort.Strings(actual.AlternativeIdentifiers.UUIDs)

	assert.EqualValues(t, expected, actual)
}
