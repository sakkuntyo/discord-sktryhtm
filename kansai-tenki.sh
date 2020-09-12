curl https://tenki.jp/forecast/6/ | grep forecast-comment -A 90 | grep "</div>" -B 20 -m 1 | sed -e "s/.*forecast-comment\">//g" -e "s/<\/div>.*//g" -e "s/<br>//g"
