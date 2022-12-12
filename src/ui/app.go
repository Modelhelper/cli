package ui

import "fmt"

func Logo(version string) string {
	var logo = `
888b     d888               888          888 888    888          888                           
8888b   d8888               888          888 888    888          888                           
88888b.d88888               888          888 888    888          888                           
888Y88888P888  .d88b.   .d88888  .d88b.  888 8888888888  .d88b.  888 88888b.   .d88b.  888d888 
888 Y888P 888 d88""88b d88" 888 d8P  Y8b 888 888    888 d8P  Y8b 888 888 "88b d8P  Y8b 888P"   
888  Y8P  888 888  888 888  888 88888888 888 888    888 88888888 888 888  888 88888888 888     
888   "   888 Y88..88P Y88b 888 Y8b.     888 888    888 Y8b.     888 888 d88P Y8b.     888     
888       888  "Y88P"   "Y88888  "Y8888  888 888    888  "Y8888  888 88888P"   "Y8888  888     
                                                                     888                       
                                                                     888                       
                                                                     888   %v                  
`
	vt := fmt.Sprintf("CLI v%s", version)

	if version == "" {
		vt = ""
	}
	return fmt.Sprintf(logo, vt)
}
