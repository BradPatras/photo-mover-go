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
# Assuming you have go installed already
go install
```

## Run it

```sh
# call it without arguments/flags to get the interactive experience
photo-mover-go

# or pass in the source/destination via flags
photo-mover-go move -source DCIM/ -destination ~/Pictures
```

I mainly just wanted to try making something with Go. So far I like it!
