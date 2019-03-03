#!/usr/bin/env bash

#####
# Run milly-menu on Amazon EC2 instance
#####

# Create database directory and start MongoDB
echo "Setting up MongoDB"
mkdir ~/db
touch ~/db/mongodb.log
mongod --fork --dbpath ~/db --logpath ~/db/mongodb.log

# Pull most recent data from GitHub
echo "Pulling latest version of code and data from GitHub"
cd /home/ec2-user/go/src/github.com/rmill040/milly-menu
git pull origin master

# Clean and then rebuild
go clean
go build

# Push .json data to MongoDB
echo "Pushing recipes data into MongoDB"
mongoimport --db recipes --collection all assets/data.json --jsonArray

# Run Go script
echo "Running Go script"
./milly-menu

# Kill mongod process and clean up
echo "Finished, cleaning up and killing mongod"
pkill mongod
rm -rf ~/db