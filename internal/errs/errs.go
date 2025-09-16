package errs

import "errors"

var (
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrWrongTokenType        = errors.New("wrong token type")
	ErrUserNotFound          = errors.New("user not found")
	ErrEmailNotVerify        = errors.New("email not verify")
	ErrTokenNotFound         = errors.New("token not found")
	ErrCategoryAlreadyExists = errors.New("category already axists")
	ErrNotAdmin              = errors.New("user does not admin")
	ErrFailedToCash          = errors.New("failed to cashed data")
)
