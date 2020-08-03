# If we're connecting over ssh, tell gpg to use the remote socket
if [ -n "$SSH_CLIENT" ] || [ -n "$SSH_TTY" ]; then
    export GPG_SOCKET_SUFFIX=.remote
fi
