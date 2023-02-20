# CryoUtilities
Scripts and utilities to enhance the Steam Deck experience.

[![Watch me on YouTube:](https://img.shields.io/youtube/channel/subscribers/UCJ2wc4hCWI8bEki48Zv45fQ?color=%23FF0000&label=Subscribe%20on%20YouTube&style=flat-square)](https://www.youtube.com/@cryobyte33) [![Support me on Patreon](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Fshieldsio-patreon.vercel.app%2Fapi%3Fusername%3Dcryobyte33%26type%3Dpatrons&style=flat)](https://patreon.com/cryobyte33)

## Update
Major update! The entire project has been rewritten in Go, and now contains more functionality!

If you're interested, please see the [announcement video here](https://youtu.be/C9EjXYZUqUs), where I go over all the new features and how they work.

## Functionality
* **NEW** One-click set-to-recommended settings
* **NEW** One-click revert-to-stock settings
* Swap Tuner
  * Swap File Resizer
  * Swappiness Changer
* **NEW** Memory Parameter Tuning
  * **NEW** HugePages Toggle
  * **NEW** Compaction Proactiveness Changer
  * **NEW** HugePage Defragmentation Toggle
  * **NEW** Page Lock Unfairness Changer
  * **NEW** Shared Memory (shmem) Toggle
* **NEW** Storage Manager
  * **NEW** Sync shadercache and compatdata to the same location the game is installed
  * **NEW** Delete shadercache and compatdata for whichever games you select
* **NEW** Full CLI mode

Look below for common questions and answers, or go check out my [YouTube Channel](https://www.youtube.com/@cryobyte33) 
for examples on how to use this and what performance you can expect.

## Install
[Download this link](https://raw.githubusercontent.com/CryoByte33/steam-deck-utilities/main/InstallCryoUtilities.desktop) 
to your desktop (right click and save file) on your Steam Deck, then double-click it. 

This will install the program, create desktop icons, and create menu entries.

## Usage
**NOTE:** This **REQUIRES** a password set on the Steam Deck. That can be done with the `passwd` command.
### GUI
After installation, just run the "CryoUtilities" icon on the desktop or the application menu under "Utilities".

### CLI
The latest version has a full CLI handler, which can be used to perform all tweaks, but doesn't perform game data operations.
```
sudo ~/.cryo_utilities/cryo_utilities <command> [parameter]
```
If you want to see the available commands and accepted values, you can use:
```
sudo ~/.cryo_utilities/cryo_utilities help
```
**Note:** You _need_ to use sudo for the tweaks to work, otherwise it can't write to the necessary locations on disk.

## Upgrade
Double-click the "Update CryoUtilities" icon on the desktop, you will get a dialog box when the update is complete.

## Uninstall
Double-click the "Uninstall CryoUtilities" icon on the desktop, you will be asked if you're sure, then asked if you want 
to revert the tweaks that have been made.

## Revert To Default Settings
To revert to the Steam Deck defaults, do one of the following:
* Boot CryoUtilities and click "Stock" on the homepage.
* Uninstall CryoUtilities, you'll be asked if you want to revert to stock settings. Choose yes.

After choosing these options, the Deck will be identical to an unmodified version.

## Known Issues
* If the drive becomes full during the swap file resize, you can trigger a known SteamOS bug that causes boot loops.
  * CryoUtilities is programmed in such a way to not allow this, but in the very worst cases it's still possible if something is operating/downloading in the background, at the same time CryoUtilities resizes the swap file.
  * In the event that it happens, you need to either get into a live environment and delete some files, or reinstall SteamOS with the non-destructive method.
* While using CLI mode, it is possible that the swap file resize takes long enough that the sudo credentials will time out.
  * This does not occur in GUI mode, due to how I was able to implement authentication, and will be patched out of CLI-only mode soon.

## FAQ
See [the FAQ page](https://github.com/CryoByte33/steam-deck-utilities/docs/faq.md).

## Troubleshooting
### Right-clicking the link to save it doesn't open a dialog box
Reboot the Deck or restart desktop mode, afterwards the link should work.

### The .desktop file just opens with KWrite after I download it.
Make sure the `CryoUtilities.desktop` file is on the desktop when you run it. If that still doesn't work, try one of the following:
* Run `chmod +x ~/Downloads/InstallCryoUtilities.desktop` and then try again.
* Add `CryoUtilities.desktop` as a Non-Steam game and run it from Steam.

### The swap resize times out
Go to Game Mode, navigate to Settings > System, then press "Run storage device maintenance tasks" at the very bottom.
After it's completed, you should be able to resize the swap file easily.

### Trying to do anything crashes the program
Make sure that you installed using the installer. If you can't, then run:
```bash
chown -R deck:deck ~/.cryo_utilities
chmod -R 777 ~/.cryo_utilities
```
These permissions are more open than necessary, though, so only do it as a last resort.
