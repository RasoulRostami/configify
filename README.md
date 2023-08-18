# Configify

The Config File Generator App is a versatile tool that automates the process of creating and updating configuration files. It embraces a modular architecture, allowing developers to easily integrate new modules to expand the range of config file series that can be generated.

This application serves as a config file generator, designed with modularity in mind. Developers have the flexibility to introduce new modules, enhancing the capacity to generate a series of configuration files. This utility eliminates the need for manual intervention in updating and creating configuration files. The app efficiently retrieves configuration data from a key-value storage solution like Redis, enabling seamless generation and continuous updates of configuration files. By automating these processes, the app streamlines the configuration management workflow.

## Features

- Modular Design: Integrate new modules effortlessly to support the generation of various configuration file series.
- Automated Updates: Eliminate manual updates by fetching configuration data from a key-value storage, such as Redis.
- Streamlined Workflow: Simplify configuration management with automated file generation and continuous updates.

## Usage

1. Clone the repository.
2. Install the required dependencies.
3. Configure the app to connect to your chosen key-value storage.
4. Add new modules as needed to cater to different configuration requirements.
5. Run the app to initiate the automated config file generation process.

## Getting Started

To start using the Config File Generator App, follow these steps:

1. **Clone the Repository**: Clone this repository to your local machine using the following command:
   ```
   git clone https://github.com/RasoulRostami/configify.git
   ```

2. **Build**:

   ```bash
   cd src
   go build
   ```

3. **Configuration**: Open the `config.py` file and configure the app to connect to your chosen key-value storage (e.g., Redis).
   config file can be YAML, JSON and etc.
   default path: $HOME/.configify.yaml

   ```yaml
   publisher: "redis"  
   decoder: "json"     
   
   redis_config:
     address: "localhost:6379"
     password: ""
     db: 0
   
   services:
     powerdns:
       name: "powerdns"
       prefix: "dns_config_*"
       config_template: "/home/rasoul/configify/templates/pdns.txt"
       config_file_dir: "/home/rasoul/configify/powerdns"
       config_file_name: "db.%s"
       reload_command: "bash /home/rasoul/configify/bash/dev_pdns.sh"
       tbr: 10
   logging:
     error_log_path: "/home/rasoul/error.log"
     success_log_path: "/home/rasoul/success.log"
     log_level: 4
     rotation_log_time: 30
     rotation_log_size: 100
     rotation_log_backup: 5
   ```

   

4. **Run the App**: Execute the app using the following commands:

   ```
   # Bootstrap command - get all keys from key-value store, create files and exit 
   configify bootstrap
   # or
   configify bootstrap --config ./config.py
   
   # Run app to listion redis pub-sub event and update or create files
   configify stream
   
   ```

## Structure

### Database

New only redis support. you can add your own database and use Database interface. Adapter design pattern is implemented.

### Decoder

New only JSON support you can add more like cipher and use Decoder interface. Adapter design pattern is implemented.

### Services

Now it is used to generate powerDNS. You can add more service and use Service interface. each services have some process to run. 
chain of responsibility is implemented.

## Contribution

We welcome contributions to enhance the configify App. If you have ideas, bug fixes, or new module suggestions, feel free to open an issue or submit a pull request.
