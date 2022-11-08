package service

import (
	"context"
	"cryptocurrency/internal/models"
	"cryptocurrency/internal/models/mongo_models"
	"cryptocurrency/internal/models/response_models/slack_bot"
	"cryptocurrency/internal/utils/telegram_api"
	"strconv"
	"time"
)

type managerService struct {
	ManagerRepository models.ManagerRepository
	AccountRepository models.AccountRepository
}

type MSConfig struct {
	ManagerRepository models.ManagerRepository
	AccountRepository models.AccountRepository
}

func NewManagerService(c *MSConfig) models.ManagerService {
	return &managerService{
		ManagerRepository: c.ManagerRepository,
		AccountRepository: c.AccountRepository,
	}
}

func (m managerService) CreatePromoCode(ctx context.Context, promoCode *mongo_models.PromoCode) error {
	err := m.ManagerRepository.CreatePromoCode(ctx, promoCode)
	if err != nil {
		return err
	}

	return nil
}

func (m managerService) ViewPromoCodes(ctx context.Context) ([]*slack_bot.PromoCodesView, error) {
	data, err := m.ManagerRepository.ViewPromoCodes(ctx)
	if err != nil {
		return []*slack_bot.PromoCodesView{}, err
	}

	return data, nil
}

func (m managerService) AddHoursToUserSubscription(ctx context.Context, telegramUserId int64, hours int) error {
	err := m.AccountRepository.ChangeUserSubscription(ctx, telegramUserId, "premium")
	if err != nil {
		return err
	}

	err = m.AccountRepository.AddTimeToUserSubscription(ctx, telegramUserId, time.Duration(hours*1000000000*60*60))
	if err != nil {
		return err
	}

	return nil
}

func (m managerService) ViewOnlineUsersStats(ctx context.Context) ([]*slack_bot.OnlineUserStats, error) {
	onlineUsersStats, err := m.ManagerRepository.ViewOnlineUsersStats(ctx)
	if err != nil {
		return nil, err
	}

	return onlineUsersStats, err
}

func (m managerService) ViewCountRequests(ctx context.Context) ([]*slack_bot.LogRequestStats, error) {
	logRequestStats, err := m.ManagerRepository.ViewCountRequests(ctx)
	if err != nil {
		return nil, err
	}

	return logRequestStats, err
}

func (m managerService) ViewUserActions(ctx context.Context, telegramUserId int64) ([]*mongo_models.LogRequest, error) {
	userLogRequests, err := m.ManagerRepository.ViewUserActions(ctx, telegramUserId)
	if err != nil {
		return nil, err
	}

	return userLogRequests, nil
}

func (m managerService) ViewUsersInfo(ctx context.Context) ([]*mongo_models.User, error) {
	usersInfo, err := m.ManagerRepository.ViewUsersInfo(ctx)
	if err != nil {
		return nil, err
	}

	return usersInfo, nil
}

func (m managerService) ViewUserInfo(ctx context.Context, telegramUserId int64) (*slack_bot.UserInfo, error) {
	userInfo, err := m.ManagerRepository.ViewUserInfo(ctx, telegramUserId)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (m managerService) ViewActiveObservers(ctx context.Context) ([]*mongo_models.Observer, error) {
	activeObservers, err := m.ManagerRepository.ViewActiveObservers(ctx)
	if err != nil {
		return nil, err
	}

	return activeObservers, nil
}

func (m managerService) SendMessage(ctx context.Context, filter string, message string) error {
	if filter == "all" {
		users, err := m.AccountRepository.GetUsersByFilter(ctx, "all")
		if err != nil {
			return err
		}

		for _, user := range users {
			err := telegram_api.SendMessageToUser(user.TelegramUserID, message)
			if err != nil {
				return err
			}
		}
	} else if filter == "free" {
		users, err := m.AccountRepository.GetUsersByFilter(ctx, "free")
		if err != nil {
			return err
		}

		for _, user := range users {
			err := telegram_api.SendMessageToUser(user.TelegramUserID, message)
			if err != nil {
				return err
			}
		}
	} else if filter == "premium" {
		users, err := m.AccountRepository.GetUsersByFilter(ctx, "premium")
		if err != nil {
			return err
		}

		for _, user := range users {
			err := telegram_api.SendMessageToUser(user.TelegramUserID, message)
			if err != nil {
				return err
			}
		}
	} else {
		userId, err := strconv.Atoi(filter)
		if err != nil {
			return err
		}
		err = telegram_api.SendMessageToUser(int64(userId), message)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m managerService) AddUserToAdminGroup(ctx context.Context, telegramUserId int64) error {
	err := m.ManagerRepository.AddUserToAdminGroup(ctx, telegramUserId)
	if err != nil {
		return err
	}

	return nil
}
