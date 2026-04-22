#!/bin/bash
tail -n 20 <(docker logs || journalctl -n 20 || dmesg)
