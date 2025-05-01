package server

import "fmt"

type userID uint64

type userInfo struct {
	id   userID
	name *string
}

type chatMessage struct {
	user    *userInfo
	count   int
	message *string
}

type chatManager struct {
	snapshot        []*chatMessage
	visibleLength   int
	activeThreshold int
	messageCount    int
	activeUsers     map[*userInfo]int
}

func newChatManager(activeThreshold int, visibleLength int) *chatManager {
	return &chatManager{
		snapshot:        make([]*chatMessage, 0, activeThreshold+1),
		visibleLength:   visibleLength,
		activeThreshold: activeThreshold,
		messageCount:    0,
		activeUsers:     make(map[*userInfo]int),
	}
}

func newChatMessage(user *userInfo, count int, msg *string) *chatMessage {
	return &chatMessage{
		user:    user,
		count:   count,
		message: msg,
	}
}

// Posts a new message to the chat and updates the active user list by activating the poster
// if necessary, and deactivating the last poster if they haven't posted a message ever since.
// This function will result in at most one user being activated and one user being deactivated.
// The activated and the deactivated user will be returned.
func (cm *chatManager) post(u *userInfo, msg *string) (activated *userInfo, deactivated *userInfo) {
	activated = nil
	deactivated = nil

	if _, ok := cm.activeUsers[u]; !ok {
		activated = u
	}

	cm.messageCount++
	cm.activeUsers[u] = cm.messageCount
	cm.snapshot = append(cm.snapshot, newChatMessage(u, cm.messageCount, msg))

	// Update active users by deleting inactive ones
	lastUser := cm.snapshot[0].user
	if last, ok := cm.activeUsers[lastUser]; ok && last < cm.messageCount-cm.activeThreshold {
		fmt.Println(last, cm.messageCount-cm.activeThreshold)
		delete(cm.activeUsers, lastUser)
		cm.snapshot = cm.snapshot[1:]
		deactivated = lastUser
	}

	return activated, deactivated
}

func (cm *chatManager) contains(u *userInfo) bool {
	_, ok := cm.activeUsers[u]
	return ok
}

func (cm *chatManager) visibleMessages() []*chatMessage {
	return cm.snapshot[maxInt(0, len(cm.snapshot)-cm.visibleLength):len(cm.snapshot)]
}
