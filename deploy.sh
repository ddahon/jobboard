#!/bin/bash
git checkout main && git pull
make server && make scrape
sudo systemctl restart jobboard.service
exit 0
