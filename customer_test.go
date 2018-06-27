package qvo

import (
	"testing"

	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCustomer(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	//Use test token and playground
	token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJjb21tZXJjZV9pZCI6ImNvbV9NeWxDWXg1YklUbkxaUjhTcmdFbzJ3IiwiYXBpX3Rva2VuIjp0cnVlfQ.IrqOpU5fw-TtZMrKg-JkXGL4KCll-ekvqcJL4LHep8w"
	Convey("Given valid token a client should be created", t, func() {
		c := NewClient(token, true)

		Convey("After listing customers", func() {

			var where = make(map[string]map[string]interface{})
			customers, err := ListCustomers(c, 0, 0, where, "")
			So(err, ShouldBeNil)

			//Wipe all customers for testing.
			Convey("We should be able to delete them all to test fresh", func() {

				for _, customer := range customers {
					delErr := DeleteCustomer(c, customer.ID)
					So(delErr, ShouldBeNil)
				}

				Convey("We should be able to create a couple of customers", func() {

					customer1, err := CreateCustomer(c, "Ignacio Gómez", "test@manglar.cl")
					So(err, ShouldBeNil)
					So(customer1.Email, ShouldResemble, "test@manglar.cl")

					customer2, err := CreateCustomer(c, "Jere Díaz", "test2@manglar.cl")
					So(err, ShouldBeNil)
					So(customer2.Email, ShouldResemble, "test2@manglar.cl")

					log.Debugf("\n\n***\n\ncreated at: %s\n\n***\n\n", customer1.CreatedAt.String())

					Convey("So a customer should be retreivable and updatable", func() {
						retrieved, err := GetCustomer(c, customer1.ID)
						So(err, ShouldBeNil)
						So(retrieved.Email, ShouldResemble, "test@manglar.cl")

						uCustomer, err := UpdateCustomer(c, customer1.ID, "Ignacio Gómez R", "test@manglar.cl", retrieved.DefaultPaymentMethod.ID)
						So(err, ShouldBeNil)
						So(retrieved.CreatedAt, ShouldResemble, uCustomer.CreatedAt)

						Convey("So we shouldn't be able to create a new customer with an existing email", func() {
							_, err := CreateCustomer(c, "Ignacio Gómez", "test@manglar.cl")
							So(err, ShouldNotBeNil)
							log.Debugf("error: %s", err)

							Convey("Listing them with inverse orders should render inverse lists", func() {

								customersAsc, err := ListCustomers(c, 0, 0, where, "created_at ASC")
								So(err, ShouldBeNil)
								customersDesc, err := ListCustomers(c, 0, 0, where, "created_at DESC")
								So(err, ShouldBeNil)
								So(customersAsc[0].Email, ShouldResemble, customersDesc[len(customersDesc)-1].Email)
								So(customersAsc[len(customersAsc)-1].Email, ShouldResemble, customersDesc[0].Email)

								Convey("Filtering by email should work", func() {

									where["email"] = make(map[string]interface{})
									where["email"]["like"] = "%test%"

									customersLike, err := ListCustomers(c, 0, 0, where, "")
									So(err, ShouldBeNil)
									So(customersLike, ShouldHaveLength, 2)

									where["name"] = make(map[string]interface{})
									where["name"]["="] = "Ignacio Gómez"

									customersEq, err := ListCustomers(c, 0, 0, where, "")
									So(err, ShouldBeNil)
									So(customersEq, ShouldHaveLength, 1)

									Convey("Finally, listing and deleting each client should work", func() {
										where := make(map[string]map[string]interface{})
										customers, err := ListCustomers(c, 0, 0, where, "")
										So(err, ShouldBeNil)
										for _, customer := range customers {
											delErr := DeleteCustomer(c, customer.ID)
											So(delErr, ShouldBeNil)
										}
									})

								})

							})

						})

					})

				})

			})

		})

	})
}
