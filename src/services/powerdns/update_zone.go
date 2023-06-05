package powerdns

import (
	"configify/databases"
	"configify/services"
	"log"
	"os"
	"text/template"
)

type UpdateZone struct {
	configFileName string
	configFileDir  string
	configTemplate string
	next           services.Process
}

func NewUpdateZone(
	configFileName string,
	configFileDir string,
	configTemplate string,
	next services.Process) *UpdateZone {
	return &UpdateZone{
		configFileName: configFileName,
		configTemplate: configTemplate,
		configFileDir:  configFileDir,
		next:           next,
	}
}

// write zone config in a file
func (z *UpdateZone) Update(message *databases.Message) bool {
	filePath := configFileName(z.configFileDir, z.configFileName, message.Key)

	template, err := template.ParseFiles(z.configTemplate)
	if err != nil {
		log.Printf("ERROR PowerDNS Update Zone can not parse template %s \n", err)
		return false
	}
	// Create the output file
	outputFile, err := os.Create(filePath)
	if err != nil {
		log.Printf("ERROR PowerDNS Update Zone can not create output file %s \n", err)
		return false
	}
	defer outputFile.Close()
	// Execute the template with the data and save it to the file
	err = template.Execute(outputFile, message.Value)
	if err != nil {
		log.Printf("ERROR PowerDNS Update Zone can not execute template %s \n", err)
		return false
	}
	log.Printf("DEBUG PowerDNS Update Zone %s \n", message.Key)
	return z.next.Update(message)
}

// remove config file
func (z *UpdateZone) Reverse(message *databases.Message) bool {
	filePath := configFileName(z.configFileDir, z.configFileName, message.Key)
	err := os.Remove(filePath)
	if err != nil {
		log.Printf("ERROR PowerDNS Reverse Zone (%s) %s \n", message.Key, err)
		return false
	} else {
		log.Printf("DEBUG powerdns update zone %s \n", message.Key)
	}
	return z.next.Reverse(message)
}
