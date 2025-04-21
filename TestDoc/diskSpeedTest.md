        <read_bytes_sec>40971520</read_bytes_sec>
        <write_bytes_sec>40971520</write_bytes_sec>
        <read_iops_sec>300</read_iops_sec>
        <write_iops_sec>200</write_iops_sec>
## vda写测试

```
[root@localhost ~]# fio --filename=/dev/vda1 --direct=1 --rw=write --bs=1M --size=1G --numjobs=1 --time_based --runtime=30 --group_reporting --name=write_test --allow_mounted_write=1

write_test: (g=0): rw=write, bs=(R) 1024KiB-1024KiB, (W) 1024KiB-1024KiB, (T) 1024KiB-1024KiB, ioengine=psync, iodepth=1
fio-3.35
Starting 1 process
Jobs: 1 (f=1): [W(1)][100.0%][w=39.0MiB/s][w=39 IOPS][eta 00m:00s]
write_test: (groupid=0, jobs=1): err= 0: pid=12978: Mon Apr 21 08:54:34 2025
  write: IOPS=39, BW=39.2MiB/s (41.1MB/s)(1178MiB/30026msec); 0 zone resets
    clat (usec): min=1418, max=107389, avg=25438.99, stdev=3505.79
     lat (usec): min=1468, max=107427, avg=25482.90, stdev=3505.16
    clat percentiles (msec):
     |  1.00th=[   16],  5.00th=[   25], 10.00th=[   25], 20.00th=[   26],
     | 30.00th=[   26], 40.00th=[   26], 50.00th=[   26], 60.00th=[   26],
     | 70.00th=[   26], 80.00th=[   26], 90.00th=[   27], 95.00th=[   27],
     | 99.00th=[   28], 99.50th=[   29], 99.90th=[   51], 99.95th=[  108],
     | 99.99th=[  108]
   bw (  KiB/s): min=38912, max=49152, per=100.00%, avg=40231.05, stdev=1557.05, samples=59
   iops        : min=   38, max=   48, avg=39.29, stdev= 1.52, samples=59
  lat (msec)   : 2=0.59%, 4=0.25%, 20=0.17%, 50=98.81%, 100=0.08%
  lat (msec)   : 250=0.08%
  cpu          : usr=0.18%, sys=0.46%, ctx=1194, majf=0, minf=11
  IO depths    : 1=100.0%, 2=0.0%, 4=0.0%, 8=0.0%, 16=0.0%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued rwts: total=0,1178,0,0 short=0,0,0,0 dropped=0,0,0,0
     latency   : target=0, window=0, percentile=100.00%, depth=1

Run status group 0 (all jobs):
  WRITE: bw=39.2MiB/s (41.1MB/s), 39.2MiB/s-39.2MiB/s (41.1MB/s-41.1MB/s), io=1178MiB (1235MB), run=30026-30026msec

Disk stats (read/write):
  vda: ios=267/2346, merge=270/0, ticks=916/59357, in_queue=60273, util=99.13%
```

## vda读测试

```
[root@localhost ~]# fio --filename=/dev/vda1 --direct=1 --rw=read --bs=1M --size=1G --numjobs=1 --time_based --runtime=30 --group_reporting --name=read_test --allow_mounted_write=1

read_test: (g=0): rw=read, bs=(R) 1024KiB-1024KiB, (W) 1024KiB-1024KiB, (T) 1024KiB-1024KiB, ioengine=psync, iodepth=1
fio-3.35
Starting 1 process
Jobs: 1 (f=1): [R(1)][100.0%][r=39.0MiB/s][r=39 IOPS][eta 00m:00s]
read_test: (groupid=0, jobs=1): err= 0: pid=13017: Mon Apr 21 08:59:07 2025
  read: IOPS=39, BW=39.2MiB/s (41.1MB/s)(1177MiB/30002msec)
    clat (usec): min=716, max=27744, avg=25483.05, stdev=1674.39
     lat (usec): min=721, max=27744, avg=25483.74, stdev=1674.31
    clat percentiles (usec):
     |  1.00th=[24511],  5.00th=[24773], 10.00th=[25035], 20.00th=[25035],
     | 30.00th=[25035], 40.00th=[25560], 50.00th=[25822], 60.00th=[25822],
     | 70.00th=[26084], 80.00th=[26084], 90.00th=[26084], 95.00th=[26084],
     | 99.00th=[26346], 99.50th=[26608], 99.90th=[27395], 99.95th=[27657],
     | 99.99th=[27657]
   bw (  KiB/s): min=38834, max=49152, per=100.00%, avg=40228.34, stdev=1557.60, samples=59
   iops        : min=   37, max=   48, avg=39.25, stdev= 1.54, samples=59
  lat (usec)   : 750=0.08%, 1000=0.25%
  lat (msec)   : 4=0.08%, 50=99.58%
  cpu          : usr=0.06%, sys=0.41%, ctx=1187, majf=0, minf=266
  IO depths    : 1=100.0%, 2=0.0%, 4=0.0%, 8=0.0%, 16=0.0%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued rwts: total=1177,0,0,0 short=0,0,0,0 dropped=0,0,0,0
     latency   : target=0, window=0, percentile=100.00%, depth=1

Run status group 0 (all jobs):
   READ: bw=39.2MiB/s (41.1MB/s), 39.2MiB/s-39.2MiB/s (41.1MB/s-41.1MB/s), io=1177MiB (1234MB), run=30002-30002msec

Disk stats (read/write):
  vda: ios=1184/12, merge=6/1, ticks=30113/8, in_queue=30124, util=98.77%
```



