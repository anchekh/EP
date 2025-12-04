``` markdown
> Кластер на GO: Controller + Agents + Payload <
```
```
Данный кластер — отказоустойчивая среда выполнения приложений без контейнеризации. Контроллер управляет и масштабирует реплики сервисов, агенты запускают их на хостах, а payload выполняет рабочую нагрузку. Статус кластера доступен через `/status`, поддерживается восстановление реплик и добавление новых нод.
```
``` markdown
> Структура проекта <
```
``` markdown
EP/
├── agent/
│   └── main.go
├── controller/
│   └── main.go
├── payload/
│   └── main.go
├── scripts/
│   └── install.sh
├── services/
│   ├── agent.service
│   └── controller.service
├── go.mod
└── go.sum
```
``` markdown
> Сборка <
```
``` markdown
### Клонирование репозитория
```
``` bash
git clone https://github.com/anchekh/EP
cd EP
```
``` markdown
### Controller
```
``` bash
cd ~/EP/controller
go build -o controller_bin main.go
```
``` markdown
### Agent
```
``` bash
cd ~/EP/agent
go build -o agent_bin main.go
```
``` markdown
### Payload
```
``` bash
cd ~/EP/payload
go build -o payload_bin main.go
```
``` markdown
> Запуск кластера <
```
``` markdown
### Controller (вм 1)
```
``` bash
cd ~/EP/controller
./controller_bin
```
``` markdown
### Agent (вм 2)
```
``` bash
cd ~/EP/agent
./agent_bin
```
``` markdown
### Agent(2) (вм 3) (дополнительный агент)
```
``` bash
cd ~/EP/agent
./agent_bin
```
``` markdown
> Проверка работы <
```
``` markdown
### Статус кластера
```
``` bash
curl http://CONTROLLER_HOST:9000/status
```
``` markdown
### Масштабирование реплик
```
``` bash
curl -X POST -d '{"replicas":3}' http://CONTROLLER_HOST:9000/scale  # 3 реплики
curl -X POST -d '{"replicas":1}' http://CONTROLLER_HOST:9000/scale  # 1 реплика
```
``` markdown
### Остановка реплики (автовосстановление)
```
``` bash
ps aux | grep payload_bin
sudo kill -9 <PID реплики>
curl http://CONTROLLER_HOST:9000/status
```
``` markdown
### Проверка работы через прокси
```
``` bash
curl http://CONTROLLER_HOST:9000/proxy
curl http://CONTROLLER_HOST:9000/proxy
```
