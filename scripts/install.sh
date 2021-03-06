#!/usr/bin/env bash

#####
# Install milly-menu on Amazon EC2 instance
#####

# General updates
sudo yum -y update

# Install go
sudo yum install -y golang

# Install MongoDB
sudo mv assets/mongodb-org-4.0.repo /etc/yum.repos.d/mongodb-org-4.0.repo
sudo yum install -y mongodb-org

# Clone milly-menu repo
cd ~/
mkdir -p /home/ec2-user/go/src/github.com/rmill040
mv milly-menu /home/ec2-user/go/src/github.com/rmill040/milly-menu
cd /home/ec2-user/go/src/github.com/rmill040/milly-menu

# Install dependencies and build
go get
go build
./milly-menu --configure
chmod 755 scripts/*.sh

# Define cronjob for every Saturday at 8:00AM
# First write out current crontab
crontab -l 2>/dev/null

# Echo new cron into cron file
echo "0 8 * * 6 sh /home/ec2-user/go/src/github.com/rmill040/milly-menu/scripts/run.sh" >> mycron

# Install new cron file and then remove
crontab mycron
rm mycron