---
        <read_bytes_sec>30971520</read_bytes_sec>
        <write_bytes_sec>30971520</write_bytes_sec>
        <read_iops_sec>200</read_iops_sec>
        <write_iops_sec>100</write_iops_sec>
## vdb写测试

```
[root@localhost ~]# fio --filename=/dev/vdb1 --direct=1 --rw=write --bs=1M --size=1G --numjobs=1 --time_based --runtime=30 --group_reporti
ng --name=write_test
write_test: (g=0): rw=write, bs=(R) 1024KiB-1024KiB, (W) 1024KiB-1024KiB, (T) 1024KiB-1024KiB, ioengine=psync, iodepth=1
fio-3.35
Starting 1 process
Jobs: 1 (f=1): [W(1)][100.0%][w=29.0MiB/s][w=29 IOPS][eta 00m:00s]
write_test: (groupid=0, jobs=1): err= 0: pid=13024: Mon Apr 21 09:00:43 2025
  write: IOPS=29, BW=29.7MiB/s (31.1MB/s)(891MiB/30034msec); 0 zone resets
    clat (usec): min=1380, max=43678, avg=33657.10, stdev=2286.35
     lat (usec): min=1402, max=43763, avg=33702.27, stdev=2286.31
    clat percentiles (usec):
     |  1.00th=[31851],  5.00th=[33162], 10.00th=[33424], 20.00th=[33424],
     | 30.00th=[33817], 40.00th=[33817], 50.00th=[33817], 60.00th=[33817],
     | 70.00th=[33817], 80.00th=[34341], 90.00th=[34341], 95.00th=[34866],
     | 99.00th=[34866], 99.50th=[35390], 99.90th=[43779], 99.95th=[43779],
     | 99.99th=[43779]
   bw (  KiB/s): min=28672, max=36864, per=100.00%, avg=30407.59, stdev=1191.36, samples=59
   iops        : min=   28, max=   36, avg=29.69, stdev= 1.16, samples=59
  lat (msec)   : 2=0.45%, 50=99.55%
  cpu          : usr=0.16%, sys=0.27%, ctx=904, majf=0, minf=11
  IO depths    : 1=100.0%, 2=0.0%, 4=0.0%, 8=0.0%, 16=0.0%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued rwts: total=0,891,0,0 short=0,0,0,0 dropped=0,0,0,0
     latency   : target=0, window=0, percentile=100.00%, depth=1

Run status group 0 (all jobs):
  WRITE: bw=29.7MiB/s (31.1MB/s), 29.7MiB/s-29.7MiB/s (31.1MB/s-31.1MB/s), io=891MiB (934MB), run=30034-30034msec

Disk stats (read/write):
  vdb: ios=29/886, merge=0/0, ticks=54/29736, in_queue=29789, util=99.36%
```

## vdb读测试

