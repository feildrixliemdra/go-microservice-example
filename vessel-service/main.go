package main

import (
	"context"
	"errors"
	pb "shippy-microservice/vessel-service/proto/vessel"
	 "github.com/micro/go-micro"
	 
)

type Repository interface {
	FindAvailable(spec *pb.Specification) (*pb.Vessel, error)
}

type VesselRepository struct {
	vessels []*pb.Vessel
}

// FindAvailable - checks a specification against a map of vessels,
// if capacity and max weight are below a vessels capacity and max weight,
// then return that vessel.
func (repo *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	for _, vessel := range repo.vessels {
		if spec.Capacity <= vessel.Capacity && spec.MaxWeight <= vessel.MaxWeight {
			return vessel, nil
		}
	}
	return nil, errors.New("No vessel found by that spec")
}

type service struct {
	repo Repository
}

func (service *service) FindAvailable(ctx context.Context, req *pb.Specification, res *pb.Response) error {
	// Find the next available vessel
	vessel, err := service.repo.FindAvailable(req)
	if err != nil {
		return err
	}

	res.Vessel = vessel
	return nil
}
func main() {
	vessels := []*pb.Vessel{
		&pb.Vessel{
			Id:        "vessel001",
			Name:      "Boat",
			MaxWeight: 20000,
			Capacity:  500,
		},
	}
	repo := &VesselRepository{vessels}

	serve := micro.NewService(
		micro.Name("vessel.service")
	)
	serve.Init()

	// Register our implementation with 
	pb.RegisterVesselServiceHandler(serve.Server(), &service{repo})

	if err := serve.Run(); err != nil {
		fmt.Println(err)
	}
}
