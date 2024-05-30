package services

import (
	"CustomerOrderApi/Customer/entities"
	"CustomerOrderApi/Customer/repositories"
	sharedentities "CustomerOrderApi/shared/entities"
	"CustomerOrderApi/shared/helpers"
	"fmt"
	"github.com/erenkaratas99/COApiCore/pkg"
)

type Service struct {
	repository *repositories.Repository
	client     *pkg.RestClient
}

func NewService(repo *repositories.Repository) *Service {
	return &Service{repository: repo}
}

func (s *Service) CreateCustomerService(customerReq *entities.CustomerRequestModel) (*string, error) {
	insertedID, err := s.repository.InsertCustomer(customerReq)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return insertedID, nil
}

func (s *Service) GetAllCustomersService(l, o int64) (*sharedentities.ResponseModel, error) {
	customers, totalCount, err := s.repository.GetAllCustomers(l, o, true)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	length := len(customers)
	resp := helpers.NewResponseModel(totalCount, &length, customers)
	return resp, nil
}

// todo : bu yapı kullanımı
//func (s *Service) GetAllCustomers(ctx context.Context, skip int, limit int) ([]*entities.Customer, int64, error) {
//	return s.repo.GetAllCustomers(ctx, skip, limit)
//}

func (s *Service) CustomerGetById(id string) (*sharedentities.ResponseModel, error) {
	customers, totalCount, err := s.repository.CustomerGetById(id, true)
	len := 1
	customerResp := helpers.NewResponseModel(totalCount, &len, customers)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return customerResp, nil
}
