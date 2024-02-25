package notification

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"monitoring/spec"
	"net/http"
	"os"
)

func PostWebhook(params *map[string]string, msg string) error {
	if nil == params {
		return errors.New("Params must be defined for webhook")
	}

	url := (*params)["url"]
	text := (*params)["text"]

	formatted := fmt.Sprintf("{\"%s\": \"%s\"}", text, msg)

	_, err := http.Post(url, "application/json", bytes.NewReader([]byte(formatted)))
	if nil != err {
		log.Println(err.Error())
	}

	return nil
}

func SendNotification(error error, check *spec.CheckSpec) {
	for _, notification := range check.Notifications {
		switch notification.Type {
		case spec.WEBHOOK:
			msg := fmt.Sprintf("[%s] type: %s, method: %s, error: %s\nMonitored by: %s",
				check.ID,
				check.Type,
				check.Method,
				error.Error(),
				os.Getenv("MONITOR_ID"),
			)
			go PostWebhook(notification.Params, msg)
		}
	}
}

func SendNotificationResolved(msg string, check *spec.CheckSpec) {
	for _, notification := range check.Notifications {
		switch notification.Type {
		case spec.WEBHOOK:
			text := fmt.Sprintf("[%s] type: %s, method: %s, RESOLVED: %s", check.ID, check.Type, check.Method, msg)
			go PostWebhook(notification.Params, text)
		}
	}

}

func SendSuccessNotification(params *map[string]string, check *spec.CheckSpec) {
	if nil == params {
		return
	}
	msg := ""
	for k, v := range *params {
		msg += fmt.Sprintf("%s: %s\n", k, v)
	}
	for _, notification := range check.Notifications {
		switch notification.Type {
		case spec.WEBHOOK:
			text := fmt.Sprintf("[%s] type: %s, method: %s, Info: %s", check.ID, check.Type, check.Method, msg)
			go PostWebhook(notification.Params, text)
		}
	}

}
