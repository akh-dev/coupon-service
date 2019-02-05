package couponservice

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/akh-dev/coupons-service/api"
	"github.com/pkg/errors"
)

func parseBaseRequest(r *http.Request) (*api.Request, error) {
	log.Println("parsing request")

	dec := json.NewDecoder(r.Body)
	baseRequest := &api.Request{}
	err := dec.Decode(baseRequest)
	if err != nil {
		err = errors.Wrap(err, "failed to parse request")
		log.Println(err.Error())
		return nil, err
	}

	return baseRequest, nil
}

func respondBadRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	respObj := &api.Response{Error: []string{msg}}
	writeResponse(w, respObj)
}

func respondForbidden(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	respObj := &api.Response{Error: []string{"Forbidden"}}
	writeResponse(w, respObj)
}

func extractCouponsFromRequest(r *api.Request) (*api.CouponCollection, error) {

	if r == nil {
		err := errors.Errorf("Coupons data must be provided")
		log.Println(err.Error())
		return nil, err
	}

	cpnCollection := &api.CouponCollection{}
	err := json.Unmarshal(r.Data, cpnCollection)
	if err != nil {
		err = errors.Wrap(err, "failed to parse coupons request")
		log.Println(err.Error())
		return nil, err
	}

	return cpnCollection, nil
}

func extractCouponFilterFromRequest(r *api.Request) (*api.CouponFilter, error) {

	if r == nil {
		err := errors.Errorf("Request data must be provided")
		log.Println(err.Error())
		return nil, err
	}

	filter := &api.CouponFilter{}
	err := json.Unmarshal(r.Data, filter)
	if err != nil {
		err = errors.Wrap(err, "failed to parse filter from the request")
		log.Println(err.Error())
		return nil, err
	}

	return filter, nil
}

func writeResponse(w http.ResponseWriter, respObj *api.Response) {
	response, err := json.Marshal(respObj)
	if err != nil {
		log.Println(err.Error())
		return
	}

	w.Write(response)
}

func respondWithCoupons(w http.ResponseWriter, coupons []api.Coupon) {
	respObj := &api.Response{Result: coupons}
	writeResponse(w, respObj)
}
