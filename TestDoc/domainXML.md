填充的信息
```
<domain type="kvm">
  <name>testVM</name>
  <uuid>df112352-08e5-47b5-aa28-fdc0affff02c</uuid>
  <memory unit="KiB">2048000</memory>
  <currentMemory unit="KiB">1024000</currentMemory>
  <vcpu placement="static">2</vcpu>
  <os>
    <type arch="x86_64" machine="pc-q35-9.2">hvm</type>
    <boot dev="hd"/>
  </os>
  <features>
    <acpi/>
    <apic/>
  </features>
  <cpu mode="host-passthrough" check="none" migratable="on"/>
  <clock offset='utc'>
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
      <source file='/data/images/rocky9.qcow2'/>
      <target dev="vda" bus="virtio"/>
      <address type="pci" domain="0x0000" bus="0x04" slot="0x00" function="0x0"/>
      <iotune>
        <read_bytes_sec>40971520</read_bytes_sec>    <!-- 1MB/s = 1048576 字节/秒 -->
        <write_bytes_sec>40971520</write_bytes_sec>
        <read_iops_sec>300</read_iops_sec>
        <write_iops_sec>200</write_iops_sec>
      </iotune> 
    </disk>
    <disk type="file" device="disk">    <!-- 数据盘 -->
      <driver name="qemu" type="qcow2" discard="unmap"/>
      <source file='/var/lib/libvirt/images/testVol.qcow2'/>
      <target dev="vdb" bus="virtio"/>
      <address type="pci" domain="0x0000" bus="0x09" slot="0x00" function="0x0"/>
      <iotune>
        <read_bytes_sec>30971520</read_bytes_sec>
        <write_bytes_sec>30971520</write_bytes_sec>
        <read_iops_sec>200</read_iops_sec>
        <write_iops_sec>100</write_iops_sec>
      </iotune> 
    </disk>
    <disk type="file" device="cdrom">
      <driver name="qemu" type="raw"/>
      <source file='/data/images/Rocky-9.2-x86_64-minimal.iso'/>
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
      <mac address='08:1F:2F:80:97:8D'/>
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
    <graphics type="vnc" port='5910' autoport='no' listen="0.0.0.0" passwd='000000'>
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
</domain>
```

---

经过virt-manager之后的最终定义
```
<domain type="kvm">
  <name>testVM</name>
  <uuid>df112352-08e5-47b5-aa28-fdc0affff02c</uuid>
  <memory unit="KiB">2048000</memory>
  <currentMemory unit="KiB">1024000</currentMemory>
  <vcpu placement="static">2</vcpu>
  <os>
    <type arch="x86_64" machine="pc-q35-9.2">hvm</type>
    <boot dev="hd"/>
  </os>
  <features>
    <acpi/>
    <apic/>
  </features>
  <cpu mode="host-passthrough" check="none" migratable="on"/>
  <clock offset="utc">
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
    <disk type="file" device="disk">
      <driver name="qemu" type="qcow2" discard="unmap"/>
      <source file="/data/images/rocky9.qcow2"/>
      <target dev="vda" bus="virtio"/>
      <iotune>
        <read_bytes_sec>40971520</read_bytes_sec>
        <write_bytes_sec>40971520</write_bytes_sec>
        <read_iops_sec>300</read_iops_sec>
        <write_iops_sec>200</write_iops_sec>
      </iotune>
      <address type="pci" domain="0x0000" bus="0x04" slot="0x00" function="0x0"/>
    </disk>
    <disk type="file" device="disk">
      <driver name="qemu" type="qcow2" discard="unmap"/>
      <source file="/var/lib/libvirt/images/testVol.qcow2"/>
      <target dev="vdb" bus="virtio"/>
      <iotune>
        <read_bytes_sec>30971520</read_bytes_sec>
        <write_bytes_sec>30971520</write_bytes_sec>
        <read_iops_sec>200</read_iops_sec>
        <write_iops_sec>100</write_iops_sec>
      </iotune>
      <address type="pci" domain="0x0000" bus="0x09" slot="0x00" function="0x0"/>
    </disk>
    <disk type="file" device="cdrom">
      <driver name="qemu" type="raw"/>
      <source file="/data/images/Rocky-9.2-x86_64-minimal.iso"/>
      <target dev="sda" bus="sata"/>
      <readonly/>
      <address type="drive" controller="0" bus="0" target="0" unit="0"/>
    </disk>
    <controller type="usb" index="0" model="qemu-xhci" ports="15">
      <address type="pci" domain="0x0000" bus="0x02" slot="0x00" function="0x0"/>
    </controller>
    <controller type="pci" index="0" model="pcie-root"/>
    <controller type="pci" index="1" model="pcie-root-port">
      <model name="pcie-root-port"/>
      <target chassis="1" port="0x10"/>
      <address type="pci" domain="0x0000" bus="0x00" slot="0x02" function="0x0" multifunction="on"/>
    </controller>
    <controller type="pci" index="2" model="pcie-root-port">
      <model name="pcie-root-port"/>
      <target chassis="2" port="0x11"/>
      <address type="pci" domain="0x0000" bus="0x00" slot="0x02" function="0x1"/>
    </controller>
    <controller type="pci" index="3" model="pcie-root-port">
      <model name="pcie-root-port"/>
      <target chassis="3" port="0x12"/>
      <address type="pci" domain="0x0000" bus="0x00" slot="0x02" function="0x2"/>
    </controller>
    <controller type="pci" index="4" model="pcie-root-port">
      <model name="pcie-root-port"/>
      <target chassis="4" port="0x13"/>
      <address type="pci" domain="0x0000" bus="0x00" slot="0x02" function="0x3"/>
    </controller>
    <controller type="pci" index="5" model="pcie-root-port">
      <model name="pcie-root-port"/>
      <target chassis="5" port="0x14"/>
      <address type="pci" domain="0x0000" bus="0x00" slot="0x02" function="0x4"/>
    </controller>
    <controller type="pci" index="6" model="pcie-root-port">
      <model name="pcie-root-port"/>
      <target chassis="6" port="0x15"/>
      <address type="pci" domain="0x0000" bus="0x00" slot="0x02" function="0x5"/>
    </controller>
    <controller type="pci" index="7" model="pcie-root-port">
      <model name="pcie-root-port"/>
      <target chassis="7" port="0x16"/>
      <address type="pci" domain="0x0000" bus="0x00" slot="0x02" function="0x6"/>
    </controller>
    <controller type="pci" index="8" model="pcie-root-port">
      <model name="pcie-root-port"/>
      <target chassis="8" port="0x17"/>
      <address type="pci" domain="0x0000" bus="0x00" slot="0x02" function="0x7"/>
    </controller>
    <controller type="pci" index="9" model="pcie-root-port">
      <model name="pcie-root-port"/>
      <target chassis="9" port="0x18"/>
      <address type="pci" domain="0x0000" bus="0x00" slot="0x03" function="0x0"/>
    </controller>
    <controller type="sata" index="0">
      <address type="pci" domain="0x0000" bus="0x00" slot="0x1f" function="0x2"/>
    </controller>
    <controller type="virtio-serial" index="0">
      <address type="pci" domain="0x0000" bus="0x03" slot="0x00" function="0x0"/>
    </controller>
    <interface type="network">
      <mac address="08:1f:2f:80:97:8d"/>
      <source network="default"/>
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
    <graphics type="vnc" port="5910" autoport="no" listen="0.0.0.0" passwd="000000">
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
</domain>
```


