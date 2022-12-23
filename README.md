# CryoUtilities
Scripts and utilities to enhance the Steam Deck experience, particularly performance.

## Update
As of the news of SteamOS 3.4 containing TRIM support, the script has been updated.

**If you are on SteamOS lower than 3.4**, you will still get an option to enable my TRIM-on schedule.

**If you are on SteamOS 3.4 or higher**, my automatic TRIM will be disabled automatically.

**Immediate TRIM is still an option and runs on /home only**

**Current Functionality:**
* Swap File Resizer
* Swappiness Changer
* Enable/Disable TRIM on schedule (SteamOS versions lower than 3.4 ONLY)
* Immediately run TRIM on all mounted volumes

Look below for common questions and answers, or go check out my [YouTube Channel](https://www.youtube.com/channel/UCJ2wc4hCWI8bEki48Zv45fQ) for examples on how to use this and what performance you can expect.

## Support Me
[![Watch me on YouTube:](https://img.shields.io/youtube/channel/subscribers/UCJ2wc4hCWI8bEki48Zv45fQ?color=%23FF0000&label=Subscribe%20on%20YouTube&style=flat-square)](https://www.youtube.com/@cryobyte33)

[![Support me on Patreon](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Fshieldsio-patreon.vercel.app%2Fapi%3Fusername%3Dcryobyte33%26type%3Dpatrons&style=flat)](https://patreon.com/cryobyte33)

## Usage
**NOTE:** This **REQUIRES** a password set on the Steam Deck. That can be done with the `passwd` command.

After installation, just run the "CryoUtilities" icon on the desktop.

## Install
### Easy
Download InstallCryoUtilities.desktop from this repository with [this link](https://raw.githubusercontent.com/CryoByte33/steam-deck-utilities/main/InstallCryoUtilities.desktop) on your Steam Deck, then run it. (Right click and save file.)

**NOTE:** Do NOT set the downloaded desktop file as executable! This will cause the installation to fail. 

This will install a script and create a few desktop icons for the swap resizer tool.

### Local Storage / Clone
This method lets you download the script locally to have on hand. You can also modify it if you'd like, but 
I don't recommend that unless you know what you're doing!

```bash
git clone https://github.com/CryoByte33/steam-deck-utilities.git
cd steam-deck-utilities
chmod +x cryo_utilities.sh
./cryo_utilities.sh
```

## Upgrade
Double-click the "Update CryoUtilities" icon on the desktop.

## Uninstall
Double-click the "Uninstall CryoUtilities" icon on the desktop.

## Revert To Default Settings
To revert to the Steam Deck defaults, do the following:
* 1GB Swap File
* 100 Swappiness
* Disable TRIM on schedule (if applicable)

After choosing these options, the Deck will be identical to an unmodified version.

## Common Questions
#### What games benefit from the swap fix?
Any game that uses a combined 16GB for RAM and VRAM will likely benefit from the swap fix. How much the swap fix helps will depend on a few factors including asset compression, asset size, and churn of data in memory.

#### Do the swap changes work for all games?
Yes, as a matter of fact, they work for any program on SteamOS loaded from any location. That means that even Google Chrome installed on an external hard drive would benefit from these changes (If the trigger conditions are met).

#### Will an update revert the swap fix?
No, but reformatting the Deck will revert it. I tested with several updates and both the swap file and swap tendency settings were left in place.

#### Does the swap fix work on Windows?
No, but Windows does have an equivalent called a page file, and its size is automatically managed by Windows. I haven't tested performance implications yet, but I believe that it won't bring quite the same amount of performance boost as the NTFS filesystem is abysmally slow.

#### Are there any situations where the swap fix isn't a good option?
Prior to the swap tendency update, extremely memory-heavy CPU-bound processes could take a performance hit from the fix. I hadn't found an instance of it in gaming, but I confirmed that KDEnLive and Blender would both cause this to happen and worked better without the swap fix.

The latest version with the swap tendency tweak doesn't have this issue, and I've been unable to trigger worse performance in any situation during testing.

#### Does the swap fix hurt emulators?
No, but I haven't done exhaustive testing yet. That said, many viewers have reported better performance in newer emulators like RPCS3, Yuzu and CEMU. I hope to do an emulation-focused video soon since I've been emulating for 13-ish years.

#### Will the swap fix hurt my SSD?
No, setting the swap tendency to 1 will reduce wear on the SSD over the stock configuration.

#### How much SSD wear does the stock configuration cause?
In general, a swap tendency of 100 will try to put data into swap prior to RAM. Because of that, a large percentage of all memory operations would wear your SSD. I don't have a spare SSD to test this, but I would guess that the default wear pattern would reduce the SSD's life to about 4 years of average gameplay on modern games.

Setting the swap tendency to 1 won't quite prevent data from making it into swap at all times, but it will reduce the frequency to only when the Deck overflows its total memory.

#### Do differently sized SSDs have a performance difference for swapping?
Yes, the read and write speeds both make a difference in the speed of swapping, and therefore performance. That said, Valve did a pretty good job at choosing good quality parts with relatively similar speeds so we'd all have a consistent experience.

In general, the larger the SSD the faster it is, but there's a lot more to it than that. But that does mean that by upgrading your SSD size you could upgrade the performance during swap scenarios. Something to keep in mind, but it will likely be a very small difference.

#### Do these fixes affect battery life?
This doesn't directly affect either the fan speeds or battery life, but it could lift a bottleneck and allow the CPU and GPU to work harder, thus raising fan speeds or lowering battery life by technicality. 

#### Are there any downsides to these fixes?
Not that I'm aware of. At worst, performance breaks even but still extends SSD life.

#### Does this script support SteamOS 3.4+'s TRIM?
Yes, it automatically disabled my TRIM-on-schedule if it detects a version 3.4 or higher. The Immediate TRIM option is still valuable and supported in all SteamOS versions.

## Troubleshooting
### Right-clicking the link to save it doesn't open a dialog box
Reboot the Deck or restart desktop mode, afterwards the link should work.

### The .desktop file just opens with KWrite after I download it.
Run `chmod +x ~/Downloads/InstallCryoUtilities.desktop` and then try again. If that still doesn't work, add it as a Non-Steam game and run it from Steam.
