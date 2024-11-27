# EC2 Management CLI with AWS CLI Integration

This Go application allows you to manage AWS EC2 instances using AWS CLI. You can list instances, and start or stop instances by their `Name` tag.

## Prerequisites
- **AWS CLI** installed and configured with appropriate credentials (`aws configure`).
- Go installed on your machine.
- Permissions to describe, start, and stop EC2 instances in your AWS account.

## Installation
1. Clone the repository or copy the code into a file named `main.go`.
2. Build the application (optional for production use):
   ```bash
   go build  ec2-cli.go
   ```

## Usage
Run the application directly with `go run` or use the built executable.

### Commands

#### 1. List EC2 Instances
Displays a table of EC2 instances with their `Name`, private IP, public IP, availability zone (AZ), and state.

```bash
./ec2-cli list
```

Sample Output:
```
------------------------------------------------------------
Name                          Private IP         Public IP          AZ                  State         
------------------------------------------------------------
web-server                    10.0.0.1          54.123.45.67       us-east-1a          running       
db-server                     10.0.0.2                           us-east-1b          stopped       
------------------------------------------------------------
```

#### 2. Start an EC2 Instance
Starts an EC2 instance by its `Name` tag.

```bash
go run main.go start <instance-name>
```

Example:
```bash
./ec2-cli start web-server
```

#### 3. Stop an EC2 Instance
Stops an EC2 instance by its `Name` tag.

```bash
./ec2-cli stop <instance-name>
```

Example:
```bash
./ec2-cli stop web-server
```

#### 4. Display Help
Displays the help message with usage instructions.

```bash
./ec2-cli --help
```

## Error Handling
- If an instance with the specified `Name` is not found, the application will notify you.
- Errors during AWS CLI execution (e.g., permissions or configuration issues) will be displayed in the output.

## Notes
- The application uses AWS CLI under the hood, so ensure that your AWS CLI configuration is correct and functional.
- Public IP is displayed only for instances with an associated Elastic IP or auto-assigned public IP.
- The `Name` tag is case-sensitive.

## Future Enhancements
- Support for additional EC2 actions like termination or reboot.
- Improved filtering and sorting options.

## License
This project is licensed under the MIT License.

