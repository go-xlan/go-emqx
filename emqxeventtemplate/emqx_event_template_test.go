package emqxeventtemplate_test

import (
	"testing"

	"github.com/go-xlan/go-emqx/emqxeventtemplate"
	"github.com/go-xlan/go-emqx/emqxgo"
)

func TestGetEmqxEventFieldSQL(t *testing.T) {
	t.Log(emqxeventtemplate.GetEmqxEventFieldSQL(emqxgo.ConnectedEvent{}, "$events/client/connected"))
	t.Log(emqxeventtemplate.GetEmqxEventTemplate(emqxgo.ConnectedEvent{}))
}

func TestGetEmqxEventTemplate(t *testing.T) {
	t.Log(emqxeventtemplate.GetEmqxEventFieldSQL(emqxgo.DisconnectedEvent{}, "$events/client/disconnected"))
	t.Log(emqxeventtemplate.GetEmqxEventTemplate(emqxgo.DisconnectedEvent{}))
}
