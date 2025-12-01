#!/usr/bin/env zsh
set +x

# cd to the stored pwd
cd ${_mittens_PWD}
# and unset our variables if this isn't a nested _post calls
if [[ _mittens_nested ]]; then
  unset _mittens_nested PS4
  return
fi
unset _mittens_git_root _mittens_PWD script_dir PS4
