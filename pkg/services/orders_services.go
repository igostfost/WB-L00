package services

import "WB_L00/types"

func (s *Service) CreateOrder(newOrder types.Order) error {
	return s.Repo.CreateOrder(newOrder)
}

func (s *Service) GetAllOrdersFromDB() ([]types.Order, error) {
	return s.Repo.GetAllOrdersFromDB()
}
