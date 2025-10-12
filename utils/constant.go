package utils

const (
	MaxTotalTags = 5

	// Notification Comment Constant
	CREATE_COMMENT_NOTIFICATION_TYPE    = "COMMENT_NOTIFICATION"
	COMMENT_NOTIFICATION_REFERENCE_TYPE = "THREAD"
)

var CommentNotificationPriority = map[string]string{
	"COMMENT_NOTIFICATION": "medium",
}

var NotificationIdempotencyKey = map[string]string{
	"COMMENT_NOTIFICATION": "notif:comment",
}

var CommentNotificationTitle = map[string]string{
	"COMMENT_NOTIFICATION_TITLE": `berkomentar pada postinganmu yang berjudul "[TITLE]"`,
}
