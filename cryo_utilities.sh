#!/usr/bin/env bash
# Author: CryoByte33
# I am in no way responsible to damage done to any device this
# is executed on, all liability lies with the runner.

# Add vendored `gum` to the front of the path
# Note: this will be affective only within this script
export PATH="$HOME/.cryo_utilities:$PATH"

show_disclaimer() {
    DISCLAIMER_TEXT="This script was made by CryoByte33 to resize the swapfile on a Steam Deck.\n\n<b>Disclaimer: I am in no way responsible to damage done to any device this is executed on, all liability lies with the runner.</b>\n\nDo you accept these terms?"

    if [ -z "$SSH_TTY" ]; then
        zenity --question --title="Disclaimer" --text="$DISCLAIMER_TEXT" --width=600 2> /dev/null
    else
        gum format "$(printf "%s" "$DISCLAIMER_TEXT")"
        gum confirm 
    fi
}

show_error() {
    ERROR_TITLE=$1
    ERROR_TEXT=$2

    if [ -z "$SSH_TTY" ]; then
        zenity --error --title="$ERROR_TITLE" --text="$ERROR_TEXT" --width=300 2> /dev/null
    else
        gum format "**$(gum style --foreground 196 "$ERROR_TEXT")**"
    fi
}

ask_pass() {
    if [ -z "$SSH_TTY" ]; then
        zenity --password --title="Enter Password" --text="Enter Deck User Password (not Steam account!)" 2>/dev/null
    else
        gum input --password --prompt "Enter Deck User Password (not Steam account!): "
    fi
}

ask_question() {
    QUESTION_TITLE=$1
    QUESTION_TEXT=$2

    if [ -z "$SSH_TTY" ]; then
        zenity --question --title="$QUESTION_TITLE" --text="$QUESTION_TEXT" --width=300 2> /dev/null
    else
        gum format "$(printf "%s" "$QUESTION_TEXT")"
        gum confirm " "
    fi
}

choose_swap_file_size() {
    AVAILABLE=$1
    PROMPT_TEXT=$(printf "You have %q space available, what size would you like the swap file (in GB)?" "$AVAILABLE")

    if [ -z "$SSH_TTY" ]; then
        zenity --list --radiolist --text "$PROMPT_TEXT" --hide-header --column "Selected" --column "Size" TRUE "1" FALSE "2" FALSE "4" FALSE "8" FALSE "12" FALSE "16" FALSE "32" --height=400 2> /dev/null
    else
        gum format "$PROMPT_TEXT"
        gum choose "1" "2" "4" "8" "12" "16" "32"
    fi
}

choose_swappiness_value() {
    PROMPT_TEXT="What value would you like to use for swappiness? (Default: 100)"

    if [ -z "$SSH_TTY" ]; then
        zenity --list --title "Swappiness Value" --text "$PROMPT_TEXT" --column="vm.swappiness" "1" "10" "25" "50" "70" "100" --width=100 --height=300 2> /dev/null
    else
        gum format "$PROMPT_TEXT"
        gum choose "1" "10" "25" "50" "70" "100"
    fi
}

enable_or_disable_trim() {
    PROMPT_TEXT="Would you like to enable or disable TRIM?"

    if [ -z "$SSH_TTY" ]; then
        zenity --list --title "TRIM Choice" --text "$PROMPT_TEXT" --column="TRIM" "Enable" "Disable" --width=100 --height=300 2> /dev/null
    else
        gum format "$PROMPT_TEXT"
        gum choose "Enable" "Disable"
    fi
}

