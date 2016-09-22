#!/bin/bash
# Script used to setup elasticsearch. Can be run as a regular user (needs sudo)

ES_USER="elasticsearch"
ES_GROUP="$ES_USER"
ES_HOME="/etc/elasticsearch"
ES_CLUSTER="latticephere_dev"
ES_DATA_PATH="/var/data/elasticsearch"
ES_LOG_PATH="/var/log/elasticsearch"
ES_HEAP_SIZE=`expr $(free|awk '/^Mem:/{print $2}') / 2` # set to hlf of free mem
ES_MAX_OPEN_FILES=32000

# Path to main config
CONFIG="$ES_HOME/elasticsearch.yml"

# Path to service wrapper config
SERVICE_CONFIG="etc/sysconfig/elasticsearch"

# Add group and user (without creating the homedir)
echo "Add user: $ES_USER"
sudo useradd -d $ES_HOME -M -s /bin/bash -U $ES_USER

# Bump max open files for the user
sudo sh -c "echo '$ES_USER soft nofile $ES_MAX_OPEN_FILES' >> /etc/security/limits.conf"
sudo sh -c "echo '$ES_USER hard nofile $ES_MAX_OPEN_FILES' >> /etc/security/limits.conf"

cd ~

# echo "Update system"
# sudo yum update -y

echo "Install JRE"
#in hopes of preventing future problems lets stick with a fixed version and think about updates later
#sudo yum install jre -y
cd ~
wget --no-cookies --no-check-certificate --header "Cookie: gpw_e24=http%3A%2F%2Fwww.oracle.com%2F; oraclelicense=accept-securebackup-cookie" "http://download.oracle.com/otn-pub/java/jdk/8u73-b02/jdk-8u73-linux-x64.rpm"
sudo yum -y localinstall jdk-8u73-linux-x64.rpm
rm ~/jdk-8u73-linux-x64.rpm

echo "Downloading elasticsearch"
#wget https://github.com/downloads/elasticsearch/elasticsearch/elasticsearch-0.19.7.tar.gz -O elasticsearch.tar.gz

#tar -xf elasticsearch.tar.gz
#rm elasticsearch.tar.gz
#mv elasticsearch-* elasticsearch
#sudo mkdir -p $ES_HOME
#sudo mv elasticsearch/* $ES_HOME
#rm -rf elasticsearch
sudo rpm --import http://packages.elastic.co/GPG-KEY-elasticsearch
echo '[elasticsearch-2.x]
name=Elasticsearch repository for 2.x packages
baseurl=http://packages.elastic.co/elasticsearch/2.x/centos
gpgcheck=1
gpgkey=http://packages.elastic.co/GPG-KEY-elasticsearch
enabled=1
' | sudo tee /etc/yum.repos.d/elasticsearch.repo
sudo yum -y install elasticsearch


echo "Install service wrapper"
curl -L http://github.com/elasticsearch/elasticsearch-servicewrapper/tarball/master | tar -xz
sudo mv *servicewrapper*/service $ES_HOME/bin/
rm -rf *servicewrapper*
sudo $ES_HOME/bin/service/elasticsearch install

echo "Fix configuration files"
sudo sed -i "s|^# bootstrap.mlockall:.*$|bootstrap.mlockall: true|" $CONFIG
sudo sh -c "echo 'path.logs: $ES_LOG_PATH' >> $CONFIG"
sudo sh -c "echo 'path.data: $ES_DATA_PATH' >> $CONFIG"
sudo sh -c "echo 'cluster.name: $ES_CLUSTER' >> $CONFIG"
sudo sh -c "echo 'bootstrap.mlockall: true' >> $CONFIG"
sudo sh -c "echo 'ES_HEAP_SIZE=$ES_HEAP_SIZE' >> $SERVICE_CONFIG"
sudo sh -c "echo 'MAX_LOCKED_MEMORY=unlimited' >> $SERVICE_CONFIG"

sudo sh -c "echo 'http.port: 8001' >> $CONFIG"
sudo sh -c "echo 'node.master: true' >> $CONFIG"
sudo sh -c "echo 'http.cors.allow-origin: "*"' >> $CONFIG"
sudo sh -c "echo 'http.cors.enabled: true ' >> $CONFIG"
sudo sh -c "echo 'network.host: ["0.0.0.0"]' >> $CONFIG"
sudo sh -c "echo 'discovery.zen.ping.unicast.hosts: ["0.0.0.0"]' >> $CONFIG"

sudo sh -c "echo 'LimitMEMLOCK=infinity' >> /usr/lib/systemd/system/elasticsearch.service"

# Fix these two in $CONFIG if your network does not do multicast
# network.host: <ip of current node>
# discovery.zen.ping.unicast.hosts: ["<ip of other node in the cluster>"]

sudo sed -i "s|set\.default\.ES_HOME=.*$|set.default.ES_HOME=$ES_HOME|" $SERVICE_CONFIG
sudo sed -i "s|set\.default\.ES_HEAP_SIZE=[0-9]\+|set.default.ES_HEAP_SIZE=$ES_HEAP_SIZE|" $SERVICE_CONFIG
sudo sed -i "s|set\.default\.MAX_LOCKED_MEMORY=*\+|set.default.MAX_LOCKED_MEMORY=unlimited|" $SERVICE_CONFIG

sudo sed -i "s|^#RUN_AS_USER=.*$|RUN_AS_USER=$ES_USER|" $ES_HOME/bin/service/elasticsearch
sudo sed -i "s|^#ULIMIT_N=.*$|ULIMIT_N=$ES_MAX_OPEN_FILES|" $ES_HOME/bin/service/elasticsearch

echo "Create data and log directories and fix permissions"
sudo mkdir -p $ES_LOG_PATH $ES_DATA_PATH
sudo chown -R $ES_USER:$ES_GROUP $ES_LOG_PATH $ES_DATA_PATH $ES_HOME

echo "Install plugins"
sudo /usr/share/elasticsearch/bin/plugin install karmi/elasticsearch-paramedic # Paramedic = http://<node-ip>:<port>/_plugin/paramedic/index.html
sudo /usr/share/elasticsearch/bin/plugin install mobz/elasticsearch-head       # Head = http://<node-ip>:<port>/_plugin/head
sudo /usr/share/elasticsearch/bin/plugin install royrusso/elasticsearch-HQ     # HQ = http://<node-ip>:<port>/_plugin/hq
sudo systemctl restart elasticsearch.service
sudo systemctl status elasticsearch.service

# Start the daemon
sudo /etc/init.d/elasticsearch start
sudo systemctl enable elasticsearch

# create indecies (if they dont already exist)
# curl -XPUT 'http://localhost:9200/users/'
