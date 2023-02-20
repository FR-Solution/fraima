## About Project

[![Telegram](https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)](https://t.me/Dobry_kot)
[![Habrahabr](https://img.shields.io/badge/Хабр-5D6D7E?style=for-the-badge&logo=habr&logoColor=white)](https://habr.com/ru/users/dobry-kot/posts)

# Skipping phases

Для того чтоб исключить один из этапов необходимо запустить программу с флагом `skip-phases` после чего перечислить фазы через запятую

Доступные фазы

```
    containerd
    kubelet
    modprob
    sysctl
```

go build -o fraimactl cmd/fraimactl/main.go