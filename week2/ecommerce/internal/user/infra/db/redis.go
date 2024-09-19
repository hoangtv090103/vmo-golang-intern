package db

// import (
// 	"context"
// 	"ecommerce/cmd/config"
// 	"ecommerce/internal/user/domain"
// 	"strconv"
// )

// type UserRepoRedis struct {
// 	redis *config.Redis
// }

// func NewUserRepoRedis(redis *config.Redis) *UserRepoRedis {
// 	return &UserRepoRedis{
// 		redis: redis,
// 	}
// }

// func (u *UserRepoRedis) Create(user domain.User) error {
// 	err := u.redis.GetClient().Set(context.Background(), user.Username, user.Email, 0).Err()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (u *UserRepoRedis) GetByUsername(username string) (domain.User, error) {
// 	var user domain.User
// 	email, err := u.redis.GetClient().Get(context.Background(), username).Result()
// 	if err != nil {
// 		return domain.User{}, err
// 	}
// 	user.Username = username
// 	user.Email = email
// 	return user, nil
// }

// func (u *UserRepoRedis) GetByEmail(email string) (domain.User, error) {
// 	var user domain.User
// 	username, err := u.redis.GetClient().Get(context.Background(), email).Result()
// 	if err != nil {
// 		return domain.User{}, err
// 	}
// 	user.Username = username
// 	user.Email = email
// 	return user, nil
// }

// func (u *UserRepoRedis) Update(user domain.User) error {
// 	err := u.redis.GetClient().Set(context.Background(), user.Username, user.Email, 0).Err()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (u *UserRepoRedis) Delete(id int) error {
// 	idStr := strconv.Itoa(id)
// 	err := u.redis.GetClient().Del(context.Background(), idStr).Err()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
