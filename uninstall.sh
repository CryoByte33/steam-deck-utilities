#!/bin/bash
# Author: CryoByte33

# This is nested for good reason, zenity won't exit the entire script if the 'x' button is pressed.
# Nesting it forces execution only if an option is selected.
if zenity --question --title="Disclaimer" --text="This script will uninstall CryoUtilities.\n\n<b>Disclaimer:</b> Do you want to proceed?" --width=600 2>/dev/null; then
  if zenity --question --title="Revert" --text="Do you want to revert the tweaks made by CryoUtilities?\n\n<b>Note:</b> This does NOT move the game data to the original location on the SSD." --width=600 2>/dev/null; then
    # Ask for password
    hasPass=$(passwd -S "$USER" | awk -F " " '{print $2}')
    if [[ $hasPass != "P" ]]; then
      zenity --error --title="Password Error" --text="Password is not set, please set one in the terminal with the <b>passwd</b> command, then run this again." --width=400 2>/dev/null
      exit 1
    fi
    PASSWD="$(zenity --password --title="Enter Password" --text="Enter Deck User Password (not Steam account!)" 2>/dev/null)"
    echo "$PASSWD" | sudo -v -S
    ans=$?
    if [[ $ans == 1 ]]; then
      zenity --error --title="Password Error" --text="Incorrect password provided, please run this command again and provide the correct password." --width=400 2>/dev/null
      exit 1
    fi
    # Revert everything to stock
    sudo bash "$HOME"/.cryo_utilities/cryo_utilities stock
  fi
  # Delete install directory
  rm -rf "$HOME/.cryo_utilities"

  # Remove Desktop icons
  rm -rf "$HOME"/Desktop/CryoUtilitiesUninstall.desktop 2>/dev/null
  rm -rf "$HOME"/Desktop/CryoUtilities.desktop 2>/dev/null
  rm -rf "$HOME"/Desktop/UpdateCryoUtilities.desktop 2>/dev/null

  # Remove Start Menu shortcuts
  rm -rf "$HOME"/.local/share/applications/CryoUtilitiesUninstall.desktop 2>/dev/null
  rm -rf "$HOME"/.local/share/applications/CryoUtilities.desktop 2>/dev/null
  rm -rf "$HOME"/.local/share/applications/UpdateCryoUtilities.desktop 2>/dev/null
  update-desktop-database ~/.local/share/applications

  # Remove icon from KDE
  xdg-icon-resource uninstall cryo-utilities 2>/dev/null
fi
