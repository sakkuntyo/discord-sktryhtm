curl https://tenki.jp/forecast/8/ | grep forecast-comment -A 99 -m 1 | grep "</div>" -B 99 -m 1 | sed -e "s/.*forecast-comment\">//g" -e "s/<\/div>.*//g" -e "s/。/。\\n/g" | sed "/<br>/d"
