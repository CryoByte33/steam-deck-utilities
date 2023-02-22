#!/bin/bash
if [ "$(xrandr | grep ' connected' | wc -l)" -eq 1 ]; then
  export FYNE_SCALE=0.25
fi

"$HOME"/.cryo_utilities/cryo_utilities gui
