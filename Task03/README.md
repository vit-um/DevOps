# Практичне завдання 3. Створення специфікації контейнера та процес налаштування мережі.

1. Дії виконуються в [google cloud shell](https://shell.cloud.google.com/) 
2. Встановлюємо утиліту для запису дій в термінальній сесії. Та починаємо запис.  
```bash
mkdir demo && cd demo
sudo apt-get install asciinema
asciinema rec -i 1
```  
3. Створюємо файл специфікації для ініціалізації контейнеру `runc spec`

```bash
runc spec
ls -l
nano config.json
# "path": "/var/run/netns/runc"
sudo bash
brctl addbr runc0
ip link set runc0 up
ip addr add 192.168.0.1/24 dev runc0
ip a show runc0
ip link add name veth-host type veth peer name veth-guest
ip a show veth-host
ip link set veth-host up
brctl show runc0
brctl addif runc0 veth-host
brctl show runc0

ip netns add runc
ip netns ls

# та namespace runc, що ми вказували у файлі специфікації
 ip link set veth-guest netns runc

# за допомогою 'netns exec' виконуємо налаштування саме в namespace. Вкажемо ім'я та налаштуємо інтерфейс eth1
 ip netns exec runc ip link set veth-guest name eth1

# призначимо інтерфейсу в контейнері IP адресу 
 ip netns exec runc ip addr add 192.168.0.2/24 dev eth1
# піднімаємо link
 ip netns exec runc ip link set eth1 up
# додамо маршрут за замовчуванням через хост інтерфейс 
 ip netns exec runc ip route add default via 192.168.0.1
 exit
```