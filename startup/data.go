package startup

import (
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var users = []*domain.RegisteredUser{
	{
		Id:           getObjectId("723b0cc3a34d25d8567f9f82"),
		FirstName:    "Srdjan",
		LastName:     "Sukovic",
		Email:        "srdjansukovic@gmail.com",
		PhoneNumber:  "0649459562",
		Gender:       0,
		DateOfBirth:  time.Time{},
		Biography:    "Entrepreneur, investor, and business magnate. Next I’m buying Coca-Cola to put the cocaine back in.",
		IsPrivate:    false,
		IsActive:     true,
		Username:     "srdjansukovic",
		BlockedUsers: []string{},
		Experiences: []domain.Experience{
			{
				Id:                 getObjectId("723b0cc3a34d25d8567f9d72"),
				Description:        "Senior web engineer in charge of automotive project",
				StartDate:          time.Time{},
				EndDate:            time.Time{},
				Industry:           "Software",
				IsCurrentlyWorking: false,
				Location:           "Los Angeles",
				CompanyName:        "Google",
				EmploymentType:     0,
				Title:              "Full stack engineer",
			},
			{
				Id:                 getObjectId("723b0cc3a34d25d8567f9d77"),
				Description:        "Junior web engineer in charge of automotive project",
				StartDate:          time.Time{},
				EndDate:            time.Time{},
				Industry:           "Software",
				IsCurrentlyWorking: false,
				Location:           "Las Vegas",
				CompanyName:        "Facebook",
				EmploymentType:     0,
				Title:              "Full stack engineer",
			},
		},
		Educations: []domain.Education{
			{
				Id:           getObjectId("723b0cc3a34d25d8567f9d74"),
				StartDate:    time.Time{},
				EndDate:      time.Time{},
				Description:  "Graduated first in class",
				FieldOfStudy: "Computer science",
				School:       "Harvard",
				Degree:       1,
			},
		},
		Skills:      []string{"AWS", "Docker"},
		Interests:   []primitive.ObjectID{getObjectId("623b0cc3a34d25d8567f9f82")},
		Connections: []string{"stefanljubovic", "anagavrilovic"},
	},
	{
		Id:          getObjectId("723b0cc3a34d25d8567f9f83"),
		FirstName:   "Stefan",
		LastName:    "Ljubovic",
		Email:       "ljubovicstefan@gmail.com",
		PhoneNumber: "0654324995",
		Username:    "stefanljubovic",
		Gender:      0,
		DateOfBirth: time.Time{},
		Biography:   "biography sample",
		IsPrivate:   true,
		IsActive:    true,
		Experiences: []domain.Experience{
			{
				Id:                 getObjectId("723b0cc3a34d25d8567f9d73"),
				Description:        "Junior web engineer in charge of automotive project",
				StartDate:          time.Time{},
				EndDate:            time.Time{},
				Industry:           "Software",
				IsCurrentlyWorking: false,
				Location:           "Novi Sad",
				CompanyName:        "Synechron",
				EmploymentType:     0,
				Title:              "Full stack engineer",
			},
			{
				Id:                 getObjectId("723b0cc3a34d25d8567f9d79"),
				Description:        "Junior web engineer in charge of automotive project",
				StartDate:          time.Time{},
				EndDate:            time.Time{},
				Industry:           "Software",
				IsCurrentlyWorking: false,
				Location:           "Novi Sad",
				CompanyName:        "Symphony",
				EmploymentType:     0,
				Title:              "DevOps engineer",
			},
		},
		Educations:    []domain.Education{},
		Skills:        []string{"AWS", "Java"},
		Interests:     []primitive.ObjectID{},
		Connections:   []string{"anagavrilovic", "srdjansukovic"},
		BlockedUsers:  []string{"marijakljestan"},
		Notifications: false,
	},
	{
		Id:          getObjectId("723b0cc3a34d25d8567f9f84"),
		FirstName:   "Ana",
		LastName:    "Gavrilovic",
		Email:       "anagavrilovic@gmail.com",
		PhoneNumber: "0642152",
		Username:    "anagavrilovic",
		Gender:      1,
		DateOfBirth: time.Time{},
		Biography:   "biography sample",
		IsPrivate:   false,
		Experiences: []domain.Experience{
			{
				Id:                 getObjectId("723b0cc3a34d25d8567f9d74"),
				Description:        "Junior web engineer in charge of automotive project",
				StartDate:          time.Time{},
				EndDate:            time.Time{},
				Industry:           "Software",
				IsCurrentlyWorking: false,
				Location:           "Novi Sad",
				CompanyName:        "Synechron",
				EmploymentType:     0,
				Title:              "Full stack engineer",
			},
			{
				Id:                 getObjectId("723b0cc3a34d25d8567f9d78"),
				Description:        "Junior web engineer in charge of automotive project",
				StartDate:          time.Time{},
				EndDate:            time.Time{},
				Industry:           "Software",
				IsCurrentlyWorking: false,
				Location:           "Novi Sad",
				CompanyName:        "Smart Cat",
				EmploymentType:     0,
				Title:              "DevOps engineer",
			},
		},
		Educations:    []domain.Education{},
		Skills:        []string{"Java", "Docker"},
		Interests:     []primitive.ObjectID{},
		Connections:   []string{"stefanljubovic", "srdjansukovic"},
		BlockedUsers:  []string{},
		Notifications: true,
	},
	{
		Id:          getObjectId("723b0cc3a34d25d8567f9f85"),
		FirstName:   "Marija",
		LastName:    "Kljestan",
		Email:       "marijakljestan@gmail.com",
		PhoneNumber: "0642152643",
		Username:    "marijakljestan",
		Gender:      1,
		DateOfBirth: time.Time{},
		Biography:   "biography sample",
		IsPrivate:   false,
		IsActive:    true,
		Experiences: []domain.Experience{
			{
				Id:                 getObjectId("723b0cc3a34d25d8567f9d75"),
				Description:        "Junior web engineer in charge of automotive project",
				StartDate:          time.Time{},
				EndDate:            time.Time{},
				Industry:           "Software",
				IsCurrentlyWorking: false,
				Location:           "Novi Sad",
				CompanyName:        "Levi9",
				EmploymentType:     0,
				Title:              "Java engineer",
			},
			{
				Id:                 getObjectId("723b0cc3a34d25d8567f9d76"),
				Description:        "Junior web engineer in charge of automotive project",
				StartDate:          time.Time{},
				EndDate:            time.Time{},
				Industry:           "Software",
				IsCurrentlyWorking: false,
				Location:           "Novi Sad",
				CompanyName:        "Symphony",
				EmploymentType:     0,
				Title:              "DevOps engineer",
			},
		},
		Educations:    []domain.Education{},
		Skills:        []string{"AWS", "Docker", "Java"},
		Interests:     []primitive.ObjectID{},
		Connections:   []string{"srdjansukovic"},
		BlockedUsers:  []string{"stefanljubovic"},
		Notifications: true,
	},
	{
		Id:            getObjectId("723b0cc3a34d25d8567f9f86"),
		FirstName:     "Lenka",
		LastName:      "Aleksic",
		Email:         "lenka@gmail.com",
		PhoneNumber:   "064364364",
		Username:      "lenka",
		Gender:        1,
		DateOfBirth:   time.Time{},
		Biography:     "biography sample",
		IsPrivate:     false,
		IsActive:      true,
		Experiences:   []domain.Experience{},
		Educations:    []domain.Education{},
		Skills:        []string{"Java", "C#"},
		Interests:     []primitive.ObjectID{},
		Connections:   []string{},
		BlockedUsers:  []string{},
		Notifications: false,
	},
}

var connectionRequests = []*domain.ConnectionRequest{
	{
		Id: getObjectId("62b89e802697fd8b2ce82138"),
		From: domain.RegisteredUser{
			Id:          getObjectId("723b0cc3a34d25d8567f9f86"),
			FirstName:   "Lenka",
			LastName:    "Aleksic",
			Email:       "lenka@gmail.com",
			PhoneNumber: "064364364",
			Username:    "lenka",
			Gender:      1,
			DateOfBirth: time.Time{},
			Biography:   "biography sample",
			IsPrivate:   false,
			IsActive:    true,
			Experiences: []domain.Experience{},
			Educations:  []domain.Education{},
			Skills:      []string{"s1", "s2"},
			Interests:   []primitive.ObjectID{},
			Connections: []string{},
		},
		To: domain.RegisteredUser{
			Id:          getObjectId("723b0cc3a34d25d8567f9f83"),
			FirstName:   "Stefan",
			LastName:    "Ljubovic",
			Email:       "ljubovicstefan@gmail.com",
			PhoneNumber: "0654324995",
			Username:    "stefanljubovic",
			Gender:      0,
			DateOfBirth: time.Time{},
			Biography:   "biography sample",
			IsPrivate:   true,
			IsActive:    true,
			Experiences: []domain.Experience{},
			Educations:  []domain.Education{},
			Skills:      []string{"s1", "s2"},
			Interests:   []primitive.ObjectID{},
			Connections: []string{"marijakljestan", "anagavrilovic", "srdjansukovic"},
		},
		RequestTime: time.Time{},
	},
}

var companies = []*domain.Company{
	{
		Id:          getObjectId("623b0cc3a34d25d8567f9f82"),
		CompanyName: "Levi9",
		Username:    "levi9",
		Email:       "levi9@levi9.com",
		PhoneNumber: "0651234567",
		Location:    "ns",
		Description: "Technology services",
		Website:     "www.levi9.com",
		CompanySize: "1000",
		Industry:    "IT",
		IsActive:    true,
	},
	{
		Id:          getObjectId("623b0cc3a34d25d8567f9f83"),
		CompanyName: "Symphony",
		Username:    "Symphony",
		Email:       "symphony@symphony.com",
		PhoneNumber: "06517654321",
		Location:    "ns",
		Description: "Technology services",
		Website:     "www.symphony.com",
		CompanySize: "1000",
		Industry:    "IT",
		IsActive:    true,
	},
}

var jobs = []*domain.JobOffer{
	{
		Id:             getObjectId("623b0cc3a34d25d8567f9f92"),
		EmploymentType: 0,
		Position:       "DevOps Engineer",
		Prerequisites:  "AWS, Docker, Linux",
		Company:        *companies[1],
		JobDescription: "Great chance for self development and work with experts",
		Published:      time.Now().Add(24 * time.Hour),
	},
	{
		Id:             getObjectId("623b0cc3a34d25d8567f9f93"),
		EmploymentType: 0,
		Position:       "Full stack Engineer",
		Prerequisites:  "Java",
		Company:        *companies[0],
		JobDescription: "Great opportunity for self development and work with experts",
		Published:      time.Now().Add(24 * time.Hour),
	},
	{
		Id:             getObjectId("623b0cc3a34d25d8567f9f93"),
		EmploymentType: 0,
		Position:       "Experienced golang developer",
		Prerequisites:  "Golang",
		Company:        *companies[0],
		JobDescription: "Great opportunity for self development and learning from experts",
		Published:      time.Now().Add(24 * time.Hour),
	},
}

var events = []*domain.Event{
	{
		Id:        getObjectId("623b0cc3a34d25d8567f9f92"),
		Action:    "Blocked user marijakljestan",
		User:      "stefanljubovic",
		Published: time.Now().Add(24 * time.Hour),
	},
}

func getObjectId(id string) primitive.ObjectID {
	if objectId, err := primitive.ObjectIDFromHex(id); err == nil {
		return objectId
	}
	return primitive.NewObjectID()
}
