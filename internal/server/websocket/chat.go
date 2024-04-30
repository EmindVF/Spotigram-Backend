package websocket

import (
	"encoding/json"

	"spotigram/internal/service/models"
	"spotigram/internal/service/usecases"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type responceContext struct {
	MessageType  int
	Message      string
	ReceiverUUID string
}

type errorResponce struct {
	Status  string `json:"status"`
	Content string `json:"content"`
}

var (
	connections = make(map[string]*websocket.Conn)
	register    = make(chan *websocket.Conn)
	broadcast   = make(chan responceContext)
	unregister  = make(chan *websocket.Conn)
)

func RunChatHub() {
	for {
		select {
		// Connect
		case connection := <-register:
			if _, ok := connections[connection.Locals("user_uuid").(string)]; ok {
				connection.WriteMessage(websocket.CloseMessage, []byte{})
				connection.Close()
			} else {
				connections[connection.Locals("user_uuid").(string)] = connection
			}

		// Send message
		case ctx := <-broadcast:
			var connection *websocket.Conn
			var ok bool
			if connection, ok = connections[ctx.ReceiverUUID]; !ok {
				continue
			}
			if err := connection.WriteMessage(ctx.MessageType, []byte(ctx.Message)); err != nil {

				unregister <- connection
				connection.WriteMessage(websocket.CloseMessage, []byte{})
				connection.Close()
			}
		// Disconnect
		case connection := <-unregister:
			delete(connections, connection.Locals("user_uuid").(string))
		}
	}
}

func WebsocketChatUpgradeHandler(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func sendError(message string, uuid string) {
	errorMessage, _ := json.Marshal(errorResponce{
		Status:  "fail",
		Content: message,
	})
	broadcast <- responceContext{
		MessageType:  websocket.TextMessage,
		Message:      string(errorMessage),
		ReceiverUUID: uuid,
	}
}

func sendOk(uuid string) {
	errorMessage, _ := json.Marshal(errorResponce{
		Status: "ok",
	})
	broadcast <- responceContext{
		MessageType:  websocket.TextMessage,
		Message:      string(errorMessage),
		ReceiverUUID: uuid,
	}
}

func sendNotification(status, message string, uuid string) {
	errorMessage, _ := json.Marshal(errorResponce{
		Status:  status,
		Content: message,
	})
	broadcast <- responceContext{
		MessageType:  websocket.TextMessage,
		Message:      string(errorMessage),
		ReceiverUUID: uuid,
	}
}

func WebsocketChatLoop(c *websocket.Conn) {
	defer func() {
		unregister <- c
		c.Close()
	}()
	register <- c

	userUUID := c.Locals("user_uuid").(string)

	var (
		mt  int
		msg []byte
		err error
	)

	for {
		mt, msg, err = c.ReadMessage()
		if err != nil {
			return
		}

		if mt == websocket.TextMessage {
			// Unmarshal payload
			var payload models.WebsocketPayload
			err := json.Unmarshal(msg, &payload)
			if err != nil {
				sendError("invalid json", userUUID)
				continue
			}
			if payload.Action == "" || payload.Content == nil {
				sendError("invalid \"action\" or \"content\" field", userUUID)
				continue
			}

			// Action cases
			switch payload.Action {
			case "send-friend-request":
				var input = models.AddFriendRequestInput{
					SenderUUID: c.Locals("user_uuid").(string),
				}
				err := json.Unmarshal(payload.Content, &input)
				if err != nil {
					sendError("invalid \"content\"", userUUID)
					continue
				}
				if input.RecipientUUID == "" {
					sendError("invalid \"id\"", userUUID)
					continue
				}

				err = usecases.AddFriendRequest(input)
				if err != nil {
					sendError(err.Error(), userUUID)
					continue
				}

				// Notify user
				notificationMessage, _ := json.Marshal(models.FriendRequest{
					SenderId:    userUUID,
					RecipientId: input.RecipientUUID,
					IsIgnored:   false,
				})
				sendNotification(
					"friend-request-received",
					string(notificationMessage),
					input.RecipientUUID)

				sendOk(userUUID)

			case "delete-friend-request":
				var input = models.DeleteFriendRequestInput{
					SenderUUID: c.Locals("user_uuid").(string),
				}
				err := json.Unmarshal(payload.Content, &input)
				if err != nil {
					sendError("invalid \"content\"", userUUID)
					continue
				}
				if input.RecipientUUID == "" {
					sendError("invalid \"id\"", userUUID)
					continue
				}

				err = usecases.DeleteFriendRequest(input)
				if err != nil {
					sendError(err.Error(), userUUID)
					continue
				}

				// Notify user
				notificationMessage, _ := json.Marshal(models.FriendRequest{
					SenderId:    userUUID,
					RecipientId: input.RecipientUUID,
					IsIgnored:   false,
				})
				sendNotification(
					"friend-request-deleted",
					string(notificationMessage),
					input.RecipientUUID)

				sendOk(userUUID)

			case "update-friend-request":
				var input = models.UpdateFriendRequestInput{
					RecipientUUID: c.Locals("user_uuid").(string),
				}
				err := json.Unmarshal(payload.Content, &input)
				if err != nil {
					sendError("invalid \"content\"", userUUID)
					continue
				}
				if input.SenderUUID == "" {
					sendError("invalid \"id\"", userUUID)
					continue
				}

				err = usecases.UpdateFriendRequest(input)
				if err != nil {
					sendError(err.Error(), userUUID)
					continue
				}

				// Notify user
				notificationMessage, _ := json.Marshal(models.FriendRequest{
					SenderId:    input.SenderUUID,
					RecipientId: userUUID,
					IsIgnored:   input.IsIgnored,
				})
				sendNotification(
					"friend-request-updated",
					string(notificationMessage),
					input.SenderUUID)

				sendOk(userUUID)

			case "accept-friend-request":
				var input = models.AcceptFriendRequestInput{
					RecipientUUID: c.Locals("user_uuid").(string),
				}
				err := json.Unmarshal(payload.Content, &input)
				if err != nil {
					sendError("invalid \"content\"", userUUID)
					continue
				}
				if input.SenderUUID == "" {
					sendError("invalid \"id\"", userUUID)
					continue
				}

				newFriend, err := usecases.AcceptFriendRequest(input)
				if err != nil {
					sendError(err.Error(), userUUID)
					continue
				}

				// Notify user
				notificationMessage, _ := json.Marshal(newFriend)
				sendNotification(
					"friend-request-accepted",
					string(notificationMessage),
					input.SenderUUID)

				sendOk(userUUID)

			case "delete-friend":
				var input = models.DeleteFriendInput{
					User1UUID: c.Locals("user_uuid").(string),
				}
				err := json.Unmarshal(payload.Content, &input)
				if err != nil {
					sendError("invalid \"content\"", userUUID)
					continue
				}
				if input.User2UUID == "" {
					sendError("invalid \"id\"", userUUID)
					continue
				}

				err = usecases.DeleteFriend(input)
				if err != nil {
					sendError(err.Error(), userUUID)
					continue
				}

				// Notify user
				notificationMessage, _ := json.Marshal(models.Friend{
					Id1:    input.User2UUID,
					Id2:    userUUID,
					ChatId: "",
				})
				sendNotification(
					"friend-deleted",
					string(notificationMessage),
					input.User2UUID)

				sendOk(userUUID)

			case "send-message":
				var input = models.Message{
					UserId: c.Locals("user_uuid").(string),
				}
				err := json.Unmarshal(payload.Content, &input)
				if err != nil {
					sendError("invalid \"content\"", userUUID)
					continue
				}
				if input.ChatId == "" {
					sendError("invalid chat \"id\"", userUUID)
					continue
				}

				recipientUUID, err := usecases.SendMessage(input)
				if err != nil {
					sendError(err.Error(), userUUID)
					continue
				}

				// Notify user
				notificationMessage, _ := json.Marshal(input)
				sendNotification(
					"message-received",
					string(notificationMessage),
					recipientUUID)

				sendOk(userUUID)

			case "delete-message":
				var input = models.DeleteMessageInput{
					UserId: c.Locals("user_uuid").(string),
				}
				err := json.Unmarshal(payload.Content, &input)
				if err != nil {
					sendError("invalid \"content\"", userUUID)
					continue
				}
				if input.ChatId == "" {
					sendError("invalid chat \"id\"", userUUID)
					continue
				}

				recipientUUID, err := usecases.DeleteMessage(input)
				if err != nil {
					sendError(err.Error(), userUUID)
					continue
				}

				// Notify user
				notificationMessage, _ := json.Marshal(input)
				sendNotification(
					"message-deleted",
					string(notificationMessage),
					recipientUUID)

				sendOk(userUUID)

			default:
				sendError("invalid \"action\"", userUUID)
			}
		} else {
			sendError("invalid websocket message type", userUUID)
		}
	}
}
