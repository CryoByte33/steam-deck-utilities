#!/bin/bash
# Author: CryoByte33

INSTALL_DIR="$HOME/.cryo_utilities"

# Uninstall swap resizer if present
# Delete old install directory
rm -rf "$HOME/.swap_resizer" &>/dev/null

# Remove old Desktop icons
rm -rf ~/Desktop/SwapResizerUninstall.desktop &>/dev/null
rm -rf ~/Desktop/SwapResizer.desktop &>/dev/null

# Remove directory if present
rm -rf "$INSTALL_DIR" &>/dev/null

# Create a hidden directory for the script
mkdir -p "$INSTALL_DIR" &>/dev/null

# Install Script
curl https://raw.githubusercontent.com/CryoByte33/steam-deck-utilities/main/cryo_utilities.sh --silent --output "$INSTALL_DIR/cryo_utilities.sh"

# Install gum. The current version is 0.8.0
curl https://github.com/charmbracelet/gum/releases/download/v0.8.0/gum_0.8.0_linux_x86_64.tar.gz --location --silent --output "$INSTALL_DIR/gum_0.8.0_linux_x86_64.tar.gz"
(cd $INSTALL_DIR && tar xzf gum_0.8.0_linux_x86_64.tar.gz && chmod u+x gum)

# Install Icon
curl https://raw.githubusercontent.com/CryoByte33/steam-deck-utilities/main/icon.png --silent --output "$INSTALL_DIR/cryo-utilities.png"
cd $INSTALL_DIR
xdg-icon-resource install cryo-utilities.png --size 64

# Create Desktop icons
rm -rf ~/Desktop/CryoUtilitiesUninstall.desktop 2>/dev/null
echo '#!/usr/bin/env xdg-open
[Desktop Entry]
Name=Uninstall CryoUtilities
Exec=curl https://raw.githubusercontent.com/CryoByte33/steam-deck-utilities/main/uninstall.sh | bash -s --
Icon=delete
Terminal=true
Type=Application
StartupNotify=false' > ~/Desktop/CryoUtilitiesUninstall.desktop
chmod +x ~/Desktop/CryoUtilitiesUninstall.desktop

rm -rf ~/Desktop/CryoUtilities.desktop 2>/dev/null
echo '#!/usr/bin/env xdg-open
[Desktop Entry]
Name=CryoUtilities
Exec=bash $HOME/.cryo_utilities/cryo_utilities.sh
Icon=cryo-utilities
Terminal=true
Type=Application
StartupNotify=false' > ~/Desktop/CryoUtilities.desktop
chmod +x ~/Desktop/CryoUtilities.desktop

rm -rf ~/Desktop/UpdateCryoUtilities.desktop 2>/dev/null
echo '#!/usr/bin/env xdg-open
[Desktop Entry]
Name=Update CryoUtilities
Exec=curl https://raw.githubusercontent.com/CryoByte33/steam-deck-utilities/main/install.sh | bash -s --
Icon=bittorrent-sync
Terminal=true
Type=Application
StartupNotify=false' > ~/Desktop/UpdateCryoUtilities.desktop
chmod +x ~/Desktop/UpdateCryoUtilities.desktop
