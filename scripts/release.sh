#!/bin/bash

set -euo pipefail

HERE="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

cd $HERE/..

BUMP=${1:-patch}
current_branch=$(git rev-parse --abbrev-ref HEAD)

go get -u github.com/flood-io/gitsem

reset="\033[0m"
red="\033[31m"
green="\033[32m"
yellow="\033[33m"
blue="\033[34m"
echo -n -e $reset

echo -e " ${blue}~~ Flood CLI Release ~~${reset}"
echo -e "${blue}Previewing your ${green}$BUMP${blue} bump:${reset}"

if [[ $current_branch != "master" ]]; then
  echo -e "${red}You are releasing from ${yellow}$current_branch${red} rather than master$reset"
  echo Are you sure this is really what you want to do?
  echo
fi

echo -e ${blue}bump preview:${reset}
gitsem -preview $BUMP

echo "if you're happy with that hit enter, otherwise ctrl-c"
read

echo -e ${green}ok, here we go${reset}

echo -e "- ${blue}bumping${reset}"
gitsem $BUMP

echo -e "- ${blue}pushing${reset}"
git push origin HEAD --tags

echo -e "${green}done${reset}"
