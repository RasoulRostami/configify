package powerdns

import (
	"configify/databases"
	"configify/services"
	"fmt"
	"log"
	"os"
	"text/template"
)

type UpdateZone struct {
	config_file_name string
	config_file_dir  string
	config_template  string
	next             services.Process
}

func NewUpdateZone(
	config_file_name string,
	config_file_dir string,
	config_template string,
	next services.Process) *UpdateZone {
	return &UpdateZone{
		config_file_name: config_file_name,
		config_template:  config_template,
		config_file_dir:  config_file_dir,
		next:             next,
	}
}

func (z *UpdateZone) Update(message *databases.Message) {

	file_name := fmt.Sprintf(z.config_file_name, message.Key)
	file_path := fmt.Sprintf("%s/%s", z.config_file_dir, file_name)
	// get tempalte
	template, err := template.ParseFiles(z.config_template)
	if err != nil {
		panic(err)
	}
	// Create the output file
	outputFile, err := os.Create(file_path)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()
	// Execute the template with the data and save it to the file
	err = template.Execute(outputFile, message.Value)
	if err != nil {
		log.Fatal(err)
	}
	z.next.Update(message)
}

func (t *UpdateZone) Reverse(message *databases.Message) {
	//t.Next.Reverse(message)
}
