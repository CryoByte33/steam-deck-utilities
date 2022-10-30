# steam-deck-utilities
Scripts and utilities to enhance the Steam Deck experience, particularly performance.

**Current Functionality:**
* Swap File Resizer
* Swappiness Changer
* Enable/Disable TRIM
* Immediately run TRIM on all mounted volumes

Look below for common questions and answers, or go check out my [YouTube Channel](https://www.youtube.com/channel/UCJ2wc4hCWI8bEki48Zv45fQ) for examples on how to use this and what performance you can expect.

## Usage
**NOTE:** This **REQUIRES** a password set on the Steam Deck. That can be done with the `passwd` command.

After installation, just run the "CryoUtilities" icon on the desktop.

## Install
### Easy
Download InstallCryoUtilities.desktop from this repository with [this link](https://raw.githubusercontent.com/CryoByte33/steam-deck-utilities/main/InstallCryoUtilities.desktop) on your Steam Deck, then run it. (Right click and save file)

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
Double-clock the "Update CryoUtilities" icon on the desktop.

## Uninstall
Double-click the "Uninstall CryoUtilities" icon on the desktop.

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


## Troubleshooting
### Right-clicking the link to save it doesn't open a dialog box
Reboot the Deck or restart desktop mode, afterwards the link should work.

### The .desktop file just opens with KWrite after I download it.
Run `chmod +x ~/Downloads/InstallCryoUtilities.desktop` and then try again. If that still doesn't work, add it as a Non-Steam game and run it from Steam.
