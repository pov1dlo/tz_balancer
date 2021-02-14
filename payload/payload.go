package payload

// Payload ...
type Payload struct {
	Price    int `json:price`
	Quantity int `json:quantity`
	Amount   int `json:amount`
	Object   int `json:object`
	Method   int `json:method`
}
