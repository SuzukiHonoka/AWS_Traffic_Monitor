# AWS Traffic Monitor

A lightweight monitoring tool for AWS services that tracks network traffic usage and executes predefined actions when usage limits are reached.

## Overview

AWS Traffic Monitor uses the official AWS Go SDK to check traffic usage for your AWS resources. When usage exceeds your configured limits, the tool can automatically execute custom commands or perform API actions like shutting down instances to prevent additional charges.

## Features

- Real-time monitoring of AWS network traffic
- Configurable usage limits with custom actions
- Direct AWS API integration (no AWS CLI required)
- Built-in shutdown functionality via AWS API
- Periodic checking with adjustable intervals
- Simple JSON configuration format

## Requirements

- Properly configured AWS credentials
- Go 1.16+ (for building from source)

## Usage

```bash
# Basic usage
atm -c config.json
```

### Command Line Arguments

- `-c` Path to config file (JSON format)

## Configuration

Create a JSON configuration file with your instance details, limits, and actions:

```json
[
  {
    "Name": "Ubuntu-1",
    "Limit": {
      "Unit": "GB", 
      "Value": 1000
    },
    "Command": ["shutdown"]
  },
  {
    "Name": "WebServer",
    "Limit": {
      "Unit": "GB",
      "Value": 500
    },
    "Command": ["sudo service nginx stop"]
  }
]
```

Each entry in the configuration array represents a monitored instance with:
- `Name`: Instance identifier
- `Limit`: Traffic limit with unit and value
- `Command`: Array of commands to execute when limit is reached
    - Use `"shutdown"` to force shutdown the instance via AWS API

## Supported Services

Currently, AWS Traffic Monitor supports:
- **Amazon Lightsail**

Support for additional AWS services is planned for future releases.

## License

This project is licensed under the GNU General Public License v3.0 - see the LICENSE file for details.

## Author

Created by SuzukiHonoka (starx)