# photo-mover-go

Transfer and organize photos into photo library

```
SD card/
├─ DSC001.jpg
├─ DSC002.jpg
├─ DSC003.jpg
└─ DSC004.jpg
```

becomes

```
photos/
├─ 2026-5-12/
│ ├─ DSC001.jpg
│ └─ DSC002.jpg
└─ 2026-2-9/
  ├─ DSC003.jpg
  └─ DSC004.jpg
```

## Build it

```sh
go build
```

## Run it

```sh
./photo-mover-go -source DCIM/ -destination ~/Pictures
```

I mainly just wanted to try making something with Go. So far I like it!
