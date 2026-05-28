package usecase

import (
	"context"
	"strings"

	"samrs-backend/pkg/util"
)

type NotificationProvider interface {
	Channel() string
	Send(ctx context.Context, req NotificationRequest) error
}

type NotificationUsecase interface {
	Send(ctx context.Context, req NotificationRequest) error
}

type NotificationRequest struct {
	Channel  string                 `json:"channel"`
	To       string                 `json:"to"`
	Subject  string                 `json:"subject,omitempty"`
	Message  string                 `json:"message,omitempty"`
	Template string                 `json:"template,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`
}

type notificationUsecase struct {
	providers map[string]NotificationProvider
}

func NewNotificationUsecase(providers []NotificationProvider) NotificationUsecase {
	registry := make(map[string]NotificationProvider)
	for _, provider := range providers {
		registry[strings.ToLower(provider.Channel())] = provider
	}
	return &notificationUsecase{providers: registry}
}

func (u *notificationUsecase) Send(ctx context.Context, req NotificationRequest) error {
	channel := strings.ToLower(strings.TrimSpace(req.Channel))
	if channel == "" {
		return util.ErrValidation("channel wajib diisi")
	}

	provider, ok := u.providers[channel]
	if !ok {
		return util.ErrValidation("channel tidak didukung")
	}

	return provider.Send(ctx, req)
}

type noopNotificationProvider struct {
	channel string
}

func NewNoopNotificationProvider(channel string) NotificationProvider {
	return &noopNotificationProvider{channel: channel}
}

func (p *noopNotificationProvider) Channel() string {
	return p.channel
}

func (p *noopNotificationProvider) Send(ctx context.Context, req NotificationRequest) error {
	return nil
}
