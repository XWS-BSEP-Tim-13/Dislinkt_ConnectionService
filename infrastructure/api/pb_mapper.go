package api

import (
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/domain"
	pb "github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/infrastructure/grpc/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func mapCompanyToPB(company *domain.Company) *pb.Company {
	companyPb := &pb.Company{
		Id:          company.Id.Hex(),
		CompanyName: company.CompanyName,
		Username:    company.Username,
		Description: company.Description,
		Location:    company.Location,
		Website:     company.Website,
		CompanySize: company.CompanySize,
		Industry:    company.Industry,
	}
	return companyPb
}

func mapJobOfferToPB(job *domain.JobOffer) *pb.JobOffer {
	jobPb := &pb.JobOffer{
		Id:             job.Id.Hex(),
		JobDescription: job.JobDescription,
		Position:       job.Position,
		Prerequisites:  job.Prerequisites,
		EmploymentType: pb.EmploymentType(job.EmploymentType),
		Company: &pb.Company{
			Id:          job.Company.Id.Hex(),
			CompanyName: job.Company.CompanyName,
			Username:    job.Company.Username,
			Description: job.Company.Description,
			Location:    job.Company.Location,
			Website:     job.Company.Website,
			CompanySize: job.Company.CompanySize,
			Industry:    job.Company.Industry,
		},
	}
	return jobPb
}

func mapUserToPB(user *domain.RegisteredUser) *pb.User {
	userPb := &pb.User{
		Id:          user.Id.Hex(),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Gender:      pb.User_Gender(user.Gender),
		DateOfBirth: timestamppb.New(user.DateOfBirth),
		Biography:   user.Biography,
		IsPrivate:   user.IsPrivate,
		Username:    user.Username,
	}

	for _, experience := range user.Experiences {
		userPb.Experiences = append(userPb.Experiences, &pb.Experience{
			Id:                 experience.Id.Hex(),
			Title:              experience.Title,
			EmploymentType:     pb.Experience_EmploymentType(experience.EmploymentType),
			CompanyName:        experience.CompanyName,
			Location:           experience.Location,
			IsCurrentlyWorking: experience.IsCurrentlyWorking,
			StartDate:          timestamppb.New(experience.StartDate),
			EndDate:            timestamppb.New(experience.EndDate),
			Industry:           experience.Industry,
			Description:        experience.Description,
		})
	}

	for _, education := range user.Educations {
		userPb.Educations = append(userPb.Educations, &pb.Education{
			Id:           education.Id.Hex(),
			School:       education.School,
			Degree:       pb.Education_Degree(education.Degree),
			FieldOfStudy: education.FieldOfStudy,
			StartDate:    timestamppb.New(education.StartDate),
			EndDate:      timestamppb.New(education.EndDate),
			Description:  education.Description,
		})
	}

	for _, skill := range user.Skills {
		userPb.Skills = append(userPb.Skills, skill)
	}

	for _, interest := range user.Interests {
		userPb.Interests = append(userPb.Interests, interest.Hex())
	}

	for _, connection := range user.Connections {
		userPb.Connections = append(userPb.Connections, connection)
	}

	return userPb
}

func mapConnectionRequestToPB(request *domain.ConnectionRequest) *pb.ConnectionRequest {
	connectionPb := &pb.ConnectionRequest{
		Id:          request.Id.Hex(),
		From:        mapUserToPB(&request.From),
		To:          mapUserToPB(&request.To),
		RequestTime: timestamppb.New(request.RequestTime),
	}
	return connectionPb
}

func mapEventToPB(request *domain.Event) *pb.Event {
	eventPb := &pb.Event{
		Id:        request.Id.Hex(),
		Action:    request.Action,
		User:      request.User,
		Published: timestamppb.New(request.Published),
	}
	return eventPb
}
