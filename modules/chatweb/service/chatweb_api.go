package service

import (
	"github.com/cool-team-official/cool-admin-go/cool"
)

type ChatwebApiService struct {
	*cool.Service
}

func NewChatwebApiService() *ChatwebApiService {
	return &ChatwebApiService{
		&cool.Service{
			//    Model: model.NewChatwebApi(),
		},
	}
}
