#!/bin/sh

read -r _input || true

printf '%s\n' '{
  "permission": "ask",
  "user_message": "This git command can publish or close out mission work. Confirm that shipping was explicitly requested.",
  "agent_message": "The ship gate requires explicit user approval before git merge, git push, or git tag."
}'
