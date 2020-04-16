#!/bin/sh
# make it in crontab
# # m h  dom mon dow   command
# 0 5 * * 0 ~/work/acdailybot/reset_turnip.sh

CHANNELID=-436800666
echo "remove channelid $CHANNELID turnip:" `date` >> ~/work/acdailybot/log.txt
echo "delete from turnip where channelid = $CHANNELID" | sqlite3 ~/work/acdailybot/acbot.db >> ~/work/acdailybot/log.txt
