package dto

type RegisterRequest struct {
	Email      string  `json:"email"`
	Password   string  `json:"password"`
	Role       string  `json:"role"`
	EmployeeId *string `json:"employeeId,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string           `json:"token"`
	User  UserInfoResponse `json:"user"`
}

type UserInfoResponse struct {
	ID       string            `json:"id"`
	Email    string            `json:"email"`
	Role     string            `json:"role"`
	Employee *EmployeeResponse `json:"employee,omitempty"`
}
