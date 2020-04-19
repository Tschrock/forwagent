# Sample additions to ~/.bashrc
# Modify interface and port below as desired.

# Check if forwagent isn't running, in which case, run it.
if ! pgrep -u $USER -f "forwagent" > /dev/null
then
  forwagent 192.168.137.1:4711 2>&1 > /dev/null &
fi

# Make SSH use the correct socket.
export SSH_AUTH_SOCK="$HOME/.gnupg/S.gpg-agent.ssh"