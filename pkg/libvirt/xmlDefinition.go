package libvirt

// domainTemplate 保存域定义的XML模板
const domainTemplate = `<domain type="kvm">
  <name>{{.Name}}</name>
  <uuid>{{.UUID}}</uuid>
  <memory unit="KiB">{{.MaxMem}}</memory>
  <currentMemory unit="KiB">{{.CurrentMem}}</currentMemory>
  <vcpu placement="static">{{.VCPU}}</vcpu>
  <os>
    <type arch="{{.Arch}}" machine="pc-q35-9.2">hvm</type>
    <boot dev="hd"/>
  </os>
  <features>
    <acpi/>
    <apic/>
  </features>
  <cpu mode="host-passthrough" check="none" migratable="on"/>
  <clock offset='{{.ClockOffset}}'>
    <timer name="rtc" tickpolicy="catchup"/>
    <timer name="pit" tickpolicy="delay"/>
    <timer name="hpet" present="no"/>
  </clock>
  <on_poweroff>destroy</on_poweroff>
  <on_reboot>restart</on_reboot>
  <on_crash>destroy</on_crash>
  <pm>
    <suspend-to-mem enabled="no"/>
    <suspend-to-disk enabled="no"/>
  </pm>
  <devices>
    <emulator>/usr/bin/qemu-system-x86_64</emulator>
    <disk type="file" device="disk">    <!-- 系统盘 -->
      <driver name="qemu" type="qcow2" discard="unmap"/>
      <source file='{{.OsDiskSource}}'/>
      <target dev="vda" bus="virtio"/>
      <address type="pci" domain="0x0000" bus="0x04" slot="0x00" function="0x0"/>
      <iotune>
        <read_bytes_sec>{{.OsDiskRrate}}</read_bytes_sec>    <!-- 1MB/s = 1048576 字节/秒 -->
        <write_bytes_sec>{{.OsDiskWrate}}</write_bytes_sec>
        <read_iops_sec>{{.OsDiskRIOPS}}</read_iops_sec>
        <write_iops_sec>{{.OsDiskWIOPS}}</write_iops_sec>
      </iotune> 
    </disk>
    <disk type="file" device="disk">    <!-- 数据盘 -->
      <driver name="qemu" type="qcow2" discard="unmap"/>
      <source file='{{.DataDiskSource}}'/>
      <target dev="vdb" bus="virtio"/>
      <address type="pci" domain="0x0000" bus="0x09" slot="0x00" function="0x0"/>
      <iotune>
        <read_bytes_sec>{{.DataDiskRrate}}</read_bytes_sec>
        <write_bytes_sec>{{.DataDiskWrate}}</write_bytes_sec>
        <read_iops_sec>{{.DataDiskRIOPS}}</read_iops_sec>
        <write_iops_sec>{{.DataDiskWIOPS}}</write_iops_sec>
      </iotune> 
    </disk>
    <disk type="file" device="cdrom">
      <driver name="qemu" type="raw"/>
      <source file='{{.CDRomSource}}'/>
      <target dev="sda" bus="sata"/>
      <readonly/>
      <!-- boot order='2'/ --> <!-- 表示第几个启动设备 -->
      <address type="drive" controller="0" bus="0" target="0" unit="0"/>
    </disk>
    <!-- 控制器，待研究哪些是必须的，以及其中的关系 -->
    <controller type="usb" index="0" model="qemu-xhci" ports="15">
      <address type="pci" domain="0x0000" bus="0x02" slot="0x00" function="0x0"/>
    </controller>
    <controller type="pci" index="0" model="pcie-root"/>
    <controller type="sata" index="0">
      <address type="pci" domain="0x0000" bus="0x00" slot="0x1f" function="0x2"/>
    </controller>
    <controller type="virtio-serial" index="0">
      <address type="pci" domain="0x0000" bus="0x03" slot="0x00" function="0x0"/>
    </controller>
    <!-- 网卡 -->
    <interface type="network">
      <mac address='{{.NatMac}}'/>
      <source network="default"/>  <!-- 默认NAT -->
      <model type="virtio"/>
      <address type="pci" domain="0x0000" bus="0x07" slot="0x00" function="0x0"/>
    </interface>
    <serial type="pty">
      <target type="isa-serial" port="0">
        <model name="isa-serial"/>
      </target>
    </serial>
    <console type="pty">
      <target type="serial" port="0"/>
    </console>
    <channel type="unix">
      <target type="virtio" name="org.qemu.guest_agent.0"/>
      <address type="virtio-serial" controller="0" bus="0" port="1"/>
    </channel>
    <input type="mouse" bus="ps2"/>
    <input type="keyboard" bus="ps2"/>
    <!-- VNC配置 -->
    <!-- <graphics type="vnc" port="-1" autoport="yes">
      <listen type="address"/>
    </graphics> -->
    <graphics type="vnc" port='{{.VncPort}}' autoport='{{.IsVncAutoPort}}' listen="0.0.0.0" passwd='{{.VncPasswd}}'>
        <listen type="address" address="0.0.0.0"/>
    </graphics>
    <audio id="1" type="none"/>
    <video>
      <model type="vga" vram="16384" heads="1" primary="yes"/>
      <address type="pci" domain="0x0000" bus="0x00" slot="0x01" function="0x0"/>
    </video>
    <watchdog model="itco" action="reset"/>
    <memballoon model="virtio">
      <address type="pci" domain="0x0000" bus="0x05" slot="0x00" function="0x0"/>
    </memballoon>
    <rng model="virtio">
      <backend model="random">/dev/urandom</backend>
      <address type="pci" domain="0x0000" bus="0x06" slot="0x00" function="0x0"/>
    </rng>
  </devices>
</domain>`

// DomainTemplateParams 定义了虚拟机XML模板中的所有变量
type DomainTemplateParams struct {
	Name        string // 虚拟机名称
	UUID        string // 虚拟机UUID
	MaxMem      int    // 单位KiB
	CurrentMem  int    // 单位KiB
	VCPU        int    // 虚拟CPU数量
	Arch        string // 例如: "x86_64"
	ClockOffset string // 例如: "utc"

	// 系统盘配置
	OsDiskSource string
	OsDiskRrate  int64 // 读取限制，字节/秒
	OsDiskWrate  int64 // 写入限制，字节/秒
	OsDiskRIOPS  int64 // 读取IOPS限制
	OsDiskWIOPS  int64 // 写入IOPS限制

	// 数据盘配置
	DataDiskSource string
	DataDiskRrate  int64
	DataDiskWrate  int64
	DataDiskRIOPS  int64
	DataDiskWIOPS  int64

	// CDROM配置
	CDRomSource string

	// 网络配置
	NatMac string

	// VNC配置
	VncPort       int
	IsVncAutoPort bool
	VncPasswd     string
}
