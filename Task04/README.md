# Практичне завдання 4. Застосування утиліти `Dive` в процесі CI/CD тестування

Потрібно записати демонстрацію процесу використання утиліти з наступним параметром:  
  `dive --ci --lowestEfficiency=0.9 <image_name>`
для подальшого її використання в процесах автоматизованого тестування образу.

## Сценарій
```bash
asciinema rec -i 1 
make image
dive {sha256}
dive --ci --lowestEfficiency=0.99 {sha256}
# Result:FAIL [Total:3] [Passed:1] [Failed:1] [Warn:0] [Skipped:1]
wget https://raw.githubusercontent.com/wagoodman/dive/main/.data/.dive-ci
dive --ci {sha256}
# Result:PASS [Total:3] [Passed:3] [Failed:0] [Warn:0] [Skipped:0]  {0.99  300kb 0.03}
nano .dive-ci 
dive --ci {sha256}
# Result:FAIL [Total:3] [Passed:0] [Failed:3] [Warn:0] [Skipped:0]
rm .dive-ci
docker images 
docker rmi
nano Makefile 
make dive
# Result:FAIL [Total:3] [Passed:1] [Failed:1] [Warn:0] [Skipped:1]
nano Dockerfile
make dive
# Result:PASS [Total:3] [Passed:2] [Failed:0] [Warn:0] [Skipped:1]
```
## Запис 
![Image](./622097.gif)  
