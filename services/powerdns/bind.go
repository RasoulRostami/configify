package powerdns

import (
	"configify/databases"
	"fmt"
	"log"
	"os"
	"text/template"
)

type Bind struct {
	config_file_name string
	config_file_dir  string
	template         string
}

func NewBind(config_file_dir string) *Bind {
	template := "\nzone \"{{.hostname}}\" {\n    type master;\n    file \"{{.db_path}}\"; \n};"
	config_file_name := fmt.Sprintf("%s/named.conf.local", config_file_dir)
	return &Bind{
		config_file_name: config_file_name,
		config_file_dir:  config_file_dir,
		template:         template,
	}
}

func (b *Bind) Update(message *databases.Message) {
	data := make(map[string]string)
	data["hostname"] = message.Value["hostname"].(string)
	data["db_path"] = fmt.Sprintf("%s/db.%s", b.config_file_dir, message.Key)
	// template
	tmpl, err := template.New("bind").Parse(b.template)
	if err != nil {
		log.Fatal(" error 1")
		log.Fatal(err)
	}

	// Execute the template with the data
	outputFile, err := os.OpenFile(b.config_file_name, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(" error 2")
		log.Fatal(err)
	}
	defer outputFile.Close()

	// Execute the template with the data and append it to the file
	err = tmpl.Execute(outputFile, data)
	if err != nil {
		log.Fatal(" error 3")
		log.Fatal(err)
	}

}

func (b *Bind) Reverse(message *databases.Message) {
}

//
//class BindRoot(Processor): #
//
//    def __init__(self, **kwargs) -> None:
//        self.template = Template(self.raw_template)
//        super().__init__(**kwargs)
//
//    def handle(self, key: str, value: dict) -> dict:
//        zone = self.get_zone(key, value['hostname'])
//        template_str = self.template.render(zone)
//
//        root_config = self.output_dir + '/named.conf.local'
//        with open(root_config, 'r+') as myfile:
//            if f"""zone "{zone['hostname']}" """ in myfile.read():
//                return value
//
//        with open(root_config, 'a') as myfile:
//            myfile.write(template_str)
//
//        return value
//
//    def reverse(self, key: str) -> None:
//        pattern = fr'''zone ".+" {{\s+type master;\s+file "{self.output_file_name(key)}";\s+}};'''
//        root_config = self.output_dir + '/named.conf.local'
//
//        with open(root_config, 'r') as myfile:
//            contents = re.sub(pattern=pattern, repl="", string=myfile.read())
//
//        FileHelper().save_file(root_config, contents)
//
//    def get_zone(self, key: str, hostname: str):
//        return {
//            'hostname': hostname,
//            'db_path': self.output_file_name(key)
//        }
