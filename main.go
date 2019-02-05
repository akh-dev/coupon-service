package main

import (
	"fmt"
	"log"

	"github.com/akh-dev/coupons-service/config"
	"github.com/akh-dev/coupons-service/couponservice"
)

func main() {

	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("Failed to load config: %+v", err)
	}

	couponService, err := couponservice.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialise coupon service: %+v", err)
	}

	couponService.ListenAndServe()

	log.Println("Coupon-Service started, press <ENTER> to exit")
	fmt.Scanln()

}
