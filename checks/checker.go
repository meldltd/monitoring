package checks

import (
	"fmt"
	"log"
	"monitoring/notification"
	"monitoring/spec"
	"slices"
	"time"
)

type CheckHandler struct {
	Spec *spec.CheckSpec

	CurrentFailCount    int32
	CurrentSuccessCount int32
	IsDown              bool
	FirstFailure        time.Time
	LastCheck           time.Time
	LastNotification    time.Time
	LastExpect          string
	WaitChan            chan time.Time
}

func NewHandler(spec spec.CheckSpec) CheckHandler {
	return CheckHandler{Spec: &spec, IsDown: false}
}

func (c *CheckHandler) Run() {
	log.Printf("[-] [%s] Starting handler\n", c.Spec.ID)
	for {
		waitFor := time.Duration(c.Spec.Interval)

		if c.IsDown && c.Spec.IntervalErrored != 0 {
			waitFor = time.Duration(c.Spec.IntervalErrored)
		}

		// <-c.WaitChan
		<-time.After(waitFor * time.Second)

		switch c.Spec.Type {
		case spec.HTTPS:
			go runCheck(c.CheckHTTPS, c)
		case spec.HTTP:
			go runCheck(c.CheckHTTP, c)
		case spec.IHTTPS:
			go runCheck(c.CheckIHTTPS, c)
		case spec.PG:
			go runCheck(c.CheckPG, c)
		case spec.CH:
			log.Fatalln("Clickhouse Not implemented yet.")
		case spec.TLS:
			go runCheck(c.CheckTLS, c)
		case spec.SSH:
			go runCheck(c.CheckSSH, c)
		case spec.AWSCOST:
			go runCheck(c.CheckAWSCosts, c)
		case spec.K8SSPEC:
			go runCheck(c.CheckK8S, c)
		case spec.AWSINSTANCECOUNT:
			go runCheck(c.CheckAWSInstanceCount, c)
		case spec.GITHUB:
			go runCheck(c.CheckGitHub, c)
		}

	}
}

func (c *CheckHandler) CheckStarted() {
	c.LastCheck = time.Now()
}

func (c *CheckHandler) ResetCounters() {
	c.FirstFailure = time.Time{}
	c.CurrentFailCount = 0
	c.CurrentSuccessCount = 0
	c.IsDown = false
}

func (c *CheckHandler) IncrementFailures() {
	c.IsDown = true
	if c.FirstFailure.IsZero() {
		c.FirstFailure = time.Now()
	}

	c.CurrentFailCount++
}

func (c *CheckHandler) IncrementSuccess() {
	c.CurrentSuccessCount++
}

func (c *CheckHandler) SetExpect(value string) {
	c.LastExpect = value
}

func runCheck(checkFunc func(*spec.CheckSpec) (*map[string]string, error), handler *CheckHandler) {
	start := time.Now()
	log.Printf("[-] [%s] Starting check\n", handler.Spec.ID)
	handler.CheckStarted()
	params, err := checkFunc(handler.Spec)
	if nil != err {
		log.Printf("[X] [%s] [%dms] Check failed: %s\n", handler.Spec.ID, time.Now().Sub(start).Milliseconds(), err.Error())
		handler.IncrementFailures()
		if handler.CurrentFailCount >= handler.Spec.FailLimit {
			notification.SendNotification(err, handler.Spec)
		}
		return
	}

	handler.IncrementSuccess()
	if handler.CurrentSuccessCount >= handler.Spec.SuccessLimit {
		durationLen := time.Now().Sub(handler.FirstFailure)

		if handler.IsDown {
			handler.IsDown = false
			msg := fmt.Sprintf("Alarm cleared. Duration %d ms", durationLen.Milliseconds())
			go notification.SendNotificationResolved(msg, handler.Spec)
		}
		handler.ResetCounters()
	}

	if slices.Contains(handler.Spec.Tags, spec.TAG_SUCESS) {
		log.Println("SUCESS")
		go notification.SendSuccessNotification(params, handler.Spec)
	}
	log.Printf("[+] [%s] [%dms] Check passed\n", handler.Spec.ID, time.Now().Sub(start).Milliseconds())
}
