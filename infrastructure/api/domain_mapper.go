package api

import (
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/domain"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/domain/enum"
	pb "github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/infrastructure/grpc/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func mapCompanyPBToDomain(companyPb *pb.Company) *domain.Company {
	company := &domain.Company{
		CompanyName: companyPb.CompanyName,
		Username:    companyPb.Username,
		Email:       companyPb.Email,
		PhoneNumber: companyPb.PhoneNumber,
		Description: companyPb.Description,
		Location:    companyPb.Location,
		Website:     companyPb.Website,
		CompanySize: companyPb.CompanySize,
		Industry:    companyPb.Industry,
	}
	return company
}

func mapJobOfferToDomain(jobOfferDto *pb.JobOffer) *domain.JobOffer {
	jobOffer := &domain.JobOffer{
		Id: primitive.NewObjectID(),
		Company: domain.Company{
			CompanyName: jobOfferDto.Company.CompanyName,
			Username:    jobOfferDto.Company.Username,
			Email:       jobOfferDto.Company.Email,
			PhoneNumber: jobOfferDto.Company.PhoneNumber,
			Description: jobOfferDto.Company.Description,
			Location:    jobOfferDto.Company.Location,
			Website:     jobOfferDto.Company.Website,
			CompanySize: jobOfferDto.Company.CompanySize,
			Industry:    jobOfferDto.Company.Industry,
		},
		JobDescription: jobOfferDto.JobDescription,
		Position:       jobOfferDto.Position,
		Prerequisites:  jobOfferDto.Prerequisites,
		EmploymentType: enum.EmploymentType(jobOfferDto.EmploymentType),
		Published:      time.Now(),
	}
	return jobOffer
}

func mapUserToDomain(userPb *pb.User) *domain.RegisteredUser {
	user := &domain.RegisteredUser{
		Username:    (*userPb).Username,
		FirstName:   (*userPb).FirstName,
		LastName:    (*userPb).LastName,
		Email:       (*userPb).Email,
		PhoneNumber: (*userPb).PhoneNumber,
		Gender:      enum.Gender((*userPb).Gender),
		DateOfBirth: (*((*userPb).DateOfBirth)).AsTime(),
		Biography:   (*userPb).Biography,
		IsPrivate:   (*userPb).IsPrivate,
	}

	user.Experiences = []domain.Experience{}
	for _, experience := range (*userPb).Experiences {
		id, err := primitive.ObjectIDFromHex(experience.Id)
		if err != nil {
			continue
		}

		user.Experiences = append(user.Experiences, domain.Experience{
			Id:                 id,
			Title:              experience.Title,
			EmploymentType:     enum.EmploymentType(experience.EmploymentType),
			CompanyName:        experience.CompanyName,
			Location:           experience.Location,
			IsCurrentlyWorking: experience.IsCurrentlyWorking,
			StartDate:          experience.StartDate.AsTime(),
			EndDate:            experience.EndDate.AsTime(),
			Industry:           experience.Industry,
			Description:        experience.Description,
		})
	}

	user.Educations = []domain.Education{}
	for _, education := range (*userPb).Educations {
		id, err := primitive.ObjectIDFromHex(education.Id)
		if err != nil {
			continue
		}

		user.Educations = append(user.Educations, domain.Education{
			Id:           id,
			School:       education.School,
			Degree:       enum.Degree(education.Degree),
			FieldOfStudy: education.FieldOfStudy,
			StartDate:    education.StartDate.AsTime(),
			EndDate:      education.EndDate.AsTime(),
			Description:  education.Description,
		})
	}

	user.Skills = []string{}
	for _, skill := range (*userPb).Skills {
		user.Skills = append(user.Skills, skill)
	}

	user.Interests = []primitive.ObjectID{}
	for _, interest := range (*userPb).Interests {
		interestId, err := primitive.ObjectIDFromHex(interest)
		if err != nil {
			continue
		}

		user.Interests = append(user.Interests, interestId)
	}

	user.Connections = []primitive.ObjectID{}
	for _, connection := range (*userPb).Connections {
		connectionId, err := primitive.ObjectIDFromHex(connection)
		if err != nil {
			continue
		}
		user.Connections = append(user.Connections, connectionId)
	}

	return user
}
