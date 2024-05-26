package usecases

import (
	"spotigram/internal/customerrors"
	"spotigram/internal/service/abstractions"
	"spotigram/internal/service/models"
	"spotigram/internal/utility"
)

// A use case to get user's sent friend request.
// Expects access token deserialization beforehand.
// Validates the passed uuid.
// May return ErrInvalidInput, ErrInternal, ErrNotFound on failure.
func GetFriendRequestsSent(gfrsi models.GetFriendRequestsSentInput) ([]models.FriendRequest, error) {
	if check := utility.IsValidUUID(gfrsi.UserUUID); !check {
		return nil, &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}

	friendsRequests, err :=
		abstractions.FriendRequestRepositoryInstance.
			GetFriendRequestsSent(gfrsi.UserUUID, gfrsi.Offset)
	if err != nil {
		return nil, err
	}

	return friendsRequests, nil
}

// A use case to get user's received friend request.
// Expects access token deserialization beforehand.
// Validates the passed uuid.
// May return ErrInvalidInput, ErrInternal, ErrNotFound on failure.
func GetFriendRequestsReceived(gfrri models.GetFriendRequestsReceivedInput) ([]models.FriendRequest, error) {
	if check := utility.IsValidUUID(gfrri.UserUUID); !check {
		return nil, &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}

	friendsRequests, err :=
		abstractions.FriendRequestRepositoryInstance.
			GetFriendRequestsReceived(gfrri.UserUUID, gfrri.Offset)
	if err != nil {
		return nil, err
	}

	return friendsRequests, nil
}

// A use case to add a friend request.
// Expects access token deserialization beforehand.
// Validates the passed uuid.
// May return ErrInvalidInput, ErrInternal on failure.
func AddFriendRequest(afri models.AddFriendRequestInput) error {
	if check := utility.IsValidUUID(afri.SenderUUID); !check {
		return &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}
	if check := utility.IsValidUUID(afri.RecipientUUID); !check {
		return &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}

	if afri.SenderUUID == afri.RecipientUUID {
		return &customerrors.ErrInvalidInput{
			Message: "cannot befriend yourself"}
	}

	check, err := abstractions.FriendRepositoryInstance.
		DoesFriendExist(afri.SenderUUID, afri.RecipientUUID)
	if err != nil {
		return err
	}
	if check {
		return &customerrors.ErrInvalidInput{
			Message: "this friendship already exists"}
	}

	check, err = abstractions.FriendRequestRepositoryInstance.
		DoesFriendRequestExist(afri.SenderUUID, afri.RecipientUUID)
	if err != nil {
		return err
	}
	if check {
		return &customerrors.ErrInvalidInput{
			Message: "this friend request already exists"}
	}

	check, err = abstractions.FriendRequestRepositoryInstance.
		DoesFriendRequestExist(afri.RecipientUUID, afri.SenderUUID)
	if err != nil {
		return err
	}
	if check {
		return &customerrors.ErrInvalidInput{
			Message: "this user already wants to be friends"}
	}

	err = abstractions.FriendRequestRepositoryInstance.
		AddFriendRequest(models.FriendRequest{
			SenderId:    afri.SenderUUID,
			RecipientId: afri.RecipientUUID,
			IsIgnored:   false,
		})

	return err
}

// A use case to update a friend request's ignore status.
// Expects access token deserialization beforehand.
// Validates the passed uuid.
// May return ErrInvalidInput, ErrInternal, ErrNotFound on failure.
func UpdateFriendRequest(ufri models.UpdateFriendRequestInput) error {
	if check := utility.IsValidUUID(ufri.SenderUUID); !check {
		return &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}
	if check := utility.IsValidUUID(ufri.RecipientUUID); !check {
		return &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}

	if ufri.SenderUUID == ufri.RecipientUUID {
		return &customerrors.ErrInvalidInput{
			Message: "cannot befriend yourself"}
	}

	err := abstractions.FriendRequestRepositoryInstance.
		UpdateIsIgnored(ufri.SenderUUID, ufri.RecipientUUID, ufri.IsIgnored)
	return err
}

// A use case to delete a friend request.
// Expects access token deserialization beforehand.
// Validates the passed uuid.
// May return ErrInvalidInput, ErrInternal, ErrNotFound on failure.
func DeleteFriendRequest(dfri models.DeleteFriendRequestInput) error {
	if check := utility.IsValidUUID(dfri.SenderUUID); !check {
		return &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}
	if check := utility.IsValidUUID(dfri.RecipientUUID); !check {
		return &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}

	err := abstractions.FriendRequestRepositoryInstance.
		DeleteFriendRequest(dfri.SenderUUID, dfri.RecipientUUID)
	return err
}

// A use case to accept a friend request.
// Expects access token deserialization beforehand.
// Validates the passed uuid.
// Returns a model of a generated friend.
// May return ErrInvalidInput, ErrInternal, ErrNotFound on failure.
func AcceptFriendRequest(afri models.AcceptFriendRequestInput) (*models.Friend, error) {
	if check := utility.IsValidUUID(afri.SenderUUID); !check {
		return nil, &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}
	if check := utility.IsValidUUID(afri.RecipientUUID); !check {
		return nil, &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}

	if afri.SenderUUID == afri.RecipientUUID {
		return nil, &customerrors.ErrInvalidInput{
			Message: "cannot befriend yourself"}
	}

	err := abstractions.FriendRequestRepositoryInstance.
		DeleteFriendRequest(afri.SenderUUID, afri.RecipientUUID)
	if err != nil {
		return nil, err
	}

	newFriend := models.Friend{
		Id1:    afri.SenderUUID,
		Id2:    afri.RecipientUUID,
		ChatId: utility.GenerateUUID(),
	}

	err = abstractions.FriendRepositoryInstance.
		AddFriend(newFriend)
	if err != nil {
		return nil, err
	}

	err = abstractions.ReadTimeRepositoryInstance.
		AddReadTime(&models.ReadTime{
			UserId: afri.SenderUUID,
			ChatId: newFriend.ChatId,
			TimeId: 0,
		})
	if err != nil {
		return nil, err
	}

	err = abstractions.ReadTimeRepositoryInstance.
		AddReadTime(&models.ReadTime{
			UserId: afri.RecipientUUID,
			ChatId: newFriend.ChatId,
			TimeId: 0,
		})
	if err != nil {
		return nil, err
	}

	return &newFriend, err
}
