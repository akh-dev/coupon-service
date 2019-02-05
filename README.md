# coupon-service

Sample create:
curl -X POST -d '{"apiKey":"Valid API Key","data":{"coupons":[{"name":"Save £1 at Tesco","brand":"Tesco","value":1,"expiry":"2019-03-01T00:00:00Z"},{"name":"Save £2 at Boots","brand":"Boots","value":2,"expiry":"2019-04-01T00:00:00Z"}]}}' -H "Content-Type:application/json" localhost:8080

Sample response:
{"result":[{"id":"5c58ea1afaa48016746e59b9","name":"Save £1 at Tesco","brand":"Tesco","value":1,"expiry":"2019-03-01T00:00:00Z","createdAt":"2019-02-05T01:42:50.667Z"},{"id":"5c58ea1afaa48016746e59ba","name":"Save £2 at Boots","brand":"Boots","value":2,"expiry":"2019-04-01T01:00:00+01:00","createdAt":"2019-02-05T01:42:50.667Z"}]}osyan@andrew-ryzen ~/go/src $ 




Sample update:
curl -X POST -d '{"apiKey":"Valid API Key","data":{"coupons":[{"name":"Save £3 at Tesco","brand":"Tesco","value":3,"expiry":"2019-03-01T00:00:00Z"}]}}' -H "Content-Type:application/json" localhost:8080

Sample response:
{"result":[{"id":"5c58f1e10f468a8b68c814ca","name":"Save £3 at Tesco","brand":"Tesco","value":3,"expiry":"2019-03-01T00:00:00Z","createdAt":"2019-02-05T02:16:01.549Z"}]}osyan@andrew-ryzen ~/go/src $ 





Sample list all:
curl -X GET -d '{"apiKey":"Valid API Key","data":{}}' -H "Content-Type:application/json" localhost:8080

Sample list with filters (any combination):
curl -X GET -d '{"apiKey":"Valid API Key","data":{"idIn":["5c58f1e10f468a8b68c814ca"],"nameContains":"£3","brandEqual":"Tesco","valueFrom":1,"valueTo":3,"expiryFrom":"2019-03-01T00:00:00Z","expiryTo":"2019-03-02T00:00:00Z","createdAtFrom":"2018-03-02T00:00:00Z","createdAtTo":"2019-03-02T00:00:00Z"}}' -H "Content-Type:application/json" localhost:8080

Sample respopnse:
{"result":[{"id":"5c58f1e10f468a8b68c814ca","name":"Save £3 at Tesco","brand":"Tesco","value":3,"expiry":"2019-03-01T00:00:00Z","createdAt":"2019-02-05T02:16:01.549Z"}]}
