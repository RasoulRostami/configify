package powerdns

import (
	"configify/databases"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"text/template"
)

type Bind struct {
	config_file_name string
	config_file_dir  string
	config_file      string
	template         string
}

func NewBind(config_file_dir string, config_file_name string) *Bind {
	template := "zone \"{{.hostname}}\" {\n    type master;\n    file \"{{.db_path}}\"; \n};\n"
	config_file := fmt.Sprintf("%s/named.conf.local", config_file_dir)
	return &Bind{
		config_file_name: config_file_name,
		config_file_dir:  config_file_dir,
		config_file:      config_file,
		template:         template,
	}
}

func (b *Bind) Update(message *databases.Message) {
	data := make(map[string]string)
	data["hostname"] = message.Value["hostname"].(string)
	data["db_path"] = configFileName(b.config_file_dir, b.config_file_name, message.Key)
	// template
	tmpl, err := template.New("bind").Parse(b.template)
	if err != nil {
		log.Fatal(err)
	}

	// Execute the template with the data
	outputFile, err := os.OpenFile(b.config_file, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// Execute the template with the data and append it to the file
	err = tmpl.Execute(outputFile, data)
	if err != nil {
		log.Fatal(err)
	}

}

// Remove zone config from the files
func (b *Bind) Reverse(message *databases.Message) {
	// Read the file content
	content, err := ioutil.ReadFile(b.config_file)
	if err != nil {
		log.Printf("ERROR PowerDNS can not open %s", b.config_file)
	}

	// Define the regular expression pattern to search
	zoneFile := configFileName(b.config_file_dir, b.config_file_name, message.Key)
	pattern := fmt.Sprintf(`zone\s+"([^"]+)"\s+{\s+type\s+master;\s+file\s+"%s";\s+};`, zoneFile)

	// Define the replacement value
	replacement := ""

	// Create a regular expression object
	regex := regexp.MustCompile(pattern)

	// Perform the replacement
	newContent := regex.ReplaceAllString(string(content), replacement)

	// Write the updated content back to the file
	err = ioutil.WriteFile(b.config_file, []byte(newContent), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("")
}
