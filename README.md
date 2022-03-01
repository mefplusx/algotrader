```
python3 manage.py runserver
```
Первичная настройка сервера
===========
```
ifconfig | grep "inet " | grep -Fv 127.0.0.1 | awk '{print $2}' 
sudo timedatectl set-timezone Europe/Minsk
```
Синхронизация
===========
Частичная синхронизация
```
python3 main.py --sync ETH-USDT
```
Полная синхронизация всех пар с биржей
```
./syncall.sh start
./syncall.sh stop
```
Боевой бот
===========
```
./bot.sh start BTC-USDT
./bot.sh stop
```
```
./bot.sh startall
./bot.sh stopall
```
Неиронка
===========
Запуск сервера нод
```
./server.sh start
./server.sh stop
```
Запуск нод на машинах
```
./node.sh start
./node.sh stop
```
Генерация диапазонов
- для начала необходимо запустить демон генерации
```
./gentasks.sh --daemon
```
- генерация в 5 шагов с авто записью в конфиг
```
./gentasks.sh {currency pair} {short/long}
./gentasks.sh --all
```
- генерация с указанием шага без авто записи в конфиг
```
./gentasks.sh --manual {step} {currency pair}:{short/long}
```
Детали торговли по истории
- сложный процент
```
go run *.go --detail ETH-USDT
```
- помесячно
```
go run *.go --detail ETH-USDT --m
```
Проверка целосности истории
```
python3 main.py --check
```
Вывод сравнения валют
```
go run *.go --compare --m
```
Загрузить best в конфиг
```
go run *.go --update
```
Подобрать крипту, которая:
- ALLOW_MIN_VALUE = 1500000 "Минимальный объем за 24 часа"
- MIN_DAYS_PRICES = 500 "Есть минимум дней истории"
```
python3 main.py --collect-tasks
```
Перебросить баланс main->trade и наоборот
```
python3 main.py --transfer-balance
```
