package usecases

import (
	"spotigram/internal/customerrors"
	"spotigram/internal/service/abstractions"
	"spotigram/internal/service/models"
	"spotigram/internal/utility"
)

// A use case to get user's friends.
// Expects access token deserialization beforehand.
// Validates the passed uuid.
// May return ErrInvalidInput, ErrInternal, ErrNotFound on failure.
func GetFriends(gfi models.GetFriendsInput) ([]models.Friend, error) {
	if check := utility.IsValidUUID(gfi.UserUUID); !check {
		return nil, &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}

	friends, err :=
		abstractions.FriendRepositoryInstance.GetFriends(gfi.UserUUID, gfi.Offset)
	if err != nil {
		return nil, err
	}

	return friends, nil
}

// A use case to delete a user's friend.
// Expects access token deserialization beforehand.
// Validates the passed uuid.
// May return ErrInvalidInput, ErrInternal, ErrNotFound on failure.
func DeleteFriend(dfi models.DeleteFriendInput) error {
	if check := utility.IsValidUUID(dfi.User1UUID); !check {
		return &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}
	if check := utility.IsValidUUID(dfi.User2UUID); !check {
		return &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}

	chatId, err := abstractions.FriendRepositoryInstance.
		GetChatIdByFriend(dfi.User1UUID, dfi.User2UUID)
	if err != nil {
		return err
	}

	err = abstractions.ChatRepositoryInstance.
		DeleteChat(chatId)
	if err != nil {
		return err
	}

	err = abstractions.FriendRepositoryInstance.
		DeleteFriend(dfi.User1UUID, dfi.User2UUID)

	return err
}
