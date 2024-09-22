package user

func userToDTO(user User) RegisterResponse {
	response := RegisterResponse{
		UserID:      user.ID.String(),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Address:     user.Address,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt.String(),
	}
	return response
}

func userUpdataeToDTO(user User) UpdateResponse {
	response := UpdateResponse{
		UserID:      user.ID.String(),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Address:     user.Address,
		PhoneNumber: user.PhoneNumber,
		Updated_at:  user.UpdatedAt.String(),
	}
	return response
}