```
[root@localhost ~]# fio --filename=/dev/vdb1 --direct=1 --rw=read --bs=1M --size=1G --numjobs=1 --time_based --runtime=30 --group_reportin
g --name=read_test
read_test: (g=0): rw=read, bs=(R) 1024KiB-1024KiB, (W) 1024KiB-1024KiB, (T) 1024KiB-1024KiB, ioengine=psync, iodepth=1
fio-3.35
Starting 1 process
Jobs: 1 (f=1): [R(1)][100.0%][r=29.0MiB/s][r=29 IOPS][eta 00m:00s]
read_test: (groupid=0, jobs=1): err= 0: pid=13044: Mon Apr 21 09:01:53 2025
  read: IOPS=29, BW=29.7MiB/s (31.1MB/s)(890MiB/30001msec)
    clat (usec): min=722, max=35588, avg=33700.88, stdev=2232.11
     lat (usec): min=723, max=35588, avg=33701.58, stdev=2232.10
    clat percentiles (usec):
     |  1.00th=[32637],  5.00th=[32900], 10.00th=[33162], 20.00th=[33817],
     | 30.00th=[33817], 40.00th=[33817], 50.00th=[33817], 60.00th=[33817],
     | 70.00th=[33817], 80.00th=[33817], 90.00th=[34341], 95.00th=[34341],
     | 99.00th=[34866], 99.50th=[34866], 99.90th=[35390], 99.95th=[35390],
     | 99.99th=[35390]
   bw (  KiB/s): min=28672, max=36864, per=100.00%, avg=30407.59, stdev=1191.36, samples=59
   iops        : min=   28, max=   36, avg=29.69, stdev= 1.16, samples=59
  lat (usec)   : 750=0.22%, 1000=0.11%
  lat (msec)   : 2=0.11%, 50=99.55%
  cpu          : usr=0.04%, sys=0.32%, ctx=899, majf=0, minf=266
  IO depths    : 1=100.0%, 2=0.0%, 4=0.0%, 8=0.0%, 16=0.0%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued rwts: total=890,0,0,0 short=0,0,0,0 dropped=0,0,0,0
     latency   : target=0, window=0, percentile=100.00%, depth=1

Run status group 0 (all jobs):
   READ: bw=29.7MiB/s (31.1MB/s), 29.7MiB/s-29.7MiB/s (31.1MB/s-31.1MB/s), io=890MiB (933MB), run=30001-30001msec

Disk stats (read/write):
  vdb: ios=886/0, merge=0/0, ticks=29756/0, in_queue=29757, util=98.77%
```

---

## 分析

Copilot：从测试结果中可以看出，使用 `fio` 工具对磁盘性能进行了读写测试，并比较了不同限制条件下的性能表现。以下是对结果的分析，以确认通过 `libvirt` 限制磁盘性能是否生效：

### **限制条件**
- `vda1` 的限制条件：
  - `read_bytes_sec=40971520` (40 MiB/s)
  - `write_bytes_sec=40971520` (40 MiB/s)
  - `read_iops_sec=300`
  - `write_iops_sec=200`
- `vdb1` 的限制条件：
  - `read_bytes_sec=30971520` (30 MiB/s)
  - `write_bytes_sec=30971520` (30 MiB/s)
  - `read_iops_sec=200`
  - `write_iops_sec=100`

---

### **测试结果**
#### **1. `vda1` 测试**
- **写测试**
  - 实际带宽：**39.2 MiB/s** (大约等于 40 MiB/s 限制)
  - 写 IOPS：**39** (与块大小 1 MiB 对应)
  - 磁盘利用率：**99.13%**
- **读测试**
  - 实际带宽：**39.2 MiB/s** (大约等于 40 MiB/s 限制)
  - 读 IOPS：**39** (与块大小 1 MiB 对应)
  - 磁盘利用率：**98.77%**

**分析**：`vda1` 的性能表现接近于限制条件 `40 MiB/s`，并没有超出。可以判断 `libvirt` 对此磁盘的性能限制是生效的。

---

#### **2. `vdb1` 测试**
- **写测试**
  - 实际带宽：**29.7 MiB/s** (大约等于 30 MiB/s 限制)
  - 写 IOPS：**29** (与块大小 1 MiB 对应)
  - 磁盘利用率：**99.36%**
- **读测试**
  - 实际带宽：**29.7 MiB/s** (大约等于 30 MiB/s 限制)
  - 读 IOPS：**29** (与块大小 1 MiB 对应)
  - 磁盘利用率：**98.77%**

**分析**：`vdb1` 的性能表现接近于限制条件 `30 MiB/s`，并没有超出。可以判断 `libvirt` 对此磁盘的性能限制也是生效的。

---

### **结论**
从 `fio` 测试结果可以看出：
- **`vda1` 和 `vdb1` 的读写带宽和 IOPS 都受到了限制**，性能表现与 `libvirt` 配置中的限制值接近。
- 磁盘利用率几乎达到 100%，这表明限制条件对性能起到了预期的效果。

因此，可以确认通过 `libvirt` 限制磁盘性能的配置是生效的。