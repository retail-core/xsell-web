package handlers

import "net/http"

type NotificationPref struct {
	Key, Label, Desc string
	Enabled          bool
}

func getSettingsData() map[string]any {
	return map[string]any{
		"Notifications": []NotificationPref{
			{"low_stock", "Low Stock Alerts", "Get notified when items run low", true},
			{"sale_completed", "New Sale Completed", "Notify when a sale is processed", true},
			{"debt_reminders", "Debt Reminders", "Reminders for outstanding customer credit", true},
			{"daily_summary", "Daily Sales Summary", "End-of-day summary of your sales", false},
		},
	}
}

func (h *AccountHandler) SettingsPage(w http.ResponseWriter, r *http.Request) {
	h.Tmpl.ExecuteTemplate(w, "settings-page", getSettingsData())
}