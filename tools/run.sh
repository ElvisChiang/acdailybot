#!/bin/sh

ssh -n -f $BOTHOST "sh -c 'cd $BOTPATH; nohup ./acdailybot >> log.txt 2>&1 &'"
