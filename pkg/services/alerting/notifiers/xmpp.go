package notifiers

import (
	/*"github.com/grafana/grafana/pkg/bus"*/
	"github.com/grafana/grafana/pkg/log"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/alerting"
	"github.com/mattn/go-xmpp"
	"fmt"
)

func init() {
	alerting.RegisterNotifier("xmpp", NewXmppNotifier)
}

func NewXmppNotifier(model *m.AlertNotification) (alerting.Notifier, error) {
	user := model.Settings.Get("user").MustString()
	if user == "" {
		return nil, alerting.ValidationError{Reason: "Could not find user property in settings"}
	}
	password := model.Settings.Get("password").MustString()
	if password == "" {
		return nil, alerting.ValidationError{Reason: "Could not find password property in settings"}
	}
	server := model.Settings.Get("server").MustString()
	if server == "" {
		return nil, alerting.ValidationError{Reason: "Could not find server property in settings"}
	}
	room := model.Settings.Get("room").MustString()
	if room == "" {
		return nil, alerting.ValidationError{Reason: "Could not find room property in settings"}
	}

	logger := log.New("alerting.notifier.xmpp")

	opt := xmpp.Options{
		Host:     server,
		User:     user,
		Password: password,
		Resource: "Grafana",
		NoTLS:    true,
		Debug:    false,
		Session:  false,
	}

	client, err := opt.NewClient()

	if err != nil {
		fmt.Printf("%s", err)
	}


	//client.JoinMUCNoHistory(room, opt.Resource)

	return &XmppNotifier{
		NotifierBase: NewNotifierBase(model.Id, model.IsDefault, model.Name, model.Type, model.Settings),
		Room:         room,
		Client:       client,
		Logger:       logger,
	}, nil
}

type XmppNotifier struct {
	NotifierBase
	Room       string
	Client     *xmpp.Client
	Logger     log.Logger
}

func (this *XmppNotifier) Notify(evalContext *alerting.EvalContext) error {
	this.Logger.Info("Sending xmpp notification")

	this.Client.Send(xmpp.Chat{Remote: this.Room, Type: "chat", Text: evalContext.GetNotificationTitle()})

	return nil
}

