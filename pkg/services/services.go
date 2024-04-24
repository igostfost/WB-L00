package services

import "WB_L00/pkg/repository"

type Service struct {
	repo *repository.Repostitory
}

func NewService(repo *repository.Repostitory) *Service {
	return &Service{repo: repo}
}
