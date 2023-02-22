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