modify_swap() {
    SIZE=$1
    SWAPFILE=$2

    if [ -z "$SSH_TTY" ]; then
        (
        echo 0
        echo "# Disabling swap..."
        sudo swapoff -a
        echo 25
        echo "# Creating new $SIZE GB swapfile (be patient, this can take between 10 seconds and 30 minutes)..."
        sudo dd if=/dev/zero of="$SWAPFILE" bs=1G count="$SIZE" status=none
        echo 50
        echo "# Setting permissions on swapfile..."
        sudo chmod 0600 "$SWAPFILE"
        echo 75
        echo "# Initializing new swapfile..."
        sudo mkswap "$SWAPFILE"  
        sudo swapon "$SWAPFILE" 
        echo 100
        echo "# Process completed! You can verify the file is resized by doing 'ls -lash $SWAPFILE' or using 'swapon -s'."
        ) | zenity --title "Resizing Swap File" --progress --no-cancel --width=800 2> /dev/null
    else
        gum spin --title "Disabling swap" sudo swapoff -a
        gum spin --title "Creating new $SIZE GB swapfile (be patient, this can take between 10 seconds and 30 minutes)..." sudo dd if=/dev/zero of="$SWAPFILE" bs=1G count="$SIZE" status=none
        gum spin --title "Setting permissions on swapfile" sudo chmod 0600 "$SWAPFILE"
        gum spin --title "Initializing new swapfile" sudo mkswap "$SWAPFILE"
        gum spin --title "Initializing new swapfile" sudo swapon "$SWAPFILE"
        gum format "Process completed! You can verify the file is resized by doing 'ls -lash $SWAPFILE' or using 'swapon -s'."
        fi
    }

    run_trim() {
        if [ -z "$SSH_TTY" ]; then
            echo -e "\nTRIM Debug:"
            echo "-------------"
            (
            echo 50
            echo "# Running TRIM, please be patient (this can take up to 30 minutes)..."
            sudo fstrim -v /home
            echo 100
            echo "# TRIM Complete!"
            ) | zenity --title "Running TRIM" --progress --no-cancel --width=800 2> /dev/null
        else
            gum spin --title "Running TRIM, please be patient (this can take up to 30 minutes)" sudo fstrim -v /home
            gum format "TRIM Complete!"
        fi
    }

    if ! (show_disclaimer); then
        show_error "Terms Denied" "Terms were denied, cannot proceed."
        exit 1
    fi
    hasPass=$(passwd -S "$USER" | awk -F " " '{print $2}')
    if [[ ! $hasPass == "P" ]]; then
        show_error "Password Error" "Password is not set, please set one in the terminal with the <b>passwd</b> command, then run this again."
        exit 1
    fi
    PASSWD="$(ask_pass)"
    echo "$PASSWD" | sudo -v -S
    ans=$?
    if [[ $ans == 1 ]]; then
        show_error "Password Error" "Incorrect password provided, please run this command again and provide the correct password."
        exit 1
    fi
    echo -e "\nDebugging Information:"
    echo "----------------------"

    MACHINE_CURRENT_SWAP_SIZE=$(stat --format='%s' /home/swapfile)
    CURRENT_SWAP_SIZE=$(( MACHINE_CURRENT_SWAP_SIZE / 1024 / 1024 / 1024 ))
    CURRENT_VM_SWAPPINESS=$(sysctl vm.swappiness | awk '{print $3}')
    STEAMOS_VERSION=$(sudo cat /etc/os-release | grep VERSION_ID | sed 's/VERSION_ID=//g')

# Swapfile Size Changer
if (ask_question "Change Swap Size?" "Do you want to change the swap file size?\n\nCurrent Size: $CURRENT_SWAP_SIZE\nRecommended: 16"); then
    AVAILABLE=$(df --output="avail" -lh --sync /home | grep -v "Avail" | sed -e 's/^[ \t]*//')
    MACHINE_AVAILABLE=$(( $(df --output="avail" -l --sync /home | grep -v "Avail" | sed -e 's/^[ \t]*//') * 1024 ))
    SIZE=$(choose_swap_file_size "$AVAILABLE")
    MACHINE_SIZE=$(( SIZE * 1024 * 1024 ))
    TOTAL_AVAILABLE=$(( MACHINE_AVAILABLE + MACHINE_CURRENT_SWAP_SIZE ))
    echo "Swap Debug:"
    echo "-----------"
    echo "Bytes Available: $MACHINE_AVAILABLE"
    echo "Chosen Size: $MACHINE_SIZE"
    echo "Current Swap Size in Bytes: $MACHINE_CURRENT_SWAP_SIZE"
    echo "Total Size Available: $TOTAL_AVAILABLE"

    if [ "$MACHINE_SIZE" -lt $TOTAL_AVAILABLE ]; then
        modify_swap "$SIZE" "/home/swapfile"
    else
        show_error "Invalid Size" "You selected a size greater than the space you have available, cannot proceed." 
    fi
fi
# Swappiness Change
# Thank you to protosam for the idea and some of the code here.
if (ask_question "Change Swappiness?" "Would you like to change swappiness?\n\nCurrent value: $CURRENT_VM_SWAPPINESS\nRecommended: 1"); then
    SWAPPINESS_ANSWER=$(choose_swappiness_value)
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

    if (ask_question "Toggle TRIM?" "Would you like to enable or disable TRIM running on a schedule?\n\nCurrent value: $TRIM_STATUS\nRecommended: Enabled"); then
        TRIM_CHOICE=$(enable_or_disable_trim)
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
if (ask_question "Run TRIM Now?" "Would you like to run TRIM right now?\n\n<b>Note:</b> This can take up to 30 minutes or so."); then
    run_trim
fi
