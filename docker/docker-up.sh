#!/bin/bash

usage() {
  echo "usage:"
  echo   $0 --storage=db    
  echo   $0 --storage=inmemory
  exit 1
}

if [ $# -eq 0 ]; then
  usage
fi

while [[ $# -gt 0 ]]; do
  case "$1" in
    --storage=db)
      echo "run docker compose..."
      docker-compose up --build
      shift
      ;;
    --storage=inmemory)
      echo "run docker..."
      cd ..
      docker build -t csapp -f docker/Dockerfile .
      docker run -p 8080:8080 csapp  
      shift
      ;;
    *)
      echo "unknown: $1"
      usage
      ;;
  esac
done
