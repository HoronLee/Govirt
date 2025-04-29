package xmlDefine

const PoolTemplate = `<pool type="{{.Type}}">
  <name>{{.Name}}</name>
  <uuid>{{.UUID}}</uuid>
  <target>
    <path>{{.Path}}</path>
  </target>
</pool>`

type PoolTemplateParams struct {
	Type      string `default:"dir"` // 存储池类型 dir
	Name      string // 存储池名称
	UUID      string // 存储池UUID
	Autostart bool   `default:"1"` // 是否自启动
	Path      string // 存储池路径
}

const VolumeTemplate = `<volume>
  <name>{{.Name}}</name>
  <capacity>{{.Capacity}}</capacity>
  <allocation>{{.Allocation}}</allocation>
  <target>
    <format type="{{.Type}}"/>
  </target>
  {{if notEmpty .BackingStorePath}}
  <backingStore>
    <path>{{.BackingStorePath}}</path>
  </backingStore>
  {{end}}
</volume>`

type VolumeTemplateParams struct {
	Name             string // 卷名称
	Capacity         uint64
	Allocation       int64  `default:"0"`     // 默认不立刻分配卷
	Type             string `default:"qcow2"` // 卷格式 qcow2
	BackingStorePath string // 卷的后备存储路径
}

const NetworkTemplate = `<network>
  <name>{{.Name}}</name>
  <uuid>{{.UUID}}</uuid>
  <forward mode="{{.ForwardMode}}"/>
  <domain name="{{.DomainName}}"/>
  <ip address="{{.IPAddress}}" netmask="{{.NetMask}}">
    <dhcp>
      <range start="{{.DhcpStart}}" end="{{.DhcpEnd}}"/>
    </dhcp>
  </ip>
</network>`

type NetworkTemplateParams struct {
	Name        string // 网络名称
	UUID        string // 网络UUID
	DomainName  string // 域名，一般与Name相同
	Autostart   bool   `default:"1"`   // 是否自启动
	ForwardMode string `default:"nat"` // 转发模式
	IPAddress   string // IP地址
	NetMask     string `default:"255.255.255.0"` // 子网掩码
	DhcpStart   string // DHCP起始地址
	DhcpEnd     string // DHCP结束地址
}

// DomainTemplate 保存域定义的XML模板
const DomainTemplate = `<domain type="kvm">
  <name>{{.Name}}</name>
  <uuid>{{.UUID}}</uuid>
  <memory unit="KiB">{{.MaxMem}}</memory>
  <currentMemory unit="KiB">{{.CurrentMem}}</currentMemory>
  <vcpu placement="static">{{.VCPU}}</vcpu>
  <os>
    <type arch="{{.Arch}}" machine="pc-q35-9.2">hvm</type>
    {{- if notEmpty .CDRomSource -}}
        <boot dev="{{if notEmpty .BootDev}}{{.BootDev}}{{else}}cdrom{{end}}"/>
    {{- else -}}
        <boot dev="{{.BootDev}}"/>
    {{- end -}}
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
      <iotune>
        <read_bytes_sec>{{.OsDiskRrate}}</read_bytes_sec>    <!-- 1MB/s = 1048576 字节/秒 -->
        <write_bytes_sec>{{.OsDiskWrate}}</write_bytes_sec>
        <read_iops_sec>{{.OsDiskRIOPS}}</read_iops_sec>
        <write_iops_sec>{{.OsDiskWIOPS}}</write_iops_sec>
      </iotune> 
    </disk>
    <!-- 数据盘是可选的 -->
    {{if notEmpty .DataDiskSource}}
    <disk type="file" device="disk">    <!-- 数据盘 -->
      <driver name="qemu" type="qcow2" discard="unmap"/>
      <source file='{{.DataDiskSource}}'/>
      <target dev="vdb" bus="virtio"/>
      <iotune>
        <read_bytes_sec>{{.DataDiskRrate}}</read_bytes_sec>
        <write_bytes_sec>{{.DataDiskWrate}}</write_bytes_sec>
        <read_iops_sec>{{.DataDiskRIOPS}}</read_iops_sec>
        <write_iops_sec>{{.DataDiskWIOPS}}</write_iops_sec>
      </iotune> 
    </disk>
    {{end}}
    <!-- CDROM是可选的 -->
    {{if notEmpty .CDRomSource}}
    <disk type="file" device="cdrom">
      <driver name="qemu" type="raw"/>
      <source file='{{.CDRomSource}}'/>
      <target dev="sda" bus="sata"/>
      <readonly/>
      <address type="drive" controller="0" bus="0" target="0" unit="0"/>
    </disk>
    {{end}}
    <!-- 控制器，待研究哪些是必须的，以及其中的关系 -->
    <controller type="usb" index="0" model="qemu-xhci" ports="15">
    </controller>
    <controller type="pci" index="0" model="pcie-root"/>
    <controller type="sata" index="0">
    </controller>
    <controller type="virtio-serial" index="0">
    </controller>
    <!-- 网卡 -->
    <interface type="network">
      <source network="{{.ExterName}}"/>
      <mac address='{{.ExterMac}}'/>
      <model type="virtio"/>
    </interface>
    <interface type="network">
      <source network="{{.InterName}}"/>
      <mac address='{{.InterMac}}'/>
      <model type="virtio"/>
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
    <graphics type="vnc" {{if ne .VncPort "-1"}}port='{{.VncPort}}'{{else}}autoport='yes'{{end}} listen="0.0.0.0"{{if notEmpty .VncPasswd}} passwd='{{.VncPasswd}}'{{end}}>
        <listen type="address" address="0.0.0.0"/>
    </graphics>
    <audio id="1" type="none"/>
    <video>
      <model type="vga" vram="16384" heads="1" primary="yes"/>
      <address type="pci" domain="0x0000" bus="0x00" slot="0x01" function="0x0"/>
    </video>
    <watchdog model="itco" action="reset"/>
    <memballoon model="virtio">
    </memballoon>
    <rng model="virtio">
      <backend model="random">/dev/urandom</backend>
    </rng>
  </devices>
</domain>
`

