package controller

import "menta-backend/db"

type ChatController struct {
}

type SimpleMessageAuthor struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type SimpleMessage struct {
	ID           string              `json:"id"`
	CreationDate db.DateTime         `json:"creationDate"`
	EditionDate  db.DateTime         `json:"editionDate"`
	Author       SimpleMessageAuthor `json:"author"`
	Content      string              `json:"content"`
}

func (c *ChatController) GetMessages(group string, limit int) ([]SimpleMessage, *ErrorResponse) {
	res, err := db.DB.ChatMessage.FindMany(
		db.ChatMessage.Room.Where(
			db.ChatRoom.ID.Equals(group),
		),
	).Take(limit).With(
		db.ChatMessage.Author.Fetch(),
	).Exec(ctx)
	if err != nil {
		return nil, &ErrorResponse{
			Code:    500,
			Message: "db error",
		}
	}

	simple := make([]SimpleMessage, len(res))
	for k, v := range res {
		simple[k] = SimpleMessage{
			ID:           v.ID,
			CreationDate: v.CreatedAt,
			EditionDate:  v.RefreshedAt,
			Author: SimpleMessageAuthor{
				ID:       v.Author().ID,
				Username: v.Author().Username,
			},
			Content: v.Content,
		}
	}

	return simple, nil
}
