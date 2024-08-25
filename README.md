# mkmyfat

「中身がホスト側から触らずに変わっていくUSBメモリ」のためのFATイメージを作成するツールです。

```
git clone git@github.com:gpioblink/mkmyfat.git
cd mkmyfat
go run main.go create test1.img 128MiB mp4 5 1MiB
```

## Usage

### 「中身がホスト側から触らずに変わっていくUSBメモリ」用のFATイメージの作成

「中身がホスト側から触らずに変わっていくUSBメモリ」に使用するためのイメージを作成できます。
このツールで作成したFATは、全てのファイルが必ず連続したクラスタに保管されます。

そのため、各種管理用のパラメータを一切変えずに、ファイルのコンテンツのアドレスを直接書き換えることが可能です。

このツールでフォーマットしたイメージを、今後作成予定の、書き換えに対応した「Mass Storage Gadget」で用いることで、「中身がホスト側から触らずに変わっていくUSBメモリ」を作成できるようになります。


```
go run main.go create <filename> <size> <fileExt> <eachFileSize> <numberOfFiles>
```

実行イメージ

```
$ go run main.go create test1.img 128MiB mp4 5 1MiB 
imagePath test1.img, fileSize 134217728, fileExt mp4, numOfFiles 5, eachFileSize 1048576 

...

***** Root File List *****
IAM0RDVEMP4[1048576bytes]: 0x102400-0x202400 clus=3
IAM1RDVEMP4[1048576bytes]: 0x202400-0x302400 clus=1027
IAM2RDVEMP4[1048576bytes]: 0x302400-0x402400 clus=2051
IAM3RDVEMP4[1048576bytes]: 0x402400-0x502400 clus=3075
IAM4RDVEMP4[1048576bytes]: 0x502400-0x602400 clus=4099
```

### 新しいファイルの作成

ファイルのない、新規のFAT32のイメージを作成します

```
go run main.go create <filename> <size>
```

例

```
$ go run main.go create test5.img 128MiB
imagePath test5.img, fileSize 134217728 

***** BPB *****
BS_jmpBoot([3]bytes): �X�
BS_OEMName([8]bytes): mkmy.fat
BPB_BytsPerSec(uint16):  0x200
BPB_SecPerClus(uint8):  0x2
BPB_RsvdSecCnt(uint16):  0x20
BPB_NumFATs(uint8):  0x2
BPB_RootEntCnt(uint16):  0x0
BPB_TotSec16(uint16):  0x0
BPB_Media(uint8):  0xf8
BPB_FATSz16(uint16):  0x0
BPB_SecPerTrk(uint16):  0x20
BPB_NumHeads(uint16):  0x8
BPB_HiddSec(uint32):  0x0
BPB_TotSec32(uint32):  0x40000
BPB_FATSz32(uint32):  0x3f8
BPB_ExtFlags(uint16):  0x0
BPB_FSVer(uint16):  0x0
BPB_RootClus(uint32):  0x2
BPB_FSInfo(uint16):  0x1
BPB_BkBootSec(uint16):  0x6
BPB_Reserved([12]bytes): 
BS_DrvNum(uint8):  0x80
BS_Reserved(uint8):  0x0
BS_BootSig(uint8):  0x29
BS_VolID(uint32):  0x8ee5ba36
BS_VolLab([11]bytes): NO NAME    
BS_FilSysType([8]bytes): FAT32   
BS_BootCode32([420]bytes): 
BS_Sign(uint16):  0xaa55

***** FSInfo *****
FSI_LeadSig(uint32):  0x41615252
FSI_Reserved1([480]bytes): 
FSI_StrucSig(uint32):  0x61417272
FSI_Free_Count(uint32):  0xffffffff
FSI_Nxt_Free(uint32):  0xffffffff
FSI_Reserved2([12]bytes): 
FSI_TrailSig(uint32):  0xaa550000

***** FAT *****
models.FAT{0x0:0xffffff8, 0x1:0xfffffff, 0x2:0xffffff8}

***** DIR *****
cluster(slice):  []models.DirectoryEntry{}

```

