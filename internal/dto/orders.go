package dto

import "time"

type Order struct {
	OrderID        int64     `json:"order_id" validate:"required"`
	OrderDate      time.Time `json:"order_date" validate:"required"`
	Status         string    `json:"status" validate:"required,oneof='pending' 'shipped' 'delivered' 'cancelled'"`
	TotalAmount    float64   `json:"total_amount" validate:"required,gt=0"`
	Products       []Product `json:"products" validate:"required,dive"`
	Customer       Customer  `json:"customer" validate:"required"`
	PaymentMethod  string    `json:"payment_method" validate:"required,oneof='credit_card' 'paypal' 'bank_transfer'"`
	PaymentStatus  string    `json:"payment_status" validate:"required,oneof='paid' 'unpaid' 'refunded'"`
	DeliveryDate   time.Time `json:"delivery_date,omitempty"`
	TrackingNumber int64     `json:"tracking_number,omitempty"`
	CreatedAt      time.Time `json:"created_at" validate:"required"`
	UpdatedAt      time.Time `json:"updated_at" validate:"required"`
}

type Product struct {
	ProductID int64   `json:"product_id" validate:"required"`
	Name      string  `json:"name" validate:"required"`
	Quantity  int     `json:"quantity" validate:"required,gt=0"`
	Price     float64 `json:"price" validate:"required,gt=0"`
}

type Customer struct {
	CustomerID int64  `json:"customer_id" validate:"required"`
	FirstName  string `json:"first_name" validate:"required"`
	LastName   string `json:"last_name" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Phone      string `json:"phone" validate:"required"`
	Address    string `json:"address" validate:"required"`
}

type OrderMessageResponse struct {
	RequestId string `json:"RequestId" validate:"required"`
	Status    string `json:"Status" validate:"required"`
}
