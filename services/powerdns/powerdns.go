package powerdns

import (
	"configify/databases"
	"configify/services"
	"log"

	"github.com/spf13/cast"
)

type PowerDNS struct {
	Tbr            int
	Name           string
	Reload_command string
	update_process services.Process
}

func NewPowerDNS(config map[string]interface{}) *PowerDNS {
	// creating needed process
	bind_process := NewBind(cast.ToString(config["config_file_dir"]))
	zone_process := NewUpdateZone(
		cast.ToString(config["config_file_name"]),
		cast.ToString(config["config_file_dir"]),
		cast.ToString(config["config_template"]),
		bind_process,
	)
	txt_splite_process := NewTxtRecordSplit(zone_process)
	check_message := NewCheckRelatedmessage(
		cast.ToString(config["prefix"]), txt_splite_process,
	)
	// end creating process
	return &PowerDNS{
		Tbr:            cast.ToInt(config["tbr"]),
		Name:           cast.ToString(config["name"]),
		Reload_command: cast.ToString(config["reload_command"]),
		update_process: check_message,
	}
}

func (p *PowerDNS) Update(message *databases.Message) {
	p.update_process.Update(message)
	log.Printf("DEBUG PowerDNS update: (%s) \n", message.Key)
}

func (P *PowerDNS) Reverse(message *databases.Message) {
}

func (P *PowerDNS) Reload() {
}
