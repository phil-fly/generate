package osinfo

import (
	"fmt"
	"github.com/StackExchange/wmi"
	"net"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

var (
	advapi = syscall.NewLazyDLL("Advapi32.dll")
	kernel = syscall.NewLazyDLL("Kernel32.dll")
)

func GetOSInformation() string {

	osInformation := "OSInfo: "
	osInformation += "\n" + " Boot time: " + getStartTime()
	osInformation += "\n" + " Current user: " + getUserName()
	osInformation += "\n" + " OS: " + runtime.GOOS
	osInformation += "\n" + " System version: " + getSystemVersion()
	osInformation += "\n" + " Motherboard: " + getMotherboardInfo()
	osInformation += "\n" + " Bios info: " +  getBiosInfo()

	osInformation += "\n" + " CPU: " + getCpuInfo()
	osInformation += "\n" + " Memory: " + fmt.Sprintf("%v", getMemory())
	osInformation += "\n" + " Disk: " + fmt.Sprintf("%v", getDiskInfo())
	osInformation += "\n" + " Interfaces: " + fmt.Sprintf("%v", getIntfs())

	return osInformation
}

//开机时间
func getStartTime() string {
	GetTickCount := kernel.NewProc("GetTickCount")
	r, _, _ := GetTickCount.Call()
	if r == 0 {
		return ""
	}
	ms := time.Duration(r * 1000 * 1000)
	return ms.String()
}

//当前用户名
func getUserName() string {
	var size uint32 = 128
	var buffer = make([]uint16, size)
	user := syscall.StringToUTF16Ptr("USERNAME")
	domain := syscall.StringToUTF16Ptr("USERDOMAIN")
	r, err := syscall.GetEnvironmentVariable(user, &buffer[0], size)
	if err != nil {
		return ""
	}
	buffer[r] = '@'
	old := r + 1
	if old >= size {
		return syscall.UTF16ToString(buffer[:r])
	}
	r, err = syscall.GetEnvironmentVariable(domain, &buffer[old], size-old)
	return syscall.UTF16ToString(buffer[:old+r])
}

//系统版本
func getSystemVersion() string {
	version, err := syscall.GetVersion()
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%d.%d (%d)", byte(version), uint8(version>>8), version>>16)
}

type diskusage struct {
	Path  string `json:"path"`
	Total uint64 `json:"total"`
	Free  uint64 `json:"free"`
}

func usage(getDiskFreeSpaceExW *syscall.LazyProc, path string) (diskusage, error) {
	lpFreeBytesAvailable := int64(0)
	var info = diskusage{Path: path}
	diskret, _, err := getDiskFreeSpaceExW.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(info.Path))),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&(info.Total))),
		uintptr(unsafe.Pointer(&(info.Free))))
	if diskret != 0 {
		err = nil
	}
	return info, err
}

//硬盘信息
func getDiskInfo() (infos []diskusage) {
	GetLogicalDriveStringsW := kernel.NewProc("GetLogicalDriveStringsW")
	GetDiskFreeSpaceExW := kernel.NewProc("GetDiskFreeSpaceExW")
	lpBuffer := make([]byte, 254)
	diskret, _, _ := GetLogicalDriveStringsW.Call(
		uintptr(len(lpBuffer)),
		uintptr(unsafe.Pointer(&lpBuffer[0])))
	if diskret == 0 {
		return
	}
	for _, v := range lpBuffer {
		if v >= 65 && v <= 90 {
			path := string(v) + ":"
			if path == "A:" || path == "B:" {
				continue
			}
			info, err := usage(GetDiskFreeSpaceExW, string(v)+":")
			if err != nil {
				continue
			}
			infos = append(infos, info)
		}
	}
	return infos
}

type memoryStatusEx struct {
	cbSize                  uint32
	dwMemoryLoad            uint32
	ullTotalPhys            uint64 // in bytes
	ullAvailPhys            uint64
	ullTotalPageFile        uint64
	ullAvailPageFile        uint64
	ullTotalVirtual         uint64
	ullAvailVirtual         uint64
	ullAvailExtendedVirtual uint64
}

//内存信息
func getMemory() string {
	GlobalMemoryStatusEx := kernel.NewProc("GlobalMemoryStatusEx")
	var memInfo memoryStatusEx
	memInfo.cbSize = uint32(unsafe.Sizeof(memInfo))
	mem, _, _ := GlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memInfo)))
	if mem == 0 {
		return ""
	}
	return fmt.Sprint(memInfo.ullTotalPhys / (1024 * 1024))
}

type intfInfo struct {
	Name string
	Ipv4 []string
	Ipv6 []string
}

//网卡信息
func getIntfs() []intfInfo {
	var intf, err = net.Interfaces()
	if err != nil {
		return []intfInfo{}
	}
	var is = make([]intfInfo, len(intf))
	for i, v := range intf {
		ips, err := v.Addrs()
		if err != nil {
			continue
		}
		is[i].Name = v.Name
		for _, ip := range ips {
			if strings.Contains(ip.String(), ":") {
				is[i].Ipv6 = append(is[i].Ipv6, ip.String())
			} else {
				is[i].Ipv4 = append(is[i].Ipv4, ip.String())
			}
		}
	}
	return is
}

//主板信息
func getCpuInfo() (string) {
	var s = []struct {
		Name string
	}{}

	var cpuinfo = ""
	var err = wmi.Query("SELECT  Name  FROM Win32_Processor WHERE (Name IS NOT NULL)", &s)
	if err != nil {
		return cpuinfo
	}
	cpuinfo = cpuinfo + s[0].Name

	return cpuinfo
}

//主板信息
func getMotherboardInfo() (string) {
	var s = []struct {
		Product string
	}{}
	var s1 = []struct {
		SerialNumber string
	}{}
	var boardinfo = ""
	var err = wmi.Query("SELECT  Product  FROM Win32_BaseBoard WHERE (Product IS NOT NULL)", &s)
	if err != nil {
		return boardinfo
	}
	boardinfo = boardinfo + "Product:" + s[0].Product
	err = wmi.Query("SELECT  SerialNumber  FROM Win32_BaseBoard WHERE (SerialNumber IS NOT NULL)", &s1)
	if err != nil {
		return boardinfo
	}
	boardinfo = boardinfo + " SerialNumber:" + s1[0].SerialNumber

	return boardinfo
}

//BIOS信息
func getBiosInfo() string {
	var s = []struct {
		Name string
	}{}
	err := wmi.Query("SELECT Name FROM Win32_BIOS WHERE (Name IS NOT NULL)", &s) // WHERE (BIOSVersion IS NOT NULL)
	if err != nil {
		return ""
	}
	return s[0].Name
}