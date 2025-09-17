#!/usr/bin/env bash
set -euo pipefail

if [ "$#" -lt 1 ]; then
  echo "usage: $0 host:port [host:port ...] -- command [args...]" >&2
  exit 1
fi

# collect services until --
services=()
while [[ $# -gt 0 ]]; do
  case "$1" in
    --)
      shift
      break
      ;;
    *)
      services+=("$1")
      shift
      ;;
  esac
done

if [ "$#" -eq 0 ]; then
  echo "no command specified after --" >&2
  exit 1
fi

for target in "${services[@]}"; do
  host="${target%%:*}"
  port="${target##*:}"
  echo "Waiting for $host:$port ..."
  for i in {1..120}; do
    if (echo > "/dev/tcp/$host/$port") >/dev/null 2>&1; then
      echo "$host:$port is up"
      break
    fi
    sleep 1
    if [ "$i" -eq 120 ]; then
      echo "Timed out waiting for $host:$port" >&2
      exit 1
    fi
  done
done

exec "$@"
