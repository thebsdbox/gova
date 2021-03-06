package main

import (
	"encoding/xml"
	"fmt"
)

// These define the Virtual Hardware types as per the OVF specifcation
const (
	HardwareCPU    = "3"
	HardwareMEM    = "4"
	HardwareIDE    = "5"
	HardwareSCSI   = "6"
	HardwareNet    = "10"
	HardwareFloppy = "14"
	HardwareCDROM  = "15" //16 is also supported o_O
	HardwareDisk   = "17"
	HardwareUSB    = "23"
)

// Envelope : is the parent XML and holds all information about a VM
type Envelope struct {
	XMLName xml.Name       `xml:"Envelope"`
	BuildID string         `xml:"vmwbuildId,attr"`
	XMLNS   string         `xml:"xmlns,attr"`
	CIM     string         `xml:"xmlns:cim,attr"`
	OVF     string         `xml:"xmlns:ovf,attr"`
	RASD    string         `xml:"xmlns:rasd,attr"`
	VMW     string         `xml:"xmlns:vmw,attr"`
	VSSD    string         `xml:"xmlns:vssd,attr"`
	XSI     string         `xml:"xmlns:xsi,attr"`
	File    []References   `xml:"References>File"`
	Disk    []*DiskSection `xml:"DiskSection,omitempty"`
	Net     *NetSection    `xml:"NetworkSection,omitempty"`
	VM      VirtualSystem  `xml:"VirtualSystem"`
}

// References : Contains the references to additional files
type References struct {
	OVFHREF string `xml:"ovf:href,attr"`
	OVFID   string `xml:"ovf:id,attr"`
	OVFSIZE string `xml:"ovf:size,attr"`
}

// DiskSection : Defines the disks attached to the VM
type DiskSection struct {
	DiskInfo string `xml:"Info"`
	Disk     []DiskDetails
}

// DiskDetails : Details all of the parts of a disk
type DiskDetails struct {
	OVFCAP     string `xml:"ovf:capacity,attr"`
	OVCAPUNITS string `xml:"ovf:capacityAllocationUnits,attr"`
	OVFDISKID  string `xml:"ovf:diskId,attr"`
	OVFFILEREF string `xml:"ovf:fileRef,attr"`
	OVFFORMAT  string `xml:"ovf:format,attr"`
	OVFPOPSIZE string `xml:"ovf:populatedSize,attr"`
}

// NetSection : Details the networking configuration
type NetSection struct {
	NetInfo string `xml:"Info"`
	Network struct {
		NetName string `xml:"ovf:name,attr"`
		NetDesc string `xml:"Description"`
	} `xml:"Network"`
}

// VirtualSystem : The overall struct that details the Virtual Machine
type VirtualSystem struct {
	VID        string `xml:"ovf:id,attr"`
	VInfo      string `xml:"Info"`
	VName      string `xml:"Name"`
	VOSSection struct {
		VOSID   string `xml:"ovf:id,attr"`
		VOSType string `xml:"ovf:osType,attr"`
		VOSInfo string `xml:"Info"`
		VOSDesc string `xml:"Description,omitempty"`
	} `xml:"OperatingSystemSection"`
	VHardware VirtualHardware `xml:"VirtualHardwareSection"`
	// AnnotationSection struct {
	// 	AnnotationRequired string `xml:"ovf:required,attr,omitempty"`
	// 	AnnotationInfo     string `xml:"Info"`
	// 	AnnotationText     string `xml:"Annotation"`
	// } `xml:"AnnotationSection,omitempty"`
}

// VirtualHardware : The overall struct that details the Virtual Machine
type VirtualHardware struct {
	VHWInfo   string `xml:"Info"`
	VHWSystem struct {
		VHWSystemName string `xml:"vssd:ElementName"`
		VHWInstanceID string `xml:"vssd:InstanceID"`
		VHWSystemID   string `xml:"vssd:VirtualSystemIdentifier"`
		VHWSystemType string `xml:"vssd:VirtualSystemType"`
	} `xml:"System"`
	VHWItem []VirtualHardwareItem `xml:"Item"`
}

// VirtualHardwareItem : The overall struct that details the Virtual Machine
type VirtualHardwareItem struct {
	VHWRequired            string `xml:"ovf:required,attr,omitempty"`
	VHWAddress             string `xml:"rasd:Address,omitempty"`
	VHWAddressOnParent     string `xml:"rasd:AddressOnParent,omitempty"`
	VHWAllocationUnits     string `xml:"rasd:AllocationUnits,omitempty"`
	VHWAutomaticAllocation string `xml:"rasd:AutomaticAllocation,omitempty"`
	VHWConnection          string `xml:"rasd:Connection,omitempty"`
	VHWDescription         string `xml:"rasd:Description,omitempty"`
	VHWElementName         string `xml:"rasd:ElementName,omitempty"`
	VHWHostResource        string `xml:"rasd:HostResource,omitempty"`
	VHWInstanceID          string `xml:"rasd:InstanceID,omitempty"`
	VHWParent              string `xml:"rasd:Parent,omitempty"`
	VHWResourceSubType     string `xml:"rasd:ResourceSubType,omitempty"`
	VHWResourceType        string `xml:"rasd:ResourceType,omitempty"`
	VHWVirtualQuantity     string `xml:"rasd:VirtualQuantity,omitempty"`
}

func newDMTFEnvelope() Envelope {
	var envelope Envelope
	envelope.BuildID = "linuxkitOVF"
	envelope.XMLNS = "http://schemas.dmtf.org/ovf/envelope/1"
	envelope.CIM = "http://schemas.dmtf.org/wbem/wscim/1/common"
	envelope.OVF = "http://schemas.dmtf.org/ovf/envelope/1"
	envelope.RASD = "http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_ResourceAllocationSettingData"
	envelope.VMW = "http://www.vmware.com/schema/ovf"
	envelope.VSSD = "http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_VirtualSystemSettingData"
	envelope.XSI = "http://www.w3.org/2001/XMLSchema-instance"
	return envelope
}

