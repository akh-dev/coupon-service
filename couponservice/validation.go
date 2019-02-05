package couponservice

import (
	"fmt"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"

	"github.com/akh-dev/coupons-service/api"
)

const COUPON_MIN_EXPIRY_DATE string = "2010-01-01T00:00:00Z"

//func type: validator for a single coupon data
type cpnValidatorFunc func(cpnCollection *api.Coupon) (validationSuccess bool, errors []string)

//validates a coupon collection before performing an insert
func (s *CouponService) validateManyForInsert(cpnCollection *api.CouponCollection) (validationSuccess bool, errors []string) {
	return validateMany(cpnCollection, validateOneForInsert)
}

//validates a coupon collection before performing an update
func (s *CouponService) validateManyForUpdate(cpnCollection *api.CouponCollection) (validationSuccess bool, errors []string) {
	return validateMany(cpnCollection, validateOneForUpdate)
}

//generic validation for a coupon collection (actual validator is passed as parameter)
func validateMany(cpnCollection *api.CouponCollection, validator cpnValidatorFunc) (validationSuccess bool, errors []string) {
	validationSuccess = true

	if cpnCollection == nil || len(cpnCollection.Coupons) == 0 {
		validationSuccess = false
		errors = append(errors, "No coupon data provided")
	} else {
		for _, cpn := range cpnCollection.Coupons {
			//ok, e := (&cpn)
			ok, e := validator(&cpn)
			validationSuccess = validationSuccess && ok
			errors = append(errors, e...)
		}
	}

	return validationSuccess, errors
}

//validates one coupon before inserting
func validateOneForInsert(cpn *api.Coupon) (validationSuccess bool, errors []string) {

	if cpn == nil {
		return false, []string{"ValidateNewCoupon: coupon data needs to be provided"}
	}

	validationSuccess = true
	errors = []string{}

	if cpn.Name == "" {
		validationSuccess = false
		errors = append(errors, "ValidateNewCoupon: Coupon name must be provided")
	}

	if cpn.Brand == "" {
		validationSuccess = false
		errors = append(errors, "ValidateNewCoupon: Coupon brand must be provided")
	}

	if cpn.Value <= 0 {
		validationSuccess = false
		errors = append(errors, "ValidateNewCoupon: A positive coupon value must be provided")
	}

	bot, _ := time.Parse(time.RFC3339, COUPON_MIN_EXPIRY_DATE)
	if cpn.Expiry.IsZero() || cpn.Expiry.Before(bot) {
		validationSuccess = false
		errors = append(errors, fmt.Sprintf("ValidateNewCoupon: Coupon expiry date must be after %s", COUPON_MIN_EXPIRY_DATE))
	}

	return validationSuccess, errors
}

//validates one coupon before updating
func validateOneForUpdate(cpn *api.Coupon) (ok bool, errors []string) {

	if cpn == nil {
		return false, []string{"ValidateUpdateCoupon: coupon data needs to be provided"}
	}

	ok, errors = true, []string{}

	if cpn.Id.IsZero() {
		ok = false
		errors = append(errors, "ValidateUpdateCoupon: Coupon id must be provided")
	}

	if cpn.Value < 0 {
		ok = false
		errors = append(errors, "ValidateUpdateCoupon: A positive coupon value must be provided")
	}

	if !cpn.CreatedAt.IsZero() {
		ok = false
		errors = append(errors, "ValidateUpdateCoupon: CreatedAt is a read-only field")
	}

	return ok, errors
}

func validateCouponId(id interface{}) (ok bool, errors []string) {

	ok = true
	errors = []string{}

	switch id.(type) {
	case primitive.ObjectID:
		cast := id.(primitive.ObjectID)
		if cast.IsZero() {
			ok = false
			errors = append(errors, "ValidateUpdateCoupon: Coupon id must be provided")
		}
	case int:
		cast := id.(int)
		if cast < 1 {
			ok = false
			errors = append(errors, "ValidateUpdateCoupon: Coupon id must be provided")
		}
	case string:
		cast := id.(string)
		if cast == "" {
			ok = false
			errors = append(errors, "ValidateUpdateCoupon: Coupon id must be provided")
		}
	default:
		ok = false
		errors = append(errors, "ValidateUpdateCoupon: Coupon id is of unknown type")
	}

	return ok, errors
}
