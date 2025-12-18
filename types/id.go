package types

type ID string

const (
	BotUserID            ID = "01J00000000000000000000BOT"
	ManagerServiceUserID ID = "01J0000000000000000MANAGER"
	ChatServiceUserID    ID = "01J0000000000000000000CHAT"
	NotificationUserID   ID = "01J000000000000000000NOTIF"
	AdminUserID          ID = "01J000000000000000000ADMIN"
	StorageUserID        ID = "01J0000000000000000STORAGE"
	StoryUserID          ID = "01J000000000000000000STORY"
	RAGUserID            ID = "01J00000000000000000000RAG"
)

const (
	ManagerServiceUsername = "manager-service"
	ChatServiceUsername    = "chat-service"
	NotificationUsername   = "notification-service"
	AdminServiceUsername   = "admin-service"
	StorageServiceUsername = "storage-service"
	StoryUsername          = "story-service"
	RAGServiceUsername     = "rag-service"
)
