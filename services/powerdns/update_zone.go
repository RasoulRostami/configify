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
	file_path := configFileName(z.config_file_dir, z.config_file_name, message.Key)
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
	log.Printf("DEBUG PowerDNS Update Zone %s", message.Key)
	z.next.Update(message)
}

func (z *UpdateZone) Reverse(message *databases.Message) {
	file_path := configFileName(z.config_file_dir, z.config_file_name, message.Key)
	err := os.Remove(file_path)
	if err != nil {
		fmt.Printf("ERROR PowerDNS Update Zone (%s) %s \n", message.Key, err)
	} else {
		log.Printf("DEBUG powerdns update zone %s \n", message.Key)
	}
	z.next.Reverse(message)
}
