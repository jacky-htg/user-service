package service

import (
	"context"
	"database/sql"
	"regexp"
	users "user-service/pb"

	"user-service/internal/model"
	"user-service/internal/pkg/app"
	"user-service/internal/pkg/db/redis"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Company struct
type Company struct {
	Db    *sql.DB
	Cache *redis.Cache
}

// Registration new company
func (u *Company) Registration(ctx context.Context, in *users.CompanyRegistration) (*users.CompanyRegistration, error) {
	var output users.CompanyRegistration
	var err error

	// validate company request
	{
		if len(in.GetCompany().GetName()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid company name")
		}

		if len(in.GetCompany().GetCode()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid company code")
		}

		if len(in.GetCompany().GetAddress()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid company address")
		}

		if len(in.GetCompany().GetCity()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid company city")
		}

		if len(in.GetCompany().GetProvince()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid company province")
		}

		if len(in.GetCompany().GetPhone()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid company phone")
		}

		if len(in.GetCompany().GetPic()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid company pic")
		}

		if len(in.GetCompany().GetPicPhone()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid company pic phone")
		}

		// company code validation
		{
			companyModel := model.Company{}
			companyModel.Pb.Code = in.GetCompany().GetCode()
			err = companyModel.GetByCode(ctx, u.Db)
			if err != nil {
				if st, ok := status.FromError(err); ok && st.Code() != codes.NotFound {
					return &output, err
				}
			}

			if len(companyModel.Pb.GetId()) > 0 {
				return &output, status.Error(codes.AlreadyExists, "code must be unique")
			}
		}

		switch in.GetCompany().GetPackageOfFeature().String() {
		case "ALL":
			var featurePackageModel model.FeaturePackage
			featurePackageModel.Pb.Name = in.GetCompany().GetPackageOfFeature()
			err = featurePackageModel.GetByName(ctx, u.Db)
			if err != nil {
				return &output, err
			}

			in.GetCompany().Features = featurePackageModel.Pb.Features
		case "SIMPLE":
			var featurePackageModel model.FeaturePackage
			featurePackageModel.Pb.Name = in.GetCompany().GetPackageOfFeature()
			err = featurePackageModel.GetByName(ctx, u.Db)
			if err != nil {
				return &output, err
			}

			in.GetCompany().Features = featurePackageModel.Pb.Features
		case "CUSTOME":
			if len(in.GetCompany().GetFeatures()) == 0 {
				return &output, status.Error(codes.InvalidArgument, "Please supply valid company features")
			}

			for _, feature := range in.GetCompany().GetFeatures() {
				var featureModel model.Feature
				featureModel.Pb.Id = feature.GetId()
				err = featureModel.Get(ctx, u.Db)
				if err != nil {
					return &output, err
				}
			}
		}
	}

	// user validation
	{

		if len(in.GetUser().GetEmail()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid email")
		}

		if len(in.GetUser().GetName()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid name")
		}

		// username validation
		{
			if len(in.GetUser().GetUsername()) == 0 {
				return &output, status.Error(codes.InvalidArgument, "Please supply valid username")
			}

			if len(in.GetUser().GetUsername()) < 8 {
				return &output, status.Error(codes.InvalidArgument, "username min 8 character")
			}

			userModel := model.User{}
			userModel.Pb.Username = in.GetUser().GetUsername()
			err = userModel.GetByUsername(ctx, u.Db)
			if err != nil {
				if st, ok := status.FromError(err); ok && st.Code() != codes.NotFound {
					return &output, err
				}
			}

			if len(userModel.Pb.GetId()) > 0 {
				return &output, status.Error(codes.AlreadyExists, "username must be unique")
			}
		}

		// email validation
		{
			var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
			if valid := func(e string) bool {
				if len(e) < 3 && len(e) > 254 {
					return false
				}
				return emailRegex.MatchString(e)
			}(in.GetUser().GetEmail()); !valid {
				return &output, status.Error(codes.InvalidArgument, "Please supply valid email")
			}

			userModel := model.User{}
			userModel.Pb.Email = in.GetUser().GetEmail()
			err = userModel.GetByEmail(ctx, u.Db)
			if err != nil {
				if st, ok := status.FromError(err); ok && st.Code() != codes.NotFound {
					return &output, err
				}
			}

			if len(userModel.Pb.GetId()) > 0 {
				return &output, status.Error(codes.AlreadyExists, "email must be unique")
			}
		}
	}

	companyRegisterModel := model.CompanyRegister{}
	companyRegisterModel.Pb = users.CompanyRegistration{
		Company: in.GetCompany(),
		User:    in.GetUser(),
	}
	companyRegisterModel.Password = generateRandomPassword()
	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return &output, status.Errorf(codes.Internal, "begin tx: %v", err)
	}

	err = companyRegisterModel.Registration(ctx, u.Db, tx)
	if err != nil {
		tx.Rollback()
		return &output, err
	}

	tx.Commit()

	// TODO : currenly we accept logo url, next we should accept logo file and process to upload log
	// TODO : send email to inform username and password
	return &companyRegisterModel.Pb, err
}

