#!/bin/bash

# execute openvpn with it's profile file and also set authentication to it

user=$1
config=$2
vpnProfile=$3

result=$(content-plus-totp -c $config)
if [ $? != 0 ]; then
  echo "cannot be proceed, read above read message for more details"
  exit 1
fi

sudo bash -c "openvpn --config ${vpnProfile} --auth-user-pass <(echo -e \"${user}${result}\")"
