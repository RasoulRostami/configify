package powerdns

import (
	"configify/databases"
	"configify/services"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cast"
)

func configFileName(dir string, name string, key string) string {
	file_name := fmt.Sprintf(name, key)
	return fmt.Sprintf("%s/%s", dir, file_name)
}

type PowerDNS struct {
	tbr           int
	Name          string
	need_reload   bool
	last_reload   time.Time
	reloadCommand string
	process       services.Process
}

func NewPowerDNS(config map[string]interface{}) *PowerDNS {
	// creating needed process
	bind_process := NewBind(
		cast.ToString(config["config_file_dir"]),
		cast.ToString(config["config_file_name"]),
	)
	zone_process := NewUpdateZone(
		cast.ToString(config["config_file_name"]),
		cast.ToString(config["config_file_dir"]),
		cast.ToString(config["config_template"]),
		bind_process,
	)
	txt_splite_process := NewTxtRecordSplit(zone_process)
	check_message := NewCheckMessage(
		cast.ToString(config["prefix"]), txt_splite_process,
	)
	// end creating process
	return &PowerDNS{
		tbr:           cast.ToInt(config["tbr"]),
		Name:          cast.ToString(config["name"]),
		reloadCommand: cast.ToString(config["reload_command"]),
		process:       check_message,
		need_reload:   false,
		last_reload:   time.Now().Add(-10 * time.Second),
	}
}

// updating config when new one was created or old one was updated
func (p *PowerDNS) Update(message *databases.Message) {
	result := p.process.Update(message)
	if result {
		p.need_reload = true
		log.Printf("DEBUG PowerDNS was updated: (%s) \n", message.Key)
	}
}

// Reversing config when old one was removed
func (p *PowerDNS) Reverse(message *databases.Message) {
	result := p.process.Reverse(message)
	if result {
		p.need_reload = true
		log.Printf("DEBUG PowerDNS was reversed: (%s) \n", message.Key)
	}
}

// Reload service after modifying
func (p *PowerDNS) Reload() {
	duration := time.Since(p.last_reload)
	if duration.Seconds() > float64(p.tbr) && p.need_reload {
		// generate executable command
		split_command := strings.Split(p.reloadCommand, " ")
		commond := split_command[0]
		args := split_command[1:]
		cmd := exec.Command(commond, args...)
		// Run the command and wait for it to complete
		err := cmd.Run()
		if err != nil {
			log.Printf("ERROR PowerDNS can not run realod command %s \n", err)
		}
		p.last_reload = time.Now()
		log.Println("SUCCESS PowerDNS was realoded.")
	}
}
