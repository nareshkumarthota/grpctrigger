package grpctrigger

import "github.com/nareshkumarthota/grpctrigger/pb"

var employees = []pb.Employee{
	pb.Employee{
		Id:                  1,
		BadgeNumber:         2080,
		FirstName:           "testfn1",
		LastName:            "testln1",
		VacationAccrualRate: 2,
		VacationAccrued:     30,
	},
	pb.Employee{
		Id:                  2,
		BadgeNumber:         2082,
		FirstName:           "testfn2",
		LastName:            "testln2",
		VacationAccrualRate: 2.3,
		VacationAccrued:     30.43,
	},
	pb.Employee{
		Id:                  3,
		BadgeNumber:         2083,
		FirstName:           "testfn3",
		LastName:            "testln3",
		VacationAccrualRate: 2.23,
		VacationAccrued:     30.56,
	},
}
