#!/usr/bin/env bash

#####
# Install milly-menu on Amazon EC2 instance
#####

# General updates
sudo yum update

# Install go
sudo yum install -y go

# Install MongoDB
sudo mv mongodb-org-4.0.repo /etc/yum.repos.d/mongodb-org-4.0.repo
sudo yum install -y mongodb-org

# Clone milly-menu repo
go get -d ./...
go build
./milly-menu --configure
chmod 755 run.sh

# Define cronjob for every Saturday at 8:00AM
# First write out current crontab
crontab -l > mycron

# Echo new cron into cron file
echo "0 8 * * 6 ./milly-menu/run.sh" >> mycron

# Install new cron file and then remove
crontab mycron
rm mycron