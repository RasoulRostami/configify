package services

import (
	"configify/databases"
	"log"
	"regexp"
	"time"
)

type Render struct {
	Prefix        string
	UpdateProcess Process
	DeleteProcess Process
	ReloadCommand string
	need_reload   bool
	prefix_regex  *regexp.Regexp
	last_update   time.Time
}

func NewRender(
	Prefix string,
	UpdateProcess Process,
	DeleteProcess Process,
	ReloadCommand string,
) *Render {
	regex := regexp.MustCompile(Prefix)
	return &Render{
		Prefix:        Prefix,
		UpdateProcess: UpdateProcess,
		DeleteProcess: DeleteProcess,
		ReloadCommand: ReloadCommand,
		prefix_regex:  regex,
		need_reload:   false,
	}
}

// Get message and decide what to do
func (p *Render) Handel(messasge databases.Message) {
	if !p.IsRelatedMessage(messasge.Key) {
		log.Printf("%s is not related to PoweDNS service", messasge.Key)
	} else {
		if messasge.Type == databases.Set {
			p.UpdateProcess.Execute(messasge)
		} else {
			p.DeleteProcess.Execute(messasge)
		}
		p.reload()
	}
}

// Get message key and return bool which is related message or not
func (p *Render) IsRelatedMessage(key string) bool {
	return p.prefix_regex.MatchString(key)
}

// reload service
func (p *Render) reload() {

}
