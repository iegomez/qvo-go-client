package qvo

import (
	"testing"
)

func TestSubscription(t *testing.T) {

	//Can't test subscriptions if there are no cards.

	/*
		log.SetLevel(log.DebugLevel)
		//Use test token and playground
		token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJjb21tZXJjZV9pZCI6ImNvbV9NeWxDWXg1YklUbkxaUjhTcmdFbzJ3IiwiYXBpX3Rva2VuIjp0cnVlfQ.IrqOpU5fw-TtZMrKg-JkXGL4KCll-ekvqcJL4LHep8w"
		Convey("Given valid token a client should be created", t, func() {
			c := NewClient(token, true)

			Convey("After listing subscriptions", func() {

				var where = make(map[string]map[string]interface{})
				subscriptions, err := ListSubscriptions(c, 0, 0, where, "")
				So(err, ShouldBeNil)

				//Wipe all subscriptions for testing.
				Convey("We should be able to cancel them all to test fresh", func() {

					for _, subscription := range subscriptions {
						delErr := CancelSubscription(c, subscription.ID, false)
						So(delErr, ShouldBeNil)
					}

					Convey("We should be able to create a plan, a customer, a subscription with default values and one with more info", func() {

						customer1, err := CreateCustomer(c, "Ignacio Gómez", "test@manglar.cl")
						So(err, ShouldBeNil)
						So(customer1.Email, ShouldResemble, "test@manglar.cl")

						testPlan1 := Plan{
							ID:                "test-plan-1-id",
							Name:              "Test Plan 1",
							Price:             "19990",
							Currency:          "CLP",
							Interval:          "month",
							IntervalCount:     1,
							TrialPeriodDays:   0,
							DefaultCycleCount: 3,
						}
						plan1, err := CreatePlan(c, testPlan1)
						So(err, ShouldBeNil)
						So(plan1.Name, ShouldResemble, testPlan1.Name)

						subscription1, err := CreateSubscription(c, customer1.ID, plan1.ID, "", 0, 0, nil)
						So(err, ShouldBeNil)
						//Cycle count should resemble the default one.
						So(subscription1.CycleCount, ShouldResemble, plan1.DefaultCycleCount)

						start := time.Now().AddDate(0, 0, 2)
						subscription2, err := CreateSubscription(c, customer1.ID, plan1.ID, "IVA", 19.0, 2, &start)
						So(err, ShouldBeNil)
						So(subscription2.TaxName, ShouldResemble, "IVA")
						So(subscription2.TaxPercent, ShouldResemble, "19.0")
						So(subscription2.CycleCount, ShouldResemble, 2)
						So(subscription2.Start, ShouldHappenAfter, time.Now())

						Convey("So a subscription should be retrievable and updatable", func() {

							//Create another plan with a greater price to reassign the subscription.
							testPlan2 := Plan{
								ID:                "test-plan-2-id",
								Name:              "Test Plan 2",
								Price:             "29990",
								Currency:          "CLP",
								Interval:          "month",
								IntervalCount:     1,
								TrialPeriodDays:   10,
								DefaultCycleCount: 1,
							}
							plan2, err := CreatePlan(c, testPlan2)
							So(err, ShouldBeNil)
							So(plan2.Name, ShouldResemble, testPlan2.Name)

							retrieved, err := GetSubscription(c, subscription1.ID)
							So(err, ShouldBeNil)
							So(retrieved.CycleCount, ShouldResemble, subscription1.CycleCount)

							uSubscription, err := UpdateSubscription(c, subscription1.ID, plan2.ID)
							So(err, ShouldBeNil)
							So(uSubscription.Debt, ShouldBeGreaterThan, 0)
							So(uSubscription.Debt, ShouldResemble, 10000)

							Convey("Listing them with inverse orders should render inverse lists", func() {

								subsAsc, err := ListSubscriptions(c, 0, 0, where, "created_at ASC")
								So(err, ShouldBeNil)
								subsDesc, err := ListSubscriptions(c, 0, 0, where, "created_at DESC")
								So(err, ShouldBeNil)
								So(subsAsc[0].Debt, ShouldResemble, subsDesc[len(subsDesc)-1].Debt)
								So(subsAsc[len(subsAsc)-1].Debt, ShouldResemble, subsDesc[0].Debt)

								Convey("Filtering by debt should work", func() {

									where["debt"] = make(map[string]interface{})
									where["debt"][">"] = 0

									subsWithDebt, err := ListSubscriptions(c, 0, 0, where, "")
									So(err, ShouldBeNil)
									So(subsWithDebt, ShouldHaveLength, 1)

									where["debt"] = make(map[string]interface{})
									where["debt"]["="] = 0

									subsWithoutDebt, err := ListSubscriptions(c, 0, 0, where, "")
									So(err, ShouldBeNil)
									So(subsWithoutDebt, ShouldHaveLength, 1)

									Convey("Finally, listing and canceling subscriptions should work", func() {
										where := make(map[string]map[string]interface{})
										subs, err := ListSubscriptions(c, 0, 0, where, "")
										So(err, ShouldBeNil)
										for _, subscription := range subs {
											delErr := CancelSubscription(c, subscription.ID, false)
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
	*/
}