func addMemoryToVM(hardware *VirtualHardware, memorySize string) {
	var memHardware VirtualHardwareItem
	memHardware.VHWResourceType = HardwareMEM
	memHardware.VHWAllocationUnits = "byte * 2^20"
	memHardware.VHWDescription = "Memory Size"
	memHardware.VHWElementName = fmt.Sprintf("%sMB of memory", memorySize)
	memHardware.VHWVirtualQuantity = memorySize
	// Add an additional count to the InstanceID as 0 is the System section
	memHardware.VHWInstanceID = fmt.Sprintf("%d", len(hardware.VHWItem)+1)
	hardware.VHWItem = append(hardware.VHWItem, memHardware)
}

func addCPUtoVM(hardware *VirtualHardware, cpuCount string) {
	var cpuHardware VirtualHardwareItem
	cpuHardware.VHWResourceType = HardwareCPU
	cpuHardware.VHWAllocationUnits = "hertz * 10^6"
	cpuHardware.VHWDescription = "Number of Virtual CPUs"
	cpuHardware.VHWElementName = fmt.Sprintf("%s virtual CPU(s)", cpuCount)
	cpuHardware.VHWVirtualQuantity = cpuCount
	// Add an additional count to the InstanceID as 0 is the System section
	cpuHardware.VHWInstanceID = fmt.Sprintf("%d", len(hardware.VHWItem)+1)
	hardware.VHWItem = append(hardware.VHWItem, cpuHardware)
}

func addIDEControllerToVM(hardware *VirtualHardware) (controllerID int) {
	var controllerHardware VirtualHardwareItem
	controllerHardware.VHWResourceType = HardwareIDE
	controllerHardware.VHWAddress = "1"
	controllerHardware.VHWDescription = "IDE Controller"
	controllerHardware.VHWElementName = "ideController1"
	// This needs returning so other devices can be attached to it
	controllerID = len(hardware.VHWItem) + 1
	controllerHardware.VHWInstanceID = fmt.Sprintf("%d", controllerID)
	hardware.VHWItem = append(hardware.VHWItem, controllerHardware)
	return controllerID
}

func addCDToController(hardware *VirtualHardware, controllerID int, ovfFileID string) {

	var cdHardware VirtualHardwareItem
	cdHardware.VHWResourceType = HardwareCDROM
	cdHardware.VHWAddressOnParent = "0"
	cdHardware.VHWAutomaticAllocation = "true"
	cdHardware.VHWElementName = "cdrom0"
	// File needs adding to references and disksection
	// Should become something like file1
	cdHardware.VHWParent = fmt.Sprintf("%d", controllerID)
	cdHardware.VHWHostResource = fmt.Sprintf("ovf:/file/%s", ovfFileID)
	cdHardware.VHWInstanceID = fmt.Sprintf("%d", len(hardware.VHWItem)+1)
	hardware.VHWItem = append(hardware.VHWItem, cdHardware)
}

func addSCSIControllertoVM(hardware *VirtualHardware) (controllerID int) {
	var controllerHardware VirtualHardwareItem
	controllerHardware.VHWResourceType = HardwareSCSI
	// Perhaps support more controller types in the future
	controllerHardware.VHWResourceSubType = "lsilogic"
	controllerHardware.VHWAddress = "0"
	controllerHardware.VHWDescription = "SCSI Controller"
	controllerHardware.VHWElementName = "scsiController0"
	// This needs returning so other devices can be attached to it
	controllerID = len(hardware.VHWItem) + 1
	controllerHardware.VHWInstanceID = fmt.Sprintf("%d", controllerID)
	hardware.VHWItem = append(hardware.VHWItem, controllerHardware)
	return controllerID
}

func addNicToVM(hardware *VirtualHardware, networkName string) {
	var networkHardware VirtualHardwareItem
	networkHardware.VHWResourceType = HardwareNet
	// Add support for VMXNET3 etc. in future
	networkHardware.VHWResourceSubType = "E1000"
	networkHardware.VHWAddressOnParent = "2"
	networkHardware.VHWAutomaticAllocation = "true"
	networkHardware.VHWConnection = networkName
	networkHardware.VHWDescription = fmt.Sprintf("E1000 ethernet adapter on %s", networkName)
	networkHardware.VHWElementName = "ethernet0"
	networkHardware.VHWInstanceID = fmt.Sprintf("%d", len(hardware.VHWItem)+1)
	hardware.VHWItem = append(hardware.VHWItem, networkHardware)
}

func appendFilesToReferences(references *Envelope, ovfHref string, ovfID string, ovfSize string) {
	var newFile References
	newFile.OVFHREF = ovfHref
	newFile.OVFID = ovfID
	newFile.OVFSIZE = ovfSize
	references.File = append(references.File, newFile)
}

func appdendDiskToDiskSection(disksection *DiskSection, ovfCapacity string, ovfCapUnits string, ovfDiskID string, ovfFileRef string, ovfFormat string, ovfPopSize string) {
	var newDisk DiskDetails
	newDisk.OVFCAP = ovfCapacity
	newDisk.OVCAPUNITS = ovfCapUnits
	newDisk.OVFDISKID = ovfDiskID
	newDisk.OVFFILEREF = ovfFileRef
	newDisk.OVFFORMAT = ovfFormat
	newDisk.OVFPOPSIZE = ovfPopSize
	disksection.Disk = append(disksection.Disk, newDisk)
}
