package brands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/Financial-Times/tme-reader/tmereader"
	log "github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/jaytaylor/html2text"
)

const (
	cacheBucket = "brand"
)

// BrandService - interface for retrieving v1 brands
type BrandService interface {
	getBrands() (io.PipeReader, error)
	getBrandLinks() (io.PipeReader, error)
	getBrandUUIDs() (io.PipeReader, error)
	getBrandByUUID(uuid string) (brand, bool, error)
	getCount() (int, error)
	isInitialised() bool
	isDataLoaded() bool
	reloadDB() error
	Shutdown() error
	loadCuratedBrands([]berthaBrand) error
}

type brandServiceImpl struct {
	sync.RWMutex
	repository    tmereader.Repository
	baseURL       string
	taxonomyName  string
	maxTmeRecords int
	initialised   bool
	dataLoaded    bool
	cacheFileName string
	db            *bolt.DB
	berthaURL     string
}

// NewBrandService - create a new BrandService
func NewBrandService(repo tmereader.Repository, baseURL string, taxonomyName string, maxTmeRecords int, cacheFileName string, berthaURL string) BrandService {
	s := &brandServiceImpl{repository: repo, baseURL: baseURL, taxonomyName: taxonomyName, maxTmeRecords: maxTmeRecords, initialised: true, cacheFileName: cacheFileName, berthaURL: berthaURL}
	go func(service *brandServiceImpl) {
		err := service.loadDB()
		if err != nil {
			log.Errorf("Error while creating BrandService: [%v]", err.Error())
		}
		var bBrands []berthaBrand

		bBrands, err = getBerthaBrands(service.berthaURL)
		if err != nil {
			log.Errorf("Error on Bertha load: [%v]", err.Error())
		} else {
			err = service.loadCuratedBrands(bBrands)
			if err != nil {
				log.Errorf("Error while loading in the curated brands: [%v]", err.Error())
			}
		}

	}(s)
	return s
}

func (s *brandServiceImpl) isInitialised() bool {
	s.RLock()
	defer s.RUnlock()
	return s.initialised
}

func (s *brandServiceImpl) setInitialised(val bool) {
	s.Lock()
	s.initialised = val
	s.Unlock()
}

func (s *brandServiceImpl) isDataLoaded() bool {
	s.RLock()
	defer s.RUnlock()
	return s.dataLoaded
}

func (s *brandServiceImpl) setDataLoaded(val bool) {
	s.Lock()
	s.dataLoaded = val
	s.Unlock()
}

func (s *brandServiceImpl) Shutdown() error {
	log.Info("Shuting down...")
	s.Lock()
	defer s.Unlock()
	s.initialised = false
	s.dataLoaded = false
	if s.db == nil {
		return errors.New("DB not open")
	}
	return s.db.Close()
}

func (s *brandServiceImpl) getCount() (int, error) {
	s.RLock()
	defer s.RUnlock()
	if !s.isDataLoaded() {
		return 0, nil
	}

	var count int
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(cacheBucket))
		if bucket == nil {
			return fmt.Errorf("Bucket %v not found!", cacheBucket)
		}
		count = bucket.Stats().KeyN
		return nil
	})
	return count, err
}

func (s *brandServiceImpl) getBrands() (io.PipeReader, error) {
	s.RLock()
	pv, pw := io.Pipe()
	go func() {
		defer s.RUnlock()
		defer pw.Close()
		s.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(cacheBucket))
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				if _, err := pw.Write(v); err != nil {
					return err
				}
				io.WriteString(pw, "\n")
			}
			return nil
		})
	}()
	return *pv, nil
}

func (s *brandServiceImpl) getBrandUUIDs() (io.PipeReader, error) {
	s.RLock()
	pv, pw := io.Pipe()
	go func() {
		defer s.RUnlock()
		defer pw.Close()
		s.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(cacheBucket))
			c := b.Cursor()
			encoder := json.NewEncoder(pw)
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				if k == nil {
					break
				}
				pl := brandUUID{UUID: string(k[:])}
				if err := encoder.Encode(pl); err != nil {
					return err
				}
			}
			return nil
		})
	}()
	return *pv, nil
}

