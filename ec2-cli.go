package main

import (
 "encoding/json"
 "fmt"
 "os"
 "os/exec"
 "strings"
)

type EC2Instance struct {
 InstanceID string `json:"InstanceId"`
 State struct {
  Name string `json:"Name"`
 } `json:"State"`
 Tags []struct {
  Key   string `json:"Key"`
  Value string `json:"Value"`
 } `json:"Tags"`
 PrivateIPAddress string `json:"PrivateIpAddress"`
 PublicIPAddress  string `json:"PublicIpAddress"`
 Placement        struct {
  AvailabilityZone string `json:"AvailabilityZone"`
 } `json:"Placement"`
}

type EC2Response struct {
 Reservations []struct {
  Instances []EC2Instance `json:"Instances"`
 } `json:"Reservations"`
}

func printTableRow(columns []string, widths []int) {
 for i, col := range columns {
  fmt.Printf("%-*s", widths[i], col)
 }
 fmt.Println()
}

func main() {
 if len(os.Args) > 1 {
  command := os.Args[1]
  if command == "--help" || command == "help" {
   printHelp()
   return
  }
  if command == "list" {
   displayEC2Instances()
   return
  }
  if len(os.Args) < 3 {
   fmt.Println("Please provide the instance Name.")
   return
  }
  name := os.Args[2]
  switch command {
  case "start":
   manageInstance(name, "start-instances")
  case "stop":
   manageInstance(name, "stop-instances")
  default:
   fmt.Println("Invalid command. Use 'list', 'start', 'stop', or '--help'.")
  }
  return
 }

 // Default to displaying help if no arguments are provided
 printHelp()
}

func displayEC2Instances() {
 // AWS CLI command to describe EC2 instances
 cmd := exec.Command("aws", "ec2", "describe-instances", "--output", "json")
 output, err := cmd.Output()
 if err != nil {
  fmt.Printf("Error executing AWS CLI command: %v\n", err)
  return
 }

 // Parse the JSON output
 var ec2Response EC2Response
 err = json.Unmarshal(output, &ec2Response)
 if err != nil {
  fmt.Printf("Error parsing JSON output: %v\n", err)
  return
 }

 // Table headers and column widths
 headers := []string{"Name", "Private IP", "Public IP", "AZ", "State"}
 widths := []int{30, 20, 20, 20, 15}
 border := strings.Repeat("-", sum(widths)+len(widths)-1)

 // Print table header
 fmt.Println(border)
 printTableRow(headers, widths)
 fmt.Println(border)

 // Extract and display the desired information
 for _, reservation := range ec2Response.Reservations {
  for _, instance := range reservation.Instances {
   name := ""
   for _, tag := range instance.Tags {
    if tag.Key == "Name" {
     name = tag.Value
     break
    }
   }

   columns := []string{
    name,
    instance.PrivateIPAddress,
    instance.PublicIPAddress,
    instance.Placement.AvailabilityZone,
    instance.State.Name,
   }
   printTableRow(columns, widths)
  }
 }
 fmt.Println(border)
}

func manageInstance(name, action string) {
 // AWS CLI command to describe EC2 instances
 cmd := exec.Command("aws", "ec2", "describe-instances", "--output", "json")
 output, err := cmd.Output()
 if err != nil {
  fmt.Printf("Error executing AWS CLI command: %v\n", err)
  return
 }

 // Parse the JSON output
 var ec2Response EC2Response
 err = json.Unmarshal(output, &ec2Response)
 if err != nil {
  fmt.Printf("Error parsing JSON output: %v\n", err)
  return
 }

 // Find the instance ID by name
 instanceID := ""
 for _, reservation := range ec2Response.Reservations {
  for _, instance := range reservation.Instances {
   for _, tag := range instance.Tags {
    if tag.Key == "Name" && tag.Value == name {
     instanceID = instance.InstanceID
     break
    }
   }
   if instanceID != "" {
    break
   }
  }
  if instanceID != "" {
   break
  }
 }

 if instanceID == "" {
  fmt.Printf("Instance with Name '%s' not found.\n", name)
  return
 }

 // Execute start or stop action
 actionCmd := exec.Command("aws", "ec2", action, "--instance-ids", instanceID)
 actionOutput, err := actionCmd.CombinedOutput()
 if err != nil {
  fmt.Printf("Error executing action '%s' on instance '%s': %v\nOutput: %s\n", action, name, err, string(actionOutput))
  return
 }

 fmt.Printf("Action '%s' executed on instance '%s'. Response: %s\n", action, name, string(actionOutput))
}

func sum(numbers []int) int {
 total := 0
 for _, number := range numbers {
  total += number
 }
 return total
}

func printHelp() {
 helpText := `Usage:
  ./ec2-cli list                  # Display EC2 instance list
  ./ec2-cli start <instance-name>  # Start an EC2 instance by Name tag
  ./ec2-cli stop <instance-name>   # Stop an EC2 instance by Name tag

Options:
  --help                  # Display this help message
`
 fmt.Println(helpText)
}
