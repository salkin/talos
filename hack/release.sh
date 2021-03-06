#!/bin/bash

set -e

function changelog {
  if [ $# -ne 1 ]; then
    echo 1>&2 "Usage: $0 changelog <file>"
    exit 1
  fi

  conventional-changelog -p angular -i $1 -s
}

function cherry-pick {
  if [ $# -ne 2 ]; then
    echo 1>&2 "Usage: $0 cherry-pick <commit> <branch>"
    exit 1
  fi

  git checkout $2
  git fetch
  git rebase upstream/$2
  git cherry-pick -x $1
}

function commit {
  if [ $# -ne 1 ]; then
    echo 1>&2 "Usage: $0 commit <tag>"
    exit 1
  fi

  git commit -s -m "chore: prepare release $1" -m "This is the official $1 release."
}

if declare -f "$1" > /dev/null
then
  cmd="$1"
  shift
  $cmd "$@"
else
  cat <<EOF
Usage:
  commit:       Create the official release commit message.
  cherry-pick:   Cherry-pick a commit into a release branch.
  changelog:    Update the specified CHANGELOG.
EOF

  exit 1
fi
