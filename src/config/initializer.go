package config

import "fmt"

// Initialize builds the configuration
func (c *Config) Initialize() error {

	fmt.Println("Initialize stuff from config here")

	fmt.Println("Clear screen between questions")
	fmt.Println("Use ASCII color")
	fmt.Println("Set template location")
	fmt.Println("Set language location")
	fmt.Println("Set developer info")
	fmt.Println("Create connections - if not now - how")
	fmt.Println("Create code section - if not now - how")

	return c.Save(Location())

}
