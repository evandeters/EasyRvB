package service

type Service interface {
	ReadConfig(string) error
}

type GenericService struct {
	Name     string
	Port     int16
	Protocol string
	Limit    int
}

type Dependency struct {
	Name string
}
