#!/bin/bash

if [[ $UID -ne 0 ]]; then
    echo "Run as sudoer"
    exit 0
fi

SERIAL=""

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --serial)
      SERIAL="$2"
      shift # past argument
      shift # past value
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

# Check if serial is empty
if [ -z "$SERIAL" ]; then
  echo "Error: --serial argument required"
  exit 1
fi

install_dir="/usr/bin"
service_dir="/etc/systemd/system"
service_name="thoth-agent.service"
config_dir="/etc/thoth"

if [[ ! -f "$install_dir/thoth-agent" ]]; then
    cp thoth-agent "$install_dir/thoth-agent"
    echo "Copied thoth-agent to $install_dir"
fi

chmod +x "$install_dir/thoth-agent"

if [[ ! -f "$service_dir/$service_name" ]]; then
    cp thoth-agent.service "$service_dir/$service_name"
    echo "Copied $service_name to $service_dir"
fi

mkdir -p "$config_dir"

if [[ ! -f "$config_dir/serial.id" ]]; then
  echo $SERIAL > "$config_dir/serial.id"
fi

systemctl is-enabled "$service_name" &>/dev/null || {
    systemctl enable "$service_name"
    echo "Enabled $service_name"
}

systemctl is-active "$service_name" &>/dev/null || {
    systemctl start "$service_name"
    echo "Started $service_name"
}

# Loading animation function
animate_loading() {
    local chars="/-\|"
    local delay=0.1

    while true; do
        for (( i = 0; i < ${#chars}; i++ )); do
            echo -en "\r[${chars:$i:1}] Installing Thoth Agent..."
            sleep "$delay"
        done
    done
}

# Start the loading animation in the background
animate_loading &

# Simulate installation process (sleep for demonstration)
sleep 5

# Stop the loading animation by killing the background process
kill $!

# Print a newline to clear the loading animation line
echo

echo "Thoth Agent installation completed!"