func (s *brandServiceImpl) getBrandLinks() (io.PipeReader, error) {
	s.RLock()
	pv, pw := io.Pipe()
	go func() {
		defer s.RUnlock()
		defer pw.Close()
		io.WriteString(pw, "[")
		s.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(cacheBucket))
			c := b.Cursor()
			encoder := json.NewEncoder(pw)
			var k []byte
			k, _ = c.First()
			for {
				if k == nil {
					break
				}
				pl := brandLink{APIURL: s.baseURL + "/" + string(k[:])}
				if err := encoder.Encode(pl); err != nil {
					return err
				}
				if k, _ = c.Next(); k != nil {
					io.WriteString(pw, ",")
				}
			}
			return nil
		})
		io.WriteString(pw, "]")
	}()
	return *pv, nil
}

func (s *brandServiceImpl) getBrandByUUID(uuid string) (brand, bool, error) {
	s.RLock()
	defer s.RUnlock()
	var cachedValue []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(cacheBucket))
		if bucket == nil {
			return fmt.Errorf("Bucket %v not found!", cacheBucket)
		}
		cachedValue = bucket.Get([]byte(uuid))
		return nil
	})

	if err != nil {
		log.Errorf("ERROR reading from cache file for [%v]: %v", uuid, err.Error())
		return brand{}, false, err
	}
	if len(cachedValue) == 0 {
		log.Infof("INFO No cached value for [%v].", uuid)
		return brand{}, false, nil
	}

	var cachedBrand brand
	if err := json.Unmarshal(cachedValue, &cachedBrand); err != nil {
		log.Errorf("ERROR unmarshalling cached value for [%v]: %v.", uuid, err.Error())
		return brand{}, true, err
	}
	log.Info(cachedBrand)
	return cachedBrand, true, nil
}

func (s *brandServiceImpl) openDB() error {
	s.Lock()
	defer s.Unlock()
	log.Infof("Opening database '%v'.", s.cacheFileName)
	if s.db == nil {
		var err error
		if s.db, err = bolt.Open(s.cacheFileName, 0600, &bolt.Options{Timeout: 1 * time.Second}); err != nil {
			log.Errorf("ERROR opening cache file for init: %v.", err.Error())
			return err
		}
	}
	return s.createCacheBucket()
}

func (s *brandServiceImpl) reloadDB() error {
	s.setDataLoaded(false)

	err := s.loadDB()
	if err != nil {
		log.Errorf("Error while creating BrandService: [%v]", err.Error())
		return err
	}
	var bBrands []berthaBrand

	bBrands, err = getBerthaBrands(s.berthaURL)
	if err != nil {
		log.Errorf("Error on Bertha load: [%v]", err.Error())
	} else {
		err = s.loadCuratedBrands(bBrands)
		if err != nil {
			log.Errorf("Error while loading in the curated brands: [%v]", err.Error())
		}
	}

	return nil
}

func (s *brandServiceImpl) loadDB() error {
	var wg sync.WaitGroup
	log.Info("Loading DB...")
	c := make(chan []brand)
	go s.processBrands(c, &wg)
	defer func(w *sync.WaitGroup) {
		close(c)
		w.Wait()
	}(&wg)

	if err := s.openDB(); err != nil {
		s.setInitialised(false)
		return err
	}

	responseCount := 0
	for {
		terms, err := s.repository.GetTmeTermsFromIndex(responseCount)
		if err != nil {
			return err
		}
		if len(terms) < 1 {
			log.Info("Finished fetching brands from TME. Waiting subroutines to terminate.")
			break
		}

		wg.Add(1)
		s.processTerms(terms, c)
		responseCount += s.maxTmeRecords
	}
	return nil
}

func (s *brandServiceImpl) processTerms(terms []interface{}, c chan<- []brand) {
	log.Info("Processing terms...")
	var cacheToBeWritten []brand
	for _, iTerm := range terms {
		t := iTerm.(term)
		cacheToBeWritten = append(cacheToBeWritten, transformBrand(t, s.taxonomyName))
	}
	c <- cacheToBeWritten
}

