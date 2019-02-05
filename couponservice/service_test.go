package couponservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/akh-dev/coupons-service/api"

	"github.com/akh-dev/coupons-service/config"
)

const (
	dbHostExpected string = "127.0.0.1"
	dbPortExpected string = "11111"
	dbNameExpected string = "mockdb"

	svcContextTimeoutExpected int    = 5
	svcPortExpected           string = "80"
	svcDebugExpected          bool   = false
)

func TestNew(t *testing.T) {
	//mock config

	svc, err := getNewSvc()
	if err != nil {
		t.Error(err)
	}

	if svc == nil {
		t.Errorf("Expected a valid service object, but got nil")
	} else {

		if svc.db == nil {
			t.Error("Expected to have a db layer, but got nil")
		}

		if svc.port != svcPortExpected {
			t.Errorf("Expected service port to be set to %s, but got %s", svcPortExpected, svc.port)
		}

		toDuration := time.Duration(svcContextTimeoutExpected) * time.Second
		if svc.timeout != toDuration {
			t.Errorf("Expected service ctx timeout to be set to %s, but got %s", toDuration, svc.timeout)
		}

		if svc.debug != svcDebugExpected {
			t.Errorf("Expected service debug flag to be set to %t, but got %t", svcDebugExpected, svc.debug)
		}
	}
}

func TestHandleCouponsRequest(t *testing.T) {
	s, err := getNewSvc()
	if err != nil {
		t.Log(err)
		return
	}

	r, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%s", s.port), bytes.NewBuffer([]byte("{}")))
	if err != nil {
		t.Error(err)
		return
	}
	r.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.handleCouponsRequest(w, r)

	data, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error(err)
		return
	}
	//t.Logf("got response: %s", string(data))

	resp := &api.Response{}
	if err := json.Unmarshal(data, resp); err != nil {
		t.Error(err)
		return
	}

}

func TestHandleListCoupons(t *testing.T) {
	s, err := getNewSvc()
	if err != nil {
		t.Log(err)
		return
	}

	s.db = newDbMock()

	payload := `{}`
	r := &api.Request{
		ApiKey: "dont care",
		Data:   []byte(payload),
	}

	w := httptest.NewRecorder()
	s.handleListCoupons(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("unexpected http status %d", w.Code)
		return
	}

	data, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error(err)
		return
	}

	resp := &api.Response{}
	if err := json.Unmarshal(data, resp); err != nil {
		t.Error(err)
		return
	}

	if len(resp.Error) > 0 {
		t.Errorf("Service returned unexpected errors: %s", strings.Join(resp.Error, ":"))
	}

}

func TestHandleCreateCoupon(t *testing.T) {
	s, err := getNewSvc()
	if err != nil {
		t.Log(err)
		return
	}

	s.db = newDbMock()

	payload := `{"coupons":[{"name":"Save £1 at Tesco","brand":"Tesco","value":1,"expiry":"2019-03-01T00:00:00Z"},{"name":"Save £2 at Boots","brand":"Boots","value":2,"expiry":"2019-04-01T00:00:00Z"}]}`
	r := &api.Request{
		ApiKey: "dont care",
		Data:   []byte(payload),
	}

	w := httptest.NewRecorder()
	s.handleCreateCoupon(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("unexpected http status %d", w.Code)
		return
	}

	data, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error(err)
		return
	}

	resp := &api.Response{}
	if err := json.Unmarshal(data, resp); err != nil {
		t.Error(err)
		return
	}

	if len(resp.Error) > 0 {
		t.Errorf("Service returned unexpected errors: %s", strings.Join(resp.Error, ":"))
	}

}

func TestHandleUpdateCoupon(t *testing.T) {
	s, err := getNewSvc()
	if err != nil {
		t.Log(err)
		return
	}

	s.db = newDbMock()

	payload := `{"coupons":[{"id":"5c58ea1afaa48016746e59b9","name":"Save. Tesco. $1"},{"id":"5c58ea1afaa48016746e59ba","expiry":"2020-12-31T23:59:59Z"}]}`
	r := &api.Request{
		ApiKey: "dont care",
		Data:   []byte(payload),
	}

	w := httptest.NewRecorder()
	s.handleListCoupons(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("unexpected http status %d", w.Code)
		return
	}

	data, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error(err)
		return
	}

	resp := &api.Response{}
	if err := json.Unmarshal(data, resp); err != nil {
		t.Error(err)
		return
	}

	if len(resp.Error) > 0 {
		t.Errorf("Service returned unexpected errors: %s", strings.Join(resp.Error, ":"))
	}
}

func TestAuthenticate(t *testing.T) {
	s, err := getNewSvc()
	if err != nil {
		t.Log(err)
		return
	}

	s.db = newDbMock()

	payload := `{}`
	req1 := &api.Request{
		ApiKey: "some invalid key",
		Data:   []byte(payload),
	}
	if authenticated := s.authenticate(req1); authenticated {
		t.Error("authentication failed: the call was able to authenticate with an invalid key")
	}

	req2 := &api.Request{
		ApiKey: "Valid API Key",
		Data:   []byte(payload),
	}
	if authenticated := s.authenticate(req2); !authenticated {
		t.Error("authentication failed: the call was not able to authenticate with a valid key")
	}
}

func newMockConfig() *config.Config {
	cfg := &config.Config{}

	cfg.DB = config.DBConf{
		Host: dbHostExpected,
		Port: dbPortExpected,
		Name: dbNameExpected,
	}

	cfg.Service = config.ServiceConf{
		CtxTimeout: svcContextTimeoutExpected,
		Port:       svcPortExpected,
		Debug:      svcDebugExpected,
	}

	return cfg
}

func getNewSvc() (*CouponService, error) {
	cfg := newMockConfig()

	return New(cfg)
}

type DbMock struct {
	mongoClient *mongo.Client
	dbName      string
	timeout     time.Duration
}

func (mock *DbMock) CreateCoupons(coupons []api.Coupon) (*mongo.InsertManyResult, error) {
	insRes := &mongo.InsertManyResult{
		InsertedIDs: []interface{}{},
	}

	return insRes, nil
}

func (mock *DbMock) UpdateCoupons(coupons []api.Coupon) (int64, error) {
	return 0, nil
}

func (mock *DbMock) FindByIds(ids []interface{}) ([]api.Coupon, error) {
	return []api.Coupon{}, nil
}

func (mock *DbMock) SearchFromRequest(reqFilter *api.CouponFilter) ([]api.Coupon, error) {
	return []api.Coupon{}, nil
}

func newDbMock() *DbMock {
	return &DbMock{
		mongoClient: nil,
		timeout:     time.Duration(1) * time.Second,
		dbName:      "test",
	}
}
