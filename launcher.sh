!/bin/bash
if [ "$(xrandr --current --verbose | grep '*' | xargs | cut -d ' ' -f1)" = "800x1280" ]; then
  export FYNE_SCALE=0.25
fi

"$HOME"/.cryo_utilities/cryo_utilities gui