func (s *brandServiceImpl) processBrands(c <-chan []brand, wg *sync.WaitGroup) {
	for brands := range c {
		log.Infof("Processing batch of %v brands.", len(brands))
		if err := s.db.Batch(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(cacheBucket))
			if bucket == nil {
				return fmt.Errorf("Cache bucket [%v] not found!", cacheBucket)
			}
			for _, anBrand := range brands {
				marshalledBrand, err := json.Marshal(anBrand)
				if err != nil {
					return err
				}
				err = bucket.Put([]byte(anBrand.UUID), marshalledBrand)
				if err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			log.Errorf("ERROR storing to cache: %+v.", err)
		}
		wg.Done()
	}

	log.Info("Finished processing all brands.")
	if s.isInitialised() {
		s.setDataLoaded(true)
	}
}

func (s *brandServiceImpl) createCacheBucket() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		if tx.Bucket([]byte(cacheBucket)) != nil {
			log.Infof("Deleting bucket '%v'.", cacheBucket)
			if err := tx.DeleteBucket([]byte(cacheBucket)); err != nil {
				log.Warnf("Cache bucket [%v] could not be deleted.", cacheBucket)
			}
		}
		log.Infof("Creating bucket '%s'.", cacheBucket)
		_, err := tx.CreateBucket([]byte(cacheBucket))
		return err
	})
}

func getBerthaBrands(berthaURL string) ([]berthaBrand, error) {
	res, err := http.Get(berthaURL)
	if err != nil {
		return []berthaBrand{}, err
	}

	var bBrands []berthaBrand
	err = json.NewDecoder(res.Body).Decode(&bBrands)
	return bBrands, err
}

func (s *brandServiceImpl) loadCuratedBrands(bBrands []berthaBrand) error {
	s.Lock()
	defer s.Unlock()
	err := s.db.Batch(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(cacheBucket))
		if bucket == nil {
			return fmt.Errorf("Cache bucket [%v] not found!", cacheBucket)
		}

		for _, b := range bBrands {
			cachedBrand := bucket.Get([]byte(b.UUID))
			var a brand
			if b.TmeIdentifier == "" && b.UUID != "dbb0bdae-1f0c-11e4-b0cb-b2227cce2b54" {
				// We've put this check in because editorial sometimes forget the TME Identifier.
				// The UUID is for the FT, which is a special case (no TME Identifier but we still want it)
				log.Warnf("No TME Identifier, ignoring curated brand %s [%s]", b.PrefLabel, b.UUID)
				continue
			} else if cachedBrand == nil {
				log.Warnf("Curated brand %s [%s] was not found in cache.  Adding without V1 information.", b.PrefLabel, b.UUID)
				a, _ = berthaToBrand(b)
			} else {
				json.Unmarshal(cachedBrand, &a)
				a, _ = addBerthaInformation(a, b)

				bucket.Delete([]byte(b.UUID))
			}
			newCachedVersion, _ := json.Marshal(a)
			err := bucket.Put([]byte(a.UUID), newCachedVersion)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func addBerthaInformation(a brand, b berthaBrand) (brand, error) {
	plainDescription, err := html2text.FromString(b.DescriptionXML)
	if err != nil {
		return a, err
	}
	a.UUID = b.UUID
	a.PrefLabel = b.PrefLabel
	a.Strapline = b.Strapline
	a.ParentUUID = b.ParentUUID
	a.Description = plainDescription
	a.DescriptionXML = b.DescriptionXML
	a.ImageURL = b.ImageURL
	a.Type = "Brand"

	return a, nil
}

func berthaToBrand(a berthaBrand) (brand, error) {
	plainDescription, err := html2text.FromString(a.DescriptionXML)

	if err != nil {
		return brand{}, err
	}

	altIds := alternativeIdentifiers{
		UUIDs: []string{a.UUID},
		TME:   []string{a.TmeIdentifier},
	}

	p := brand{
		UUID:                   a.UUID,
		ParentUUID:             a.ParentUUID,
		PrefLabel:              a.PrefLabel,
		Type:                   "Brand",
		Strapline:              a.Strapline,
		Description:            plainDescription,
		DescriptionXML:         a.DescriptionXML,
		ImageURL:               a.ImageURL,
		AlternativeIdentifiers: altIds,
	}

	return p, err
}
