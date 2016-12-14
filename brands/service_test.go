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

	"github.com/Financial-Times/tme-reader/tmereader"
	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type testSuiteForBrands struct {
	name  string
	uuid  string
	found bool
	err   error
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
	assert.Equal(t, "28d66fcc-bb56-363d-80c1-f2d957ef58cf", res[0].UUID)
	assert.Equal(t, "be2e7e2b-0fa2-3969-a69b-74c46e754032", res[1].UUID)
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
	assert.Equal(t, "28d66fcc-bb56-363d-80c1-f2d957ef58cf", res[0].UUID)
	assert.Equal(t, "be2e7e2b-0fa2-3969-a69b-74c46e754032", res[1].UUID)
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
	assert.Equal(t, "/base/url/28d66fcc-bb56-363d-80c1-f2d957ef58cf", res[0].APIURL)
	assert.Equal(t, "/base/url/be2e7e2b-0fa2-3969-a69b-74c46e754032", res[1].APIURL)
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
		{"Success", "28d66fcc-bb56-363d-80c1-f2d957ef58cf", true, nil},
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

func assertCount(t *testing.T, s BrandService, expected int) {
	count, err := s.getCount()
	assert.NoError(t, err)
	assert.Equal(t, expected, count)
}

func createTestBrandService(repo tmereader.Repository, cacheFileName string) BrandService {
	return NewBrandService(repo, "/base/url", "taxonomy_string", 1, cacheFileName, "http://bertha/url")
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

func TestBerthaToBrand(t *testing.T) {
	testBrand := berthaBrand{
		Active:         true,
		PrefLabel:      "Financial Times",
		Strapline:      "Make the right connections",
		DescriptionXML: "<p>The Financial Times (FT) is one of the world’s leading business news and information organisations.</p>",
		ImageURL:       "http://aboutus.ft.com/files/2010/11/ft-logo.gif",
		UUID:           "dbb0bdae-1f0c-11e4-b0cb-b2227cce2b54",
		ParentUUID:     "dbb0bdae-1f0c-11e4-b0cb-846947257459",
		TmeIdentifier:  "1234567890",
	}
	expectedBrand := brand{
		UUID:           "dbb0bdae-1f0c-11e4-b0cb-b2227cce2b54",
		ParentUUID:     "dbb0bdae-1f0c-11e4-b0cb-846947257459",
		PrefLabel:      "Financial Times",
		Type:           "Brand",
		Strapline:      "Make the right connections",
		Description:    "The Financial Times (FT) is one of the world’s leading business news and information organisations.",
		DescriptionXML: "<p>The Financial Times (FT) is one of the world’s leading business news and information organisations.</p>",
		ImageURL:       "http://aboutus.ft.com/files/2010/11/ft-logo.gif",
		AlternativeIdentifiers: alternativeIdentifiers{
			UUIDs: []string{"dbb0bdae-1f0c-11e4-b0cb-b2227cce2b54"},
			TME:   []string{"1234567890"},
		},
	}

	actualBrand, err := berthaToBrand(testBrand)
	assert.Equal(t, expectedBrand, actualBrand)
	assert.Nil(t, err)
}

// func TestBadAddBertha(t *testing.T) {
// 	testBrand := berthaBrand{
// 		Active:         true,
// 		PrefLabel:      "Financial Times",
// 		Strapline:      "Make the right connections",
// 		DescriptionXML: "<p>The Financial Times (FT) is one of the world’s leading business news and information organisations.</p>",
// 		ImageURL:       "http://aboutus.ft.com/files/2010/11/ft-logo.gif",
// 		UUID:           "dbb0bdae-1f0c-11e4-b0cb-b2227cce2b54",
// 		ParentUUID:     "dbb0bdae-1f0c-11e4-b0cb-846947257459",
// 		TmeIdentifier:  "1234567890",
// 	}
// 	emptyBrand := brand{}

// 	_, err := addBerthaInformation(emptyBrand, testBrand)
// 	assert.EqualError(t, err, "Bertha UUID doesn't match brand UUID")
// }

func TestGoodAddBertha(t *testing.T) {
	testBrand := berthaBrand{
		Active:         true,
		PrefLabel:      "Financial Times",
		Strapline:      "Make the right connections",
		DescriptionXML: "<p>The Financial Times (FT) is one of the world’s leading business news and information organisations.</p>",
		ImageURL:       "http://aboutus.ft.com/files/2010/11/ft-logo.gif",
		UUID:           "dbb0bdae-1f0c-11e4-b0cb-b2227cce2b54",
		ParentUUID:     "dbb0bdae-1f0c-11e4-b0cb-846947257459",
		TmeIdentifier:  "1234567890",
	}
	emptyBrand := brand{
		UUID:      "dbb0bdae-1f0c-11e4-b0cb-b2227cce2b54",
		PrefLabel: "Fred Black",
		AlternativeIdentifiers: alternativeIdentifiers{
			UUIDs: []string{"dbb0bdae-1f0c-11e4-b0cb-b2227cce2b54"},
			TME:   []string{"1234567890"},
		},
	}
	expectedBrand := brand{
		UUID:           "dbb0bdae-1f0c-11e4-b0cb-b2227cce2b54",
		ParentUUID:     "dbb0bdae-1f0c-11e4-b0cb-846947257459",
		PrefLabel:      "Financial Times",
		Type:           "Brand",
		Strapline:      "Make the right connections",
		Description:    "The Financial Times (FT) is one of the world’s leading business news and information organisations.",
		DescriptionXML: "<p>The Financial Times (FT) is one of the world’s leading business news and information organisations.</p>",
		ImageURL:       "http://aboutus.ft.com/files/2010/11/ft-logo.gif",
		AlternativeIdentifiers: alternativeIdentifiers{
			UUIDs: []string{"dbb0bdae-1f0c-11e4-b0cb-b2227cce2b54"},
			TME:   []string{"1234567890"},
		},
	}

	actualBrand, err := addBerthaInformation(emptyBrand, testBrand)
	assert.Equal(t, expectedBrand, actualBrand)
	assert.Nil(t, err)
}

func TestLoadingCuratedBrands(t *testing.T) {
	tmpfile := getTempFile(t)
	defer os.Remove(tmpfile.Name())

	brandService := NewBrandService(&dummyRepo{}, "/base/url", "taxonomy", 1, tmpfile.Name(), "/bertha/url")
	log.Info(brandService)
	input := []berthaBrand{
		berthaBrand{
			Active:         true,
			PrefLabel:      "Financial Times",
			Strapline:      "Make the right connections",
			DescriptionXML: "<p>The Financial Times (FT) is one of the world’s leading business news and information organisations.</p>",
			ImageURL:       "http://aboutus.ft.com/files/2010/11/ft-logo.gif",
			UUID:           "dbb0bdae-1f0c-11e4-b0cb-b2227cce2b54",
			ParentUUID:     "dbb0bdae-1f0c-11e4-b0cb-846947257459",
			TmeIdentifier:  "1234567890",
		},
	}
	expectedBrand := brand{
		UUID:           "dbb0bdae-1f0c-11e4-b0cb-b2227cce2b54",
		ParentUUID:     "dbb0bdae-1f0c-11e4-b0cb-846947257459",
		PrefLabel:      "Financial Times",
		Type:           "Brand",
		Strapline:      "Make the right connections",
		Description:    "The Financial Times (FT) is one of the world’s leading business news and information organisations.",
		DescriptionXML: "<p>The Financial Times (FT) is one of the world’s leading business news and information organisations.</p>",
		ImageURL:       "http://aboutus.ft.com/files/2010/11/ft-logo.gif",
		AlternativeIdentifiers: alternativeIdentifiers{
			UUIDs: []string{"dbb0bdae-1f0c-11e4-b0cb-b2227cce2b54"},
			TME:   []string{"1234567890"},
		},
	}
	waitTillInit(t, brandService)
	waitTillDataLoaded(t, brandService)

	brandService.loadCuratedBrands(input)
	actualOutput, found, err := brandService.getBrandByUUID("dbb0bdae-1f0c-11e4-b0cb-b2227cce2b54")
	assert.Equal(t, true, found)
	assert.EqualValues(t, expectedBrand, actualOutput)
	assert.Nil(t, err)
}
