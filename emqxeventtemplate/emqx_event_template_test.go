package emqxeventtemplate_test

import (
	"testing"

	"github.com/go-xlan/go-emqx/emqxeventtemplate"
	"github.com/go-xlan/go-emqx/emqxgo"
)

func TestGetEmqxEventTemplate(t *testing.T) {
	t.Log(emqxeventtemplate.GetEmqxEventTemplate(emqxgo.ConnectedEvent{}))
	t.Log(emqxeventtemplate.GetEmqxEventTemplate(emqxgo.DisconnectedEvent{}))
}
