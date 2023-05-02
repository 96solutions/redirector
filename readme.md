# TODO
- Use concurrency + pool of workers for handling Clicks after processing redirect request
  - save to storage (via query builder & multiple ORMs)
    - mongodb
    - mysql
    - pg
    - clickhouse
  - save to queue
    - kafka
    - RabbitMQ
- Add https://github.com/spf13/viper
- Add https://github.com/spf13/cobra
  - serve
  - status
  - update number of workers
  - shutdown
  - log
  - cache / cache reset
- Think about https://github.com/c-bata/go-prompt or try to implement same commands as listed above
- Create multiple server implementations
  - HTTP controller
    - net/http
    - gin
    - echo
    - fiber
    - beego
  - gRPC controller
- Add logs
- Add tracing
- Add health check


## see https://tproger.ru/articles/shift-to-golang-ozon-roadmap/