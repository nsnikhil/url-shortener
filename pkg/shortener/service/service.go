package service

type Service struct {
	shortner ShortenerService
}

func NewService(shortner ShortenerService) *Service {
	return &Service{
		shortner: shortner,
	}
}

func (s *Service) GetShortenerService() ShortenerService {
	return s.shortner
}
