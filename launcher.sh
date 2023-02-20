#!/bin/bash
if [ "$(xdpyinfo | grep dimension | awk '{print $2}' | cut -d 'x' -f2)" -eq "800" ]; then
  export FYNE_SCALE=0.25
fi

"$HOME"/.cryo_utilities/cryo_utilities gui
