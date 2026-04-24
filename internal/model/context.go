package model

type ContextKey string

const (
	ContextKeyUserID ContextKey = "userID"
	ContextKeyEmail  ContextKey = "email"
	ContextKeyRole   ContextKey = "role"
)
