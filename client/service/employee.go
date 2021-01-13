package service

import (
	"context"
	"io"
	"log"
	users "user-service/pb"
)

// CreateEmployee service client
func CreateEmployee(ctx context.Context, employee users.EmployeeServiceClient) {
	response, err := employee.Create(setMetadata(ctx), &users.Employee{
		Code:     "RAS",
		Name:     "Rijal Asepnugroho",
		Address:  "Jl Address",
		City:     "Lampung",
		Province: "Lampung",
		Jabatan:  "Head of IT",
		User: &users.User{
			Id: "362e5eb5-51b5-412d-8d5f-081c5aa494ce",
		},
	})
	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}

// UpdateEmployee service client
func UpdateEmployee(ctx context.Context, employee users.EmployeeServiceClient) {
	response, err := employee.Update(setMetadata(ctx), &users.Employee{
		Id:       "382a011d-8c2f-4e6b-8e8f-7b41f1332dca",
		Name:     "Rijal Asepnugroho",
		Address:  "Jl Address",
		City:     "Lampung",
		Province: "Lampung",
		Jabatan:  "CTO",
	})
	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}

// ViewEmployee service client
func ViewEmployee(ctx context.Context, employee users.EmployeeServiceClient) {
	response, err := employee.View(setMetadata(ctx), &users.Id{Id: "382a011d-8c2f-4e6b-8e8f-7b41f1332dca"})
	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}

// DeleteEmployee service client
func DeleteEmployee(ctx context.Context, employee users.EmployeeServiceClient) {
	response, err := employee.Delete(setMetadata(ctx), &users.Id{Id: "382a011d-8c2f-4e6b-8e8f-7b41f1332dca"})
	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}

// ListEmployee service client
func ListEmployee(ctx context.Context, employee users.EmployeeServiceClient) {
	stream, err := employee.List(setMetadata(ctx), &users.ListEmployeeRequest{})
	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Fatal("end stream")
			break
		}
		if err != nil {
			log.Fatalf("cannot receive %v", err)
		}
		log.Printf("Resp received: %s : %v", resp.GetEmployee(), resp.GetPagination())
	}
}
