#!/usr/bin/env bash
version_file="version"

# return default version if the version file is not found
[[ ! -f version ]] && echo "0.0.1" > ${version_file} && echo "0.0.1" && exit

ver=`cat ${version_file}`
major=0
minor=0
build=0

# break down the version number into it's components
regex="([0-9]+).([0-9]+).([0-9]+)"
if [[ $ver =~ $regex ]]; then
  major="${BASH_REMATCH[1]}"
  minor="${BASH_REMATCH[2]}"
  build="${BASH_REMATCH[3]}"
fi

# check paramater to see which number to increment
if [[ "$1" == "build" ]]; then
  build=$(echo $build + 1 | bc)
elif [[ "$1" == "feature" ]]; then
  minor=$(echo $minor + 1 | bc)
  build=0
elif [[ "$1" == "major" ]]; then
  major=$(echo $major+1 | bc)
  minor=0
  build=0
else
  echo "usage: ./version.sh [major/feature/build]"
  exit -1
fi

# echo the new version number
new_version=${major}.${minor}.${build}
echo ${new_version} > ${version_file} && echo ${new_version}