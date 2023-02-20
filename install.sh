#!/bin/bash
# Author: CryoByte33

# Uninstall swap resizer if present
# Delete old install directory
rm -rf "$HOME/.swap_resizer" &>/dev/null

# Remove old Desktop icons
rm -rf ~/Desktop/SwapResizerUninstall.desktop &>/dev/null
rm -rf ~/Desktop/SwapResizer.desktop &>/dev/null

# Remove directory if present
rm -rf "$HOME/.cryo_utilities" &>/dev/null

# Create a hidden directory for the script
mkdir -p "$HOME/.cryo_utilities" &>/dev/null

# Install binary
wget https://github.com/CryoByte33/steam-deck-utilities/releases/download/latest/cryoutilities -O "$HOME/.cryo_utilities/cryo_utilities"
chmod +x "$HOME/.cryo_utilities/cryo_utilities"

# Install launcher script
curl https://raw.githubusercontent.com/CryoByte33/steam-deck-utilities/main/launcher.sh --silent --output "$HOME/.cryo_utilities/launcher.sh"
chmod +x "$HOME/.cryo_utilities/launcher.sh"

# Install Icon
curl https://raw.githubusercontent.com/CryoByte33/steam-deck-utilities/main/cmd/cryoutilities/Icon.png --silent --output "$HOME/.cryo_utilities/cryo-utilities.png"
cd ~/.cryo_utilities || exit
xdg-icon-resource install cryo-utilities.png --size 64

# Create Desktop icons
rm -rf "$HOME"/Desktop/CryoUtilitiesUninstall.desktop 2>/dev/null
echo '#!/usr/bin/env xdg-open
[Desktop Entry]
Name=Uninstall CryoUtilities
Exec=curl https://raw.githubusercontent.com/CryoByte33/steam-deck-utilities/main/uninstall.sh | bash -s --
Icon=delete
Terminal=false
Type=Application
StartupNotify=false' >"$HOME"/Desktop/CryoUtilitiesUninstall.desktop
chmod +x "$HOME"/Desktop/CryoUtilitiesUninstall.desktop

rm -rf "$HOME"/Desktop/CryoUtilities.desktop 2>/dev/null
echo "#!/usr/bin/env xdg-open
[Desktop Entry]
Name=CryoUtilities
Exec=bash $HOME/.cryo_utilities/launcher.sh
Icon=cryo-utilities
Terminal=false
Type=Application
StartupNotify=false" >"$HOME"/Desktop/CryoUtilities.desktop
chmod +x "$HOME"/Desktop/CryoUtilities.desktop

rm -rf "$HOME"/Desktop/UpdateCryoUtilities.desktop 2>/dev/null
echo "#!/usr/bin/env xdg-open
[Desktop Entry]
Name=Update CryoUtilities
Exec=curl https://raw.githubusercontent.com/CryoByte33/steam-deck-utilities/main/install.sh | bash -s --
Icon=bittorrent-sync
Terminal=false
Type=Application
StartupNotify=false" >"$HOME"/Desktop/UpdateCryoUtilities.desktop
chmod +x "$HOME"/Desktop/UpdateCryoUtilities.desktop

# Create Start Menu Icons
rm -rf "$HOME"/.local/share/applications/CryoUtilitiesUninstall.desktop 2>/dev/null
echo "#!/usr/bin/env xdg-open
[Desktop Entry]
Name=CryoUtilities - Uninstall
Exec=curl https://raw.githubusercontent.com/CryoByte33/steam-deck-utilities/main/uninstall.sh | bash -s --
Icon=delete
Terminal=false
Type=Application
Categories=Utility
StartupNotify=false" >"$HOME"/.local/share/applications/CryoUtilitiesUninstall.desktop
chmod +x "$HOME"/.local/share/applications/CryoUtilitiesUninstall.desktop

rm -rf "$HOME"/.local/share/applications/CryoUtilities.desktop 2>/dev/null
echo "#!/usr/bin/env xdg-open
[Desktop Entry]
Name=CryoUtilities
Exec=bash $HOME/.cryo_utilities/launcher.sh
Icon=cryo-utilities
Terminal=false
Type=Application
Categories=Utility
StartupNotify=false" >"$HOME"/.local/share/applications/CryoUtilities.desktop
chmod +x "$HOME"/.local/share/applications/CryoUtilities.desktop

rm -rf "$HOME"/.local/share/applications/UpdateCryoUtilities.desktop 2>/dev/null
echo "#!/usr/bin/env xdg-open
[Desktop Entry]
Name=CryoUtilities - Update
Exec=curl https://raw.githubusercontent.com/CryoByte33/steam-deck-utilities/main/install.sh | bash -s --
Icon=bittorrent-sync
Terminal=false
Type=Application
Categories=Utility
StartupNotify=false" >"$HOME"/.local/share/applications/UpdateCryoUtilities.desktop
chmod +x "$HOME"/.local/share/applications/UpdateCryoUtilities.desktop

update-desktop-database ~/.local/share/applications

zenity --info --text="Install/upgrade of CryoUtilities has been completed!" --width=300
