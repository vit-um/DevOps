# Task14 Міграція Pipeline в gitlab cicd
1. Реєструємось та експортуємо наш проект з GitHUB на [GitLAB](https://gitlab.com/vit-um/kbot)
2. Клонуємо репозиторій на локальну машину `git clone https://gitlab.com/vit-um/kbot.git`
3. На віддаленому репозиторії натискаємо Edit -> Web IDE
2. Створюємо файл `.gitlab-ci.yml` та починаємо описувати в ньому процес CI

sudo apt-get install certbot
export EMAIL="vit@i.ua"   
export DOMAIN="kbot.remote.gitlab.dev"

certbot -d "${DOMAIN}" \
  -m "${EMAIL}" \
  --config-dir ~/.certbot/config \
  --logs-dir ~/.certbot/logs \
  --work-dir ~/.certbot/work \
  --manual \
  --preferred-challenges dns certonly

  

Account registered.
Requesting a certificate for kbot.remote.gitlab.dev

- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
Please deploy a DNS TXT record under the name:

_acme-challenge.kbot.remote.gitlab.dev.

with the following value:

PFnD8RDL5sNS4wg5E34rPAIo0m_5ndPMkBFrGxlNoD4

Before continuing, verify the TXT record has been deployed. Depending on the DNS
provider, this may take some time, from a few seconds to multiple minutes. You can
check if it has finished deploying with aid of online tools, such as the Google
Admin Toolbox: https://toolbox.googleapps.com/apps/dig/#TXT/_acme-challenge.kbot.remote.gitlab.dev.


3. Доопрацювати:
- Запустити пейплайн  
https://www.youtube.com/watch?v=jAIhhULc7YA  
https://www.youtube.com/watch?v=phlsVGysQSw  

