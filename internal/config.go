// CryoUtilities
// Copyright (C) 2023 CryoByte33 and contributors to the CryoUtilities project

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package internal

import (
	"image/color"
	"os"
	"path/filepath"
)

// CurrentVersionNumber Version number to build with, Fyne can't support build flags just yet.
var CurrentVersionNumber = "2.2.1"

// Get home Directory
var HomeDirectory, _ = os.UserHomeDir()

// InstallDirectory Location the program is installed.
var InstallDirectory = filepath.Join(HomeDirectory, ".cryo_utilities")

// LogFilePath Location of the log file
var LogFilePath = filepath.Join(InstallDirectory, "cryoutilities.log")

type Tweak struct {
	Location    string
	Recommended string
	Default     string
}

////////////////////
// Tweak Settings //
////////////////////

var TweakList = map[string]Tweak{
	"swappiness":               {Location: "/proc/sys/vm/swappiness", Recommended: "1", Default: "100"},
	"page_lock_unfairness":     {Location: "/proc/sys/vm/page_lock_unfairness", Recommended: "1", Default: "5"},
	"compaction_proactiveness": {Location: "/proc/sys/vm/compaction_proactiveness", Recommended: "0", Default: "20"},
	"hugepages":                {Location: "/sys/kernel/mm/transparent_hugepage/enabled", Recommended: "always", Default: "madvise"},
	"shmem_enabled":            {Location: "/sys/kernel/mm/transparent_hugepage/shmem_enabled", Recommended: "advise", Default: "never"},
	"defrag":                   {Location: "/sys/kernel/mm/transparent_hugepage/khugepaged/defrag", Recommended: "0", Default: "1"},
}

// These are supposed to be integers, I can't store them as strings in the map above. idk what to do about these.
var RecommendedSwapSize = 16
var RecommendedSwapSizeBytes = int64(RecommendedSwapSize * GigabyteMultiplier)

//////////////////////
// Default Settings //
//////////////////////

var DefaultSwapFileLocation = "/home/swapfile"
var DefaultSwapSize = 1
var DefaultSwapSizeBytes = int64(DefaultSwapSize * GigabyteMultiplier)

var TmpFilesRoot = "/etc/tmpfiles.d"

var TemplateUnitFile = "# Path Mode UID GID Age Argument\nw PARAM - - - - VALUE"

var OldSwappinessUnitFile = "/etc/sysctl.d/zzz-custom-swappiness.conf"
var NHPTestingFile = "/proc/sys/vm/nr_hugepages"

/////////////////
// UI Settings //
/////////////////

// HeaderTextSize Header Text Size
var HeaderTextSize = float32(32)

// SubHeadingTextSize Subheader Text Size
var SubHeadingTextSize = float32(16)

// Green UI Color
var Green = color.RGBA{R: 0, G: 155, B: 0, A: 255}

// Gray UI Color
var Gray = color.RGBA{R: 155, G: 155, B: 155, A: 255}

// Red UI Color
var Red = color.RGBA{R: 255, G: 0, B: 0, A: 255}

// White UI Color
var White = color.RGBA{R: 255, G: 255, B: 255, A: 255}

//////////////////////////////////
// Swap and swappiness settings //
//////////////////////////////////

// AvailableSwapSizes A list of swap sizes available to choose from, in GB
var AvailableSwapSizes = []string{"2", "4", "6", "8", "12", "16", "20", "24", "32"}

// AvailableSwappinessOptions A list of swappiness options to choose from, valid range 0-200
var AvailableSwappinessOptions = []string{"0", "1", "10", "25", "50", "60", "75", "90", "100 (Default)", "150", "200"}

// SpaceOverhead The amount of space to keep available above the swapfile size, should prevent boot loops
var SpaceOverhead = 1 * GigabyteMultiplier // 1GB

// GigabyteMultiplier Used to convert gigabytes to bytes
var GigabyteMultiplier = 1024 * 1024 * 1024

////////////////////////
// Game Data settings //
////////////////////////

// LibraryVDFLocation The default location of Steam's library VDF
var LibraryVDFLocation = filepath.Join(HomeDirectory, ".steam/steam/steamapps/libraryfolders.vdf")

// MountDirectory The folder where all external devices are mounts
var MountDirectory = "/run/media"

// SteamDataRoot The default location where Steam keeps compatdata and shadercache
var SteamDataRoot = filepath.Join(HomeDirectory, ".local/share/Steam")

// SteamCompatRoot Generates the full path of the compatdata folder, on SSD
var SteamCompatRoot = filepath.Join(SteamDataRoot, "steamapps/compatdata")

// SteamShaderRoot Generates the full path of the shadercache folder, on SSD
var SteamShaderRoot = filepath.Join(SteamDataRoot, "steamapps/shadercache")

// ExternalDataRoot The location where I'll keep compatdata and shadercache on microSD cards
var ExternalDataRoot = "cryoutilities_steam_data"

// ExternalCompatRoot Generates the full path of the compatdata folder, on microSD
var ExternalCompatRoot = filepath.Join(ExternalDataRoot, "compatdata")

// ExternalShaderRoot Generates the full path of the shadercache folder, on microSD
var ExternalShaderRoot = filepath.Join(ExternalDataRoot, "shadercache")

// SteamApiUrl The URL for the Steam GetAppList URL
var SteamApiUrl = "https://api.steampowered.com/ISteamApps/GetAppList/v0002/"

// SteamGameMaxInteger Anything over this number is presumed to be a Proton version
// Prevents accidental removal of Proton files
var SteamGameMaxInteger = 1000000000
