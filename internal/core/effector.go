package core

// Tool represents a polymorphic effector that the agent can execute
type Tool interface {
	ID() string
	Execute(target string, kb *KnowledgeBase) string
}
