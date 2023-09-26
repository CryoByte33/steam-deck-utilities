#!/bin/bash
# Author: CryoByte33 and contributors to the CryoUtilities project

if [ "$(xrandr | grep -c ' connected')" -eq 1 ]; then
  export FYNE_SCALE=0.25
fi

"$HOME/.cryo_utilities/cryo_utilities" gui
