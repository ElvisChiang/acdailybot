#!/bin/sh
# make it in crontab
# # m h  dom mon dow   command
# 0 5 * * * ~/work/acdailybot/reset_hl.sh

CHANNELID=-436800666
echo "remove channelid $CHANNELID highlight:" `date` >> ~/work/acdailybot/log.txt
echo "delete from highlight where channelid = $CHANNELID" | sqlite3 ~/work/acdailybot/acbot.db >> ~/work/acdailybot/log.txt
