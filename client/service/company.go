package service

import (
	"context"
	"log"

	"github.com/jacky-htg/erp-proto/go/pb/users"
)

// Registration service client
func Registration(ctx context.Context, company users.CompanyServiceClient) {
	value, ok := users.EnumPackageOfFeature_value["ALL"]
	if !ok {
		log.Fatal("Invalid package of Feature")
	}

	response, err := company.Registration(ctx, &users.CompanyRegistration{
		Company: &users.Company{
			Address:          "Jalan Minangkabau",
			City:             "Jakarta Pusat",
			Code:             "SRTU",
			Logo:             "logo",
			Name:             "Sri Ratu",
			Npwp:             "npwp",
			PackageOfFeature: users.EnumPackageOfFeature(value),
			Phone:            "081244444444",
			Pic:              "Admin Pasaraya",
			PicPhone:         "081312222222",
			Province:         "Jakarta",
		},
		User: &users.User{
			Email:    "admin.sriratu@gmail.com",
			Name:     "Admin Sri Ratu",
			Username: "admin-srtu",
		},
	})
	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}

// UpdateCompany service client
func UpdateCompany(ctx context.Context, company users.CompanyServiceClient) {
	value, ok := users.EnumPackageOfFeature_value["ALL"]
	if !ok {
		log.Fatal("Invalid package of Feature")
	}

	response, err := company.Update(setMetadata(ctx), &users.Company{
		Id:               "45e5719f-6797-4463-b454-7413ce3a58f7",
		Address:          "Jalan Minangkabau",
		City:             "Jakarta Pusat",
		Logo:             "logo",
		Name:             "Pasaraya",
		Npwp:             "npwp",
		PackageOfFeature: users.EnumPackageOfFeature(value),
		Phone:            "081244444444",
		Pic:              "Admin Pasaraya",
		PicPhone:         "081312222222",
		Province:         "Jakarta",
	})
	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}

// ViewCompany service client
func ViewCompany(ctx context.Context, company users.CompanyServiceClient) {
	response, err := company.View(setMetadata(ctx), &users.Id{Id: "45e5719f-6797-4463-b454-7413ce3a58f7"})
	if err != nil {
		log.Fatalf("Error when calling grpc service: %s", err)
	}
	log.Printf("Resp received: %v", response)
}
