package powerdns

import (
	"configify/databases"
	"configify/services"
	"fmt"
	"log"

	"github.com/spf13/cast"
)

func configFileName(dir string, name string, key string) string {
	file_name := fmt.Sprintf(name, key)
	return fmt.Sprintf("%s/%s", dir, file_name)
}

type PowerDNS struct {
	Tbr            int
	Name           string
	Reload_command string
	process        services.Process
}

func NewPowerDNS(config map[string]interface{}) *PowerDNS {
	// creating needed process
	bind_process := NewBind(cast.ToString(config["config_file_dir"]), cast.ToString(config["config_file_name"]))
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
		process:        check_message,
	}
}

func (p *PowerDNS) Update(message *databases.Message) {
	p.process.Update(message)
	log.Printf("DEBUG PowerDNS update: (%s) \n", message.Key)
}

func (p *PowerDNS) Reverse(message *databases.Message) {
	p.process.Reverse(message)
	log.Printf("DEBUG PowerDNS update: (%s) \n", message.Key)
}

func (p *PowerDNS) Reload() {
}
