package qvo

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPlan(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	//Use test token and playground
	token := os.Getenv("QVO_TEST_TOKEN")
	Convey("Given valid token a client should be created", t, func() {
		c := NewClient(token, true)
		//Set log level at debug.
		c.SetLogLevel(log.DebugLevel)

		Convey("After listing plans", func() {

			var where = make(map[string]map[string]interface{})
			plans, err := ListPlans(c, 0, 0, where, "")
			So(err, ShouldBeNil)

			//Wipe all plans for testing.
			Convey("We should be able to delete them all to test fresh", func() {

				for _, plan := range plans {
					delErr := DeletePlan(c, plan.ID)
					So(delErr, ShouldBeNil)
				}

				Convey("We should be able to create a couple of plans", func() {

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

					testPlan2 := Plan{
						ID:                "test-plan-2-id",
						Name:              "Test Plan 2",
						Price:             "29990.0",
						Currency:          "UF",
						Interval:          "month",
						IntervalCount:     1,
						TrialPeriodDays:   0,
						DefaultCycleCount: 5,
					}
					plan2, err := CreatePlan(c, testPlan2)
					So(err, ShouldBeNil)
					So(plan2.Name, ShouldResemble, testPlan2.Name)

					Convey("So a plan should be retreivable and updatable", func() {

						retrieved, err := GetPlan(c, plan1.ID)
						So(err, ShouldBeNil)
						So(retrieved.Name, ShouldResemble, testPlan1.Name)

						uPlan, err := UpdatePlan(c, retrieved.ID, "Modified plan 1")
						So(err, ShouldBeNil)
						So(retrieved.CreatedAt, ShouldResemble, uPlan.CreatedAt)
						So(uPlan.Name, ShouldResemble, "Modified plan 1")

						Convey("So we shouldn't be able to create a new plan with an existing id", func() {

							testPlan3 := Plan{
								ID:                "test-plan-1-id",
								Name:              "Test Plan 2",
								Price:             "29990.0",
								Currency:          "UF",
								Interval:          "month",
								IntervalCount:     1,
								TrialPeriodDays:   0,
								DefaultCycleCount: 5,
							}

							_, err := CreatePlan(c, testPlan3)
							So(err, ShouldNotBeNil)
							log.Debugf("error: %s", err)

							Convey("Listing them with inverse orders should render inverse lists", func() {

								plansAsc, err := ListPlans(c, 0, 0, where, "created_at ASC")
								So(err, ShouldBeNil)
								plansDesc, err := ListPlans(c, 0, 0, where, "created_at DESC")
								So(err, ShouldBeNil)
								So(plansAsc[0].Name, ShouldResemble, plansDesc[len(plansDesc)-1].Name)
								So(plansAsc[len(plansAsc)-1].Name, ShouldResemble, plansDesc[0].Name)

								Convey("Filtering by name should work", func() {

									where["name"] = make(map[string]interface{})
									where["name"]["like"] = "%Test%"

									plansLike, err := ListPlans(c, 0, 0, where, "")
									So(err, ShouldBeNil)
									So(plansLike, ShouldHaveLength, 1)

									where["name"] = make(map[string]interface{})
									where["name"]["="] = "Test Plan 2"

									plansEq, err := ListPlans(c, 0, 0, where, "")
									So(err, ShouldBeNil)
									So(plansEq, ShouldHaveLength, 1)

									Convey("Finally, listing and deleting each plan should work", func() {
										where := make(map[string]map[string]interface{})
										plans, err := ListPlans(c, 0, 0, where, "")
										So(err, ShouldBeNil)
										for _, plan := range plans {
											delErr := DeletePlan(c, plan.ID)
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
