package domain

import (
	"fmt"
	"regexp"
	"time"

	core_errors "github.com/IwantHappiness/url-shortener/internal/core/errors"
)

type User struct {
	ID        int       `json:"id"`
	Version   int       `json:"version"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUser(id, version int, nickname, email string, createdAt time.Time) User {
	return User{
		ID:        id,
		Version:   version,
		Nickname:  nickname,
		Email:     email,
		CreatedAt: createdAt,
	}
}

func NewUserUninitialized(nickname, email string) User {
	return NewUser(UninitializedID, UninitializedVersion, nickname, email, UninitializedCreatedAt)
}

func (u *User) ValidateUser() error {
	nicknameLen := len([]rune(u.Nickname))
	if nicknameLen < 3 || nicknameLen > 20 {
		return fmt.Errorf("invalid `Nickname` len: %d: %w", nicknameLen, core_errors.ErrInvalidArgument)
	}

	emailLen := len([]rune(u.Email))
	if emailLen < 3 || emailLen > 255 {
		return fmt.Errorf("invalid `Email` len: %d: %w", emailLen, core_errors.ErrInvalidArgument)
	}

	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{1,}$`)
	if !re.MatchString(u.Email) {
		return fmt.Errorf("invalid `Email` format: %s: %w", u.Email, core_errors.ErrInvalidArgument)
	}

	return nil
}

type UserPatch struct {
	Nickname Nullable[string]
	Email    Nullable[string]
}

func NewUserPatch(nickname Nullable[string], email Nullable[string]) UserPatch {
	return UserPatch{
		Nickname: nickname,
		Email:    email,
	}
}

func (p *UserPatch) Validate() error {
	if p.Nickname.Set && p.Nickname.Value == nil {
		return fmt.Errorf("'Nickname' can't be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}

	if p.Email.Set && p.Email.Value == nil {
		return fmt.Errorf("'Email' can't be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate user patch: %w", err)
	}

	tmp := *u

	if patch.Nickname.Set {
		tmp.Nickname = *patch.Nickname.Value
	}

	if patch.Email.Set {
		tmp.Email = *patch.Email.Value
	}

	if err := tmp.ValidateUser(); err != nil {
		return fmt.Errorf("validate patched user: %w", err)
	}

	*u = tmp

	return nil
}
