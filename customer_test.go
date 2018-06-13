package qvo

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCustomer(t *testing.T) {
	//Use test token and playground
	token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJjb21tZXJjZV9pZCI6ImNvbV9RbnB5bkhiOHJkVWRPSl82MWFoR0ZBIiwiYXBpX3Rva2VuIjp0cnVlfQ.b0DqEiqx9rwBC70kdvSfqQ3F_gVG6jFOcs5QlNWnJHk"
	Convey("Given valid token a client should be created", t, func() {
		c := NewClient(token, true)

		//Create a customer
		/*
					//Customer struct to represent a qvo customer object.
			type Customer struct {
				ID                   string         `json:"id"`
				DefaultPaymentMethod Card           `json:"default_payment_method"`
				Name                 string         `json:"name"`
				Email                string         `json:"email"`
				Subscriptions        []Subscription `json:"subscriptions"`
				Cards                []Card         `json:"cards"`
				Transactions         []Transaction  `json:"transactions"`
				CreatedAt            time.Time      `json:"created_at"`
				UpdatedAt            time.Time      `json:"updated_at"`
			}
		*/
		Convey("So we should be able to create a couple of customers", func() {

			customer1, err := CreateCustomer(c, "Ignacio Gómez", "iegomez@manglar.cl")
			So(err, ShouldBeNil)
			So(customer1.Email, ShouldResemble, "iegomez@manglar.cl")

			customer2, err := CreateCustomer(c, "Jere Díaz", "jeremias@manglar.cl")
			So(err, ShouldBeNil)
			So(customer2.Email, ShouldResemble, "jeremias@manglar.cl")

			Convey("So a customer should be updatable and retreivable", func() {
				retrieved, err := GetCustomer(c, customer1.ID)
				So(err, ShouldBeNil)
				So(retrieved.Email, ShouldResemble, "iegomez@manglar.cl")

				uCustomer, err := UpdateCustomer(c, customer1.ID, "Ignacio Gómez R", "iegomez@manglar.cl", retrieved.DefaultPaymentMethod.ID)
				So(err, ShouldBeNil)
				So(retrieved.CreatedAt, ShouldResemble, uCustomer.CreatedAt)
			})

		})

	})
}
