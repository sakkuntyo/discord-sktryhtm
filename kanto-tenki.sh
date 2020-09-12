curl https://tenki.jp/forecast/3/16/ | grep forecast-comment | sed -e "s/.*forecast-comment\">//g" -e "s/<\/div>.*//g" -e "s/。/。\\n/g"
