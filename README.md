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

## Get it

3 options:

**Download it** from the Releases tab, unzip it and run it like any other program

**Install it** (assuming you have go installed)
```sh
go install github.com/BradPatras/photo-mover-go
```

**Build it yourself** (assuming you have go installed)
```sh
# after cloning this repo
cd photo-mover-go
go install
```

## Run it

```sh
# call it without arguments/flags to get the interactive experience
photo-mover-go

# or use the `move` argument and pass in the source/destination via flags
photo-mover-go move -source DCIM/ -destination ~/Pictures
```

I mainly just wanted to try making something with Go. So far I like it!
