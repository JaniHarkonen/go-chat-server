package server

type userID uint64

type userInfo struct {
	id   userID
	name *string
}

type chatMessage struct {
	user    *userInfo
	count   uint64
	message *string
}

type chatManager struct {
	snapshot        []*chatMessage
	activeThreshold uint64
	messageCount    uint64
	activeUsers     map[*userInfo]uint64
}

func newChatManager(activeThreshold uint64) *chatManager {
	return &chatManager{
		snapshot:        make([]*chatMessage, 0, activeThreshold+1),
		activeThreshold: activeThreshold,
		messageCount:    0,
		activeUsers:     make(map[*userInfo]uint64),
	}
}

func newChatMessage(user *userInfo, count uint64, msg *string) *chatMessage {
	return &chatMessage{
		user:    user,
		count:   count,
		message: msg,
	}
}

func (cm *chatManager) post(u *userInfo, msg *string) {
	// Append thE message to chat
	cm.messageCount++
	cm.activeUsers[u] = cm.messageCount
	cm.snapshot = append(cm.snapshot, newChatMessage(u, cm.messageCount, msg))

	// Update active users by deleting inactive ones
	lastUser := cm.snapshot[0].user
	if last, ok := cm.activeUsers[lastUser]; ok && last < cm.messageCount-cm.activeThreshold {
		delete(cm.activeUsers, lastUser)
		cm.snapshot = cm.snapshot[1:]
	}
}

func (um *chatManager) contains(u *userInfo) bool {
	_, ok := um.activeUsers[u]
	return ok
}
