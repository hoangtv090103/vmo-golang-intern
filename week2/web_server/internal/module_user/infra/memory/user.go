// internal/infra/memory/user_repo_memory.go
package memory

import (
	"errors"
	"week2-clean-architecture/internal/module_user/domain"
)

type UserRepoMemory struct {
	users  map[int]*domain.User
	nextID int
}

func NewUserRepoMemory() *UserRepoMemory {
	return &UserRepoMemory{
		users:  make(map[int]*domain.User),
		nextID: 1,
	}
}

func (r *UserRepoMemory) GetAll() ([]*domain.User, error) {
	users := make([]*domain.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepoMemory) GetByUsername(username string) (*domain.User, error) {
	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *UserRepoMemory) Create(user *domain.User) error {
	user.UserID = r.nextID
	r.users[user.UserID] = user
	r.nextID++
	return nil
}
