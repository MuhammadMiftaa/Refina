package service

import (
	"database/sql"
	"errors"
	"time"

	"server/internal/repository"
	"server/internal/types/dto"
	"server/internal/types/entity"
	helper "server/internal/utils"
)

type UsersService interface {
	Register(user dto.UsersRequest) (dto.UsersResponse, error)
	Login(user dto.UsersRequest) (*string, error)
	OAuthLogin(name string, email string) (*string, error)
	GetAllUsers() ([]dto.UsersResponse, error)
	GetUserByID(id string) (dto.UsersResponse, error)
	GetUserByEmail(email string) (dto.UsersResponse, error)
	UpdateUser(id string, userNew dto.UsersRequest) (dto.UsersResponse, error)
	VerifyUser(email string) (dto.UsersResponse, error)
	DeleteUser(id string) (dto.UsersResponse, error)
}

type usersService struct {
	userRepository repository.UsersRepository
}

func NewUsersService(usersRepository repository.UsersRepository) UsersService {
	return &usersService{usersRepository}
}

func (user_serv *usersService) Register(user dto.UsersRequest) (dto.UsersResponse, error) {
	// VALIDASI APAKAH NAME, EMAIL, PASSWORD KOSONG
	if user.Name == "" || user.Email == "" || user.Password == "" {
		return dto.UsersResponse{}, errors.New("name, email, and password cannot be blank")
	}

	// VALIDASI UNTUK FORMAT EMAIL SUDAH BENAR
	if isValid := helper.EmailValidator(user.Email); !isValid {
		return dto.UsersResponse{}, errors.New("please enter a valid email address")
	}

	// MENGECEK APAKAH EMAIL SUDAH DIGUNAKAN
	userExist, err := user_serv.userRepository.GetUserByEmail(user.Email)
	if err == nil && (userExist.Email != "") {
		return dto.UsersResponse{}, errors.New("email already exists")
	}

	// VALIDASI PASSWORD SUDAH SESUAI, MIN 8 KARAKTER, MENGANDUNG ALFABET DAN NUMERIK
	hasMinLen, hasLetter, hasDigit := helper.PasswordValidator(user.Password)
	if !hasMinLen {
		return dto.UsersResponse{}, errors.New("password must be at least 8 characters long")
	}
	if !hasLetter {
		return dto.UsersResponse{}, errors.New("password must contain at least one letter")
	}
	if !hasDigit {
		return dto.UsersResponse{}, errors.New("password must contain at least one number")
	}

	// HASHING PASSWORD MENGGUNAKAN BCRYPT
	hashedPassword, err := helper.PasswordHashing(user.Password)
	if err != nil {
		return dto.UsersResponse{}, err
	}
	user.Password = hashedPassword

	newUser, err := user_serv.userRepository.CreateUser(entity.Users{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		return dto.UsersResponse{}, err
	}

	userResponse := helper.ConvertToResponseType(newUser).(dto.UsersResponse)

	return userResponse, nil
}

func (user_serv *usersService) Login(user dto.UsersRequest) (*string, error) {
	// VALIDASI APAKAH EMAIL DAN PASSWORD KOSONG
	if user.Email == "" || user.Password == "" {
		return nil, errors.New("email and password cannot be blank")
	}

	// MENGECEK APAKAH USER SUDAH TERDAFTAR
	userExist, err := user_serv.userRepository.GetUserByEmail(user.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// VALIDASI APAKAH PASSWORD SUDAH SESUAI
	if !helper.ComparePass(userExist.Password, user.Password) {
		return nil, errors.New("password is incorrect")
	}

	token, err := helper.GenerateToken(userExist.ID.String(), userExist.Name, userExist.Email)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (user_serv *usersService) OAuthLogin(name string, email string) (*string, error) {
	token, err := helper.GenerateToken("99", name, email)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (user_serv *usersService) GetAllUsers() ([]dto.UsersResponse, error) {
	users, err := user_serv.userRepository.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var usersResponse []dto.UsersResponse
	for _, user := range users {
		userResponse, _ := helper.ConvertToResponseType(user).(dto.UsersResponse)
		usersResponse = append(usersResponse, userResponse)
	}

	return usersResponse, nil
}

func (user_serv *usersService) GetUserByID(id string) (dto.UsersResponse, error) {
	user, err := user_serv.userRepository.GetUserByID(id)
	if err != nil {
		return dto.UsersResponse{}, err
	}

	userResponse := helper.ConvertToResponseType(user)

	return userResponse.(dto.UsersResponse), nil
}

func (user_serv *usersService) GetUserByEmail(email string) (dto.UsersResponse, error) {
	user, err := user_serv.userRepository.GetUserByEmail(email)
	if err != nil {
		return dto.UsersResponse{}, err
	}

	userResponse := helper.ConvertToResponseType(user)

	return userResponse.(dto.UsersResponse), nil
}

func (user_serv *usersService) UpdateUser(id string, userNew dto.UsersRequest) (dto.UsersResponse, error) {
	// MENGAMBIL DATA YANG INGIN DI UPDATE
	user, err := user_serv.userRepository.GetUserByID(id)
	if err != nil {
		return dto.UsersResponse{}, err
	}

	// VALIDASI APAKAH FULLNAME & EMAIL KOSONG
	if userNew.Name == "" && userNew.Email == "" {
		return dto.UsersResponse{}, errors.New("fullname and email cannot be blank")
	}

	// VALIDASI APAKAH FULLNAME / EMAIL SUDAH DI INPUT
	if userNew.Name != "" {
		user.Name = userNew.Name
	}

	if userNew.Email != "" {
		// VALIDASI UNTUK FORMAT EMAIL SUDAH BENAR
		if isValid := helper.EmailValidator(userNew.Email); !isValid {
			return dto.UsersResponse{}, errors.New("please enter a valid email address")
		}
		// MENGECEK APAKAH EMAIL SUDAH DIGUNAKAN
		existingUser, err := user_serv.userRepository.GetUserByEmail(userNew.Email)
		if err == nil && existingUser.ID != user.ID {
			return dto.UsersResponse{}, errors.New("email already in use by another user")
		}
		user.Email = userNew.Email
	}

	userUpdated, err := user_serv.userRepository.UpdateUser(user)
	if err != nil {
		return dto.UsersResponse{}, err
	}

	userResponse := helper.ConvertToResponseType(userUpdated)

	return userResponse.(dto.UsersResponse), nil
}

func (user_serv *usersService) VerifyUser(email string) (dto.UsersResponse, error) {
	// MENGAMBIL DATA YANG INGIN DI UPDATE
	user, err := user_serv.userRepository.GetUserByEmail(email)
	if err != nil {
		return dto.UsersResponse{}, err
	}

	current := time.Now()
	user.EmailVerfiedAt = sql.NullTime{
		Time:  current,
		Valid: true,
	}

	userExist, err := user_serv.userRepository.UpdateUser(user)
	if err != nil {
		return dto.UsersResponse{}, err
	}

	userResponse := helper.ConvertToResponseType(userExist).(dto.UsersResponse)

	return userResponse, nil
}

func (user_serv *usersService) DeleteUser(id string) (dto.UsersResponse, error) {
	// MENGAMBIL DATA YANG INGIN DI DELETE
	user, err := user_serv.userRepository.GetUserByID(id)
	if err != nil {
		return dto.UsersResponse{}, err
	}

	userDeleted, err := user_serv.userRepository.DeleteUser(user)
	if err != nil {
		return dto.UsersResponse{}, err
	}

	userResponse := helper.ConvertToResponseType(userDeleted)

	return userResponse.(dto.UsersResponse), nil
}
