package command

type UpdateUser struct {
	ID         string
	IsAdmin    bool
	IsDisabled bool
}

type UpdateNotificationPreferences struct {
	Phone       string
	NotifyEmail bool
	NotifySMS   bool
}
