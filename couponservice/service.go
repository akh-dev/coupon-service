package couponservice

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/akh-dev/coupons-service/dblayer"

	"github.com/akh-dev/coupons-service/api"
	"github.com/akh-dev/coupons-service/config"
	"github.com/akh-dev/coupons-service/util"

	"github.com/mongodb/mongo-go-driver/mongo"
)

type CouponService struct {
	db      dblayer.Interface
	timeout time.Duration
	port    string
	debug   bool
}

func New(cfg *config.Config) (*CouponService, error) {

	client, err := mongo.NewClient(fmt.Sprintf("mongodb://%s:%s", cfg.DB.Host, cfg.DB.Port))
	if err != nil {
		log.Printf("Failed to create a mongo client: %s", err.Error())
		return nil, err
	}

	timeout := time.Duration(cfg.Service.CtxTimeout) * time.Second

	ctx, _ := context.WithTimeout(context.Background(), timeout)
	err = client.Connect(ctx)

	if err != nil {
		log.Printf("Failed to connect to mongo server: %s", err.Error())
		return nil, err
	}

	db, err := dblayer.New(client, cfg.DB.Name, timeout)
	if err != nil {
		log.Printf("Failed to create db layer object: %s", err.Error())
		return nil, err
	}

	service := &CouponService{
		db:      db,
		timeout: timeout,
		port:    cfg.Service.Port,
		debug:   cfg.Service.Debug,
	}

	return service, nil
}

func (s *CouponService) ListenAndServe() {
	http.HandleFunc("/", s.handleCouponsRequest)

	go func() {
		//err := http.ListenAndServeTLS(fmt.Sprintf(":%s", s.port), "cert.pem", "key.pem", new(util.GzipHandler))
		err := http.ListenAndServe(fmt.Sprintf(":%s", s.port), new(util.GzipHandler))
		if err != nil {
			log.Fatal(err.Error())
		}
	}()
}

func (s *CouponService) handleCouponsRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	baseRequest, err := parseBaseRequest(r)
	if err != nil {
		log.Printf("errors during handleCouponsRequest:%s", err.Error())
		respondBadRequest(w, err.Error())
	}

	if !s.authenticate(baseRequest) {
		respondForbidden(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.handleListCoupons(w, baseRequest)
	case http.MethodPost:
		s.handleCreateCoupon(w, baseRequest)
	case http.MethodPut:
		s.handleUpdateCoupon(w, baseRequest)
	default:
		respondBadRequest(w, "unknown request")
	}
}

func (s *CouponService) handleListCoupons(w http.ResponseWriter, r *api.Request) {
	filter, err := extractCouponFilterFromRequest(r)
	if err != nil {
		respondBadRequest(w, err.Error())
	}

	if s.debug {
		log.Printf("request data: %s", string(r.Data))
	}

	coupons, err := s.db.SearchFromRequest(filter)
	if err != nil {
		respObj := &api.Response{Error: []string{err.Error()}}
		writeResponse(w, respObj)
		return
	}

	respondWithCoupons(w, coupons)
	return
}

func (s *CouponService) handleCreateCoupon(w http.ResponseWriter, r *api.Request) {

	cpnCollection, err := extractCouponsFromRequest(r)
	if err != nil {
		respondBadRequest(w, err.Error())
	}

	if s.debug {
		log.Printf("coupon data: %s", string(r.Data))
	}

	if validationSuccess, errors := s.validateManyForInsert(cpnCollection); !validationSuccess {
		respObj := &api.Response{Error: errors}
		writeResponse(w, respObj)
		return
	}

	if s.debug {
		log.Println("Creating new coupons")
	}

	res, err := s.db.CreateCoupons(cpnCollection.Coupons)
	if err != nil {
		respObj := &api.Response{Error: []string{err.Error()}}
		writeResponse(w, respObj)
		return
	}
	log.Printf("%d coupons created", len(res.InsertedIDs))

	coupons, err := s.db.FindByIds(res.InsertedIDs)
	if err != nil {
		respObj := &api.Response{Error: []string{err.Error()}}
		writeResponse(w, respObj)
		return
	}

	respondWithCoupons(w, coupons)
	return
}

func (s *CouponService) handleUpdateCoupon(w http.ResponseWriter, r *api.Request) {
	cpnCollection, err := extractCouponsFromRequest(r)
	if err != nil {
		respondBadRequest(w, err.Error())
	}

	if s.debug {
		log.Printf("coupon data: %s", string(r.Data))
	}

	if validationSuccess, errors := s.validateManyForUpdate(cpnCollection); !validationSuccess {
		respObj := &api.Response{Error: errors}
		writeResponse(w, respObj)
		return
	}

	if s.debug {
		log.Println("Updating coupons")
	}

	updCount, err := s.db.UpdateCoupons(cpnCollection.Coupons)
	if err != nil {
		respObj := &api.Response{Error: []string{err.Error()}}
		writeResponse(w, respObj)
		return
	}

	log.Printf("%d coupons updated", updCount)

	cpnIDs := []interface{}{}
	for _, cpn := range cpnCollection.Coupons {
		cpnIDs = append(cpnIDs, cpn.Id)
	}

	coupons, err := s.db.FindByIds(cpnIDs)
	if err != nil {
		respObj := &api.Response{Error: []string{err.Error()}}
		writeResponse(w, respObj)
		return
	}

	respondWithCoupons(w, coupons)
	return
}

func (s *CouponService) authenticate(r *api.Request) bool {

	//TODO: replace stud with user lookup
	if r == nil {
		log.Println("No API key provided, authentication failed")
		return false
	}
	if r.ApiKey != "Valid API Key" {
		log.Printf("invalid key provided: %s - Authentication failed", r.ApiKey)
		return false
	}

	log.Println("Authentication successful")
	return true
}
