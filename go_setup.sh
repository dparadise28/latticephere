#!/bin/bash
## Install Golang 1.6.2 64Bits on Linux (Debian|Ubuntu|OpenSUSE|CentOS)
## http://www.linuxpro.com.br/2015/06/golang-aula-1-instalacao-da-linguagem-no-linux.html
## Run as root (sudo su)
## Thank's **Bruno Albuquerque bga at bug-br.org.br**


GO_URL="https://storage.googleapis.com/golang"
GO_FILE="go1.6.2.linux-amd64.tar.gz"

# Check if user has root privileges
if [[ $EUID -ne 0 ]]; then
echo "You must run the script as root or using sudo"
   exit 1
fi


GET_OS=$(cat /etc/os-release | head -n1 | cut -d'=' -f2 | awk '{ print tolower($1) }'| tr -d '"')

if [[ $GET_OS == 'debian' || $GET_OS == 'ubuntu' ]]; then
   apt-get update
   apt-get install wget git-core
fi

if [[ $GET_OS == 'opensuse' ]]; then
   zypper in -y wget git-core
fi

if [[ $GET_OS == 'centos' ]]; then
   yum install wget git-core
fi


cd /tmp
wget --no-check-certificate ${GO_URL}/${GO_FILE}
tar -xzf ${GO_FILE}
mv go /usr/local/go


echo 'export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/GO
export PATH=$PATH:$GOPATH/bin' >> /etc/profile

### You do not need to run commands with root or sudo
source /etc/profile
mkdir -p $HOME/GO

## Test if Golang is working
go version

### The output is this:
go version go1.6 linux/amd64

echo 'export PATH=$PATH:/usr/local/go/bin:/root/GO/bin' >>~/.bash_profile
