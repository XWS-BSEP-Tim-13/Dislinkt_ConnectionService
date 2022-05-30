package domain

import (
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/domain/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Company struct {
	Id          primitive.ObjectID `bson:"_id"`
	CompanyName string             `bson:"company_name" validate:"required,companyName"`
	Username    string             `bson:"username" validate:"required,username"`
	Email       string             `bson:"email" validate:"required,email"`
	PhoneNumber string             `bson:"phone_number" validate:"required,numeric,min=9,max=10"`
	Description string             `bson:"description"`
	Location    string             `bson:"location" validate:"required,max=256"`
	Website     string             `bson:"website" validate:"required,website"`
	CompanySize string             `bson:"company_size" validate:"required,companyName"`
	Industry    string             `bson:"industry" validate:"required,max=256"`
	IsActive    bool               `bson:"is_active"`
}

type JobOffer struct {
	Id             primitive.ObjectID  `bson:"_id"`
	Position       string              `bson:"position"`
	JobDescription string              `bson:"job_description"`
	Prerequisites  string              `bson:"prerequisites"`
	Company        Company             `bson:"company"`
	EmploymentType enum.EmploymentType `bson:"employment_type"`
	Published      time.Time           `bson:"published" validate:"required"`
}

type Education struct {
	Id           primitive.ObjectID `bson:"_id"`
	School       string             `bson:"school"`
	Degree       enum.Degree        `bson:"degree"`
	FieldOfStudy string             `bson:"field_of_study"`
	StartDate    time.Time          `bson:"start_date"`
	EndDate      time.Time          `bson:"end_date"`
	Description  string             `bson:"description"`
}

type Experience struct {
	Id                 primitive.ObjectID  `bson:"_id"`
	Title              string              `bson:"title"`
	EmploymentType     enum.EmploymentType `bson:"employment_type"`
	CompanyName        string              `bson:"company_name"`
	Location           string              `bson:"location"`
	IsCurrentlyWorking bool                `bson:"is_currently_working"`
	StartDate          time.Time           `bson:"start_date"`
	EndDate            time.Time           `bson:"end_date"`
	Industry           string              `bson:"industry"`
	Description        string              `bson:"description"`
}

type RegisteredUser struct {
	Id          primitive.ObjectID   `bson:"_id"`
	FirstName   string               `bson:"first_name"`
	LastName    string               `bson:"last_name"`
	Email       string               `bson:"email"`
	PhoneNumber string               `bson:"phone_number"`
	Gender      enum.Gender          `bson:"gender"`
	DateOfBirth time.Time            `bson:"date_of_birth"`
	Biography   string               `bson:"biography"`
	IsPrivate   bool                 `bson:"is_private"`
	Experiences []Experience         `bson:"experiences"`
	Educations  []Education          `bson:"educations"`
	Skills      []string             `bson:"skills"`
	Interests   []primitive.ObjectID `bson:"interests"`
	Connections []primitive.ObjectID `bson:"experiences"`
	Username    string               `bson:"username" validate:"required,username"`
}

type ConnectionRequest struct {
	Id          primitive.ObjectID `bson:"_id"`
	From        RegisteredUser     `bson:"from"`
	To          RegisteredUser     `bson:"to"`
	RequestTime time.Time          `bson:"request_time"`
}
