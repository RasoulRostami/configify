package powerdns

import (
	"configify/databases"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"
)

type Bind struct {
	configFileName string
	configFileDir  string
	configFile     string
	template       string
}

func NewBind(configFileDir string, configFileName string) *Bind {
	template := "zone \"{{.hostname}}\" {\n    type master;\n    file \"{{.db_path}}\"; \n};\n"
	configFile := fmt.Sprintf("%s/named.conf.local", configFileDir)
	return &Bind{
		configFileName: configFileName,
		configFileDir:  configFileDir,
		configFile:     configFile,
		template:       template,
	}
}

// Adding zone config file to zone list
func (b *Bind) Update(message *databases.Message) bool {
	data := make(map[string]string)
	data["hostname"] = message.Value["hostname"].(string)
	data["db_path"] = configFileName(b.configFileDir, b.configFileName, message.Key)
	// checking domain already is exists
	content, err := ioutil.ReadFile(b.configFile)
	if err != nil {
		log.Println("ERROR PowerDNS Bind Process: can not opne file.")
	}
	pattern := fmt.Sprintf(`zone "%s"`, data["hostname"])
	if strings.Contains(string(content), pattern) {
		return true
	} else {
		// Adding new config to file
		tmpl, err := template.New("bind").Parse(b.template)
		if err != nil {
			log.Printf("ERROR PowerDNS Bind Proccess can not parse template. %s \n", err)
			return false
		}

		outputFile, err := os.OpenFile(b.configFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("ERROR PowerDNS Bind Proccess can not open config file. %s \n", err)
			return false
		}
		defer outputFile.Close()

		// Execute the template with the data and append it to the file
		err = tmpl.Execute(outputFile, data)
		if err != nil {
			log.Printf("ERROR PowerDNS Bind Proccess can not execute template. %s \n", err)
			return false
		}
		return true
	}
}

// Removing zone config file from zone list
func (b *Bind) Reverse(message *databases.Message) bool {
	content, err := ioutil.ReadFile(b.configFile)
	if err != nil {
		log.Printf("ERROR PowerDNS Bind Process can not read file %s \n", err)
		return false
	}

	zoneFile := configFileName(b.configFileDir, b.configFileName, message.Key)
	pattern := fmt.Sprintf(`zone\s+"([^"]+)"\s+{\s+type\s+master;\s+file\s+"%s";\s+};`, zoneFile)

	replacement := ""
	regex := regexp.MustCompile(pattern)

	// Perform the replacement, remove config from list
	newContent := regex.ReplaceAllString(string(content), replacement)

	// Write the updated content back to the file
	err = ioutil.WriteFile(b.configFile, []byte(newContent), os.ModePerm)
	if err != nil {
		log.Printf("ERROR PowerDNS Bind Process can write file %s \n", err)
		return false
	}
	return true
}
