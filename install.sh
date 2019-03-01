#!/usr/bin/env bash

#####
# Install milly-menu on Amazon EC2 instance
#####

# General updates
sudo yum update

# Install git
sudo yum install -y git

# Install MongoDB
sudo mv mongodb-org-4.0.repo /etc/yum.repos.d/mongodb-org-4.0.repo
sudo yum install -y mongodb-org

# Clone milly-menu repo
git clone https://github.com/rmill040/milly-menu.git
cd milly-menu 
go get -d ./...
go build
./milly-menu --configure
chmod 755 milly-menu run.sh

# Define cronjob for every Saturday at 8:00AM
# First write out current crontab
crontab -l > mycron

# Echo new cron into cron file
echo "0 8 * * 6 ./milly-menu/run.sh" >> mycron

# Install new cron file and then remove
crontab mycron
rm mycron