package model

type MessageRequest struct {
	ReceiverId string `json:"receiver_id"`
	Content    string `json:"content"`
}

type NewDeviceRequest struct {
	DeviceId string `json:"device_id"`
}

type NotificationRequest struct {
	ReceiverId string `json:"receiver_id"`
	Title      string `json:"title"`
	Body       string `json:"body"`
}

// UserFollowerMilestoneNotificationRequest POST "notifications/user-milestone"
type UserFollowerMilestoneNotificationRequest struct {
	UserId     string `json:"user_id"`
	FollowerId string `json:"follower_id"`
}

// MentionNotificationRequest POST "notifications/mention"
type MentionNotificationRequest struct {
	UserId   string `json:"user_id"`
	TaggerId string `json:"tagger_id"`
}

//// TrendingTwitNotificationRequest POST "notifications/trending-twit"
//type TrendingTwitNotificationRequest struct {
//	UserID string `json:"user_id"`
//	TwitId string `json:"twit_id"`
//}
