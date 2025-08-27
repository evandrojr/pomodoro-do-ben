#!/bin/bash

# Create $HOME/.local/bin if it doesn't exist
mkdir -p $HOME/.local/bin

# Copy executable to $HOME/.local/bin
cp pomodoro-do-ben $HOME/.local/bin/

# Copy desktop file to $HOME/.local/share/applications
cp pomodoro.desktop $HOME/.local/share/applications/

# Create icons of different sizes
for size in 16 24 32 48 64 128 256 512; do
  mkdir -p $HOME/.local/share/icons/hicolor/${size}x${size}/apps
  convert pomodoro-do-ben.png -resize ${size}x${size} $HOME/.local/share/icons/hicolor/${size}x${size}/apps/pomodoro-do-ben.png
done

# Update icon cache
gtk-update-icon-cache -f -t $HOME/.local/share/icons/hicolor
