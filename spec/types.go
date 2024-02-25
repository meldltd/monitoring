package spec

type CheckType string
type NotificationType string
type CheckMethod string

const (
	TLS              CheckType = "tls"
	TCP              CheckType = "tcp"
	HTTP             CheckType = "http"
	HTTPS            CheckType = "https"
	IHTTPS           CheckType = "ihttps"
	SSH              CheckType = "ssh"
	PG               CheckType = "pg"
	CH               CheckType = "ch"
	K8SSPEC          CheckType = "k8s-spec"
	AWSCOST          CheckType = "awscost"
	AWSINSTANCECOUNT CheckType = "awsinstancecount"
	GITHUB           CheckType = "github"

	WEBHOOK NotificationType = "webhook"
	SNS     NotificationType = "sns"
	EMAIL   NotificationType = "email"
	SMS     NotificationType = "SMS"

	CONNECT  CheckMethod = "connect"
	CHANGE   CheckMethod = "change"
	STATUS   CheckMethod = "status"
	CONTAINS CheckMethod = "contains"
	EXPIRES  CheckMethod = "expires"
	QUERY    CheckMethod = "query"
)

const (
	TAG_SUCESS    = "success"
	APPENED_VALUE = "append_value"
)

type CheckTags string

type CheckSpec struct {
	ID                string                  `json:"identifier"`
	Type              CheckType               `json:"type"`
	Method            CheckMethod             `json:"method"`
	CheckParams       *map[string]string      `json:"CheckParams,omitempty"`
	DSN               string                  `json:"DSN"`
	Notifications     []Notification          `json:"notifications"`
	DSNParams         *map[string]interface{} `json:"DSNParams,omitempty"`
	Expect            *string                 `json:"expect"`
	Timeout           int32                   `json:"timeout"`
	Measure           bool                    `json:"measure"`
	FailLimit         int32                   `json:"failLimit"`
	Interval          int32                   `json:"interval"`
	IntervalErrored   int32                   `json:"intervalErrored"`
	SuccessLimit      int32                   `json:"successLimit"`
	NotificationDelay int32                   `json:"notificationDelay"`
	Tags              []CheckTags             `json:"tags"`
}

type CheckFile struct {
	Checks []CheckSpec `json:"checks"`
}

type Notification struct {
	Type   NotificationType   `json:"type"`
	Params *map[string]string `json:"params,omitempty"`
}
