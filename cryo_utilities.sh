#!/bin/bash
# Author: CryoByte33
# I am in no way responsible to damage done to any device this
# is executed on, all liability lies with the runner.

if ! (zenity --question --title="Disclaimer" --text="This script was made by CryoByte33 to resize the swapfile on a Steam Deck.\n\n<b>Disclaimer: I am in no way responsible to damage done to any device this is executed on, all liability lies with the runner.</b>\n\nDo you accept these terms?" --width=600 2> /dev/null); then
    zenity --error --title="Terms Denied" --text="Terms were denied, cannot proceed." --width=300 2> /dev/null
    exit 1
fi
hasPass=$(passwd -S "$USER" | awk -F " " '{print $2}')
if [[ ! $hasPass == "P" ]]; then
    zenity --error --title="Password Error" --text="Password is not set, please set one in the terminal with the <b>passwd</b> command, then run this again." --width=400 2> /dev/null
    exit 1
fi
PASSWD="$(zenity --password --title="Enter Password" --text="Enter Deck User Password (not Steam account!)" 2>/dev/null)"
echo "$PASSWD" | sudo -v -S
ans=$?
if [[ $ans == 1 ]]; then
    zenity --error --title="Password Error" --text="Incorrect password provided, please run this command again and provide the correct password." --width=400 2> /dev/null
    exit 1
fi
echo -e "\nDebugging Information:"
echo "----------------------"

MACHINE_CURRENT_SWAP_SIZE=$(ls -l /home/swapfile | awk '{print $5}')
CURRENT_SWAP_SIZE=$(( MACHINE_CURRENT_SWAP_SIZE / 1024 / 1024 / 1024 ))
CURRENT_VM_SWAPPINESS=$(sysctl vm.swappiness | awk '{print $3}')
STEAMOS_VERSION=$(sudo cat /etc/os-release | grep VERSION_ID | sed 's/VERSION_ID=//g')

# Swapfile Size Changer
if zenity --question --title="Change Swap Size?" --text="Do you want to change the swap file size?\n\nCurrent Size: $CURRENT_SWAP_SIZE\nRecommended: 16" --width=300 2> /dev/null; then
    AVAILABLE=$(df --output="avail" -lh --sync /home | grep -v "Avail" | sed -e 's/^[ \t]*//')
    MACHINE_AVAILABLE=$(( $(df --output="avail" -l --sync /home | grep -v "Avail" | sed -e 's/^[ \t]*//') * 1024 ))
    SIZE=$(zenity --list --radiolist --text "You have $AVAILABLE space available, what size would you like the swap file (in GB)?" --hide-header --column "Selected" --column "Size" TRUE "1" FALSE "2" FALSE "4" FALSE "8" FALSE "12" FALSE "16" FALSE "32" --height=400 2> /dev/null)
    MACHINE_SIZE=$(( SIZE * 1024 * 1024 ))
    TOTAL_AVAILABLE=$(( MACHINE_AVAILABLE + MACHINE_CURRENT_SWAP_SIZE ))
    echo "Swap Debug:"
    echo "-----------"
    echo "Bytes Available: $MACHINE_AVAILABLE"
    echo "Chosen Size: $MACHINE_SIZE"
    echo "Current Swap Size in Bytes: $MACHINE_CURRENT_SWAP_SIZE"
    echo "Total Size Available: $TOTAL_AVAILABLE"

    if [ "$MACHINE_SIZE" -lt $TOTAL_AVAILABLE ]; then
        (
            echo 0
            echo "# Disabling swap..."
            sudo swapoff -a
            echo 25
            echo "# Creating new $SIZE GB swapfile (be patient, this can take between 10 seconds and 30 minutes)..."
            sudo dd if=/dev/zero of=/home/swapfile bs=1G count="$SIZE" status=none
            echo 50
            echo "# Setting permissions on swapfile..."
            sudo chmod 0600 /home/swapfile
            echo 75
            echo "# Initializing new swapfile..."
            sudo mkswap /home/swapfile  
            sudo swapon /home/swapfile 
            echo 100
            echo "# Process completed! You can verify the file is resized by doing 'ls -lash /home/swapfile' or using 'swapon -s'."
        ) | zenity --title "Resizing Swap File" --progress --no-cancel --width=800 2> /dev/null
    else
        zenity --error --title="Invalid Size" --text="You selected a size greater than the space you have available, cannot proceed." --width=500 2> /dev/null
    fi
fi
# Swappiness Change
# Thank you to protosam for the idea and some of the code here.
if zenity --question --title="Change Swappiness?" --text="Would you like to change swappiness?\n\nCurrent value: $CURRENT_VM_SWAPPINESS\nRecommended: 1" --width=300 2> /dev/null; then
    SWAPPINESS_ANSWER=$(zenity --list --title "Swappiness Value" --text "What value would you like to use for swappiness? (Default: 100)" --column="vm.swappiness" "1" "10" "25" "50" "70" "100" --width=100 --height=300 2> /dev/null)
    echo -e "\nSwappiness Debug:"
    echo "-------------------"
    sudo sysctl -w "vm.swappiness=$SWAPPINESS_ANSWER"
    if [ "$SWAPPINESS_ANSWER" -eq "100" ]; then
        sudo rm -f -f /etc/sysctl.d/zzz-custom-swappiness.conf
    else
        echo "vm.swappiness=$SWAPPINESS_ANSWER" | sudo tee /etc/sysctl.d/zzz-custom-swappiness.conf
    fi
fi
# Whether to manipulate the TRIM timer
if (( $(echo "$STEAMOS_VERSION 3.4" | awk '{print ($1 < $2)}') )); then
    echo "SteamOS version with no native TRIM support, providing schedule option..."
    # Check for current TRIM status on SteamOS versions lower than 3.4
    systemctl list-timers | grep fstrim &>/dev/null
    if [ "$?" == "1" ]; then
        TRIM_STATUS="Disabled"
    else
        TRIM_STATUS="Enabled"
    fi
    if zenity --question --title="Toggle TRIM?" --text="Would you like to enable or disable TRIM running on a schedule?\n\nCurrent value: $TRIM_STATUS\nRecommended: Enabled" --width=300 2> /dev/null; then
        TRIM_CHOICE=$(zenity --list --title "TRIM Choice" --text "Would you like to enable or disable TRIM" --column="TRIM" "Enable" "Disable" --width=100 --height=300 2> /dev/null)
        if [ "$TRIM_CHOICE" = "Enable" ]; then
            sudo systemctl enable --now fstrim.timer &>/dev/null
        else
            sudo systemctl disable --now fstrim.timer &>/dev/null
        fi
    fi
else
    echo "SteamOS version with native TRIM support, disabling custom schedule..."
    sudo systemctl disable --now fstrim.timer &>/dev/null
fi

# Whether to execute TRIM immediately
if zenity --question --title="Run TRIM Now?" --text="Would you like to run TRIM right now?\n\n<b>Note:</b> This can take up to 30 minutes or so." --width=300 2> /dev/null; then
    echo -e "\nTRIM Debug:"
    echo "-------------"
    (
        echo 50
        echo "# Running TRIM, please be patient (this can take up to 30 minutes)..."
        sudo fstrim -v /home
        echo 100
        echo "# TRIM Complete!"
    ) | zenity --title "Running TRIM" --progress --no-cancel --width=800 2> /dev/null
fi