// DomainTemplateParams 定义了虚拟机XML模板中的所有变量
type DomainTemplateParams struct {
	// 基本配置
	Name        string `default:"vm"`
	UUID        string // UUID，留空自动生成
	Autostart   bool   `default:"1"`       // 是否自启动，默认自启动
	MaxMem      int    `default:"1048576"` // 单位KiB，默认1GB
	CurrentMem  int    `default:"1048576"` // 单位KiB，默认1GB
	VCPU        int    `default:"1"`       // 虚拟CPU数量，默认1个
	Arch        string `default:"x86_64"`  // 架构
	ClockOffset string `default:"utc"`     // 时钟偏移
	BootDev     string `default:"hd"`      // 启动设备，默认硬盘启动

	// 系统盘配置
	OsDiskSource string // 必须提供
	OsDiskRrate  int64  `default:"104857600"` // 默认100MB/s
	OsDiskWrate  int64  `default:"104857600"` // 默认100MB/s
	OsDiskRIOPS  int64  `default:"1000"`      // 默认1000 IOPS
	OsDiskWIOPS  int64  `default:"1000"`      // 默认1000 IOPS

	// 数据盘配置
	DataDiskSource string
	DataDiskRrate  int64 `default:"104857600"` // 默认100MB/s
	DataDiskWrate  int64 `default:"104857600"` // 默认100MB/s
	DataDiskRIOPS  int64 `default:"1000"`      // 默认1000 IOPS
	DataDiskWIOPS  int64 `default:"1000"`      // 默认1000 IOPS

	// CDROM配置
	CDRomSource string // 不设默认值，可为空

	// 网络配置
	ExterName string `config:"network.external.name"` // 外网网卡名称，必需
	ExterMac  string // 留空自动生成
	InterName string `config:"network.internal.name"` // 内网网卡名称，必需
	InterMac  string // 留空自动生成

	// VNC配置
	VncPort       string `default:"-1"`  // 默认VNC端口
	IsVncAutoPort string `default:"yes"` // 默认自动分配端口
	VncPasswd     string `default:""`    // 默认无密码

	// 外部设置
	OsImageID   string // 镜像标识符，可以是名称或UUID
	OsImageType string // 镜像类型，qcow2或raw
	OsCapacity  uint64 // 系统盘容量，单位GB
}
