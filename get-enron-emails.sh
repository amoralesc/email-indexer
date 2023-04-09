#!/usr/bin/bash

# Determine what do to if the "emails" directory already exists
if [ -d "emails" ]; then
  read -p "The 'emails' directory already exists. Do you want to delete it and its contents? (y/n) " confirm
  if [ "$confirm" == "y" ]; then
    rm -rf emails
  else
    echo "Aborting script."
    exit 1
  fi
fi

# Download the tar file
wget http://www.cs.cmu.edu/~enron/enron_mail_20110402.tgz

# Extract the contents to a directory called "emails"
mkdir emails
tar -xvzf enron_mail_20110402.tgz -C emails/
rm enron_mail_20110402.tgz
