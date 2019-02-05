package couponservice

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"

	"github.com/pkg/errors"

	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/akh-dev/coupons-service/api"
	"github.com/mongodb/mongo-go-driver/bson"
)

const (
	DB_COUPON_COLLECTION string = "coupons"
)

func (s *CouponService) createCoupons(coupons []api.Coupon) (*mongo.InsertManyResult, error) {

	db := s.mongoClient.Database(s.dbName)
	couponColl := db.Collection(DB_COUPON_COLLECTION)

	documents := []interface{}{}
	for _, cpn := range coupons {
		documents = append(documents, bson.M{
			"name":      cpn.Name,
			"brand":     cpn.Brand,
			"value":     cpn.Value,
			"expiry":    cpn.Expiry,
			"createdAt": time.Now(),
		})

	}
	ctx, _ := context.WithTimeout(context.Background(), s.timeout)
	res, err := couponColl.InsertMany(ctx, documents)
	if err != nil {
		err = errors.Wrap(err, "failed to write new coupons to the db")
		log.Println(err.Error())
		return nil, err
	}

	return res, nil

}

func (s *CouponService) updateCoupons(coupons []api.Coupon) (int64, error) {

	db := s.mongoClient.Database(s.dbName)
	couponColl := db.Collection(DB_COUPON_COLLECTION)

	var UpdatedCnt int64 = 0

	for _, cpn := range coupons {
		fields := bson.D{}

		if cpn.Name != "" {
			fields = append(fields, bson.E{"name", cpn.Name})
		}
		if cpn.Brand != "" {
			fields = append(fields, bson.E{"brand", cpn.Brand})
		}
		if cpn.Value > 0 {
			fields = append(fields, bson.E{"value", cpn.Value})
		}
		if !cpn.Expiry.IsZero() {
			fields = append(fields, bson.E{"expiry", cpn.Expiry})
		}

		ctx, _ := context.WithTimeout(context.Background(), s.timeout)
		res, err := couponColl.UpdateOne(
			ctx,
			bson.D{
				{"item", "paper"},
			},
			bson.D{
				{"$set", fields},
				{"$currentDate", bson.D{
					{"lastModified", true},
				}},
			},
		)

		if err != nil {
			err = errors.Wrap(err, "failed to write new coupons to the db")
			log.Println(err.Error())
			return UpdatedCnt, err
		}

		UpdatedCnt = UpdatedCnt + res.ModifiedCount
	}

	return UpdatedCnt, nil

}

func (s *CouponService) findByIds(ids []interface{}) ([]api.Coupon, error) {
	idsBsonA := bson.A{}
	for _, id := range ids {
		idsBsonA = append(idsBsonA, id)
	}

	filter := bson.D{{"_id", bson.D{{"$in", idsBsonA}}}}
	return s.findManyWithFilter(filter)
}

func (s *CouponService) findManyWithFilter(filter interface{}) ([]api.Coupon, error) {
	if s.debug {
		log.Printf("searching with filter %s", filter)
	}

	db := s.mongoClient.Database(s.dbName)
	couponColl := db.Collection(DB_COUPON_COLLECTION)
	ctx, _ := context.WithTimeout(context.Background(), s.timeout)
	cur, err := couponColl.Find(ctx, filter)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer cur.Close(ctx)

	coupons := []api.Coupon{}
	for cur.Next(ctx) {
		log.Println("found one")
		cpn := api.Coupon{}
		err := cur.Decode(&cpn)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		coupons = append(coupons, cpn)
	}
	if err := cur.Err(); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return coupons, nil
}

func (s *CouponService) searchFromRequest(reqFilter *api.CouponFilter) ([]api.Coupon, error) {
	dbFilter, err := buildFilterFromRequest(reqFilter)
	if err != nil {
		return nil, err
	}
	return s.findManyWithFilter(dbFilter)
}

func buildFilterFromRequest(reqFilter *api.CouponFilter) (bson.D, error) {
	fieldsFilter := bson.D{}

	if reqFilter == nil {
		return nil, errors.Errorf("Search criteria must be provided")
	}

	//filter by ID
	if len(reqFilter.IdIn) > 0 {
		ids := bson.A{}
		for _, id := range reqFilter.IdIn {
			objId, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				return nil, err
			}
			ids = append(ids, objId)
		}

		if len(ids) > 0 {
			fieldsFilter = append(fieldsFilter, bson.E{"_id", bson.D{{"$in", ids}}})
		}
	}

	//Filter by Value
	if reqFilter.ValueFrom != nil || reqFilter.ValueTo != nil {
		valueFilter := bson.D{}
		if reqFilter.ValueFrom != nil {
			valueFilter = append(valueFilter, bson.E{"$gte", *reqFilter.ValueFrom})
		}
		if reqFilter.ValueTo != nil {
			valueFilter = append(valueFilter, bson.E{"$lte", *reqFilter.ValueTo})
		}
		fieldsFilter = append(fieldsFilter, bson.E{"value", valueFilter})
	}

	//Filter by Name (loose search)
	if reqFilter.NameContains != "" {
		fieldsFilter = append(fieldsFilter, bson.E{"name", primitive.Regex{Pattern: fmt.Sprintf(".*%s.*", reqFilter.NameContains), Options: ""}})
	}

	//Filter by Brand (exact match)
	if reqFilter.BrandEqual != "" {
		fieldsFilter = append(fieldsFilter, bson.E{"brand", reqFilter.BrandEqual})
	}

	//Filter by expiry date
	if !reqFilter.ExpiryFrom.IsZero() || !reqFilter.ExpiryTo.IsZero() {
		expiryFilter := bson.D{}
		if !reqFilter.ExpiryFrom.IsZero() {
			expiryFilter = append(expiryFilter, bson.E{"$gte", reqFilter.ExpiryFrom})
		}
		if !reqFilter.ExpiryTo.IsZero() {
			expiryFilter = append(expiryFilter, bson.E{"$lte", reqFilter.ExpiryTo})
		}
		fieldsFilter = append(fieldsFilter, bson.E{"expiry", expiryFilter})
	}

	//Filter by createdAt date
	if !reqFilter.CreatedAtFrom.IsZero() || !reqFilter.CreatedAtTo.IsZero() {
		expiryFilter := bson.D{}
		if !reqFilter.CreatedAtFrom.IsZero() {
			expiryFilter = append(expiryFilter, bson.E{"$gte", reqFilter.CreatedAtFrom})
		}
		if !reqFilter.CreatedAtTo.IsZero() {
			expiryFilter = append(expiryFilter, bson.E{"$lte", reqFilter.CreatedAtTo})
		}
		fieldsFilter = append(fieldsFilter, bson.E{"createdAt", expiryFilter})
	}

	return fieldsFilter, nil
}
