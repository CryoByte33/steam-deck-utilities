#!/bin/bash
# Author: CryoByte33

if zenity --question --title="Disclaimer" --text="This script will uninstall CryoUtilities.\n\n<b>Disclaimer:</b> This will <b>NOT</b> remove any changes made by the script. If you want to remove the\nchanges, use the script one more time to revert to default values prior to uninstallation.\n\nProceed?" --width=600 2> /dev/null; then
    # Delete install directory
    rm -rf "$HOME/.cryo_utilities"

    # Remove Desktop icons
    rm -rf ~/Desktop/CryoUtilitiesUninstall.desktop 2>/dev/null
    rm -rf ~/Desktop/CryoUtilities.desktop 2>/dev/null

    # Remove icon from KDE
    xdg-icon-resource uninstall cryo-utilities 2>/dev/null
fi
