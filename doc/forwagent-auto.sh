#!/bin/sh

UID=$(id -u)

# Make sure the run folders are created
mkdir -p  /run/user/$UID/gnupg
mkdir -p  /run/user/$UID/forwagent

# Remove the gpg-agent sockets if they exist
rm /run/user/$UID/gnupg/S.gpg-agent
rm /run/user/$UID/gnupg/S.gpg-agent.ssh

# If the local gpg-agent is running, deleting it's socket will cause it to shutdown, which will
# cause it to remove the file we are about to write, so let's wait a bit for it to shutdown
sleep 5s

# Redirect the gnupg sockets to forwagent
echo '%Assuan%\nsocket=/run/user/'$UID'/forwagent/S.gpg-agent${GPG_SOCKET_SUFFIX}' > /run/user/$UID/gnupg/S.gpg-agent
echo '%Assuan%\nsocket=/run/user/'$UID'/forwagent/S.gpg-agent.ssh${GPG_SOCKET_SUFFIX}' > /run/user/$UID/gnupg/S.gpg-agent.ssh

# Start forwagent
exec forwagent -gpgsocket /run/user/$UID/forwagent/S.gpg-agent.remote -sshsocket /run/user/$UID/forwagent/S.gpg-agent.ssh.remote
