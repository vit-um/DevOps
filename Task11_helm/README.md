# Task 11: HELM release

1. Створіть новий Helm чарт за допомогою команди (приклад розглядається на Coding Session):

```sh
$ helm create helm
Creating helm
```

2. Підготуйте файл "values.yaml" у директорії чарту, включивши до нього блок:

```yaml
image:
  repository: umanetsvitaliy
  tag: "v1.2.0-660ae95"
  arch: amd64
# Додатково визначте секцію для токену TELE_TOKEN
secret:
  name: "kbot"
  env: "TELE_TOKEN"
  key: "token"
securityContext:
  privileged: true
```
  
3. Відредагуйте файл "deployment.yaml" у каталозі "templates" та додайте блок з посиланням на образ контейнеру:

```yaml
spec:
  template:
    spec:
      containers:
        - name: {{ .Release.Name }}
          image: {{ .Values.image.repository }}/{{ .Chart.Name }}:{{ .Values.image.tag }}-{{ .Values.image.arch | default "amd64"}}
  
# Додатково створіть блок для змінної середовища TELE_TOKEN із застосуванням Kubernetes secret
          env:
          - name: {{ .Values.secret.env }}
            valueFrom:
              secretKeyRef:
                key: {{ .Values.secret.key }}
                name: {{ .Values.secret.name }}
```
4. Запакуйте Helm чарт за допомогою команди:

```sh
$ helm lint ./helm
$ helm package ./helm
```
  
5. Створіть новий реліз GitHub за допомогою інтерактивної команди GitHub CLI (вам може знадобитися GITHUB_TOKEN):

```sh
$ gh auth login  
$ gh release create

```
  
6. Перевірте створений реліз командою:

gh release list
  
7. Додайте до релізу helm пакет:

gh release upload <RELEASE> <CHART_NAME>.tgz
  
8. Протестуйте Helm чарт, встановивши його за допомогою команди:

helm install <CHART_NAME> <CHART_URL>
  
9. Перевірте чи все необхідне вказано в інструкції та чарт встановлено і працює коректно.

10. Після виконання завдання обов'язково перегляньте і протестуйте Helm пакет, щоб переконатися, що все відповідає вимогам і функціонує коректно, додайте URL-адресу до HELM пакету релізу як відповідь.