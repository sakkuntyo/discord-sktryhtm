curl https://tenki.jp/lite/ | grep forecast-comment -A 99 -m 1 | grep "</div>" -B 99 | sed -e "s/.*forecast-comment\">//g" -e "s/<\/div>.*//g" -e "s/。/。\\n/g"
