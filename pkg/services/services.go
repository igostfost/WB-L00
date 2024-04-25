package services

import "WB_L00/pkg/repository"

type Service struct {
	Repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{Repo: repo}
}
