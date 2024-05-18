package websocket

import (
	"encoding/json"
	"time"

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

type echoResponce struct {
	Status  string `json:"status"`
	Content string `json:"content"`
}

type notificationResponce struct {
	Status  string          `json:"status"`
	Content json.RawMessage `json:"content"`
}

var (
	connections = make(map[string]*websocket.Conn)
	register    = make(chan *websocket.Conn)
	broadcast   = make(chan responceContext)
	unregister  = make(chan *websocket.Conn)
)

func RunChatHub() {
	ticker := time.NewTicker(3 * time.Minute)
	for {
		select {
		// Connect
		case connection := <-register:
			if _, ok := connections[connection.Locals("user_uuid").(string)]; ok {
				connections[connection.Locals("user_uuid").(string)].
					WriteMessage(websocket.CloseMessage, []byte{})
			}
			connections[connection.Locals("user_uuid").(string)] = connection

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

		// Send pings
		case <-ticker.C:
			for _, c := range connections {
				c.WriteMessage(websocket.TextMessage, []byte("ping"))
			}
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
	errorMessage, _ := json.Marshal(echoResponce{
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
	errorMessage, _ := json.Marshal(echoResponce{
		Status: "ok",
	})
	broadcast <- responceContext{
		MessageType:  websocket.TextMessage,
		Message:      string(errorMessage),
		ReceiverUUID: uuid,
	}
}

func sendNotification(status string, message json.RawMessage, uuid string) {
	errorMessage, _ := json.Marshal(notificationResponce{
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
	register <- c
	userUUID := c.Locals("user_uuid").(string)

	c.SetReadLimit(256 * 1024)

	var (
		mt  int
		msg []byte
		err error
	)

	for {
		c.SetReadDeadline(time.Now().Add(5 * time.Minute))
		mt, msg, err = c.ReadMessage()
		if err != nil {
			break
		}

		if mt == websocket.TextMessage && len(msg) < 10 {
			if string(msg) == "pong" {
				continue
			}
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
				var input = models.AddFriendRequestInput{}
				err := json.Unmarshal(payload.Content, &input)
				if err != nil {
					sendError("invalid \"content\"", userUUID)
					continue
				}
				if input.RecipientUUID == "" {
					sendError("invalid \"id\"", userUUID)
					continue
				}

				input.SenderUUID = c.Locals("user_uuid").(string)

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
					notificationMessage,
					input.RecipientUUID)

				sendOk(userUUID)

			case "delete-friend-request":
				var input = models.DeleteFriendRequestInput{}
				err := json.Unmarshal(payload.Content, &input)
				if err != nil {
					sendError("invalid \"const\"", userUUID)
					continue
				}
				if input.RecipientUUID == "" {
					sendError("invalid \"id\"", userUUID)
					continue
				}

				input.SenderUUID = c.Locals("user_uuid").(string)

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
					notificationMessage,
					input.RecipientUUID)

				sendOk(userUUID)

			case "update-friend-request":
				var input = models.UpdateFriendRequestInput{}
				err := json.Unmarshal(payload.Content, &input)
				if err != nil {
					sendError("invalid \"content\"", userUUID)
					continue
				}
				if input.SenderUUID == "" {
					sendError("invalid \"id\"", userUUID)
					continue
				}

				input.RecipientUUID = c.Locals("user_uuid").(string)

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
					notificationMessage,
					input.SenderUUID)

				sendOk(userUUID)

			case "accept-friend-request":
				var input = models.AcceptFriendRequestInput{}
				err := json.Unmarshal(payload.Content, &input)
				if err != nil {
					sendError("invalid \"content\"", userUUID)
					continue
				}
				if input.SenderUUID == "" {
					sendError("invalid \"id\"", userUUID)
					continue
				}
				input.RecipientUUID = c.Locals("user_uuid").(string)

				newFriend, err := usecases.AcceptFriendRequest(input)
				if err != nil {
					sendError(err.Error(), userUUID)
					continue
				}

				// Notify user
				notificationMessage, _ := json.Marshal(newFriend)
				sendNotification(
					"friend-request-accepted",
					notificationMessage,
					input.SenderUUID)

				sendOk(userUUID)

				sendNotification(
					"friend-added",
					notificationMessage,
					userUUID)

			case "delete-friend":
				var input = models.DeleteFriendInput{}
				err := json.Unmarshal(payload.Content, &input)
				if err != nil {
					sendError("invalid \"content\"", userUUID)
					continue
				}
				if input.User2UUID == "" {
					sendError("invalid \"id\"", userUUID)
					continue
				}

				input.User1UUID = c.Locals("user_uuid").(string)

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
					notificationMessage,
					input.User2UUID)

				sendOk(userUUID)

			case "send-message":
				var input = models.Message{}
				err := json.Unmarshal(payload.Content, &input)
				input.UserId = c.Locals("user_uuid").(string)
				if err != nil {
					sendError("invalid \"content\"", userUUID)
					continue
				}

				if input.ChatId == "" {
					sendError("invalid chat \"id\"", userUUID)
					continue
				}

				recipientUUID, message, err := usecases.SendMessage(input)
				if err != nil {
					sendError(err.Error(), userUUID)
					continue
				}

				// Notify user
				notificationMessage, _ := json.Marshal(message)
				sendNotification(
					"message-received",
					notificationMessage,
					recipientUUID)

				sendOk(userUUID)

				sendNotification(
					"message-received",
					notificationMessage,
					userUUID)

			case "delete-message":
				var input = models.DeleteMessageInput{}
				err := json.Unmarshal(payload.Content, &input)
				if err != nil {
					sendError("invalid \"content\"", userUUID)
					continue
				}
				if input.ChatId == "" {
					sendError("invalid chat \"id\"", userUUID)
					continue
				}
				input.UserId = c.Locals("user_uuid").(string)
				recipientUUID, err := usecases.DeleteMessage(input)
				if err != nil {
					sendError(err.Error(), userUUID)
					continue
				}

				// Notify user
				notificationMessage, _ := json.Marshal(input)
				sendNotification(
					"message-deleted",
					notificationMessage,
					recipientUUID)

				sendOk(userUUID)

			default:
				sendError("invalid \"action\"", userUUID)
			}
		} else {
			sendError("invalid websocket message type", userUUID)
		}
	}

	unregister <- c
	c.Close()
}
