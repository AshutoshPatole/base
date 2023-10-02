package config

import (
	"fmt"
	"log"
	"net"

	c "github.com/AshutoshPatole/ssh-manager/utils"
	"github.com/TwiN/go-color"
	"github.com/spf13/viper"
)

var config c.Config
var groupIndex int = -1
var envIndex int = -1

func SaveServer(hostname, user, group, env, alias string, keyAuth bool) {
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalln(err)
	}

	groupIndex = -1
	envIndex = -1

	// Find the group and environment, or create them if they don't exist
	for idx, grp := range config.Groups {
		if group == grp.Name {
			groupIndex = idx
			for idj, envs := range grp.Environment {
				if envs.Name == env {
					envIndex = idj
					break
				}
			}
		}
	}

	ip, err := IP(hostname)
	if err != nil {
		fmt.Println(color.InYellow("Could not resolve IP address"))
	}

	server := c.Server{
		HostName: hostname,
		IP:       ip,
		KeyAuth:  keyAuth,
		Alias:    alias,
	}

	// If the group doesn't exist, create a new group
	if groupIndex == -1 {
		newGroup := c.Group{
			Name: group,
			User: user,
			Environment: []c.Env{
				{
					Name:    env,
					Servers: []c.Server{server},
				},
			},
		}
		config.Groups = append(config.Groups, newGroup)
		groupIndex = len(config.Groups) - 1
	} else {
		// If the environment doesn't exist, create a new environment
		if envIndex == -1 {
			newEnvironment := c.Env{
				Name:    env,
				Servers: []c.Server{server},
			}
			config.Groups[groupIndex].Environment = append(config.Groups[groupIndex].Environment, newEnvironment)
			envIndex = len(config.Groups[groupIndex].Environment) - 1
		}

		// Check for duplicate server within the environment
		isDup := checkDuplicateServer(server, config.Groups[groupIndex].Environment[envIndex].Servers)

		// Save server information if not a duplicate
		if !isDup {
			config.Groups[groupIndex].Environment[envIndex].Servers = append(config.Groups[groupIndex].Environment[envIndex].Servers, server)
		} else {
			fmt.Println(color.InYellow("Server already exists in group"))
		}
	}

	// Save the information in the config file
	viper.Set("groups", config.Groups)
	if err := viper.WriteConfig(); err != nil {
		log.Fatalln(err)
	}
}

func IP(host string) (string, error) {
	addr := net.ParseIP(host)
	if addr == nil {
		// if addr is nil then its a domain
		ips, err := net.LookupIP(host)
		if err != nil {
			fmt.Println("Error:", err)
		}
		if len(ips) != 0 {
			return ips[0].String(), nil
		} else {
			return "", nil
		}

	} else {
		// since it is not a domain it has to be a ip
		return host, nil
	}
}

func checkDuplicateServer(s c.Server, server []c.Server) bool {
	isDuplicate := false
	for idx, server := range server {
		if server.HostName == s.HostName && server.IP == s.IP {
			isDuplicate = true
			if s.KeyAuth {
				config.Groups[groupIndex].Environment[envIndex].Servers[idx].KeyAuth = true
			}
		}
	}
	return isDuplicate
}
