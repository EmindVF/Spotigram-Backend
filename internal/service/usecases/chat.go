package usecases

import (
	"spotigram/internal/customerrors"
	"spotigram/internal/service/abstractions"
	"spotigram/internal/service/models"
	"spotigram/internal/utility"
	"time"
)

// A use case to get messages of a chat
func GetMessages(input models.GetMessagesInput) ([]models.Message, error) {
	if check := utility.IsValidUUID(input.UserId); !check {
		return nil, &customerrors.ErrInvalidInput{
			Message: "invalid user \"uuid\""}
	}
	if check := utility.IsValidUUID(input.ChatId); !check {
		return nil, &customerrors.ErrInvalidInput{
			Message: "invalid chat \"uuid\""}
	}

	friend, err := abstractions.FriendRepositoryInstance.GetFriendByChatId(input.ChatId)
	if err != nil {
		return nil, err
	}

	if input.UserId != friend.Id1 && input.UserId != friend.Id2 {
		return nil, &customerrors.ErrInvalidInput{
			Message: "user is not a member of this chat"}
	}

	if input.TimeId == 0 {
		input.TimeId = 9223372036854775807
	}

	messages, err :=
		abstractions.ChatRepositoryInstance.GetMessages(input.ChatId, input.TimeId)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

// A use case to get unread messages of a chat
func GetUnreadMessages(input models.GetUnreadMessagesInput) ([]models.Message, error) {
	if check := utility.IsValidUUID(input.UserId); !check {
		return nil, &customerrors.ErrInvalidInput{
			Message: "invalid user \"uuid\""}
	}
	if check := utility.IsValidUUID(input.ChatId); !check {
		return nil, &customerrors.ErrInvalidInput{
			Message: "invalid chat \"uuid\""}
	}

	friend, err := abstractions.FriendRepositoryInstance.GetFriendByChatId(input.ChatId)
	if err != nil {
		return nil, err
	}

	if input.UserId != friend.Id1 && input.UserId != friend.Id2 {
		return nil, &customerrors.ErrInvalidInput{
			Message: "user is not a member of this chat"}
	}

	readTime, err := abstractions.ReadTimeRepositoryInstance.
		GetReadTime(input.UserId, input.ChatId)
	if err != nil {
		return nil, err
	}

	messages, err :=
		abstractions.ChatRepositoryInstance.GetMessages(input.ChatId, readTime.TimeId)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

// A use case to send a message
// Returns the uuid of the recipient
func SendMessage(smi models.Message) (string, *models.Message, error) {
	if check := utility.IsValidUUID(smi.UserId); !check {
		return "", nil, &customerrors.ErrInvalidInput{
			Message: "invalid user \"uuid\""}
	}
	if check := utility.IsValidUUID(smi.ChatId); !check {
		return "", nil, &customerrors.ErrInvalidInput{
			Message: "invalid chat \"uuid\""}
	}
	if len(smi.Content) > 2048 || len(smi.Content) == 0 {
		return "", nil, &customerrors.ErrInvalidInput{
			Message: "message is too long"}
	}
	friend, err := abstractions.FriendRepositoryInstance.GetFriendByChatId(smi.ChatId)
	if err != nil {
		return "", nil, err
	}

	var uuidRecipient string

	if smi.UserId == friend.Id1 {
		uuidRecipient = friend.Id2
		smi.TimeId = time.Now().UTC().UnixMicro()*10 + 1
	} else if smi.UserId == friend.Id2 {
		uuidRecipient = friend.Id1
		smi.TimeId = time.Now().UTC().UnixMicro()*10 + 2
	} else {
		return "", nil, &customerrors.ErrInvalidInput{
			Message: "user is not a member of this chat"}
	}
	smi.Date = time.Now().UTC()

	err = abstractions.ChatRepositoryInstance.AddMessage(smi)
	if err != nil {
		return "", nil, nil
	}

	return uuidRecipient, &smi, nil
}

// A use case to delete a message
// Returns the uuid of the recipient
func DeleteMessage(input models.DeleteMessageInput) (string, error) {
	if check := utility.IsValidUUID(input.UserId); !check {
		return "", &customerrors.ErrInvalidInput{
			Message: "invalid user \"uuid\""}
	}
	if check := utility.IsValidUUID(input.ChatId); !check {
		return "", &customerrors.ErrInvalidInput{
			Message: "invalid chat \"uuid\""}
	}

	friend, err := abstractions.FriendRepositoryInstance.GetFriendByChatId(input.ChatId)
	if err != nil {
		return "", err
	}

	var uuidRecipient string

	if input.UserId == friend.Id1 {
		uuidRecipient = friend.Id2
	} else if input.UserId == friend.Id2 {
		uuidRecipient = friend.Id1
	} else {
		return "", &customerrors.ErrInvalidInput{
			Message: "user is not a member of this chat"}
	}

	err = abstractions.ChatRepositoryInstance.DeleteMessage(input.ChatId, input.TimeId)
	if err != nil {
		return "", nil
	}

	return uuidRecipient, nil
}

// A use case to update ReadTime.
// Returns updated models.ReadTime.
func UpdateReadTime(input models.UpdateReadTimeInput) (*models.ReadTime, string, error) {
	if check := utility.IsValidUUID(input.UserId); !check {
		return nil, "", &customerrors.ErrInvalidInput{
			Message: "invalid user \"uuid\""}
	}
	if check := utility.IsValidUUID(input.ChatId); !check {
		return nil, "", &customerrors.ErrInvalidInput{
			Message: "invalid chat \"uuid\""}
	}

	friend, err := abstractions.FriendRepositoryInstance.GetFriendByChatId(input.ChatId)
	if err != nil {
		return nil, "", err
	}

	receiverId := ""

	if friend.Id1 == input.UserId {
		receiverId = friend.Id2
	} else if friend.Id2 == input.UserId {
		receiverId = friend.Id1
	} else {
		return nil, "", &customerrors.ErrInvalidInput{
			Message: "you are not in this chat"}
	}

	input.TimeId = input.TimeId / 10
	input.TimeId = input.TimeId*10 + 3
	if input.TimeId > time.Now().UTC().UnixMicro()*10 {
		return nil, "", &customerrors.ErrInvalidInput{
			Message: "invalid timeId",
		}
	}
	err = abstractions.ReadTimeRepositoryInstance.
		UpdateReadTime(input.UserId, input.ChatId, input.TimeId)
	if err != nil {
		return nil, "", err
	}

	return &models.ReadTime{
		UserId: input.UserId,
		ChatId: input.ChatId,
		TimeId: input.TimeId,
	}, receiverId, nil
}

// A use case to get ReadTime
// Returns updated models.ReadTime
func GetReadTime(input models.GetReadTimeInput) (*models.ReadTime, error) {

	if check := utility.IsValidUUID(input.UserId); !check {
		return nil, &customerrors.ErrInvalidInput{
			Message: "invalid user \"uuid\""}
	}
	if check := utility.IsValidUUID(input.ChatId); !check {
		return nil, &customerrors.ErrInvalidInput{
			Message: "invalid chat \"uuid\""}
	}

	readTime, err := abstractions.ReadTimeRepositoryInstance.
		GetReadTime(input.UserId, input.ChatId)
	if err != nil {
		return nil, err
	}

	return readTime, nil
}
