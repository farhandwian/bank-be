package utils

import "bank-backend/module/user/entity"

func UserToDTO(user entity.User) entity.RegisterResponse {
	response := entity.RegisterResponse{
		UserID:      user.ID.String(),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Address:     user.Address,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt.String(),
	}
	return response
}

func UserUpdateToDTO(user entity.User) entity.UpdateProfileResponse {
	response := entity.UpdateProfileResponse{
		UserID:      user.ID.String(),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Address:     user.Address,
		PhoneNumber: user.PhoneNumber,
		Updated_at:  user.UpdatedAt.String(),
	}
	return response
}
