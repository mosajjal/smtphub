#!/bin/sh

# get the application name
APP_NAME="myapp"

# prompt for the application name
echo "Enter the application name: "
read APP_NAME

# replace the application name in the go files
find . -type f -name "*.go" -exec sed -i "s/myapp/$APP_NAME/g" {} \;

# replace the application name in go.mod
sed -i "s/myapp/$APP_NAME/g" go.mod

# remind the user that the package paths are still under github.com/mosajjal
echo "Don't forget to change the package paths in the go files"

# confirm and delete the setup.sh file
echo "Do you want to delete the setup.sh file? [y/n]"
read DELETE_SETUP
if [ "$DELETE_SETUP" = "y" ]; then
    rm setup.sh
fi