// Update Company
func (u *Company) Update(ctx context.Context, in *users.Company) (*users.Company, error) {
	var output users.Company
	var err error
	var companyModel model.Company

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		companyModel.Pb.Id = in.GetId()
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &output, err
	}

	// get user login
	var userLogin model.User
	userLogin.Pb.Id = ctx.Value(app.Ctx("userID")).(string)
	err = userLogin.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	if userLogin.Pb.GetCompanyId() != in.GetId() {
		return &output, status.Error(codes.Unauthenticated, "its not your company")
	}

	if len(userLogin.Pb.GetBranchId()) > 0 || len(userLogin.Pb.GetRegionId()) > 0 {
		return &output, status.Error(codes.Unauthenticated, "only user company can update the company")
	}

	err = companyModel.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	if len(in.GetName()) > 0 {
		companyModel.Pb.Name = in.GetName()
	}

	if len(in.GetAddress()) > 0 {
		companyModel.Pb.Address = in.GetAddress()
	}

	if len(in.GetCity()) > 0 {
		companyModel.Pb.City = in.GetCity()
	}

	if len(in.GetProvince()) > 0 {
		companyModel.Pb.Province = in.GetProvince()
	}

	if len(in.GetPhone()) > 0 {
		companyModel.Pb.Phone = in.GetPhone()
	}

	if len(in.GetPic()) > 0 {
		companyModel.Pb.Pic = in.GetPic()
	}

	if len(in.GetPicPhone()) > 0 {
		companyModel.Pb.PicPhone = in.GetPicPhone()
	}

	if in.GetPackageOfFeature() != companyModel.Pb.GetPackageOfFeature() {
		companyModel.UpdateFeatures = true
		switch in.GetPackageOfFeature().String() {
		case "ALL":
			var featurePackageModel model.FeaturePackage
			featurePackageModel.Pb.Name = in.GetPackageOfFeature()
			err = featurePackageModel.GetByName(ctx, u.Db)
			if err != nil {
				return &output, err
			}

			companyModel.Pb.Features = featurePackageModel.Pb.Features
		case "SIMPLE":
			var featurePackageModel model.FeaturePackage
			featurePackageModel.Pb.Name = in.GetPackageOfFeature()
			err = featurePackageModel.GetByName(ctx, u.Db)
			if err != nil {
				return &output, err
			}

			companyModel.Pb.Features = featurePackageModel.Pb.Features
		case "CUSTOME":
			if len(in.GetFeatures()) == 0 {
				return &output, status.Error(codes.InvalidArgument, "Please supply valid company features")
			}

			for _, feature := range in.GetFeatures() {
				var featureModel model.Feature
				featureModel.Pb.Id = feature.GetId()
				err = featureModel.Get(ctx, u.Db)
				if err != nil {
					return &output, err
				}
			}

			companyModel.Pb.Features = in.GetFeatures()
		}
		companyModel.Pb.PackageOfFeature = in.GetPackageOfFeature()
	} else {
		if in.GetPackageOfFeature() == companyModel.Pb.GetPackageOfFeature() && companyModel.Pb.GetPackageOfFeature().String() == "CUSTOME" {
			if len(in.GetFeatures()) > 0 {
				companyModel.UpdateFeatures = true
				for _, feature := range in.GetFeatures() {
					var featureModel model.Feature
					featureModel.Pb.Id = feature.GetId()
					err = featureModel.Get(ctx, u.Db)
					if err != nil {
						return &output, err
					}
				}

				companyModel.Pb.Features = in.GetFeatures()
			}
		}
	}

	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return &output, status.Errorf(codes.Internal, "begin tx: %v", err)
	}

	err = companyModel.Update(ctx, u.Db, tx)
	if err != nil {
		tx.Rollback()
		return &output, err
	}

	tx.Commit()

	return &companyModel.Pb, nil
}

// View Company
func (u *Company) View(ctx context.Context, in *users.Id) (*users.Company, error) {
	var output users.Company
	var err error
	var companyModel model.Company

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		companyModel.Pb.Id = in.GetId()
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &output, err
	}

	// get user login
	var userLogin model.User
	userLogin.Pb.Id = ctx.Value(app.Ctx("userID")).(string)
	err = userLogin.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	if userLogin.Pb.GetCompanyId() != in.GetId() {
		return &output, status.Error(codes.Unauthenticated, "its not your company")
	}

	err = companyModel.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	return &companyModel.Pb, err
}
