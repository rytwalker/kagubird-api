#!/bin/bash
set -eu

# ==================================================================================== #
# VARIABLES
# ==================================================================================== #

# Set the timezone for the server. A full list of available timezones can be found by running timedatectl list-timezones.
TIMEZONE=America/New_York

# Set the name of the new user to create
USERNAME=kagubird

# Prompt to enter a password for the postgreSQL kagubird user.
read -p "Enter password for kagubird DB user: " DB_PASSWORD

# Force output to be presented in en_US for the duration of the script.
export LC_ALL=en_US.UTF-8

# ==================================================================================== #
# SCRIPT LOGIC
# ==================================================================================== #

# Enable the "universe" repo.
add-apt-repository --yes universe

# Update all software packages
apt update

# Set the system timezone and install all locales.
timedatectl set-timezone ${TIMEZONE}
apt --yes install locales-all

# useradd --create-home --shell "/bin/bash" --groups sudo "${USERNAME}"

# Force a password to be set for the new user the first time they log in.
# passwd --delete "${USERNAME}"
# chage --lastday 0 ${USERNAME}

# Copy the ssh keys from the root user to the new user
# rsync --archive --chown=${USERNAME}:${USERNAME} /root/.ssh /home/${USERNAME}

# Config the firewall to allow SSH, HTTP, and HTTPS traffic.
# ufw allow 22
# ufw allow 80/tcp
# ufw allow 443/tcp
# ufw --force enable

# Install fail2ban.
apt --yes install fail2ban

# Install migrate CLI tool
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz
mv migrate.linux-amd64 /usr/local/bin/migrate

# Install postgreSQL
apt --yes install postgresql

# Set up the kagubird db and create a user account with the password entered earlier
sudo -i postgres psql -c "CREATE DATABASE kagubird"
sudo -i postgres psql -d kagubird -c "CREATE EXTENSION IF NOT EXISTS citext"
sudo -i postgres psql -d kagubird -c "CREATE ROLE kagubird WITH LOGIN PASSWORD '${DB_PASSWORD}'"

# Add a DSN for connecting to the kagubird database to the system-wide environment
# variables in the /etc/environment file
echo "KAGUBIRD_DB_DSN='postgres://kagubird:${DB_PASSWORD}@localhost/kagubird'" >> /etc/environment

# Install Caddy (see https://caddyserver.com/docs/install#debian-ubuntu-raspbian).
apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
apt update
apt --yes install caddy

# Upgrade all packages. Using the --force-confnew flag means that config files
# will be replaced if newer ones are available.
apt --yes -o Dpkg::Options::="--force-confnew" upgrade

echo "Script complete! Rebooting..."
reboot
