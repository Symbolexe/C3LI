# C3LI
Allows you to search for SSL certificates associated with a domain using crt.sh and display the details including common name, issuer, expiration dates, and more.
## Features
- [x] Fetches SSL certificates from crt.sh based on a specified domain.
- [x] Displays certificate details such as common name, issuer name, serial number, validity dates, and fingerprints.
- [x] Supports sorting certificates by issuer name or expiration date.
- [x] Optional verbose mode to display additional details like fingerprints and subject alternative names.
- [x] Option to output results to a file for easy sharing or further analysis.
## Usage
### Prerequisites
1. Go installed on your system. If not, you can download it from golang.org.
### Installation
1. Clone the repository to your local machine:

`git clone https://github.com/symbolexe/C3LI.git`

`cd C3LI`

`./C3LI --url example.com`

2. Navigate to the directory containing the tool and run it using the go run command:

`go run C3LI.go --url example.com`

Replace example.com with the domain you want to search for SSL certificates.

## Options
- [x] -url: Specify the target domain to search certificates for (required).
- [x] -v: Enable verbose output to display additional certificate details.
- [x] -sort issuer|expiration: Sort results by issuer name or expiration date.
- [x] -output filename.txt: Specify a file to save the results (default is to print to console).

## Example
Search certificates for example.com and save results to output.txt:

`go run C3LI.go --url example.com --output output.txt`
