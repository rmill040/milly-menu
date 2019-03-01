#!/usr/bin/env bash

# Create database directory and start MongoDB
echo "Setting up MongoDB"
mkdir ~/db
touch ~/db/mongodb.log
mongod --fork --dbpath ~/db --logpath ~/db/mongodb.log

# Pull most recent data from GitHub
echo "Pulling latest version of code and data from GitHub"
cd milly-menu
git pull origin master

# Push .json data to MongoDB
mongoimport --db recipes --collection all data.json --jsonArray

# Run Go script
echo "Running Go script"
./milly-menu

# Clean up directories and kill mongod process
echo "Finished, cleaning up and killing mongod"
rm -rf ../db
pkill mongod